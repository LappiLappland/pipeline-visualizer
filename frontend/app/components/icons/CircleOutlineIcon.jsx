export default function CircleOutlineIcon({className}) {
    return (
        <span
            className={`
                ${className}
                rounded-full
                flex justify-center items-center
            `}
        >
            <svg
            className="w-2/3 h-2/3"
            viewBox="0 0 20 20" 
            >
                <path d="M10,0 C4.5,0 0,4.5 0,10 C0,15.5 4.5,20 10,20 C15.5,20 20,15.5 20,10 C20,4.5 15.5,0 10,0 L10,0 Z M10,18 C5.6,18 2,14.4 2,10 C2,5.6 5.6,2 10,2 C14.4,2 18,5.6 18,10 C18,14.4 14.4,18 10,18 L10,18 Z" />
            </svg>
        </span>
        
    )
}