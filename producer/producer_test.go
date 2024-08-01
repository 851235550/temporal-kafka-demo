package producer

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
}

// SetupTest is called before every test
func (s *UnitTestSuite) SetupTest() {
	// Reset the child worker count before each test
	SetChildWorkerCnt(200)
}

func (s *UnitTestSuite) TestCronParentProducerWorkflow() {
	env := s.NewTestWorkflowEnvironment()

	// Register the workflow and activity
	env.RegisterWorkflow(CronParentProducerWorkflow)
	env.RegisterWorkflow(ChildWorkflow)
	env.RegisterActivity(ProduceMsg)

	// Mock the ProduceMsg activity
	env.OnActivity(ProduceMsg, mock.Anything, mock.Anything).Return(nil)

	// Execute the workflow
	env.ExecuteWorkflow(CronParentProducerWorkflow)

	// Check if workflow completed successfully
	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())

	// Assert that the ProduceMsg activity was called 200 times
	env.AssertExpectations(s.T())
}

func (s *UnitTestSuite) TestChildWorkflow() {
	env := s.NewTestWorkflowEnvironment()

	// Register the workflow and activity
	env.RegisterWorkflow(ChildWorkflow)
	env.RegisterActivity(ProduceMsg)

	// Mock the ProduceMsg activity
	env.OnActivity(ProduceMsg, mock.Anything, "test-msg").Return(nil)

	// Execute the workflow
	env.ExecuteWorkflow(ChildWorkflow, "test-msg")

	// Check if workflow completed successfully
	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())

	// Assert that the ProduceMsg activity was called once
	env.AssertExpectations(s.T())
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
