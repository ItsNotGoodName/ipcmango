import { useQueryClient } from "@tanstack/solid-query";
import { api } from "./data";
import { RouteLoadFuncArgs } from "~/lib/utils";

export default function (props: RouteLoadFuncArgs<{ uuid: string }>) {
  const client = useQueryClient();
  client.prefetchQuery({ ...api.devices.get(props.params.uuid) });
  client.prefetchQuery({ ...api.devices.detail(props.params.uuid) });
  client.prefetchQuery({ ...api.devices.software(props.params.uuid) });
  client.prefetchQuery({ ...api.devices.licenses(props.params.uuid) });
}
