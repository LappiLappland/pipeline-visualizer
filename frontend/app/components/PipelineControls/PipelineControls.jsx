'use client'
import { useContext } from "react";
import Button from "../Button";
import CrossIcon from "../icons/CrossIcon";
import EditorIcon from "../icons/EditorIcon";
import RestartIcon from "../icons/RestartIcon";
import StartIcon from "../icons/StartIcon";
import StatisticsIcon from "../icons/StatisticsIcon";
import { PipelineContext } from "../contexts/PipelineContext";
import ButtonLink from "../ButtonLink";
import { useParams } from "next/navigation";

export default function PipelineControls() {

    const {actionActPipeline, pipelineData, stats, permissions} = useContext(PipelineContext);
    const params = useParams()

    function getControlButtons() {
        if (!permissions.canExecute) {
            return '';
        }
        switch (stats.pipelineStatus) {
            case 'pending':
                return (
                    <>
                        <StartButton
                            onClick={() => actionActPipeline('start')}
                        />
                    </>
                )
            case 'completed':
                return (
                    <>
                        <RestartButton 
                            onClick={() => actionActPipeline('restart')}
                        />
                    </>
                )
            case 'cancelled':
            case 'failed':
                return (
                    <>
                        <RestartButton 
                            onClick={() => actionActPipeline('restart')}
                        />
                        <ContinueButton
                            onClick={() => actionActPipeline('continue')}
                        />
                    </>
                )
            case 'planned':
            case 'running':
                return (
                    <>
                        <StopButton 
                            onClick={() => actionActPipeline('stop')}
                        />
                    </>
                    
                )
            default:
                break;
        }
    }

    return (
        <div className="py-2 px-6 flex justify-between flex-wrap gap-1.5 text-sm sm:text-base">
            <div className="flex gap-2.5 flex-wrap">
                {getControlButtons()}
            </div>
            <div className="flex gap-2.5 flex-wrap">
                {/* <Button className="h-10 flex  gap-0 sm:gap-1.5 items-center justify-center">
                    <EditorIcon className="w-4 h-4 sm:w-6 sm:h-6" />
                    <span className="hidden sm:inline">
                        Editor    
                    </span>
                </Button> */}
                {!permissions.canRead || stats.pipelineStatus === 'pending' ? '' : (
                    <ButtonLink className="h-10 flex  gap-0 sm:gap-1.5 items-center justify-center"
                    href={`/pipeline/run/${params.pipelineId}/statistics`}
                    >
                        <StatisticsIcon className="w-4 h-4 sm:w-6 sm:h-6" />
                        <span className="hidden sm:inline">
                            Statistics
                        </span>
                    </ButtonLink>
                )}
                
            </div>
        </div>
    )
}

function RestartButton({onClick}) {
    return (
        <Button className="h-10 flex gap-0 sm:gap-1.5 items-center justify-center"
            onClick={onClick}
        >
            <RestartIcon className="w-4 h-4 sm:w-6 sm:h-6" />
            <span className="hidden sm:inline">
                Restart
            </span>
        </Button>
    )
}

function ContinueButton({onClick}) {
    return (
        <Button className="h-10 flex gap-0 sm:gap-1.5 items-center justify-center"
            onClick={onClick}
        >
            <StartIcon className="w-4 h-4 sm:w-6 sm:h-6" />
            <span className="hidden sm:inline">
                Continue
            </span>
        </Button>
    )
}

function StartButton({onClick}) {
    return (
        <Button className="h-10 flex gap-0 sm:gap-1.5 items-center justify-center"
            onClick={onClick}
        >
            <StartIcon className="w-4 h-4 sm:w-6 sm:h-6" />
            <span className="hidden sm:inline">
                Start  
            </span>
        </Button>
    )
}

function StopButton({onClick}) {
    return (
        <Button className="h-10 flex gap-0 sm:gap-1.5 items-center justify-center"
            onClick={onClick}
        >
            <CrossIcon className="w-4 h-4 sm:w-6 sm:h-6" />
            <span className="hidden sm:inline">
                Stop  
            </span>
        </Button>
    )
}