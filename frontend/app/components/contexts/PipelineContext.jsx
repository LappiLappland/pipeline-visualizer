import getJobStagesData from "@/app/helpers/getJobStagesData"
import getLocation from "@/app/helpers/getLocation"
import getVisualizerDataObject, { getVisualizerDataByLevel } from "@/app/helpers/getVisualizerDataObject"
import parsePermissions from "@/app/helpers/parsePermissions"
import retryCallback, { NetworkError, retryNetwork } from "@/app/helpers/retryCallback"
import Cookies from "js-cookie"
import { useParams, usePathname } from "next/navigation"
import { createContext, useCallback, useEffect, useRef, useState } from "react"

export const PipelineContext = createContext({
    pipelineData: null,
    pipelineDataByLevels: {},
    jobId: 0,
    isForbidden: false,
    jobStages: null,
    stats: {
        totalJobs: 0,
        pipelineStatus: 'pending',
        startedBy: null,
        startedAt: null,
        finishedAt: null,
        averageDuration: null,
        totalDuration: null,
        totalErrors: 0,
    },
    setCurrentStage: (n) => undefined,
    currentStage: null,
    stageLogs: [],
    permissions: {
        canWrite: false,
        canRead: false,
        canExecute: false,
        canAdmin: false,
    },

    actionActPipeline: () => undefined,
    actionActJob:  () => undefined,
})

