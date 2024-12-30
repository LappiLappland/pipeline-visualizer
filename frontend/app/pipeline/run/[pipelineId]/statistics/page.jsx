'use client'
import Button from "@/app/components/Button";
import ButtonLink from "@/app/components/ButtonLink";
import { PipelineContext } from "@/app/components/contexts/PipelineContext";
import DownloadIcon from "@/app/components/icons/DownloadIcon";
import StatisticsTable from "@/app/components/StatisticsTable/StatisticsTable";
import getLocation from "@/app/helpers/getLocation";
import Image from "next/image";
import Link from "next/link";
import { useParams } from "next/navigation";
import { useContext, useEffect, useState } from "react";

export default function Statistics() {

    const {stats} = useContext(PipelineContext)
    const params = useParams()

    const [tableData, setTableData] = useState(null);

    useEffect(() => {
        const fetchData = async () => {
            const req = await fetch(`${getLocation()}/api/pipeline/run/${params.pipelineId}/statistics`)
            const data = await req.json();
        
            setTableData(data);
        }

        fetchData();
    }, [params.pipelineId])

    function getWarningMessage() {
        let text = ''
        if (stats.pipelineStatus === 'running') {
            text = "Warning: current statistics are based on still running pipeline!"
        } else if (stats.pipelineStatus === 'pending') {
            text = "warning: current statistics are based on not started pipeline!"
        }
        return (
            <div className="text-xl font-bold text-red-500">
                {text}
            </div>
        )
    }

    return (
        <div className="">
            <div className="flex flex-col gap-1 border-2 border-gray-85 rounded-xl p-3">
                <Link className="text-gray-800 hover:underline" 
                href={`/pipeline/run/${params.pipelineId}`}
                >
                    {'<'} Back to pipeline
                </Link>
                {getWarningMessage()}
                <div>
                    <ButtonLink
                        className="py-1.5 flex items-center gap-2 w-max justify-center"
                        href={`${getLocation()}/api/pipeline/run/${params.pipelineId}/statistics/excel`}
                        download={true}
                    >
                        <DownloadIcon 
                            className="w-4 h-4 fill-gray-700"
                        />
                        Export to Excel
                    </ButtonLink>
                </div>
                {!tableData ? '' : (
                    <StatisticsTable
                        headers={tableData.headers}
                        data={tableData.data}
                    />
                )}
            </div>
        </div>
    );
}
