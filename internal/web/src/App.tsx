import { provideTheme } from "./ui/theme"
import { createQueryKeyStore } from "@lukemorales/query-key-factory";
import { lazy } from "solid-js"
import { Route, Router } from "@solidjs/router";
import { NotFound } from "./pages/404";
import { Root } from "./Root";

const Home = lazy(() => import("~/pages/Home"));
import loadHome from "~/pages/Home.data";
const Devices = lazy(() => import("~/pages/Devices"));
import loadDevices from "~/pages/Devices.data";
const Events = lazy(() => import("~/pages/Events"));
import loadEvents from "~/pages/Events.data";
const EventsLive = lazy(() => import("~/pages/EventsLive"));
import loadEventsLive from "~/pages/EventsLive.data";
const Files = lazy(() => import("~/pages/Files"));
import loadFiles from "~/pages/Files.data";

export const queries = createQueryKeyStore({
  github: {
    tanstack: null
  }
});

function App() {
  provideTheme()

  return (
    <Router root={Root} explicitLinks>
      <Route path="/" component={Home} load={loadHome} />
      <Route path="/devices" component={Devices} load={loadDevices} />
      <Route path="/events" component={Events} load={loadEvents} />
      <Route path="/events/live" component={EventsLive} load={loadEventsLive} />
      <Route path="/files" component={Files} load={loadFiles} />
      <Route path="*404" component={NotFound} />
    </Router >
  )
}

export default App
