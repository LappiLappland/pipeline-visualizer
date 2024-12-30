package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/imacks/bitflags-go"
)

type BroadcastMessage struct {
	client *Client
	data   []byte
}

type ServerReceiverConds struct {
	jobWorkerId   uint
	stageWorkerId uint
}

type ServerMessage struct {
	serverReceiverConds *ServerReceiverConds
	data                []byte
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Pipeline Worker Id
	pipelineId uint

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan BroadcastMessage

	// Last time the job progress was sent
	jobProgress map[uint]time.Time

	// Last time logs were dumped
	logDump map[uint]time.Time

	// Outbound messages to the clients
	messages chan ServerMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

type HubAction struct {
	Type string `json:"type"`
}

type HubActionJobSubscribe struct {
	ID      uint `json:"id,omitempty"`
	StageId uint `json:"stageId,omitempty"`
}

type HubActionActJob struct {
	ID  uint   `json:"id"`
	Act string `json:"act"`
}

type HubActionActPipeline struct {
	Act string `json:"act"`
}

type HubEventJobStatus struct {
	Type           string     `json:"type"`
	ID             uint       `json:"id"`
	Status         string     `json:"status"`
	PipelineStatus string     `json:"pipelineStatus,omitempty"`
	FinishedAt     *time.Time `json:"finishedAt,omitempty"`
	StartedAt      *time.Time `json:"startedAt,omitempty"`
	StartedBy      *string    `json:"startedBy,omitempty"`
}

type statusBatch struct {
	Completed []uint `json:"completed,omitempty"`
	Failed    []uint `json:"failed,omitempty"`
	Pending   []uint `json:"pending,omitempty"`
	Running   []uint `json:"running,omitempty"`
	Cancelled []uint `json:"cancelled,omitempty"`
	Planned   []uint `json:"planned,omitempty"`
}

type JobMetadata struct {
	FinishedAt *time.Time `json:"finishedAt"`
	StartedAt  *time.Time `json:"startedAt"`
	StartedBy  *string    `json:"startedBy"`
	Fields     []string   `json:"fields"`
}

type HubEventJobStatusBatch struct {
	Type           string                 `json:"type"`
	Statuses       statusBatch            `json:"statuses"`
	Metadata       *JobMetadata           `json:"metadata,omitempty"`
	PipelineStatus string                 `json:"pipelineStatus,omitempty"`
	Statistics     *PipelineRunStatistics `json:"statistics,omitempty"`
}

type stageStatusExtra struct {
	Status     string     `json:"status"`
	StartedAt  *time.Time `json:"startedAt"`
	FinishedAt *time.Time `json:"finishedAt"`
}

type HubEventStageStatus struct {
	Type  string            `json:"type"`
	ID    uint              `json:"id"`
	JobID uint              `json:"jobId"`
	Info  *stageStatusExtra `json:"info"`
}

type HubEventJobProgress struct {
	Type     string `json:"type"`
	ID       uint   `json:"id"`
	Progress uint   `json:"progress"`
}

type HubEventStageLogs struct {
	Type string         `json:"type"`
	Logs []StageLogData `json:"logs"`
}

var hubs map[uint]*Hub

func initHubs() {
	hubs = make(map[uint]*Hub)
}

func newHub(pipelineId uint) *Hub {
	hub := &Hub{
		pipelineId:  pipelineId,
		broadcast:   make(chan BroadcastMessage),
		messages:    make(chan ServerMessage),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		jobProgress: make(map[uint]time.Time),
		logDump:     make(map[uint]time.Time),
	}
	hubs[pipelineId] = hub
	go hub.run()
	return hub
}

func getHub(pipelineId uint) *Hub {
	hub, ok := hubs[pipelineId]

	if !ok {
		hub = newHub(pipelineId)
	}

	return hub
}

func (h *Hub) run() {
	updateClientsPermissions := time.NewTicker(45 * time.Second)
	defer updateClientsPermissions.Stop()

	for {
		select {
		case <-updateClientsPermissions.C:
			for client := range h.clients {
				client.updateClientPermissions()
			}
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			h.executeEvent(&message)
		case message := <-h.messages:
			for client := range h.clients {
				if !bitflags.Has(client.permissions, read) {
					continue
				}
				if h.shouldBeReported(client, &message) {
					select {
					case client.send <- message.data:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}

			}
		}
	}
}

func (h *Hub) shouldBeReported(client *Client, message *ServerMessage) bool {
	shouldReport := true

	if message.serverReceiverConds != nil {
		if message.serverReceiverConds.jobWorkerId != 0 {
			shouldReport = client.jobId == message.serverReceiverConds.jobWorkerId
			if message.serverReceiverConds.stageWorkerId != 0 && shouldReport {
				shouldReport = client.stageId == message.serverReceiverConds.stageWorkerId
			}
		}
	}

	return shouldReport
}

func (h *Hub) sendMessage(message []byte) {
	h.messages <- ServerMessage{data: message}
}

func (h *Hub) sendMessageJob(message []byte, jobWorkerId uint) {
	h.messages <- ServerMessage{
		data: message,
		serverReceiverConds: &ServerReceiverConds{
			jobWorkerId: jobWorkerId,
		},
	}
}

func (h *Hub) sendMessageStage(message []byte, jobWorkerId uint, stageWorkerId uint) {
	h.messages <- ServerMessage{
		data: message,
		serverReceiverConds: &ServerReceiverConds{
			jobWorkerId:   jobWorkerId,
			stageWorkerId: stageWorkerId,
		},
	}
}

func (h *Hub) executeEvent(message *BroadcastMessage) {
	var hubEvent HubAction
	err := json.Unmarshal(message.data, &hubEvent)
	if err != nil {
		log.Println(err)
		return
	}

	switch hubEvent.Type {
	case "actPipeline":
		var eventData HubActionActPipeline
		err := json.Unmarshal(message.data, &eventData)
		if err != nil {
			log.Println(err)
			break
		}
		h.actionActPipeline(message.client, &eventData)
	case "actJob":
		var eventData HubActionActJob
		err := json.Unmarshal(message.data, &eventData)
		if err != nil {
			log.Println(err)
			break
		}
		h.actionActJob(&eventData, message.client)
	case "jobSubscribe":
		var eventData HubActionJobSubscribe
		err := json.Unmarshal(message.data, &eventData)
		if err != nil {
			log.Println(err)
			break
		}
		h.actionJobSubscribe(&eventData, message.client)
	}
}

func (h *Hub) handleJobStatus(jobWorkerId uint, status statusOption, pipelineStatus statusOption, startedBy *string, startedAt *time.Time, finishedAt *time.Time) {
	var msg HubEventJobStatus
	msg.Type = "jobStatus"
	msg.ID = jobWorkerId
	msg.Status = status.String()
	if pipelineStatus != 0 {
		msg.PipelineStatus = pipelineStatus.String()
	}
	msg.StartedAt = startedAt
	msg.FinishedAt = finishedAt
	msg.StartedBy = startedBy

	msgData, err := json.Marshal(msg)

	if err != nil {
		log.Println(err)
		return
	}

	h.sendMessage(msgData)
}

func (h *Hub) handleStageStatus(jobWorkerId uint, stageWorkerId uint, info *stageStatusExtra) {
	var msg HubEventStageStatus
	msg.Type = "stageStatus"
	msg.ID = stageWorkerId
	msg.JobID = jobWorkerId
	msg.Info = &stageStatusExtra{
		Status:     info.Status,
		FinishedAt: info.FinishedAt,
		StartedAt:  info.StartedAt,
	}

	msgData, err := json.Marshal(msg)

	if err != nil {
		log.Println(err)
		return
	}

	h.sendMessageJob(msgData, jobWorkerId)
}

func (h *Hub) handleJobStatusBatch(statuses *statusBatch, jobMetadata *JobMetadata, pipelineStatus statusOption, statistics *PipelineRunStatistics) {
	var msg HubEventJobStatusBatch
	msg.Type = "jobStatusBatch"
	msg.Statuses = *statuses
	if pipelineStatus != 0 {
		msg.PipelineStatus = pipelineStatus.String()
	}
	msg.Statistics = statistics
	msg.Metadata = jobMetadata

	msgData, err := json.Marshal(msg)

	if err != nil {
		log.Println(err)
		return
	}

	h.sendMessage(msgData)
}

func (h *Hub) handleJobProgress(jobWorkerId uint, progress uint) {

	lastTime, ok := h.jobProgress[jobWorkerId]
	if ok {
		since := time.Since(lastTime)
		if since < time.Duration(1)*time.Second {
			return
		}
	}

	h.jobProgress[jobWorkerId] = time.Now()

	var msg HubEventJobProgress
	msg.Type = "jobProgress"
	msg.Progress = progress
	msg.ID = jobWorkerId

	msgData, err := json.Marshal(msg)

	if err != nil {
		log.Println(err)
		return
	}

	h.sendMessage(msgData)
}

func (h *Hub) handleStageLogs(jobWorkerId uint, stageWorkerId uint, logs []StageLogData, isForced bool) {

	if !isForced {
		lastTime, ok := h.logDump[jobWorkerId]
		if ok {
			since := time.Since(lastTime)
			if since < time.Duration(200)*time.Millisecond {
				return
			}
		}
	}
	h.logDump[jobWorkerId] = time.Now()

	var msg HubEventStageLogs
	msg.Type = "stageLogs"
	msg.Logs = logs

	msgData, err := json.Marshal(msg)

	if err != nil {
		log.Println(err)
		return
	}

	h.sendMessageStage(msgData, jobWorkerId, stageWorkerId)
}

func (h *Hub) actionJobSubscribe(data *HubActionJobSubscribe, client *Client) {
	client.jobId = data.ID
	client.stageId = data.StageId
	log.Printf("Client subscribed to job \"%d\" and stage \"%d\" ", data.ID, data.StageId)
}

func (h *Hub) actionActJob(data *HubActionActJob, client *Client) {
	log.Printf("hub action - actJob (%s)", data.Act)

	runner := GetRunner(h.pipelineId)
	userId := client.userId

	switch data.Act {
	case "restart":
		fallthrough
	case "start":
		runner.startJob(data.ID, uint(userId))
	case "stop":
		runner.stopJob(data.ID)
	default:
	}
}

func (h *Hub) actionActPipeline(client *Client, data *HubActionActPipeline) {

	log.Printf("hub action - actPipeline (%s)", data.Act)

	runner := GetRunner(h.pipelineId)
	userId := client.userId

	switch data.Act {
	case "restart":
		fallthrough
	case "start":
		runner.RunPipeline(client.userId)
	case "stop":
		runner.StopPipeline()
	case "continue":
		runner.ContinuePipeline(uint(userId))
	default:
	}
}
