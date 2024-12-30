'use client'
import Link from "next/link";
import JobInfo from "../../../../../components/JobInfo/JobInfo";
import ArrowIcon from "@/app/components/icons/ArrowIcon";
import { useParams } from "next/navigation";

export default function JobPage({}) {

    const params = useParams();

    return (
        <div className="border-2 border-gray-85 rounded-xl px-3 py-2">
            <Link className="text-gray-800 hover:underline" 
            href={`/pipeline/run/${params.pipelineId}`}
            >
                {'<'} Back to pipeline
            </Link>
            <JobInfo />
        </div>
    )
}