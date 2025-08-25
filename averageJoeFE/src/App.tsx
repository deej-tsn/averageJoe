import './App.css'
import Border from './components/Border'
import ChoiceComp from './components/Choice'

function App() {
  return (
    <div className='game w-full p-1'>
      <Border>
        <div id='question'>
          Up or Down?
        </div>
      </Border>

      <ChoiceComp
        choice={'Up'}
        index={1}
      />

      <ChoiceComp
        choice={'Down'}
        index={2}
      />
    </div>
  )
}

export default App
