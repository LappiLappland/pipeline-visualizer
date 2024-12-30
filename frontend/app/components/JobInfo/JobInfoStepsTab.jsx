'use client'
import { useCallback, useContext, useEffect, useRef, useState } from "react";
import ArrowIcon from "../icons/ArrowIcon";
import StatusIcon from "../StatusIcon";
import JobInfoStepsLog from "./JobInfoStepsLog";
import { PipelineContext } from "../contexts/PipelineContext";
import { formatDuration, intervalToDuration } from "date-fns";

import formatTime from "@/app/helpers/formatTime";
import { isBottomVisible, scrollToBottomOnBelowView } from "@/app/helpers/dom";

export default function JobInfoStepsTab({settings, title, status, isOpened, onClick, lines, startedAt, finishedAt}) {

    const [logHeight, setLogHeight] = useState(0);
    const [timeText, setTimeText] = useState('');
    const localTimer = useRef();
    const prevLen = useRef(0);
    const logsContainer = useRef();
    const margin = 8;

    // Set text of time
    useEffect(() => {
        // If we know both, then calculate time between them
        if (startedAt && finishedAt) {
            clearInterval(localTimer.current)
            localTimer.current = undefined;
            const duration = finishedAt - startedAt;
            setTimeText(formatTime(duration));
        // If we know only start, then we have to update time each second
        } else if (startedAt) {
            const duration = Date.now() - startedAt;
            setTimeText(formatTime(duration));
            localTimer.current = setInterval(() => {
                const duration = Date.now() - startedAt;
                setTimeText(formatTime(duration));
            }, 950);
        }
    }, [startedAt, finishedAt]);

    // Clean up timer
    useEffect(() => {
        return () => {
            if (localTimer.current) {
                clearInterval(localTimer.current);
                localTimer.current = undefined;
            }
        };
    }, []);

    // Get height of container. Needed for smooth min-height change
    const computeHeight = useCallback(
        (node) => {
            if (!node) return;
            setLogHeight(node.children[0].getBoundingClientRect().height + margin);
            logsContainer.current = node;
        },
        // eslint-disable-next-line react-hooks/exhaustive-deps
        [lines], // It is necessary !
    );

    // Scroll to the bottom of logs, if new lines appear
    useEffect(() => {
        if (lines.length > 0 && logsContainer.current) {
            if (!isBottomVisible(logsContainer.current.children[0], 10)) {
                return;
            }
            if (prevLen.current < lines.length && prevLen.current > 0) {
                scrollToBottomOnBelowView(logsContainer.current.children[0], logHeight, 10);
            }
            prevLen.current = lines.length;
        }
    }, [lines, logHeight])

    function handleClick() {
        if (status !== 'pending' && status !== 'planned') {
            onClick()
        }
    }

    return (
        <div
        >
            <button className={`w-full flex flex-row items-center justify-between py-1 px-3 rounded-lg 
                ${isOpened ? 'bg-gray-300' : 'hover:bg-gray-200'}`
            }
            onClick={handleClick}
            >
                <div className="flex flex-row items-center font-semibold">
                    {status === 'pending' ? (
                        <span className="w-6 mr-2.5" />
                    ) : (
                        <ArrowIcon direction={isOpened ? 'down' : 'right'} className="mr-2.5 transition-transform" />
                    )}
                    <StatusIcon className="w-5 h-5 mr-2.5" status={status} />
                    {title}
                </div>
                <div>
                    {timeText}
                </div>
            </button>
            <div
            ref={computeHeight}
            className={`overflow-hidden ${isOpened ? 'mb-3' : ''}
            transition-max-height duration-300 ease-in-out`}
            style={{ maxHeight: isOpened ? logHeight : 0 }}
            >
                <JobInfoStepsLog lines={lines} settings={settings} />
            </div>
        </div>
    )
}