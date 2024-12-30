import { useContext } from "react";
import { PipelineContext } from "./contexts/PipelineContext";
import JobInfoTag from "./JobInfo/JobInfoTag";
import { formatDistanceToNow } from "date-fns";
import formatRelativeTimeFromSeconds from "../helpers/formatDistanceFromSeconds";

import AutoUpdatingComponent, { getNextIntervalByDate } from "./AutoUpdatingElement";

export default function StatusBar() {
    
    const {pipelineData, stats} = useContext(PipelineContext);

    const settings = {
        addSuffix: true,
        includeSeconds: true,
        
    };
    
    const startedAt = stats.startedAt;
    const finishedAt = stats.finishedAt;
    
    const averageDuration = stats.averageDuration && formatRelativeTimeFromSeconds(stats.averageDuration)
    const totalDuration = stats.totalDuration && formatRelativeTimeFromSeconds(stats.totalDuration)
    return (
        <div className="flex flex-col lg:flex-row lg:items-center border-2 mb-3 py-2 px-3 border-gray-85 rounded-xl">
            <div className="lg:w-1/4 w-full">
                <h2 className="text-2xl font-bold mb-1">
                    {stats.title ?? ''}
                </h2>
                {!stats.startedBy ? '' : (
                    <div className="flex items-center">
                        <span className="text-gray-800">
                            Started by
                        </span>
                        <span className="ml-1 font-bold">
                            {stats.startedBy}
                        </span>
                    </div>
                )}
                
                <div className="text-sm h-5">
                    {!startedAt ? '' : (
                        <AutoUpdatingComponent 
                            initialValue={startedAt}
                            callback={(v) => 'Started at ' + formatDistanceToNow(v, settings)}
                            getNextInterval={getNextIntervalByDate}
                        />
                        
                    )}
                </div>
                <div className="text-sm h-5">
                    {!finishedAt ? '' : (
                        <AutoUpdatingComponent 
                            initialValue={finishedAt}
                            callback={(v) => 'Finished at ' + formatDistanceToNow(v, settings)}
                            getNextInterval={getNextIntervalByDate}
                        /> 
                    )}
                </div>
            </div>
            <div className="grid md:grid-cols-2 lg:grid-cols-5 gap-2 grow lg:place-items-center">
                <HeaderColumn hideOnSmall={true} title="Статус" className={"mt-1 place-self-start lg:place-self-auto col-span-2 lg:col-span-1 "}>
                    {!stats.pipelineStatus ? '' : (
                        <div>
                            <JobInfoTag
                                className="w-max"
                                status={stats.pipelineStatus}
                            />
                        </div>
                        
                    )}
                    
                </HeaderColumn>
                <HeaderColumn title="Total jobs" >
                    <div className=" font-bold text-lg py-0.5">
                        {stats.totalJobs || '-'}
                    </div>
                </HeaderColumn>
                <HeaderColumn title="Total errors" >
                    <div className=" font-bold text-lg py-0.5">
                        {stats.totalErrors ?? '-'}
                    </div>
                </HeaderColumn>
                <HeaderColumn title="Average duration" >
                    <div className=" font-bold text-lg py-0.5">
                        {averageDuration || '-'}
                    </div>
                </HeaderColumn>
                <HeaderColumn title="Total duration" >
                    <div className=" font-bold text-lg py-0.5">
                        {totalDuration || '-'}
                    </div>
                </HeaderColumn>
            </div>
        </div>
    )
}

function HeaderColumn({children, title, className, hideOnSmall}) {
    return (
        <div className={className + ' w-full'}>
            <span className={`text-gray-800 text-sm ${hideOnSmall ? 'hidden lg:block' : 'col-span-2 md:col-span-1'}`}>
                {title}
            </span>
            {children}
        </div>
    )
}