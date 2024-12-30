package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `json:"name"`
	Password string
}

type OauthProvider struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

type OauthUser struct {
	ID              uint `gorm:"primaryKey"`
	UserId          uint
	User            User
	OauthProviderId uint
	OauthProvider   OauthProvider
	ProviderUserId  string
	AccessToken     string
	RefreshToken    string
}

type UserPermissions uint8

const (
	read UserPermissions = 1 << iota
	write
	execute
	admin
)

type Branch struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `json:"name"`
}

type Commit struct {
	ID       uint `gorm:"primaryKey"`
	BranchId uint
	Branch   Branch
	UserId   uint
	User     User
	Hash     string `json:"name"`
}

type Pipeline struct {
	ID              uint `gorm:"primaryKey"`
	CommitId        uint
	Commit          Commit
	Name            string
	Description     string
	Jobs            []Job
	Users           []PipelineUser
	PipelineWorkers []PipelineWorker
	IsPrivate       bool
}

type PipelineUser struct {
	ID         uint `gorm:"primaryKey"`
	UserId     uint
	User       User
	RoleId     uint
	Role       Role
	PipelineId uint
	Pipeline   Pipeline
}

type Role struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Permissions UserPermissions
}

type Job struct {
	ID         uint `gorm:"primaryKey"`
	PipelineId uint
	Pipeline   Pipeline
	NameId     string
	Name       *string

	// TODO: these are just for testing
	WillFail  bool
	FakeTimer uint

	JobWorker *JobWorker

	Dependencies []Job `gorm:"many2many:job_dependencies;joinForeignKey:ParentId;joinReferences:ChildId"`
	Parents      []Job `gorm:"many2many:job_dependencies;joinForeignKey:ChildId;joinReferences:ParentId"`
}

type PipelineWorker struct {
	ID         uint `gorm:"primaryKey"`
	PipelineId uint
	Pipeline   Pipeline
	UserId     *uint
	User       *User
	JobWorkers []JobWorker
	StartedAt  *time.Time
	FinishedAt *time.Time
	StatusId   uint
	Status     Status
}

type Status struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

type statusOption uint

const (
	statusPending statusOption = iota + 1
	statusRunning
	statusFailed
	statusCompleted
	statusPlanned
	statusCancelled
)

func (d statusOption) String() string {
	switch d {
	case statusPending:
		return "pending"
	case statusRunning:
		return "running"
	case statusFailed:
		return "failed"
	case statusCompleted:
		return "completed"
	case statusPlanned:
		return "planned"
	case statusCancelled:
		return "cancelled"
	default:
		return "pending"
	}
}

func (d statusOption) RussianString() string {
	switch d {
	case statusPending:
		return "ожидание"
	case statusRunning:
		return "выполнение"
	case statusFailed:
		return "ошибка"
	case statusCompleted:
		return "успех"
	case statusPlanned:
		return "план"
	case statusCancelled:
		return "отмена"
	default:
		return "ожидание"
	}
}

type JobWorker struct {
	ID               uint `gorm:"primaryKey"`
	StatusId         uint
	Status           Status
	JobId            uint
	Job              Job
	UserId           *uint
	User             *User
	PipelineWorkerId uint
	PipelineWorker   PipelineWorker
	StageWorkers     []StageWorker
	StartedAt        *time.Time
	FinishedAt       *time.Time
}

type Stage struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	JobId    uint
	Job      Job
	RunsFor  uint
	WillFail bool
}

type StageWorker struct {
	ID          uint `gorm:"primaryKey"`
	StatusId    uint
	Status      Status
	StageId     uint
	Stage       Stage
	JobWorkerId uint
	JobWorker   JobWorker
	Logs        []StageLog
	StartedAt   *time.Time
	FinishedAt  *time.Time
}

type PipelineStatisticsShort struct {
	AverageDuration uint         `json:"averageDuration"`
	TotalDuration   uint         `json:"totalDuration"`
	TotalErrors     uint         `json:"totalErrors"`
	PipelineStatus  statusOption `json:"pipelineStatus"`
}

