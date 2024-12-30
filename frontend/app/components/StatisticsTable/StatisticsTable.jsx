import JobInfoTag from "../JobInfo/JobInfoTag";

export default function StatisticsTable({headers, data}) {

    const headersEl = headers.map((header) => {
        return (
            <th className="px-1 border-l last:border-r text-white bg-sky-600 border-sky-700"
            key={header.id}
            >
                {header.name}
            </th>
        )
    });
    const rowsEl = data.map((row, i) => {
        const isPipeline = i === data.length - 1;

        const columnsEl = headers.map((header) => {
            const columnData = row[header.id];
            const isDate = header.type && header.type === 'date' && columnData;
            const isStatus = header.type && header.type === 'status';
            const isString = typeof columnData === 'string' && !isDate && !isStatus;
            let extraClassName = 'text-center';
            let text = columnData;
            if (columnData == null) {
                text = '-';
            } else if (isString) {
                extraClassName = 'text-left';
            } else if (isDate) {
                text = new Date(columnData).toLocaleString();
            } else if (isStatus) {
                text = <JobInfoTag status={columnData} size="sm" className={"justify-center text-black"} />;     
                extraClassName = 'px-4';
            }

            extraClassName = isPipeline ? extraClassName + ' border-indigo-500' : extraClassName;

            return (
                <td className={`p-1 border-l last:border-r ${extraClassName}`}
                key={header.id}
                >
                    {text}
                </td>
            )
        })



        return (
            <tr className={`border-t last:border-b
                ${isPipeline ? ' font-bold bg-indigo-400 border-indigo-500 text-white' : ''}    
            `}
                key={!isPipeline ? row.id : 'pipeline'}
            >
                {columnsEl}
            </tr>
        )   
    })

    return (
        <div className="w-full overflow-auto py-2">
            <table className="w-full">
                <thead>
                    <tr className="border-t border-blue-55">
                        {headersEl}
                    </tr>
                </thead>
                <tbody>
                    {rowsEl}
                </tbody>
            </table>
        </div>
        
    )
}