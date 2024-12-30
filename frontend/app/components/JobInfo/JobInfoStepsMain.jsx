'use client';
import { useContext, useEffect, useState } from "react";
import ArrowIcon from "../icons/ArrowIcon"
import StatusIcon from "../StatusIcon"
import JobInfoStepsTab from "./JobInfoStepsTab"
import { PipelineContext } from "../contexts/PipelineContext";

export default function JobInfoStepsMain({stages, settings}) {

    const {setCurrentStage, currentStage, stageLogs} = useContext(PipelineContext)


    function handleClick(stageId) {
        if (stageId === currentStage) {
            setCurrentStage(null);
        } else {
            setCurrentStage(stageId);
        }
    }

    const itemsEl = Object.values(stages).map((stage, index) => {
        return (
            <li key={stage.id}>
                <JobInfoStepsTab
                    isOpened={stage.id === currentStage}
                    onClick={() => handleClick(stage.id)}
                    title={stage.title}
                    status={stage.status}
                    startedAt={stage.startedAt}
                    finishedAt={stage.finishedAt}
                    lines={stageLogs}
                    settings={settings}
                />
            </li>
        )
    })

    return (
        <div>
            <ul className="flex flex-col gap-1">
                {itemsEl}
            </ul>
        </div>
    )
}



