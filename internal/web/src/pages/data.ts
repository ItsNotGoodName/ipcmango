import { getApiDevices, getApiDevicesByUuidDetail, getApiDevicesByUuidLicenses, getApiDevicesByUuidSoftware, getApiDevicesByUuidStatus, getApiDevicesByUuidStorage, getApiDevicesByUuidUptime } from '~/client/services.gen';
import { createQueryKeyStore } from "@lukemorales/query-key-factory";

export const q = createQueryKeyStore({
  devices: {
    list: {
      queryKey: null,
      queryFn: getApiDevices,
      throwOnError: true,
    },
    status: (uuid: string) => ({
      queryKey: [uuid],
      queryFn: () => getApiDevicesByUuidStatus({ uuid }),
    }),
    uptime: (uuid: string) => ({
      queryKey: [uuid],
      queryFn: () => getApiDevicesByUuidUptime({ uuid }),
    }),
    detail: (uuid: string) => ({
      queryKey: [uuid],
      queryFn: () => getApiDevicesByUuidDetail({ uuid }),
    }),
    software: (uuid: string) => ({
      queryKey: [uuid],
      queryFn: () => getApiDevicesByUuidSoftware({ uuid }),
    }),
    licenses: (uuid: string) => ({
      queryKey: [uuid],
      queryFn: () => getApiDevicesByUuidLicenses({ uuid }),
    }),
    storage: (uuid: string) => ({
      queryKey: [uuid],
      queryFn: () => getApiDevicesByUuidStorage({ uuid }),
    }),
  }
})
