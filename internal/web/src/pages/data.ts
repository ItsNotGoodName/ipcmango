import {
  getApiDevices,
  getApiDevicesByUuid,
  getApiDevicesByUuidDetail,
  getApiDevicesByUuidLicenses,
  getApiDevicesByUuidSoftware,
  getApiDevicesByUuidStatus,
  getApiDevicesByUuidStorage,
  getApiDevicesByUuidUptime,
  getApiDevicesByUuidVideoInMode,
  getApiEventActions,
  getApiEventCodes,
  getApiEventRules,
  getApiEvents,
  getApiFiles,
  getApiLocations,
  getApiPagesHome,
  getApiSettings,
} from "~/client/services.gen";
import {
  createQueryKeyStore,
  createQueryKeys,
} from "@lukemorales/query-key-factory";
import { GetApiEventsData, GetApiFilesData } from "~/client";

export const api = createQueryKeyStore({
  locations: {
    list: {
      queryKey: null,
      queryFn: getApiLocations,
    },
  },
  eventRules: {
    list: {
      queryKey: null,
      queryFn: getApiEventRules,
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
  events: {
    list: (data?: GetApiEventsData) => ({
      queryKey: [data],
      queryFn: () => getApiEvents(data),
    }),
  },
  devices: {
    list: {
      queryKey: null,
      queryFn: getApiDevices,
    },
    get: (uuid: string) => ({
      queryKey: [uuid],
      queryFn: () => getApiDevicesByUuid({ uuid }),
    }),
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
  files: {
    list: (data: GetApiFilesData) => ({
      queryKey: [data],
      queryFn: () => getApiFiles(data),
    }),
  },
});

export const pages = createQueryKeys("pages", {
  home: {
    queryKey: null,
    queryFn: getApiPagesHome,
  },
});
