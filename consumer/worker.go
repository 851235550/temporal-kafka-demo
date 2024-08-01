package consumer

import (
	"context"
	"log"
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
	w := worker.New(c, parentTaskQueueName, worker.Options{})
	w.RegisterWorkflow(CronParentConsumerWorkflow)
	w.RegisterWorkflow(ChildWorkflow)
	w.RegisterActivity(ConsumeMsg)

	childWorker := worker.New(c, childTaskQueueName, worker.Options{})
	childWorker.RegisterWorkflow(ChildWorkflow)
	childWorker.RegisterActivity(ConsumeMsg)

	go func() {
		err := w.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalf("Unable to start consumer worker: %v", err)
		}
	}()

	go func() {
		err := childWorker.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalf("Unable to start consumer child worker: %v", err)
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
