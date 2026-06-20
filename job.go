package main

import (
	"os/exec"
	"sync"
)

type JobResult struct {
	RunID    uint64
	OK       bool
	ExitCode *int
	Label    string
}

type JobManager struct {
	running    bool
	runID      uint64
	resultCh   chan JobResult
	mu         sync.Mutex
}

func NewJobManager() *JobManager {
	return &JobManager{
		running:    false,
		runID:      0,
		resultCh:   make(chan JobResult, 1),
	}
}

func (j *JobManager) Start(action Action, runID uint64) {
	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		return
	}
	j.running = true
	j.runID = runID
	j.mu.Unlock()

	go func() {
		cmd := exec.Command("sh", "-c", action.Command)
		err := cmd.Run()

		var ok bool
		var exitCode *int
		if err != nil {
			ok = false
			if exitError, ok2 := err.(*exec.ExitError); ok2 {
				code := exitError.ExitCode()
				exitCode = &code
			}
		} else {
			ok = true
			code := 0
			exitCode = &code
		}

		result := JobResult{
			RunID:    runID,
			OK:       ok,
			ExitCode: exitCode,
			Label:    action.Label,
		}

		j.resultCh <- result
	}()
}

func (j *JobManager) TryRecv() *JobResult {
	select {
	case res := <-j.resultCh:
		return &res
	default:
		return nil
	}
}

func (j *JobManager) Running() bool {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.running
}

func (j *JobManager) SetRunning(v bool) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.running = v
}