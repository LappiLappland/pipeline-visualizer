import { useEffect, useRef } from "react";
import Button from "./Button";
import TickIconSelect from "./icons/TickIconSelect";

export default function MultiSelect({children, options, onClick, onClosed, onSelect, isShown}) {

    const mainContiner = useRef();

    useEffect(() => {
        const clickOutsideCallback = (e) => {
            if (mainContiner.current && !mainContiner.current.contains(e.target)) {
                onClosed();
            }
        }

        if (isShown) {
            document.addEventListener('click', clickOutsideCallback);
        }

        return () => {
            if (isShown) {
                document.removeEventListener('click', clickOutsideCallback);
            }
        }
    }, [isShown])

    function handleSelectOption(id, active) {
        onSelect(id, active)
    }

    const itemsEl = options.map((option) => {
        return (
            <li className="group"
                key={option.id}
            >
                <button className="group-first:rounded-t-lg group-last:rounded-b-lg flex items-center w-full text-left hover:bg-gray-85 px-2 py-0.5"
                    onClick={() => handleSelectOption(option.id, !option.active)}
                >
                    
                    {option.active ? (
                        <TickIconSelect 
                            className="h-4 w-4 mr-1 fill-purple-60"
                        />
                    ) : (
                        <div className="h-4 w-4 mr-1" />
                    )}
                    
                    {option.title} 
                </button> 
            </li>
        )
    })

    return (
        <div className="relative flex flex-col">
            <Button className="px-2 py-1"
                onClick={onClick}
            >
                {children}
            </Button>
            <div className={`
                w-max
                top-8 left-0 md:left-auto md:right-0 text-sm text-nowrap absolute shadow-xl mt-1 bg-white border rounded-xl
                transition-transform
                md:origin-top-right origin-top-left
                ${isShown ? '' : 'scale-0'}
            `}
                ref={mainContiner}
            >
                <ul className="">
                    {itemsEl}
                </ul>
            </div>
        </div>
    )
}