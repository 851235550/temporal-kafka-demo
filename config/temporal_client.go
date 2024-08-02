package config

import (
	"context"
	"log"

	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

func NewTemporalClient() (client.Client, error) {
	return client.Dial(client.Options{
		HostPort: "localhost:7233",
	})
}

func TerminateWorker(c client.Client, workflowID string) error {
	ctx := context.Background()
	res, err := c.ListOpenWorkflow(ctx, &workflowservice.ListOpenWorkflowExecutionsRequest{})
	if err != nil {
		log.Fatalf("List open workflow was error. err: %s", err)
	}

	if len(res.Executions) > 0 {
		for _, e := range res.Executions {
			if e.Execution.WorkflowId != workflowID {
				continue
			}

			err := c.TerminateWorkflow(ctx, workflowID, "", "Terminating existing workflow to start a new one")
			if err != nil {
				return err
			}
		}
	}

	return nil
}
