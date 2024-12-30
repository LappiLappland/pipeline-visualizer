import { useCallback, useContext, useEffect, useRef } from "react";
import CircleIcon from "../icons/CircleIcon";
import ClockIcon from "../icons/ClockIcon";
import CrossIcon from "../icons/CrossIcon";
import ExclamationIcon from "../icons/ExclamationIcon";
import RestartIcon from "../icons/RestartIcon";
import StartIcon from "../icons/StartIcon";
import TickIcon from "../icons/TickIcon";
import Link from "next/link";
import ProgressIcon from "../icons/ProgressIcon";
import VisualizerNodeHint from "./VisualizerNodeHint";
import { PipelineContext } from "../contexts/PipelineContext";
import StatusIcon from "../StatusIcon";
import getStatusColor from "@/app/helpers/getStatusColor";

const minDelay = 200;

export default function VisualizerNode({
    title,
    id,
    pipelineId = 0,
    status = 'pending',
    canExecute = false,
    startedAt,
    finishedAt,
    startedBy,
    progress,
    isUnhovered,
    isHovered,
    isMainHovered,
    onMouseEnter,
    onMouseLeave,
}) {

    const {actionActJob} = useContext(PipelineContext);

    const pressTimer = useRef(null);
    const firstPress = useRef(0);

    const unhoverClass = isUnhovered ? 'opacity-50 -z-10' : '';
    const hoverClass = isHovered ? 'shadow-lg' : '';

    useEffect(() => {
        return () => {
            if (pressTimer.current) {
                clearTimeout(pressTimer.current);
            }
        };
    }, []);

    function getActionButton() {
        if (!canExecute) {
            return '';
        }

        let icon = <></>;
        let action = () => undefined;

        switch (status) {
            case 'failed':
            case 'completed':
            case 'cancelled':
                icon = (
                    <RestartIcon 
                        className="w-4 h-4"
                    />
                );
                action = (e) => {
                    e.preventDefault();
                    if (e.timeStamp - firstPress.current > minDelay) {
                        return;
                    }
                    actionActJob(id, 'restart');
                }
                break;
            case 'running':
            case 'planned':
                icon = (
                    <CrossIcon
                        className="w-4 h-4"
                    />
                );
                action = (e) => {
                    e.preventDefault();
                    if (e.timeStamp - firstPress.current > minDelay) {
                        return;
                    }
                    actionActJob(id, 'stop')
                }
                break;
            case 'pending':
                icon = (
                    <StartIcon 
                        className="w-4 h-4 fill-gray-25"
                    />
                );
                action = (e) => {
                    e.preventDefault();
                    if (e.timeStamp - firstPress.current > minDelay) {
                        return;
                    }
                    actionActJob(id, 'start')
                }
                break;
        }

        return (
            <button className="stroke-gray-25 rounded-full p-1 hover:bg-gray-85/75"
                onClick={action}
                onMouseDown={downHandler}
                onTouchStart={downHandler}
            >
                {icon}
            </button>
        )
    }

    function downHandler(e) {
        firstPress.current = e.timeStamp;

        if (e.touches) {
            pressTimer.current = setTimeout(() => {
                onMouseEnter();
            }, minDelay);
        }
        
    }

    function upHandler(e) {
        clearTimeout(pressTimer.current);
        pressTimer.current = null;
        onMouseLeave();
    }

    function clickHandler(e) {
        if (e.timeStamp - firstPress.current > minDelay) {
            e.preventDefault();
            return;
        }
    }

    return (
        <div className={`border rounded-lg relative w-[300px]
            ease-out transition-hide items-center
            cursor-pointer
            flex flex-col
            ${unhoverClass}
            ${hoverClass}
            ${getStatusColor(status)}
        `}
        onMouseLeave={onMouseLeave}
        >
            <Link className={`select-none w-full px-3 py-3 flex flex-row justify-between items-center
            `}
            href={`/pipeline/run/${pipelineId}/job/${id}`}
            draggable="false"
            onMouseEnter={onMouseEnter}
            onClick={clickHandler}
            onMouseDown={downHandler}
            onTouchStart={downHandler}
            onMouseUp={upHandler}
            onTouchEnd={upHandler}
            >
                <div className="flex flex-row items-center">
                    <StatusIcon
                        className="w-5 h-5 mr-2"
                        status={status}
                        progress={progress}
                    />
                    <span className="font-sm">
                        {title ?? id}
                    </span>
                </div>
                {getActionButton()}
            </Link>
            <VisualizerNodeHint 
                isShown={isMainHovered}
                title={title}
                executor={startedBy}
                status={status}
                startedAt={startedAt}
                finishedAt={finishedAt}
            />
        </div>
    )
}