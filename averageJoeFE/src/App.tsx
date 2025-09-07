import React, { useEffect, useRef, useState } from 'react'
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
  const [userID, setUserID] = useState<string|undefined>(undefined)
  const [games, setGames] = useState<Record<string, string>|undefined>(undefined)
  const webSocketRef = useRef<WebSocket|null>(null)


  useEffect(()=> {
    async function getUser(){
      fetch('http://localhost:8080/user').then((response) => response.json()).then((data : {uuid : string}) => setUserID(data.uuid))
    }
    async function getGames(){
      fetch('http://localhost:8080/active-games').then((response) => response.json()).then((data : Record<string, string>) => setGames(data))
    }

    getUser()
    getGames()
    
    return () => {
      webSocketRef.current?.close()
    }
  }, [])

  if (!games || !userID) {
    return null
  }

  const game = Object.keys(games)[0]

  function joinGame(event :React.MouseEvent){
    event.preventDefault()
    webSocketRef.current = new WebSocket(`ws://localhost:8080/connect-to-game?gameID=${game}`)
    webSocketRef.current.onmessage = (ev) => {
      console.log(ev)
    } 
  }

  function sendMessage(message : string) {
    webSocketRef.current?.send(message)
  }


  return (
    <div className='game w-4xl h-lvh flex flex-col  justify-evenly items-center bg-gradient-to-t from-blue-400 to-blue-950 p-1'>
      <Border>
        <div id='question'>
          {game}
        </div>
      </Border>
      <Border>
        <button onClick={joinGame}>Join</button>
        <button onClick={() => sendMessage("hi")}>SEND HI</button>
      </Border>

    </div>
  )
}

export default App
