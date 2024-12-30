import getStatusColor from "@/app/helpers/getStatusColor";
import CircleIcon from "../icons/CircleIcon";
import ClockIcon from "../icons/ClockIcon";
import ExclamationIcon from "../icons/ExclamationIcon";
import TickIcon from "../icons/TickIcon";
import StatusIcon from "../StatusIcon";
import statusToText from "@/app/helpers/statusToText";

export default function JobInfoTag({status, size = 'sm', progress, className}) {

    return (
        <div className={`${className} flex ${size === 'sm' ? 'text-sm' : 'text-base'} font-semibold flex-row items-center px-2 py-1 rounded-lg border ${getStatusColor(status)}`}>
            <StatusIcon
                className="w-5 h-5 mr-2"
                status={status}
                progress={status > 99 ? 0 : progress}
            />
            <span>
                {statusToText(status)}
            </span>
        </div>
    )
}