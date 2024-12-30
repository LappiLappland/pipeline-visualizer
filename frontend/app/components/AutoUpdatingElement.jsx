import { useEffect, useState } from "react";

export default function AutoUpdatingComponent({ initialValue, callback, getNextInterval }) {
    const [displayValue, setDisplayValue] = useState(() => callback(initialValue));
  
    useEffect(() => {
        let timeoutId;
    
        const scheduleNextUpdate = () => {
            const interval = getNextInterval(initialValue);
    
            if (interval !== null) {
            setDisplayValue(callback(initialValue));
            timeoutId = setTimeout(scheduleNextUpdate, interval);
            }
        };
    
        scheduleNextUpdate();
    
        return () => clearTimeout(timeoutId);
    }, [initialValue, getNextInterval, callback]);

    return displayValue;
}

export function getNextIntervalByDate(initialDate) {
    const now = new Date();
    const elapsedMs = now - initialDate;
    const elapsedSeconds = elapsedMs / 1000;
  
    if (elapsedSeconds > 3600) {
      return null;
    }
  
    // Exponential increase
    const interval = Math.min(1000 * Math.pow(2, Math.floor(elapsedSeconds / 10)), 1000);
    return interval;
}