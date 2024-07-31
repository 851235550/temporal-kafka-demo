package tests

import (
	"suger/consumer"
	"testing"

	"go.temporal.io/sdk/testsuite"
)

func TestConsumerWorkflow(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	env.ExecuteWorkflow(consumer.ConsumerWorkflow)
	if !env.IsWorkflowCompleted() {
		t.Fatal("Workflow not completed")
	}
	if err := env.GetWorkflowError(); err != nil {
		t.Fatalf("Workflow failed: %v", err)
	}
}
