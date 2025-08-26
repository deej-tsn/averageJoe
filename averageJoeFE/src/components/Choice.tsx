
interface Props {
    choice : String
    index : number
}



export default function ChoiceComp({choice, index }:Props) {
    
    function handleOnClick(event : React.MouseEvent) {
        event.preventDefault()
        const form = new FormData()
        form.append('choice', String(index))
        fetch('http://localhost:8080', {
            method : 'POST',
            body : form
        }).then((response) => console.log(response))
    }


    return (
        <div onClick={handleOnClick} className="w-80 p-3 bg-red-800 flex place-content-center rounded-2xl drop-shadow-2xl border-2 border-amber-400 cursor-pointer text-3xl font-extrabold text-white">
            {choice}
        </div>
    )
}