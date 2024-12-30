import AsideMenu, { LinksGroup } from "@/app/components/AsideMenu"
import { PipelineContext } from "@/app/components/contexts/PipelineContext"
import StatusBar from "@/app/components/StatusBar"
import StatusIcon from "@/app/components/StatusIcon"
import Link from "next/link"
import { useParams } from "next/navigation"
import { useContext } from "react"
import Forbidden from "../../forbidden"
import Loader from "@/app/components/Loader"

export default function PipelineMain({children}) {

    const {pipelineDataByLevels, pipelineData, isForbidden} = useContext(PipelineContext)
    const params = useParams();

    const linksGeneral = [
        {title: 'Pipeline', href: `/pipeline/run/${params.pipelineId}`, isCurrent: params.jobId == null},
    ]
    const jobs = Object.values(pipelineDataByLevels || {}).flatMap((level) => {
        return level.map((jobId) => {
            const job = pipelineData[jobId];
            return {
                title: job.title,
                href: `/pipeline/run/${params.pipelineId}/job/${job.id}`,
                isCurrent: params.jobId && (params.jobId === job.id+''),
                icon: <StatusIcon className="w-5 h-5 mr-2.5" status={job.status} progress={job.progress} />
            }
        })
    })

    if (isForbidden) {
        return (
            <div className="flex flew-row px-12 py-6">
                <AsideMenu>
                    <Link className="text-gray-800 hover:underline pb-1 border-r-2" 
                    href={`/pipeline`}
                    >
                        {'<'} Go back
                    </Link>
                </AsideMenu>
                <div className="grow w-1/2 border-2 border-gray-85 rounded-xl rounded-b-none min-h-screen">
                        <Forbidden />
                </div>
            </div>
        )
    }

    if (!pipelineData) {
        return (
            
            <div className="flex flew-row px-12 py-6">
                <AsideMenu>
                    <Link className="text-gray-800 hover:underline pb-1 border-r-2" 
                    href={`/pipeline`}
                    >
                        {'<'} Go back
                    </Link>
                </AsideMenu>
                <div className="grow w-1/2 min-h-screen">
                    <StatusBar />
                    <div className="border-2 border-gray-85 rounded-xl rounded-b-none min-h-screen flex justify-center items-center">
                        <Loader
                            className="w-10 h-10"
                        />
                    </div>
                </div>
            </div>
        )
    }

    return (
        <div className="flex flew-row px-12 py-6">
            <AsideMenu>
                <Link className="text-gray-800 hover:underline pb-1 border-r-2" 
                href={`/pipeline`}
                >
                    {'<'} Go back
                </Link>
                <LinksGroup links={linksGeneral} />
                <LinksGroup links={jobs} title="Jobs" />
            </AsideMenu>
            <div className="grow w-1/2">
                <StatusBar />
                <div className="">
                    {children}
                </div>
            </div>
        </div>
    )
}