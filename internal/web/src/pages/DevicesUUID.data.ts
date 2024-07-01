import { useQueryClient } from "@tanstack/solid-query";
import { api } from "./data";
import { RouteLoadFuncArgs } from "@solidjs/router";

export default function (props: RouteLoadFuncArgs) {
  const client = useQueryClient();
  client.prefetchQuery({ ...api.devices.get(props.params.uuid) });
  client.prefetchQuery({ ...api.devices.detail(props.params.uuid) });
  client.prefetchQuery({ ...api.devices.software(props.params.uuid) });
  client.prefetchQuery({ ...api.devices.licenses(props.params.uuid) });
}