// Make sure to put WebsocketProvider higher up in
// the component tree than any consumer.
export const PipelineProvider = ({ children }) => {

    const [isForbidden, setIsForbidden] = useState(false); // If server denied access
    const [pipelineJobs, setPipelineJobs] = useState(null); // Stores object-map with id as key to job data 
    const [pipelineJobsByLevels, setPipelineJobsByLevel] = useState({}); // Stores object with levels as keys and array of job ids as value
    const [isOnJob, setIsOnJob] = useState(false);
    const [jobStages, setJobStages] = useState(null); // Stores array of stages
    const [socket, setSocket] = useState(null); // Stores socket
    const [currentStage, setCurrentStage] = useState(null); // Stores selected stage id
    const [stageLogs, setStageLogs] = useState([]); // Stores stage logs of selected stage
    const [permissions, setPermissions] = useState(() => ({
        canWrite: false,
        canRead: false,
        canExecute: false,
        canAdmin: false,
    })) // Stores user permissions

    const [stats, setStats] = useState(() => ({
        title: '',
        totalJobs: 0,
        pipelineStatus: 'pending',
        startedBy: null, // Convert to Date
        startedAt: null, // Convert to Date
        finishedAt: null, // Convert to Date
        averageDuration: null,
        totalDuration: null,
        totalErrors: 0,
    })); // Stores some stats that frequently change

    const params = useParams()

    const prevStage = useRef(null);

    // *** SOCKET HANDLERS START ***
    const handleJobStatus = useCallback((msg) => {
        if (pipelineJobs) {
            const newJobs = {...pipelineJobs};
            newJobs[msg.id].status = msg.status;

            // Some unique cases
            switch (msg.status) {
                case 'pending':
                    newJobs[msg.id].progress = 0;
                case 'running':
                    newJobs[msg.id].finishedAt = null;
                    newJobs[msg.id].startedAt = msg.startedAt;
                    newJobs[msg.id].startedBy = msg.startedBy;
                    break;
                case 'completed':
                case 'failed':
                    newJobs[msg.id].finishedAt = msg.finishedAt;
                    break;
                case 'pending':
                    newJobs[msg.id].finishedAt = null;
                    newJobs[msg.id].startedAt = null;
                    newJobs[msg.id].startedBy = null;
                    break;
            }

            setPipelineJobs(newJobs);
            setStats((prev) => ({
                ...prev,
                pipelineStatus: msg.pipelineStatus
            }));
        }
    }, [pipelineJobs])

    const handleJobStatusBatch = useCallback((msg) => {
        if (pipelineJobs) {
            const newJobs = {...pipelineJobs};
            Object.entries(msg.statuses).forEach(([status, jobs]) => {
                jobs.forEach((jobId) => {
                    const job = newJobs[jobId]
                    job.status = status;
                    if (job.status === 'running') {
                        job.progress = 0;
                    }
                    if (msg.metadata) {
                        // Some special cases
                        msg.metadata.fields.forEach((key) => {
                            let value = msg.metadata[key];
                            if (key === 'finishedAt') {
                                if (job.status === 'planned' || job.status === 'running') {
                                    value = null;
                                }
                            } else if (key === 'startedAt') {
                                if (job.status === 'cancelled' || job.status === 'planned') {
                                    value = null;
                                }
                            } else if (key === 'startedBy') {
                                if (job.status === 'cancelled') {
                                    value = null;
                                }
                            }
                            job[key] = value;
                        })
                    }
                });
            });
            setPipelineJobs(newJobs);

            const st = msg.statistics;
            setStats((prev) => ({
                ...prev,
                pipelineStatus: msg.pipelineStatus,
                averageDuration: st ? st.averageDuration : prev.averageDuration,
                totalDuration: st ? st.totalDuration : prev.totalDuration,
                totalErrors: st ? st.totalErrors : prev.totalErrors,
                startedBy: st ? st.startedBy : prev.startedBy,
                startedAt: st ? (st.startedAt && new Date(st.startedAt)) : prev.startedAt,
                finishedAt: st ? (st.finishedAt && new Date(st.finishedAt)) : prev.finishedAt,
            }));
        } 
    }, [pipelineJobs])

    const handleStageStatus = useCallback((msg) => {
        if (jobStages) {
            const jobId = +params.jobId;
            if (jobId === msg.jobId) {
                const newJobStages = {...jobStages};
                const stage = newJobStages[msg.id]
                const inf = msg.info;
                stage.status = inf.status;
                stage.finishedAt = inf.finishedAt ? new Date(inf.finishedAt) : inf.finishedAt;
                stage.startedAt = inf.startedAt ? new Date(inf.startedAt) : inf.startedAt;
                setJobStages(newJobStages);
            }
        }
    }, [jobStages, params.jobId])

    const handleJobProgress = useCallback((msg) => {
        if (pipelineJobs) {
            const newJobs = {...pipelineJobs}
            pipelineJobs[msg.id].progress = msg.progress
            setPipelineJobs(newJobs)
        }
    }, [pipelineJobs])

    const handleStageLogs = useCallback((msg) => {
        if (stageLogs) {
            setStageLogs((prev) => [...prev, ...msg.logs]);
        }
    }, [stageLogs])

    // *** SOCKET HANDLERS END ***

    // *** SOCKET ACTIONS START ***

    const actionActJob = useCallback((jobId, act) => {
        if (socket) {
            const msg = {
                type: 'actJob',
                act, // stop, start/restart, continue
                id: jobId
            }
            socket.send(JSON.stringify(msg));
        }
    }, [socket])

    const actionActPipeline = useCallback((act) => {
        if (socket) {
            const msg = {
                type: 'actPipeline',
                act, // stop, start/restart, continue
            }
            socket.send(JSON.stringify(msg));
        }
    }, [socket])

    const actionJobSubscribe = useCallback((jobId, stageId) => {
        if (socket) {
            const msg = {
                type: 'jobSubscribe',
                id: jobId,
                stageId: stageId
            }
            socket.send(JSON.stringify(msg));
        }
    }, [socket])

    const actionStageSubscribe = useCallback((stageId) => {
        if (params.jobId && +params.jobId) {
            actionJobSubscribe(+params.jobId, stageId);
        }
    }, [actionJobSubscribe, params.jobId])

    const actionStageUnsubscribe = useCallback(() => {
        if (params.jobId && +params.jobId) {
            actionJobSubscribe(+params.jobId, undefined);
        }
    }, [actionJobSubscribe, params.jobId])

    const actionJobUnsubscribe = useCallback(() => {
        actionJobSubscribe(undefined, undefined);
    }, [actionJobSubscribe])

    // *** SOCKET ACTIONS END ***

    // Fetches and sets initial data of pipeline and its jobs
    const fetchInitialData = useCallback(async () => {
        const headers = new Headers()
        headers.append("Authorization", "Bearer " + Cookies.get("session_token"));
        const req = await fetch(`${getLocation()}/api/pipeline/run/${params.pipelineId}`, {
            headers: headers
        })
        if (!req.ok) {
            throw new NetworkError(req.status)
        }

        const pipelineData = await req.json();

        
        setStats({
            title: pipelineData.title,
            totalJobs: pipelineData.jobs.length,
            pipelineStatus: pipelineData.pipelineStatus,
            startedBy: pipelineData.statistics.startedBy,
            startedAt: pipelineData.statistics.startedAt && new Date(pipelineData.statistics.startedAt),
            finishedAt: pipelineData.statistics.finishedAt && new Date(pipelineData.statistics.finishedAt),
            averageDuration: pipelineData.statistics.averageDuration,
            totalDuration: pipelineData.statistics.totalDuration,
            totalErrors: pipelineData.statistics.totalErrors,
        });
        setPermissions(parsePermissions(pipelineData.permissions));

        const dataObject = getVisualizerDataObject(pipelineData.jobs);
        setPipelineJobs(dataObject);
        setPipelineJobsByLevel(getVisualizerDataByLevel(dataObject));
    }, [params.pipelineId])

    // Fetches initial data and then connects to websocket
    const connectToServer = useCallback(async () => {
        const [promise, cancel] = retryCallback(fetchInitialData, 1500, retryNetwork);
        try {
            await promise
        } catch (err) {
            if (err instanceof NetworkError && (err.status === 403 || err.status === 401)) {
                setIsForbidden(true);
            }
            return null;
        }

        let socket;
        try {
            const [socketPromise, cancelSocket] = retryCallback(() => {
                const socket = new WebSocket(`${getLocation()}/api/pipeline/run/${params.pipelineId}/ws`);
                return socket;
            }, 1500, retryNetwork)
            
            socket = await socketPromise;
        } catch (err) {
            if (err instanceof NetworkError && (err.status === 403 || err.status === 401)) {
                setIsForbidden(true);
            }
            return null;
        }

        return socket;
    }, [fetchInitialData, params.pipelineId])

    // Close socket
    useEffect(() => {
        return () => {
            if (socket) {
                setSocket(null);
                socket.close();
            }
        }
    }, [socket])

    // Initiate stage information
    useEffect(() => {


        const isStageChanged = prevStage.current !== currentStage;

        if (currentStage) {
            const fetchData = async () => {
                const req = await fetch(`${getLocation()}/api/pipeline/run/${params.pipelineId}/job/${params.jobId}/stage/${currentStage}`)
                const stageData = await req.json();

                setStageLogs(stageData.logs);
                actionStageSubscribe(currentStage);
            }

            retryCallback(fetchData, 1000, 2);
            prevStage.current = currentStage;
        } else {
            if (isStageChanged) {
                actionStageUnsubscribe();
            }
        }
    }, [actionStageSubscribe, actionStageUnsubscribe, currentStage, isOnJob, params.jobId, params.pipelineId])

    // Connect to socket. Reconnect on death
    useEffect(() => {
        if (!socket) {
            const f = async () => {
                const sock = await connectToServer()
    
                setSocket(sock);
            }
            f();
        }
    }, [connectToServer, socket])


    // Initiate job information
    useEffect(() => {
        const jobId = +params.jobId;
        if (jobId && !isOnJob) {
            const fetchInitialData = async () => {
                const req = await fetch(`${getLocation()}/api/pipeline/run/${params.pipelineId}/job/${jobId}`)
                if (!req.ok) {
                    throw new NetworkError(req.status);
                }
                const jobData = await req.json();

                const jobStagesData = getJobStagesData(jobData.stages);
                setJobStages(jobStagesData);
                actionJobSubscribe(jobId);
                setIsOnJob(true);
            }
            const [promise, cancel] = retryCallback(fetchInitialData, 1500, NetworkError);
        } else if (!jobId && isOnJob) {
            actionJobUnsubscribe();
            setIsOnJob(false);
            setCurrentStage(null);
        }
    }, [actionJobSubscribe, actionJobUnsubscribe, isOnJob, params])

    // Initiate socket handlers
    useEffect(() => {
        if (socket && pipelineJobs) {
            socket.onopen = () => {
                //console.log('ws - open');
            }
            socket.onclose = async () => {
                //console.log('ws - close');
                setSocket(null)

            }
            socket.onmessage = (event) => {
                //console.log('ws - on - ', JSON.parse(event.data));
        
                const msg = JSON.parse(event.data);
        
                switch (msg.type) {
                    case 'jobStatus':
                        handleJobStatus(msg);
                        break;
                    case 'jobStatusBatch':
                        handleJobStatusBatch(msg);
                        break;
                    case 'stageStatus':
                        handleStageStatus(msg);
                        break;
                    case 'jobProgress':
                        handleJobProgress(msg);
                        break;
                    case 'stageLogs':
                        handleStageLogs(msg);
                        break;
                    default:
                        break;
                }
        
            }
        }
    }, [socket, pipelineJobs, jobStages, handleJobStatus, handleJobStatusBatch, handleStageStatus, handleJobProgress, handleStageLogs])

    const ret = {
        pipelineData: pipelineJobs,
        pipelineDataByLevels: pipelineJobsByLevels,
        jobId: params.jobId,
        jobStages: jobStages,
        stats: stats,
        setCurrentStage: setCurrentStage,
        currentStage: currentStage,
        stageLogs: stageLogs,
        permissions: permissions,
        isForbidden: isForbidden,

        actionActPipeline,
        actionActJob,
    };

    return (
        <PipelineContext.Provider value={ret}>
            {children}
        </PipelineContext.Provider>
    )
}