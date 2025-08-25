import './App.css'
import Border from './components/Border'
import ChoiceComp from './components/Choice'

function App() {
  return (
    <div className='game w-4xl h-lvh flex flex-col  justify-evenly items-center bg-gradient-to-t from-blue-400 to-blue-950 p-1'>
      <Border>
        <div id='question'>
          Up or Down?
        </div>
      </Border>
      <Border>
        <div className='w-full h-60 flex flex-col justify-evenly items-center'>
          <ChoiceComp choice={'Up'} index={1}/>

          <ChoiceComp choice={'Down'} index={2}/>
        </div>
      </Border>

    </div>
  )
}

export default App
