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

var testCalledCnt = 5

// SetupTest is called before every test
func (s *UnitTestSuite) SetupTest() {
	// Reset the child worker count before each test
	SetChildWorkerCnt(testCalledCnt)
}

// TestCronParentConsumerWorkflow tests the CronParentConsumerWorkflow
func (s *UnitTestSuite) TestCronParentConsumerWorkflow() {
	env := s.NewTestWorkflowEnvironment()

	// Register workflows and activities
	env.RegisterWorkflow(CronParentConsumerWorkflow)
	env.RegisterWorkflow(ChildWorkflow)

	// Mock the activity
	env.OnActivity(ConsumeMsg, mock.Anything, mock.Anything).Return(nil)

	// Start the workflow
	env.ExecuteWorkflow(CronParentConsumerWorkflow)

	s.True(env.AssertCalled(s.T(), "ConsumeMsg", mock.Anything, mock.Anything))
	s.True(env.AssertActivityNumberOfCalls(s.T(), "ConsumeMsg", testCalledCnt))

	// Check the workflow result
	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())

	// Assert that the expected number of child workflows were executed
	s.True(env.AssertExpectations(s.T()))
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

	s.True(env.AssertCalled(s.T(), "ConsumeMsg", mock.Anything))
	s.True(env.AssertActivityNumberOfCalls(s.T(), "ConsumeMsg", 1))

	// Check the workflow result
	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())

	s.True(env.AssertExpectations(s.T()))
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
