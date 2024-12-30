export default function parsePermissions(str) {
    return {
        canRead: str[3] === '1',
        canWrite: str[2] === '1',
        canExecute: str[1] === '1',
        canAdmin: str[0] === '1',
    };
}