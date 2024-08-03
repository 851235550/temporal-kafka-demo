package consumer

import (
	"context"
	"log"
	"suger/config"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const (
	consumerWorkflowID  = "kafka-consumer-workflow"
	parentTaskQueueName = "consumer-task-queue"
	childTaskQueueName  = "consumer-child-task-queue"

	cronEveryTwoMim = "*/2 * * * *"
)

func StartWorker(c client.Client) {
	if err := config.TerminateWorker(c, consumerWorkflowID); err != nil {
		log.Fatalf("Terminate consumer workflow was error. workflowID: %s, err: %s", consumerWorkflowID, err.Error())
		return
	}

	w := worker.New(c, parentTaskQueueName, worker.Options{})
	w.RegisterWorkflow(CronParentConsumerWorkflow)
	w.RegisterWorkflow(ChildWorkflow)
	w.RegisterActivity(ConsumeMsg)

	go func() {
		err := w.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalf("Unable to start consumer worker: %v", err)
		}
	}()

	// Start a workflow execution
	we, err := c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
		ID:                       consumerWorkflowID,
		TaskQueue:                parentTaskQueueName,
		WorkflowExecutionTimeout: time.Hour,
		CronSchedule:             cronEveryTwoMim, // Ensure this schedule is correct
	}, CronParentConsumerWorkflow)
	if err != nil {
		log.Fatalf("Failed to start workflow: %v", err)
	}

	log.Printf("Started workflow with ID: %s and RunID: %s\n", we.GetID(), we.GetRunID())
}
