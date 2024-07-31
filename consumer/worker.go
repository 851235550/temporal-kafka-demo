package consumer

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func StartWorker(c client.Client) {
	w := worker.New(c, "consumer-task-queue", worker.Options{})
	w.RegisterWorkflow(CronConsumerWorkflow)
	w.RegisterWorkflow(ConsumerWorkflow)
	err := w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("Unable to start consumer worker: %v", err)
	}
}
