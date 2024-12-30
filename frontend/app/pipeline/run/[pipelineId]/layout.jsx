'use client';
import JobInfoTag from "../../../components/JobInfo/JobInfoTag";
import { PipelineContext, PipelineProvider } from "../../../components/contexts/PipelineContext";
import { useContext, useMemo } from "react";
import StatusBar from "@/app/components/StatusBar";
import Button from "@/app/components/Button";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import AsideMenu from "@/app/components/AsideMenu";
import PipelineMain from "./main";

export default function ModuleLayout({children}) {

    return (
        <PipelineProvider>
            <PipelineMain>
                {children}
            </PipelineMain>
        </PipelineProvider>
        
    )
}

