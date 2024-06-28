/* @refresh reload */
import { render } from "solid-js/web";

import "./styles/index.css";
import "./styles/hljs.css";
import App from "./App";
import { QueryClient, QueryClientProvider } from "@tanstack/solid-query";
import { SolidQueryDevtools } from "@tanstack/solid-query-devtools";

const root = document.getElementById("root");
const client = new QueryClient();

render(
  () => (
    <QueryClientProvider client={client}>
      <SolidQueryDevtools initialIsOpen={false} />
      <App />
    </QueryClientProvider>
  ),
  root!,
);
