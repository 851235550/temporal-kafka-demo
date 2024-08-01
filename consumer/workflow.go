package consumer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"go.temporal.io/sdk/workflow"
)

const (
	kafkaTopic     = "temporal-topic"
	kafkaBrokerURL = "localhost:9092"
	groupID        = "consumer-group-1"
)

var (
	childWorkerCnt = 200
)

func SetChildWorkerCnt(cnt int) {
	childWorkerCnt = cnt
}

func ConsumeMsg(_ context.Context) error {
	// Kafka reader configuration
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBrokerURL},
		Topic:   kafkaTopic,
		GroupID: groupID,
	})

	// Ensure to close the reader in case of panic or exit
	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("Error closing Kafka reader: %v", err)
		}
	}()

	msg, err := reader.ReadMessage(context.Background())
	if err != nil {
		log.Printf("Error reading message from Kafka: %v", err)
		return err
	}
	log.Printf("Message received: key=%s value=%s offset=%d", string(msg.Key), string(msg.Value), msg.Offset)

	return nil
}

func CronParentConsumerWorkflow(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)

	childWorkflowResults := make([]workflow.Future, 0, 200)
	for i := 0; i < childWorkerCnt; i++ {
		opts := workflow.ChildWorkflowOptions{
			WorkflowID: fmt.Sprintf("consumer-child-%d", i),
			TaskQueue:  childTaskQueueName,
		}
		childCtx := workflow.WithChildOptions(ctx, opts)
		future := workflow.ExecuteChildWorkflow(childCtx, ChildWorkflow)
		childWorkflowResults = append(childWorkflowResults, future)
	}

	// Wait for all child workflows to complete
	for i, future := range childWorkflowResults {
		err := future.Get(ctx, nil)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to execute consumer child workflow %d: %v", i, err))
			return err
		}
	}

	return nil
}

func ChildWorkflow(ctx workflow.Context) error {
	// Define activity options
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 30,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Execute the activity
	err := workflow.ExecuteActivity(ctx, ConsumeMsg).Get(ctx, nil)
	if err != nil {
		logger := workflow.GetLogger(ctx)
		logger.Error("Failed to execute activity ConsumerMsg", "Error", err)
	}

	return err
}
