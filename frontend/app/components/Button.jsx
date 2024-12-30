export default function Button({children, onClick, className, type = "button"}) {
    return (
        <button
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
            onClick={onClick}
            type={type}
        >
            {children}
        </button>        
    )
}