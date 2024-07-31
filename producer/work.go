package producer

import (
	"context"
	"log"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func StartWorker(c client.Client) {
	w := worker.New(c, "producer-task-queue", worker.Options{})
	w.RegisterWorkflow(CronProducerWorkflow)
	w.RegisterActivity(ProduceMsg)

	go func() {
		err := w.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalf("Unable to start producer worker: %v", err)
		}
	}()

	// Start a workflow execution
	we, err := c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
		ID:                       "kafka-produce-workflow",
		TaskQueue:                "producer-task-queue",
		WorkflowExecutionTimeout: time.Hour,
		CronSchedule:             "*/1 * * * *",
	}, CronProducerWorkflow)
	if err != nil {
		panic(err)
	}

	log.Printf("Started workflow with ID: %s and RunID: %s\n", we.GetID(), we.GetRunID())
}
