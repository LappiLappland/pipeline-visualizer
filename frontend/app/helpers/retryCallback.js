export default function retryCallback(callback, timeInMs, retryCallback) {
    let cancelled = false;
    const promise = new Promise((resolve, reject) => {
        let attempts = 0;
        async function attempt() {
            attempts++;

            if (cancelled) {
                resolve();
            }

            try {
                const result = await callback();
                resolve(result);
            } catch(err) {
                const again = retryCallback(attempts, err);
                
                if (again) {
                    setTimeout(attempt, timeInMs);
                } else {
                    reject(err);
                }
            }
        }

        attempt();
    });

    const cancelFunc = () => {cancelled = true;};

    return [promise, cancelFunc];
}

export class NetworkError extends Error {
    constructor(status) {
        super();
        this.status = status;
    }
}

export function retryNetwork(attempt, prevErr) {
    if (attempt > 5) {
        return false;
    }

    if (prevErr instanceof NetworkError) {
        if (prevErr.status === 403 || prevErr.status === 401) {
            return false;
        }

        return true;
    }

    return false;
}