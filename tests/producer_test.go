package tests

import (
	"suger/producer"
	"testing"

	"go.temporal.io/sdk/testsuite"
)

func TestProducerWorkflow(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	env.ExecuteWorkflow(producer.ProduceMsg)
	if !env.IsWorkflowCompleted() {
		t.Fatal("Workflow not completed")
	}
	if err := env.GetWorkflowError(); err != nil {
		t.Fatalf("Workflow failed: %v", err)
	}
}
