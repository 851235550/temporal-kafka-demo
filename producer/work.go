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
	w.RegisterWorkflow(CronParentProducerWorkflow)
	w.RegisterWorkflow(ChildWorkflow)
	w.RegisterActivity(ProduceMsg)

	childWorker := worker.New(c, "child-task-queue", worker.Options{})
	childWorker.RegisterWorkflow(ChildWorkflow)
	childWorker.RegisterActivity(ProduceMsg)

	go func() {
		err := w.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalf("Unable to start producer worker: %v", err)
		}
	}()

	go func() {
		err := childWorker.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalf("Unable to start child worker: %v", err)
		}
	}()

	// Start a workflow execution
	we, err := c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
		ID:                       "kafka-produce-workflow",
		TaskQueue:                "producer-task-queue",
		WorkflowExecutionTimeout: time.Hour,
		CronSchedule:             "*/1 * * * *", // Ensure this schedule is correct
	}, CronParentProducerWorkflow)
	if err != nil {
		log.Fatalf("Failed to start workflow: %v", err)
	}

	log.Printf("Started workflow with ID: %s and RunID: %s\n", we.GetID(), we.GetRunID())
}

// func StartWorker(c client.Client) {
// 	w := worker.New(c, "producer-task-queue", worker.Options{})
// 	w.RegisterWorkflow(CronParentProducerWorkflow)
// 	w.RegisterWorkflow(ChildWorkflow)
// 	w.RegisterActivity(ProduceMsg)

// 	go func() {
// 		err := w.Run(worker.InterruptCh())
// 		if err != nil {
// 			log.Fatalf("Unable to start producer worker: %v", err)
// 		}
// 	}()

// 	// Start a workflow execution
// 	we, err := c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
// 		ID:                       "kafka-produce-workflow",
// 		TaskQueue:                "producer-task-queue",
// 		WorkflowExecutionTimeout: time.Hour,
// 		CronSchedule:             "*/1 * * * *",
// 	}, CronParentProducerWorkflow)
// 	if err != nil {
// 		panic(err)
// 	}

// 	log.Printf("Started workflow with ID: %s and RunID: %s\n", we.GetID(), we.GetRunID())
// }