type LogType struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

type logTypeOptions uint

const (
	logTypeBase logTypeOptions = iota + 1
	logTypeError
	logTypeWarn
	logTypeInfo
)

func (d logTypeOptions) String() string {
	switch d {
	case logTypeBase:
		return "base"
	case logTypeWarn:
		return "warn"
	case logTypeError:
		return "error"
	case logTypeInfo:
		return "info"
	default:
		return "base"
	}
}

type StageLog struct {
	ID            uint `gorm:"primaryKey"`
	Text          string
	Line          uint
	LogTypeId     uint
	LogType       LogType
	StageWorkerId uint
	StageWorker   StageWorker
	CreatedAt     time.Time
}

var (
	db  *gorm.DB
	err error
)

func DBConnection() error {
	host := fmt.Sprintf("host=%s ", os.Getenv("DATABASE_HOST"))
	user := fmt.Sprintf("user=%s ", os.Getenv("DATABASE_USER"))
	password := fmt.Sprintf("password=%s ", os.Getenv("DATABASE_PASSWORD"))
	name := fmt.Sprintf("dbname=%s ", os.Getenv("DATABASE_NAME"))
	port := fmt.Sprintf("port=%s ", os.Getenv("DATABASE_PORT"))
	dsn := host + user + password + name + port

	log.Printf("Connecting to %v", dsn)

	maxRetries := 8
	retryInterval := 2 * time.Second

	var err error
	for attempts := 1; attempts <= maxRetries; attempts++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			return nil
		}

		fmt.Printf("Attempt %d/%d failed: %v. Retrying in %s...\n", attempts, maxRetries, err, retryInterval)
		time.Sleep(retryInterval)
	}

	return fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}

type PipelineWorkerStatistics struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Status          string `json:"status"`
	StatusId        statusOption
	User            *string    `json:"user,omitempty"`
	NumberOfErrors  *uint      `json:"numberOfErrors"`
	StartedAt       *time.Time `json:"startedAt"`
	FinishedAt      *time.Time `json:"finishedAt"`
	TotalDuration   *uint      `json:"totalDuration"`
	AverageDuration *uint      `json:"averageDuration"`
}

func fetchPipelineWorkerJobsStatistics(pipelineWorkerId uint) ([]PipelineWorkerStatistics, error) {
	query := `
	SELECT 
		jw.job_id AS id,
		COALESCE(j.name, j.name_id) AS name,
		s.name AS status,
		s.id AS status_id,
		u.name AS user,
		CASE 
			WHEN NOT jw.status_id IN (3, 4) THEN NULL 
			ELSE COALESCE(COUNT(sl.id), 0) 
		END AS number_of_errors,
		CASE 
			WHEN NOT jw.status_id IN (3, 4) THEN NULL 
			ELSE jw.started_at 
		END AS started_at,
		CASE 
			WHEN NOT jw.status_id IN (3, 4) THEN NULL 
			ELSE jw.finished_at 
		END AS finished_at,
		CASE 
			WHEN NOT jw.status_id IN (3, 4) THEN NULL 
			ELSE ROUND(EXTRACT(EPOCH FROM (jw.finished_at - jw.started_at))) 
		END AS total_duration,
		CASE 
			WHEN NOT jw.status_id IN (3, 4) THEN NULL 
			ELSE ROUND((COALESCE(AVG(EXTRACT(EPOCH FROM (sw.finished_at - sw.started_at))) 
			FILTER (WHERE sw.started_at IS NOT NULL AND sw.finished_at IS NOT NULL)))::numeric, 0) 
		END AS average_duration
	FROM 
		job_workers jw
	JOIN 
		jobs j ON j.id = jw.job_id
	JOIN 
		statuses s ON s.id = jw.status_id
	LEFT JOIN 
		stage_workers sw ON sw.job_worker_id = jw.id
	LEFT JOIN 
		stage_logs sl ON sl.stage_worker_id = sw.id AND sl.log_type_id = 2
	LEFT JOIN 
		users u ON u.id = jw.user_id
	WHERE 
		jw.pipeline_worker_id = ?
	GROUP BY 
		jw.job_id, j.name, j.name_id, s.name, jw.started_at, jw.finished_at, u.name, jw.status_id, jw.id, s.id
	ORDER BY 
		started_at ASC, jw.id ASC;
	`

	rows, err := db.Raw(query, pipelineWorkerId).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make([]PipelineWorkerStatistics, 0)
	for rows.Next() {
		var statsCurr PipelineWorkerStatistics
		db.ScanRows(rows, &statsCurr)
		data = append(data, statsCurr)
	}

	return data, nil
}

