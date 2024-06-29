import {
  getApiDevices,
  getApiDevicesByUuidDetail,
  getApiDevicesByUuidLicenses,
  getApiDevicesByUuidSoftware,
  getApiDevicesByUuidStatus,
  getApiDevicesByUuidStorage,
  getApiDevicesByUuidUptime,
  getApiDevicesByUuidVideoInMode,
  getApiEventActions,
  getApiEventCodes,
  getApiLocations,
  getApiPagesHome,
  getApiSettings,
} from "~/client/services.gen";
import {
  createQueryKeyStore,
  createQueryKeys,
} from "@lukemorales/query-key-factory";

export const api = createQueryKeyStore({
  locations: {
    list: {
      queryKey: null,
      queryFn: getApiLocations,
    },
  },
  eventCodes: {
    list: {
      queryKey: null,
      queryFn: getApiEventCodes,
    },
  },
  eventActions: {
    list: {
      queryKey: null,
      queryFn: getApiEventActions,
    },
  },
  settings: {
    get: {
      queryKey: null,
      queryFn: getApiSettings,
    },
  },
  devices: {
    list: {
      queryKey: null,
      queryFn: getApiDevices,
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
    video_in_mode: (uuid: string) => ({
      queryKey: [uuid],
      queryFn: () => getApiDevicesByUuidVideoInMode({ uuid }),
    }),
  },
});

export const pages = createQueryKeys("pages", {
  home: {
    queryKey: null,
    queryFn: getApiPagesHome,
  },
});
