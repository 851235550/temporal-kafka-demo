package consumer

import (
	"context"
	"log"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func StartWorker(c client.Client) {
	w := worker.New(c, "consumer-task-queue", worker.Options{})
	w.RegisterWorkflow(CronParentConsumerWorkflow)
	w.RegisterWorkflow(ChildWorkflow)
	w.RegisterActivity(ConsumerMsg)

	childWorker := worker.New(c, "consumer-child-task-queue", worker.Options{})
	childWorker.RegisterWorkflow(ChildWorkflow)
	childWorker.RegisterActivity(ConsumerMsg)

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
		ID:                       "kafka-consumer-workflow",
		TaskQueue:                "consumer-task-queue",
		WorkflowExecutionTimeout: time.Hour,
		CronSchedule:             "*/2 * * * *", // Ensure this schedule is correct
	}, CronParentConsumerWorkflow)
	if err != nil {
		log.Fatalf("Failed to start workflow: %v", err)
	}

	log.Printf("Started workflow with ID: %s and RunID: %s\n", we.GetID(), we.GetRunID())
}
