import Link from "next/link";
import LogOutButton from "./LogOutButton";
import { useEffect, useState } from "react";
import getLocation from "../helpers/getLocation";

export default function AsideMenu({children}) {

    const [userName, setUserName] = useState('');

    useEffect(() => {
        const f = async () => {
            const req = await fetch(`${getLocation()}/api/user`);

            if (!req.ok) {
                return;
            }
    
            const userData = await req.json();

            setUserName(userData.nickname);
        }

        f();
    }, [])

    return (
        <aside className="flex-col min-h-screen pr-2.5 hidden md:flex">
            <nav className="flex flex-col">
                <div className="flex flex-col py-1.5 pr-1 border-r-2">
                    <span className="block text-center h-6">
                        {userName}
                    </span>
                    <LogOutButton 
                        className="w-full text-sm"
                    />
                </div>
                
                {children}
            </nav>
            <div className="grow border-r-2" 
            />
        </aside>
    )
}

export function LinksGroup({links, title}) {

    const linksEl = links.map((link, i) => {
        return (
            <li className=""
                key={link.href}
            >
                <Link   className={`flex items-center text-nowrap grow px-2.5 py-0.5 border-r-2 ${link.isCurrent ? 'border-r-2 border-r-blue-65 bg-gray-95' : 'hover:border-r-2 hover:border-r-black/50'}`}
                    href={link.href}
                >
                    {!link.icon ? '' : link.icon}
                    {link.title}
                </Link>
            </li>
        )
    })

    return (
        <div className="">
            {!title ? '' : (
                <div className="text-gray-500 text-sm border-r-2 pt-2">
                    {title}
                </div>
            )}
            
            <ul>
                {linksEl}
            </ul>
        </div>
        
    )
}