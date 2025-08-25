import "./Border.css"

interface Props {
    children : React.ReactNode
}


export default function Border({children} : Props) {
    return (
        <div className="frame">
        <div className="bulbs top"></div>
        <div className="bulbs bottom"></div>
        <div className="bulbs left"></div>
        <div className="bulbs right"></div>
        
        <div className="content">
            {children}
        </div>
</div>
    )
}