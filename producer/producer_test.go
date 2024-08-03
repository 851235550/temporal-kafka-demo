package producer

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"
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
	env.RegisterActivity(ProduceMsg)

	// Mock the ProduceMsg activity
	env.OnActivity(ProduceMsg, mock.Anything, mock.Anything).Return(nil)

	childSuccCnt := 0
	env.SetOnChildWorkflowCompletedListener(func(workflowInfo *workflow.Info, result converter.EncodedValue, err error) {
		childSuccCnt++
	})

	// Execute the workflow
	env.ExecuteWorkflow(CronParentProducerWorkflow)

	s.True(env.AssertCalled(s.T(), "ProduceMsg", mock.Anything, mock.Anything))
	s.True(env.AssertActivityNumberOfCalls(s.T(), "ProduceMsg", testCalledCnt))

	s.Equal(childSuccCnt, testCalledCnt)

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

func (s *UnitTestSuite) TestWorkflowActivityParamCorrect() {
	env := s.NewTestWorkflowEnvironment()
	env.OnActivity(ProduceMsg, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, value string) error {
			s.Equal("test_success", value)
			return nil
		})

	env.ExecuteWorkflow(ChildWorkflow, "test_success")

	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())
	s.True(env.AssertExpectations(s.T()))
}

func (s *UnitTestSuite) TestCronProduce() {
	testWorkflow := func(ctx workflow.Context) error {
		ctx1 := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
			WorkflowRunTimeout: 10 * time.Second,
			CronSchedule:       cronEveryMim,
		})

		cronFuture := workflow.ExecuteChildWorkflow(ctx1, CronParentProducerWorkflow) // cron never stop so this future won't return

		// wait 2 minutes for the cron (cron will execute 3 times)
		_ = workflow.Sleep(ctx, time.Minute*2)
		s.False(cronFuture.IsReady())
		return nil
	}

	env := s.NewTestWorkflowEnvironment()
	env.RegisterWorkflow(testWorkflow)
	env.RegisterWorkflow(CronParentProducerWorkflow)

	execTimes := make([]time.Time, 0)
	env.OnWorkflow(CronParentProducerWorkflow, mock.Anything, mock.Anything).Return(func(ctx workflow.Context) error {
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

	startTime2, err := time.Parse(time.RFC3339, "2024-08-03T10:33:00Z")
	s.NoError(err)
	s.Equal(execTimes[2], startTime2)
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
