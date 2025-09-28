import { useEffect, useRef, useState } from 'react'
import './App.css'
import Border from './components/Border'
import ChoiceComp from './components/Choice'

interface Round {
  roundID : string,
  roundData : {
    question : string
    options : string[]
  }
}

interface wsType {
  messageType : String,
  data : any
}


function App() {
  const [round, setRound] = useState<Round | undefined>(undefined)
  const [userID, setUserID] = useState<string|undefined>(undefined)
  const [games, setGames] = useState<Record<string, string>|undefined>(undefined)
  const webSocketRef = useRef<WebSocket|null>(null)

  async function getUser(){
      fetch('http://localhost:8080/user',
        {
          method : 'POST',
          body : JSON.stringify({username : 'bob'}),
          headers : {
            "Content-Type" : 'application/json'
          }
        }
      ).then((response) => response.json()).then((data : {token : string}) => setUserID(data.token))
    }

  async function getGames(){
      fetch('http://localhost:8080/games/active-games', {
        headers: new Headers({
        'Authorization': `Bearer ${userID}`, 
        'Content-Type': 'application/x-www-form-urlencoded'
      }), 
        
      }).then((response) => response.json()).then((data : Record<string, string>) => setGames(data))
    }

  function joinGame(){
    webSocketRef.current = new WebSocket(`ws://localhost:8080/games/connect-to-game?gameID=${game}`,
      `AuthToken-${userID}`
    )
    webSocketRef.current.onmessage = (ev) => {
      let jData = JSON.parse(ev.data) as wsType
      switch (jData.messageType) {
        case "ROUND":
          setRound(jData.data as Round)
          break;
        default:
          console.log(
            `unknown data : ${jData}`
          )
          break;
      }
    } 
  }

  useEffect(()=> {
    getUser()
    return () => {
      webSocketRef.current?.close()
    }
  }, [])

  useEffect(() => {
    if (userID) {
      getGames()
    }
  }, [userID])

  if (!games || !userID) {
    return null
  }

  const game = Object.keys(games)[0]

  if (game && !webSocketRef.current) {
    joinGame()
  }

  return (
    <div className='game w-4xl h-lvh flex flex-col  justify-evenly items-center bg-gradient-to-t from-blue-400 to-blue-950 p-1'>
      <Border>
        <div id='question'>
          {round?.roundData.question}
        </div>
      </Border>
      <Border>
        {round?.roundData.options.map((choiceStr, index) => <ChoiceComp key={index} choice={choiceStr} index={index} websocketRef={webSocketRef}/> )}
      </Border>

    </div>
  )
}

export default App
