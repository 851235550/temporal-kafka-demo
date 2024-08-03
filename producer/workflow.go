package producer

import (
	"context"
	"fmt"
	"log"
	"suger/config"
	"time"

	"github.com/segmentio/kafka-go"
	"go.temporal.io/sdk/workflow"
)

var (
	childWorkerCnt = 200
)

func SetChildWorkerCnt(cnt int) {
	childWorkerCnt = cnt
}

func ProduceMsg(ctx context.Context, msg string) error {
	// Kafka writer configuration
	writer := config.NewKafkaWriter()

	// Close the writer
	defer func() {
		if err := writer.Close(); err != nil {
			log.Printf("Error closing Kafka writer: %v", err)
		}
	}()

	log.Printf("Start produce msg: %s\n", msg)
	message := kafka.Message{
		Value: []byte(msg),
	}
	err := writer.WriteMessages(ctx, message)
	if err != nil {
		log.Printf("Error writing message to Kafka: %v", err)
		return err
	}
	log.Printf("End produce msg: %s\n", msg)

	return nil
}

func CronParentProducerWorkflow(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)

	childWorkflowResults := make([]workflow.Future, 0, childWorkerCnt)
	// Start children workflows asynchronously
	for i := 0; i < childWorkerCnt; i++ {
		// Define child workflow options
		childWorkflowOptions := workflow.ChildWorkflowOptions{
			WorkflowID: fmt.Sprintf("producer-child-%d", i),
			// TaskQueue:  childTaskQueueName,
		}
		childCtx := workflow.WithChildOptions(ctx, childWorkflowOptions)
		future := workflow.ExecuteChildWorkflow(childCtx, ChildWorkflow, fmt.Sprintf("child-%d-msg-%d", i, i))
		childWorkflowResults = append(childWorkflowResults, future)
	}

	// Wait for all child workflows to complete
	for i, future := range childWorkflowResults {
		err := future.Get(ctx, nil)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to execute child workflow %d: %v", i, err))
			return err
		}
	}

	return nil
}

func ChildWorkflow(ctx workflow.Context, msg string) error {
	// Define activity options
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 30,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Execute the activity
	err := workflow.ExecuteActivity(ctx, ProduceMsg, msg).Get(ctx, nil)
	if err != nil {
		logger := workflow.GetLogger(ctx)
		logger.Error("Failed to execute activity ProduceMsg", "Error", err)
	}

	return err
}
