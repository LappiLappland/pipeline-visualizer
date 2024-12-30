import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import Button from "./Button";
import getLocation from "../helpers/getLocation";

export default function LogOutButton({className}) {

    return (
        <a
            className={`
                text-gray-25
                bg-gray-93
                hover:bg-gray-95
                px-4
                rounded-lg
                border
                border-gray-85
                hover:border-gray-90
                active:border-gray-50
                active:bg-gray-85
                ${className}
            `}
            href={`${getLocation()}/api/logout`}
        >
            Sign out
        </a>
    )
}