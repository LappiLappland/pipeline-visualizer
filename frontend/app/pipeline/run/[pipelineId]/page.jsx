import Image from "next/image";
import Visualizer from "../../../components/visualizer/Visualizer";
import PipelineControls from "../../../components/PipelineControls/PipelineControls";

export default async function Pipeline() {

    return (
        <div className="">
            <div className="border-2 border-gray-85 rounded-xl rounded-b-none">
                <PipelineControls />
            </div>
            <div className="border-2 border-t-0 border-gray-85 rounded-xl rounded-t-none">
                <Visualizer 
                />
            </div>
        </div>
  );
}
