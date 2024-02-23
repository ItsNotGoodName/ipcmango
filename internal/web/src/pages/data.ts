import { makePersisted } from "@solid-primitives/storage"
import { cache } from "@solidjs/router"
import { createStore } from "solid-js/store"
import { useClient } from "~/providers/client"
import { GetConfigResp } from "~/twirp/rpc"

// HACK: this allows App.tsx to switch routes
export const [lastConfig, setLastConfig] = makePersisted(createStore<GetConfigResp>({ siteName: "", enableSignUp: false }), { name: "config" })
export const getConfig = cache(() => useClient().public.getConfig({}).then(res => {
  setLastConfig(res.response)
  return res.response
}), "getConfig")

export const getListDevices = cache(() => useClient().user.listDevices({}).then(res => res.response.devices), "listDevices")
export const getDeviceRPCStatus = cache((id: bigint) => useClient().user.getDeviceRPCStatus({ id }).then(res => res.response), "getDeviceRPCStatus")
export const getDeviceDetail = cache((id: bigint) => useClient().user.getDeviceDetail({ id }).then(res => res.response), "getDeviceDetail")
export const getListDeviceStorage = cache((id: bigint) => useClient().user.listDeviceStorage({ id }).then(res => res.response.items), "listDeviceStorage")
export const getDeviceSoftwareVersion = cache((id: bigint) => useClient().user.getDeviceSoftwareVersion({ id }).then(res => res.response), "getDeviceSoftwareVersion")
export const getListDeviceLicenses = cache((id: bigint) => useClient().user.listDeviceLicenses({ id }).then(res => res.response.items), "listDeviceLicenses")
export const getListEmailAlarmEvents = cache(() => useClient().user.listEmailAlarmEvents({}).then(res => res.response.alarmEvents), "listEmailAlarmEvents")
export const getListEventFilters = cache(() => useClient().user.listEventFilters({}).then(res => res.response), "listEventFilters")
export const getListLatestFiles = cache(() => useClient().user.listLatestFiles({}).then(res => res.response), "listLatestFiles")
