import CancelledIcon from "./icons/CancelledIcon";
import CircleIcon from "./icons/CircleIcon";
import CircleOutlineIcon from "./icons/CircleOutlineIcon";
import ClockIcon from "./icons/ClockIcon";
import ExclamationIcon from "./icons/ExclamationIcon";
import ProgressIcon from "./icons/ProgressIcon";
import TickIcon from "./icons/TickIcon";

export default function StatusIcon({className, status, progress}) {
    switch (status) {
        case 'completed':
            return (
                <TickIcon
                    className={`${className} fill-green-500`}
                />
            );
        case 'failed':
            return (
                <ExclamationIcon
                    className={`${className} fill-red-55`}
                />
            );
        case 'running':
            if (progress == null) {
                return (
                    <ClockIcon
                        className={`${className} fill-blue-500`}
                    />
                )
            } else {
                return (
                    <ProgressIcon
                        className={`${className} bg-blue-500`}
                        progress={progress}
                    />
                )
            }
        case 'cancelled':
            return (
                <CancelledIcon
                    className={`${className} fill-gray-500`}
                />
            );
        case 'planned':
            return (
                <CircleOutlineIcon
                    className={`${className} fill-white bg-amber-500`}
                />
            );
        case 'pending':
        default:
            return (
                <CircleIcon
                    className={`${className} fill-gray-400`}
                />
            );
    }
}