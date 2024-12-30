export default function Loader({className}) {
    return (
        <span 
            className={`
                ${className}
                border-4 border-gray-500
                border-b-transparent
                rounded-full
                inline-block
                animate-spin    
            `}
        />
    )
}