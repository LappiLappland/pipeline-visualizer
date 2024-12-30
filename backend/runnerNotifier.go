package main

import "time"

type RunnerNotifier interface {
	NotifyJobStatusBatch(statuses *statusBatch, jobMetadata *JobMetadata, pipelineStatus statusOption, statistics *PipelineRunStatistics)
	NotifyJobStatus(jobWorkerId uint, status statusOption, pipelineStatus statusOption, startedBy *string, startedAt *time.Time, finishedAt *time.Time)
	NotifyStageStatus(jobWorkerId uint, stageWorkerId uint, info *stageStatusExtra)
	NotifyJobProgress(jobWorkerId uint, progress uint)
	NotifyStageLogs(jobWorkerId uint, stageWorkerId uint, logs []StageLogData, isForced bool)
}
