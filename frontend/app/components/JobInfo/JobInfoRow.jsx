export default function JobInfoRow({children}) {
    return (
        <li className="flex flex-row">
            <dt className="font-bold mr-2">
                {children[0]}
            </dt>
            <dd className="">
                {children[1]}
            </dd>
        </li>
    )
}