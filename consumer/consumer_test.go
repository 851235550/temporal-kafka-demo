package consumer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"
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
	env.RegisterActivity(ConsumeMsg)

	// Mock the activity
	env.OnActivity(ConsumeMsg, mock.Anything, mock.Anything).Return(nil)

	childSuccCnt := 0
	env.SetOnChildWorkflowCompletedListener(func(workflowInfo *workflow.Info, result converter.EncodedValue, err error) {
		childSuccCnt++
	})

	// Start the workflow
	env.ExecuteWorkflow(CronParentConsumerWorkflow)

	s.True(env.AssertCalled(s.T(), "ConsumeMsg", mock.Anything, mock.Anything))
	s.True(env.AssertActivityNumberOfCalls(s.T(), "ConsumeMsg", testCalledCnt))

	s.Equal(childSuccCnt, testCalledCnt)

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

func (s *UnitTestSuite) TestCronConsume() {
	testWorkflow := func(ctx workflow.Context) error {
		ctx1 := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
			WorkflowRunTimeout: 10 * time.Second,
			CronSchedule:       cronEveryTwoMim,
		})

		cronFuture := workflow.ExecuteChildWorkflow(ctx1, CronParentConsumerWorkflow) // cron never stop so this future won't return

		// wait 4 minutes for the cron (cron will execute 3 times)
		_ = workflow.Sleep(ctx, time.Minute*4)
		s.False(cronFuture.IsReady())
		return nil
	}

	env := s.NewTestWorkflowEnvironment()
	env.RegisterWorkflow(testWorkflow)
	env.RegisterWorkflow(CronParentConsumerWorkflow)

	execTimes := make([]time.Time, 0)
	env.OnWorkflow(CronParentConsumerWorkflow, mock.Anything, mock.Anything).Return(func(ctx workflow.Context) error {
		execTimes = append(execTimes, workflow.Now(ctx))
		return nil
	})

	startTime, err := time.Parse(time.RFC3339, "2024-08-03T10:31:10Z")
	s.NoError(err)
	env.SetStartTime(startTime)

	env.ExecuteWorkflow(testWorkflow)

	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())
	s.True(env.AssertExpectations(s.T()))

	s.Len(execTimes, 3)
	s.Equal(execTimes[0], startTime)

	startTime1, err := time.Parse(time.RFC3339, "2024-08-03T10:32:00Z")
	s.NoError(err)
	s.Equal(execTimes[1], startTime1)

	startTime2, err := time.Parse(time.RFC3339, "2024-08-03T10:34:00Z")
	s.NoError(err)
	s.Equal(execTimes[2], startTime2)
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
