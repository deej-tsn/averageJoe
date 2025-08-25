import "./Border.css"

interface Props {
    children : React.ReactNode
}


export default function Border({children} : Props) {
    return (
        <div className="frame drop-shadow-2xl">
            {/* Top row */}
            <div className="content">
                {children}
            </div>
        </div>
    )
}