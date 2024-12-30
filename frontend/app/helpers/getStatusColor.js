export default function getStatusColor(status) {
    switch (status) {
        case 'completed':
            return 'bg-green-100 border-green-200';
        case 'failed':
            return 'bg-red-100 border-red-200';
        case 'running':
            return 'bg-blue-100 border-blue-200';
        case 'cancelled':
            return 'bg-gray-200 border-gray-300'
        case 'planned':
            return 'bg-amber-200 border-amber-300'
        case 'pending':
        default:
            return 'bg-gray-100 border-gray-200';
    }
}