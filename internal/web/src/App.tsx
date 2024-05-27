import { provideTheme } from "./ui/theme"

function App() {
  provideTheme()

  return (
    <>
      <h1 class="text-3xl font-bold underline">
        Hello world!
      </h1>
    </>
  )
}

export default App
