import type { RefObject } from "react"

interface Props {
    choice : String
    websocketRef : RefObject<WebSocket|null>
}



export default function ChoiceComp({choice, websocketRef }:Props) {
    
    function handleOnClick(event : React.MouseEvent) {
        event.preventDefault()
        websocketRef.current?.send(JSON.stringify({
            messageType : "VOTE",
            data: {
                choice : choice
            }
        }))
    }


    return (
        <div onClick={handleOnClick} className="w-80 p-3 bg-red-800 flex place-content-center rounded-2xl drop-shadow-2xl border-2 border-amber-400 cursor-pointer text-3xl font-extrabold text-white">
            {choice}
        </div>
    )
}