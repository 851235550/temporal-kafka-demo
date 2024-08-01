package consumer

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

// UnitTestSuite is a suite to run Temporal tests
type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
}

// TestCronParentConsumerWorkflow tests the CronParentConsumerWorkflow
func (s *UnitTestSuite) TestCronParentConsumerWorkflow() {
	env := s.NewTestWorkflowEnvironment()

	// Register workflows and activities
	env.RegisterWorkflow(CronParentConsumerWorkflow)
	env.RegisterWorkflow(ChildWorkflow)
	env.RegisterActivity(ConsumeMsg)

	// Mock the activity
	env.OnActivity(ConsumeMsg, mock.Anything, mock.Anything).Return(nil)

	// Start the workflow
	env.ExecuteWorkflow(CronParentConsumerWorkflow)

	// Check the workflow result
	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())

	// Assert that the expected number of child workflows were executed
	env.AssertExpectations(s.T())
}

// TestChildWorkflow tests the ChildWorkflow
func (s *UnitTestSuite) TestChildWorkflow() {
	env := s.NewTestWorkflowEnvironment()

	// Register workflows and activities
	env.RegisterWorkflow(ChildWorkflow)
	env.RegisterActivity(ConsumeMsg)

	// Mock the activity
	env.OnActivity(ConsumeMsg, mock.Anything).Return(nil)

	// Start the workflow
	env.ExecuteWorkflow(ChildWorkflow)

	// Check the workflow result
	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
