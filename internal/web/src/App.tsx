import { provideTheme } from "./ui/theme"
import { createQueryKeyStore } from "@lukemorales/query-key-factory";
import { lazy } from "solid-js"
import { Route, Router } from "@solidjs/router";
import { NotFound } from "./pages/404";
import { Root } from "./Root";

const Devices = lazy(() => import("~/pages/Devices"));
import loadDevices from "~/pages/Devices.data";
const EventsLive = lazy(() => import("~/pages/EventsLive"));
import loadEventsLive from "~/pages/EventsLive.data";

export const queries = createQueryKeyStore({
  github: {
    tanstack: null
  }
});

function App() {
  provideTheme()

  return (
    <Router root={Root} explicitLinks>
      <Route path="/devices" component={Devices} load={loadDevices} />
      <Route path="/events/live" component={EventsLive} load={loadEventsLive} />
      <Route path="*404" component={NotFound} />
    </Router >
  )
}

export default App
