import { formatDistanceToNow } from "date-fns";

import AutoUpdatingComponent, { getNextIntervalByDate } from "../AutoUpdatingElement";
import statusToText from "@/app/helpers/statusToText";

export default function VisualizerNodeHint({
    isShown,
    title,
    status,
    executor,
    startedAt,
    finishedAt,
}) {

    const settings = {
        addSuffix: true,
        includeSeconds: true,
        
    };

    return (
        <div className={`absolute
            flex justify-center
            top-0
            translate-y-[58px]
            ${isShown ? 'opacity-100' : 'opacity-0'}
            transition-opacity
            rounded-lg
            z-10
            pointer-events-none
            px-3
            bg-black
            `}
        >
            <div className={`absolute w-3 h-3 border-[7px]
                rotate-45 border-b-0 border-r-0 border-black
                top-[-6px]
            `} />
            <div className="flex flex-col py-2 text-xs text-white font-medium"
            >
                <div className="group">
                    <HintEntry>
                        {title}
                    </HintEntry>
                    <HintEntry>
                        {statusToText(status)}
                    </HintEntry>
                </div>
                {!executor ? '' : (
                    <span>
                        {'Started by ' + executor}
                    </span>
                )}
                {!startedAt || !isShown ? '' : (
                    <span>
                        <AutoUpdatingComponent 
                            initialValue={startedAt}
                            callback={(v) => 'Started ' + formatDistanceToNow(v, settings)}
                            getNextInterval={getNextIntervalByDate}
                        />
                    </span>
                )}
                {!finishedAt || !isShown ? '' : (
                    <span>
                        <AutoUpdatingComponent 
                            initialValue={finishedAt}
                            callback={(v) => 'Finished ' + formatDistanceToNow(v, settings)}
                            getNextInterval={getNextIntervalByDate}
                        /> 
                    </span>
                )}
            </div>
        </div>
    )
}

function HintEntry({children}) {
    return (
        <span className={`
            after:content-['-']
            after:mx-1
            last:after:content-['']
            last:after:mx-0
        `}>
            {children}
        </span>
    )
}