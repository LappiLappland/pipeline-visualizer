export default function ProgressIcon({className, progress = 0}) {

    const deg = (progress / 100) * 360;

    return (
        <span
            className={`
                ${className}
                rounded-full
                flex justify-center items-center
            `}
        >
            <span 
                className={`
                    rounded-full
                    bg-white
                    block
                    w-1/2 h-1/2
                `}
                style={{
                    background: `conic-gradient(#ffff ${deg}deg ,#0000 ${deg + 1}deg 360deg)`
                }}
            />
        </span>
    )
}