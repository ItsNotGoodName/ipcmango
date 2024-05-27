import { provideTheme } from "./ui/theme"
import { createQueryKeyStore } from "@lukemorales/query-key-factory";
import { lazy } from "solid-js"
import { Route, Router } from "@solidjs/router";
import { NotFound } from "./pages/404";

const Devices = lazy(() => import("~/pages/Devices"));
import loadDevices from "~/pages/Devices.data";

export const queries = createQueryKeyStore({
  github: {
    tanstack: null
  }
});

function App() {
  provideTheme()

  return (
    <Router explicitLinks>
      <Route path="/devices" component={Devices} load={loadDevices} />
      <Route path="*404" component={NotFound} />
    </Router >
  )
}

export default App
