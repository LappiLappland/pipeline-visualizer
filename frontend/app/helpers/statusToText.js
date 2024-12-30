
export default function statusToText(status) {
    switch (status) {
        case 'completed':
            return 'completed';
        case 'failed':
            return 'failed';
        case 'running':
            return 'running';
        case 'cancelled':
            return 'cancelled';
        case 'planned':
            return 'planned';
        case 'pending':
        default:
            return 'pending';
    }
}