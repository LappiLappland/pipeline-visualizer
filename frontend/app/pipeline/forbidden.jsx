import CancelledIcon from "../components/icons/CancelledIcon";

export default function Forbidden() {
    return (
        <div className="flex md:flex-row flex-col gap-2 grow justify-center items-center w-full h-full">
            <CancelledIcon 
                className="w-12 h-12 fill-gray-500"
            />
            <span className="text-3xl">
                Access denied
            </span>
        </div>
    )
}