import { useEffect, useState } from 'react'
import './App.css'
import Border from './components/Border'
import ChoiceComp from './components/Choice'
import Loading from './components/Loading'

interface Round {
  question : string
  options : string[]
}


function App() {
  const [round, setRound] = useState<Round | undefined>(undefined)


  useEffect(()=> {
     fetch('http://localhost:8080').then((response) => response.json()).then((data : Round) => setRound(data))
  }, [])


  return (
    <div className='game w-4xl h-lvh flex flex-col  justify-evenly items-center bg-gradient-to-t from-blue-400 to-blue-950 p-1'>
      <Border>
        {!!round ? 
        <div id='question'>
          {round.question}
        </div> : <Loading/>}
      </Border>
      <Border>
        <div className='w-full h-60 flex flex-col justify-evenly items-center'>
          {!!round? 
          <>
            {round.options.map((option, index) => <ChoiceComp choice={option} index={index}/>)}
          </>

          : <Loading/>}
        </div>
      </Border>

    </div>
  )
}

export default App