func fetchPipelineWorkerStatistics(stats *PipelineWorkerStatistics, pipelineWorkerId uint) error {
	query := `
	SELECT 
		pw.id AS id,
		p.name AS name,
		s.name AS status,
		s.id AS status_id,
		u.name AS user,
		COALESCE(SUM(CASE WHEN sl.log_type_id = 2 THEN 1 ELSE 0 END), 0) AS number_of_errors,
		pw.started_at,
		pw.finished_at,
		ROUND(EXTRACT(EPOCH FROM (pw.finished_at - pw.started_at))) AS total_duration,
		ROUND((COALESCE(AVG(EXTRACT(EPOCH FROM (jw.finished_at - jw.started_at))) FILTER (WHERE jw.started_at IS NOT NULL AND jw.finished_at IS NOT NULL)))::numeric, 0) AS average_duration
	FROM 
		pipeline_workers pw
	JOIN 
		pipelines p ON p.id = pw.pipeline_id
	JOIN 
		statuses s ON s.id = pw.status_id
	JOIN 
		job_workers jw ON jw.pipeline_worker_id = pw.id
	JOIN 
		users u ON u.id = pw.user_id
	LEFT JOIN 
		stage_workers sw ON sw.job_worker_id = jw.id
	LEFT JOIN 
		stage_logs sl ON sl.stage_worker_id = sw.id
	WHERE 
		pw.id = ?
	GROUP BY 
    p.name, s.name, u.name, pw.started_at, pw.finished_at, pw.id, s.id;
	`

	err := db.Raw(query, pipelineWorkerId).Scan(&stats).Error
	if err != nil {
		return err
	}

	return nil
}

func fetchPipelineShortStatistics(data *PipelineStatisticsShort, pipelineWorkerId uint) error {
	query := `
		SELECT 
			CASE 
				WHEN max(pipeline_workers.finished_at) IS NOT NULL THEN 
					CASE 
						WHEN ROUND(avg(EXTRACT(epoch FROM job_workers.finished_at - job_workers.started_at))) < 0 THEN NULL
						ELSE ROUND(avg(EXTRACT(epoch FROM job_workers.finished_at - job_workers.started_at)))::integer 
					END
				ELSE 
					NULL 
			END AS average_duration,

			CASE 
				WHEN ROUND(EXTRACT(EPOCH FROM (max(pipeline_workers.finished_at) - max(pipeline_workers.started_at)))) < 0 THEN NULL
				ELSE ROUND(EXTRACT(EPOCH FROM (max(pipeline_workers.finished_at) - max(pipeline_workers.started_at))))::integer 
			END AS total_duration,

			count(*) FILTER (WHERE job_workers.status_id = 3) AS total_errors,

			CASE
				WHEN count(*) FILTER (WHERE job_workers.status_id = 4) = count(*) THEN 4
				WHEN count(*) FILTER (WHERE job_workers.status_id = 3) > 0 THEN 3
				WHEN count(*) FILTER (WHERE job_workers.status_id = 2) > 0 THEN 2
				ELSE 1
			END AS pipeline_status

		FROM job_workers
		LEFT JOIN pipeline_workers ON pipeline_workers.id = job_workers.pipeline_worker_id
		WHERE pipeline_workers.id = ?;

  	`

	err := db.Raw(query, pipelineWorkerId).Scan(&data).Error
	if err != nil {
		return err
	}

	return nil
}
