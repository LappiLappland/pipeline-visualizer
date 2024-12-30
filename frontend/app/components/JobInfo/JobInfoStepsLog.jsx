export default function JobInfoStepsLog({lines, settings}) {

    lines = lines ? filterLines(lines) : [];

    function filterLines(lines) {
        return lines.filter((line) => {
            return line.type === 'base'
                || (line.type === 'warn' && settings.showWarns) 
                || (line.type === 'error' && settings.showErrors)
                || (line.type === 'info' && settings.showInformations)
        })
    }

    function getLineExtras(type) {
        switch (type) {
            case 'warn':
                return {
                    text: 'text-orange-600',
                    bg: 'bg-orange-200/50',
                    hover: 'hover:bg-orange-300/50',
                    msg: 'Warn:',
                };
            case 'info':
                return {
                    text: 'text-blue-500',
                    bg: 'bg-blue-200/50',
                    hover: 'hover:bg-blue-300/50',
                    msg: 'Info:',
                };
            case 'error':
                return {
                    text: 'text-red-600',
                    bg: 'bg-red-200/50',
                    hover: 'hover:bg-red-300/50',
                    msg: 'Error:',
                };
            case 'base':
            default:
                return {
                    text: '',
                    bg: '',
                    hover: 'hover:bg-gray-85/50',
                    msg: '',
                };
        }
    }

    const linesComponents = lines.map((line, row) => {

        const lineExtras = getLineExtras(line.type);
        const gridClassName = settings.showTimestamps ? 'grid-cols-[min-content_max-content_1fr]' : 'grid-cols-[min-content_1fr]';

        return (
            <div key={line.line} className={`${lineExtras.bg} ${lineExtras.hover} ${gridClassName} grid py-0.5 pl-3`}>
                <div className="w-6 md:w-12 text-right mr-2 md:mr-4">
                    {line.line}
                </div>
                {!settings.showTimestamps ? '' : (
                    <div className="mr-1.5">
                        {new Date(line.createdAt).toLocaleString()}
                    </div>
                )}
                
                <div className={`$`}>
                    {!lineExtras.msg ? '' : (
                        <span className={`${lineExtras.text} mr-1.5`}>
                            {lineExtras.msg}
                        </span>
                    )}
                    
                    <span>
                        {line.text}
                    </span>
                </div>
            </div>
        )
    })

    return (
        <div className="text-sm mt-2">
            {linesComponents}
        </div>
    )
}