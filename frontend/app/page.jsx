import Image from "next/image";
import Visualizer from "./components/visualizer/Visualizer";
import { permanentRedirect } from "next/navigation";

export default async function Home() {
    permanentRedirect('/login');
}
