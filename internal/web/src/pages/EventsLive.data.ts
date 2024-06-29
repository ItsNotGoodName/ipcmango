import { useQueryClient } from "@tanstack/solid-query";
import { api } from "./data";

export default function () {
  const client = useQueryClient();
  client.prefetchQuery({ ...api.devices.list });
  client.prefetchQuery({ ...api.eventCodes.list });
  client.prefetchQuery({ ...api.eventActions.list });
}
