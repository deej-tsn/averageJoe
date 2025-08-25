
interface Props {
    choice : String
    index : number
}



export default function ChoiceComp({choice, index }:Props) {
    return (
        <div className="choice ">
            {choice}
        </div>
    )
}