import { provideTheme } from "./ui/theme";
import { createQueryKeyStore } from "@lukemorales/query-key-factory";
import { lazy } from "solid-js";
import { Route, Router } from "@solidjs/router";
import { NotFound } from "./pages/404";
import { Root } from "./Root";

const Home = lazy(() => import("~/pages/Home"));
import loadHome from "~/pages/Home.data";
const Settings = lazy(() => import("~/pages/Settings"));
import loadSettings from "~/pages/Settings.data";
const Devices = lazy(() => import("~/pages/Devices"));
import loadDevices from "~/pages/Devices.data";
const Emails = lazy(() => import("~/pages/Emails"));
import loadEmails from "~/pages/Emails.data";
const Events = lazy(() => import("~/pages/Events"));
import loadEvents from "~/pages/Events.data";
const EventsLive = lazy(() => import("~/pages/EventsLive"));
import loadEventsLive from "~/pages/EventsLive.data";
const EventsRules = lazy(() => import("~/pages/EventsRules"));
import loadEventsRules from "~/pages/EventsRules.data";
const Files = lazy(() => import("~/pages/Files"));
import loadFiles from "~/pages/Files.data";

export const queries = createQueryKeyStore({
  github: {
    tanstack: null,
  },
});

function App() {
  provideTheme();

  return (
    <Router root={Root} explicitLinks>
      <Route path="/" component={Home} load={loadHome} />
      <Route path="/settings" component={Settings} load={loadSettings} />
      <Route path="/devices" component={Devices} load={loadDevices} />
      <Route path="/emails" component={Emails} load={loadEmails} />
      <Route path="/events" component={Events} load={loadEvents} />
      <Route path="/events/live" component={EventsLive} load={loadEventsLive} />
      <Route
        path="/events/rules"
        component={EventsRules}
        load={loadEventsRules}
      />
      <Route path="/files" component={Files} load={loadFiles} />
      <Route path="*404" component={NotFound} />
    </Router>
  );
}

export default App;
