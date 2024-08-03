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

var testCalledCnt = 5

// SetupTest is called before every test
func (s *UnitTestSuite) SetupTest() {
	// Reset the child worker count before each test
	SetChildWorkerCnt(testCalledCnt)
}

func (s *UnitTestSuite) TestCronParentProducerWorkflow() {
	env := s.NewTestWorkflowEnvironment()

	// Register the workflow and activity
	env.RegisterWorkflow(CronParentProducerWorkflow)
	env.RegisterWorkflow(ChildWorkflow)

	// Mock the ProduceMsg activity
	env.OnActivity(ProduceMsg, mock.Anything, mock.Anything).Return(nil)

	// Execute the workflow
	env.ExecuteWorkflow(CronParentProducerWorkflow)

	s.True(env.AssertCalled(s.T(), "ProduceMsg", mock.Anything, mock.Anything))
	s.True(env.AssertActivityNumberOfCalls(s.T(), "ProduceMsg", testCalledCnt))

	// Check if workflow completed successfully
	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())

	// Assert that the ProduceMsg activity was called 10 times
	s.True(env.AssertExpectations(s.T()))
}

func (s *UnitTestSuite) TestChildWorkflow() {
	env := s.NewTestWorkflowEnvironment()

	// Register the workflow and activity
	env.RegisterWorkflow(ChildWorkflow)
	env.RegisterActivity(ProduceMsg)

	// Mock the ProduceMsg activity
	env.OnActivity(ProduceMsg, mock.Anything, "msg").Return(nil)

	// Execute the workflow
	env.ExecuteWorkflow(ChildWorkflow, "msg")

	s.True(env.AssertCalled(s.T(), "ProduceMsg", mock.Anything, "msg"))
	s.True(env.AssertActivityNumberOfCalls(s.T(), "ProduceMsg", 1))

	// Check if workflow completed successfully
	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())

	// Assert that the ProduceMsg activity was called once
	s.True(env.AssertExpectations(s.T()))
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
