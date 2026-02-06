package worker

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

type testJob struct {
	id        string
	executed  *int32
	shouldErr bool
}

func (j *testJob) ID() string {
	return j.id
}

func (j *testJob) Execute(ctx context.Context) error {
	atomic.AddInt32(j.executed, 1)
	time.Sleep(10 * time.Millisecond) // Simulate work
	if j.shouldErr {
		return context.Canceled
	}
	return nil
}

func TestWorkerPoolBasic(t *testing.T) {
	pool := NewWorkerPool(3, 10)
	pool.Start()
	defer pool.Stop()
	
	var executed int32
	
	// Submit 5 jobs
	for i := 0; i < 5; i++ {
		job := &testJob{
			id:       string(rune('A' + i)),
			executed: &executed,
		}
		pool.Submit(job)
	}
	
	// Wait for jobs to complete
	time.Sleep(200 * time.Millisecond)
	
	if atomic.LoadInt32(&executed) != 5 {
		t.Errorf("Expected 5 jobs executed, got %d", executed)
	}
}

func TestWorkerPoolConcurrency(t *testing.T) {
	pool := NewWorkerPool(5, 100)
	pool.Start()
	defer pool.Stop()
	
	var executed int32
	jobCount := 50
	
	// Submit many jobs
	for i := 0; i < jobCount; i++ {
		job := &testJob{
			id:       string(rune(i)),
			executed: &executed,
		}
		pool.Submit(job)
	}
	
	// Wait for completion
	time.Sleep(2 * time.Second)
	
	if atomic.LoadInt32(&executed) != int32(jobCount) {
		t.Errorf("Expected %d jobs executed, got %d", jobCount, executed)
	}
}

func TestWorkerPoolShutdown(t *testing.T) {
	pool := NewWorkerPool(2, 5)
	pool.Start()
	
	var executed int32
	
	// Submit jobs
	for i := 0; i < 3; i++ {
		job := &testJob{
			id:       string(rune('X' + i)),
			executed: &executed,
		}
		pool.Submit(job)
	}
	
	// Immediate shutdown
	pool.Stop()
	
	// Should have processed at least some jobs
	if atomic.LoadInt32(&executed) == 0 {
		t.Error("Expected at least some jobs to be executed before shutdown")
	}
}

func TestWorkerPoolErrorHandling(t *testing.T) {
	pool := NewWorkerPool(2, 5)
	pool.Start()
	defer pool.Stop()
	
	var executed int32
	
	// Submit job that errors
	job := &testJob{
		id:        "error-job",
		executed:  &executed,
		shouldErr: true,
	}
	pool.Submit(job)
	
	time.Sleep(100 * time.Millisecond)
	
	// Should still execute (error is logged, not fatal)
	if atomic.LoadInt32(&executed) != 1 {
		t.Error("Job should execute even if it errors")
	}
}
