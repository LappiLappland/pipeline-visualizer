import Link from "next/link";
import JobInfoRow from "./JobInfoRow";
import JobInfoTag from "./JobInfoTag";
import JobInfoStepsMain from "./JobInfoStepsMain";
import { useContext, useState } from "react";
import { PipelineContext } from "../contexts/PipelineContext";
import { formatDistanceToNow, formatDuration, intervalToDuration } from "date-fns";

import AutoUpdatingComponent, { getNextIntervalByDate } from "../AutoUpdatingElement";
import MultiSelect from "../MultiSelect";

export default function JobInfo() {

    const {pipelineData: dataJobs, jobId, jobStages} = useContext(PipelineContext);
    const [showOptions, setShowOptions] = useState(false);

    const [settings, setSettings] = useState(() => ({
        showErrors: true,
        showWarnings: true,
        showInformations: true,
        showTimestamps: false,
    }))

    const options = [
        {id: 'showTimestamps', title: 'Show timestamps', active: settings.showTimestamps},
        {id: 'showErrors', title: 'Show errors', active: settings.showErrors},
        {id: 'showWarnings', title: 'Show warnings', active: settings.showWarnings},
        {id: 'showInformations', title: 'Show informations', active: settings.showInformations},
    ]

    const jobInfo = dataJobs ? dataJobs[jobId] : {};
    const title = jobInfo?.title ?? '';

    const settingsDistance = {
        addSuffix: true,
        includeSeconds: true,
        
    };

    function handleSelectOption(id, active) {
        const newSettings = {...settings};
        newSettings[id] = active;
        setSettings(newSettings);
    }

    return (
        <div className="flex flex-col lg:flex-row mt-1">
            <div className="grow">
                <h2 className="text-2xl font-bold mb-1">
                    {title}
                </h2>
                <div className="flex sm:items-center sm:justify-between mb-2 flex-col sm:flex-row">
                    <div className="flex md:items-center flex-col md:flex-row">
                        <div className="w-max">
                            <JobInfoTag 
                                status={jobInfo.status}
                                progress={jobInfo.progress}
                            />
                        </div>
                        
                        {!jobInfo.startedAt ? '' : (
                            <span className="ml-1.5">
                                <AutoUpdatingComponent
                                    initialValue={jobInfo.startedAt}
                                    callback={(v) => 'Started at ' + formatDistanceToNow(v, settingsDistance)}
                                    getNextInterval={getNextIntervalByDate}
                                />
                            </span>
                            
                        )}
                        {!jobInfo.startedBy ? '' : (
                            <span className="ml-1 font-bold">
                                {jobInfo.startedBy}
                            </span>
                        )}
                    </div>
                    <div className="mt-1 sm:mt-0 self-start sm:self-auto">
                        <MultiSelect
                            options={options}
                            onClick={() => setShowOptions((prev) => !prev)}
                            onClosed={() => setShowOptions(false)}
                            onSelect={handleSelectOption}
                            isShown={showOptions}
                        >
                            Фильтры
                        </MultiSelect>
                    </div>
                </div>
                <div className="w-full p-1 bg-gray-95 border-gray-85 border rounded-lg">
                    {!jobStages ? '' : (
                        <JobInfoStepsMain
                            stages={jobStages}
                            settings={settings}
                        />
                    )}
                </div>
            </div>
            <div className="w-full lg:w-auto lg:max-w-2/5 pt-3.5 lg:pt-0 lg:pl-8 pr-4">
                <ul className="flex flex-col gap-2">
                    <JobInfoRow>
                        <>Duration:</>
                        <>
                        {
                            jobInfo.startedAt && jobInfo.finishedAt ? (
                                formatDuration(intervalToDuration({ start: jobInfo.startedAt, end: jobInfo.finishedAt }), {})
                            ) : '-'
                        }
                        </>
                    </JobInfoRow>
                    <JobInfoRow>
                        <>Started at:</>
                        <>
                        {
                            jobInfo.startedAt ? (
                                jobInfo.startedAt.toLocaleString()
                            ) : '-'
                        }
                        </>
                    </JobInfoRow>
                    <JobInfoRow>
                        <>Finished at:</>
                        <>
                        {
                            jobInfo.finishedAt ? (
                                jobInfo.finishedAt.toLocaleString()
                            ) : '-'
                        }
                        </>
                    </JobInfoRow>
                </ul>
            </div>
        </div>
    )
}