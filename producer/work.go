package producer

import (
	"context"
	"log"
	"suger/config"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"go.temporal.io/api/enums/v1"
)

const (
	producerWorkflowID  = "kafka-producer-workflow"
	parentTaskQueueName = "producer-task-queue"
	childTaskQueueName  = "producer-child-task-queue"

	cronEveryMim = "*/1 * * * *"
)

func StartWorker(c client.Client) {
	if err := config.TerminateWorker(c, producerWorkflowID); err != nil {
		log.Fatalf("Terminate workflow was error. workflowID: %s, err: %s", producerWorkflowID, err.Error())
		return
	}

	w := worker.New(c, parentTaskQueueName, worker.Options{})
	w.RegisterWorkflow(CronParentProducerWorkflow)
	w.RegisterWorkflow(ChildWorkflow)
	w.RegisterActivity(ProduceMsg)

	go func() {
		err := w.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalf("Unable to start producer worker: %v", err)
		}
	}()

	// Start a workflow execution
	we, err := c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
		ID:                       producerWorkflowID,
		TaskQueue:                parentTaskQueueName,
		WorkflowExecutionTimeout: time.Hour,
		WorkflowIDReusePolicy:    enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE, // Correct usage,
		CronSchedule:             cronEveryMim,                                   // Ensure this schedule is correct
	}, CronParentProducerWorkflow)
	if err != nil {
		log.Fatalf("Failed to start workflow: %v", err)
	}

	log.Printf("Started workflow with ID: %s and RunID: %s\n", we.GetID(), we.GetRunID())
}
