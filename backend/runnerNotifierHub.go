package main

import "time"

type RunnerNotifierHub struct {
	pipelineId uint
}

func NewRunnerNotifierHub(pipelineId uint) *RunnerNotifierHub {
	return &RunnerNotifierHub{pipelineId: pipelineId}
}

func (n *RunnerNotifierHub) NotifyJobStatusBatch(statuses *statusBatch, jobMetadata *JobMetadata, pipelineStatus statusOption, statistics *PipelineRunStatistics) {
	hub, ok := hubs[n.pipelineId]
	if ok {
		hub.handleJobStatusBatch(statuses, jobMetadata, pipelineStatus, statistics)
	}
}

func (n *RunnerNotifierHub) NotifyJobStatus(jobWorkerId uint, status statusOption, pipelineStatus statusOption, startedBy *string, startedAt *time.Time, finishedAt *time.Time) {
	hub, ok := hubs[n.pipelineId]
	if ok {
		hub.handleJobStatus(jobWorkerId, status, pipelineStatus, startedBy, startedAt, finishedAt)
	}
}

func (n *RunnerNotifierHub) NotifyStageStatus(jobWorkerId uint, stageWorkerId uint, info *stageStatusExtra) {
	hub, ok := hubs[n.pipelineId]
	if ok {
		hub.handleStageStatus(jobWorkerId, stageWorkerId, info)
	}
}

func (n *RunnerNotifierHub) NotifyJobProgress(jobWorkerId uint, progress uint) {
	hub, ok := hubs[n.pipelineId]
	if ok {
		hub.handleJobProgress(jobWorkerId, progress)
	}
}

func (n *RunnerNotifierHub) NotifyStageLogs(jobWorkerId uint, stageWorkerId uint, logs []StageLogData, isForced bool) {
	hub, ok := hubs[n.pipelineId]
	if ok {
		hub.handleStageLogs(jobWorkerId, stageWorkerId, logs, isForced)
	}
}
