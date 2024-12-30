'use client'
import { useEffect, useState } from "react";
import AsideMenu, { LinksGroup } from "../components/AsideMenu";
import getLocation from "../helpers/getLocation";
import Link from "next/link";

export default function PipelineSimple() {

    const [data, setData] = useState([]);

    useEffect(() => {
        const f = async () => {
            const req = await fetch(`${getLocation()}/api/pipelines`)
            const data = await req.json();

            setData(data);
        }

        f()
    }, []);

    const itemsEl = data.map((pipe) => {
        return (
            <li className="border-2 border-black border-x-0 last:borded-t-0"
            key={pipe.id}
            >
                <Link className={`
                    flex flex-col group py-2
                `}
                    href={`${getLocation()}/pipeline/run/${pipe.id}`}
                >
                    <span className="text-2xl group-hover:underline">
                        {pipe.name}
                    </span>
                    <span className="text-gray-600">
                        {pipe.description}
                    </span>
                </Link>
            </li>
        )
    })

    return (
        <div className="flex flew-row px-12 py-6 w-screen">
            <AsideMenu />
            <div className="grow flex flex-col gap-1 border-2 border-gray-85 rounded-xl p-3">
                <h1 className="text-center text-2xl">
                    Demo examples
                </h1>
                <ul>
                    {itemsEl}
                </ul>
            </div>
        </div>
    )
}