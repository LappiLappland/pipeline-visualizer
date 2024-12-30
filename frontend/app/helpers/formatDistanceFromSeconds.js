import { formatDuration, intervalToDuration } from "date-fns";
import { ru } from "date-fns/locale/ru"

export default function formatRelativeTimeFromSeconds(seconds, locale = 'en') {
    const duration = intervalToDuration({ start: 0, end: seconds * 1000 });
    return formatDuration(duration, {});
}