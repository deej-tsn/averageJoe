
interface Props {
    choice : String
    index : number
}



export default function ChoiceComp({choice, index }:Props) {
    return (
        <div className="w-80 p-3 bg-red-800 flex place-content-center rounded-2xl drop-shadow-2xl border-2 border-amber-400 cursor-pointer text-3xl font-extrabold text-white">
            {choice}
        </div>
    )
}