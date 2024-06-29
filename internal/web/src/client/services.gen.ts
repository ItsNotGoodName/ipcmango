// This file is auto-generated by @hey-api/openapi-ts

import type { CancelablePromise } from './core/CancelablePromise';
import { OpenAPI } from './core/OpenAPI';
import { request as __request } from './core/request';
import type { GetApiDevicesResponse, PutApiDevicesData, PutApiDevicesResponse, PostApiDevicesCreateData, PostApiDevicesCreateResponse, DeleteApiDevicesByUuidData, DeleteApiDevicesByUuidResponse, GetApiDevicesByUuidData, GetApiDevicesByUuidResponse, PostApiDevicesByUuidData, PostApiDevicesByUuidResponse, GetApiDevicesByUuidCoaxialCapsData, GetApiDevicesByUuidCoaxialCapsResponse, PostApiDevicesByUuidCoaxialSpeakerData, PostApiDevicesByUuidCoaxialSpeakerResponse, GetApiDevicesByUuidCoaxialStatusData, GetApiDevicesByUuidCoaxialStatusResponse, PostApiDevicesByUuidCoaxialWhiteLightData, PostApiDevicesByUuidCoaxialWhiteLightResponse, GetApiDevicesByUuidDetailData, GetApiDevicesByUuidDetailResponse, GetApiDevicesByUuidFileData, GetApiDevicesByUuidFileResponse, GetApiDevicesByUuidGroupsData, GetApiDevicesByUuidGroupsResponse, GetApiDevicesByUuidLicensesData, GetApiDevicesByUuidLicensesResponse, GetApiDevicesByUuidPtzPresetsData, GetApiDevicesByUuidPtzPresetsResponse, PostApiDevicesByUuidRebootData, PostApiDevicesByUuidRebootResponse, DeleteApiDevicesByUuidScanCursorData, DeleteApiDevicesByUuidScanCursorResponse, GetApiDevicesByUuidScanCursorData, GetApiDevicesByUuidScanCursorResponse, PostApiDevicesByUuidScanFullData, PostApiDevicesByUuidScanFullResponse, PostApiDevicesByUuidScanManualData, PostApiDevicesByUuidScanManualResponse, PostApiDevicesByUuidScanQuickData, PostApiDevicesByUuidScanQuickResponse, GetApiDevicesByUuidSnapshotData, GetApiDevicesByUuidSnapshotResponse, GetApiDevicesByUuidSoftwareData, GetApiDevicesByUuidSoftwareResponse, GetApiDevicesByUuidStatusData, GetApiDevicesByUuidStatusResponse, GetApiDevicesByUuidStorageData, GetApiDevicesByUuidStorageResponse, GetApiDevicesByUuidUptimeData, GetApiDevicesByUuidUptimeResponse, GetApiDevicesByUuidUsersData, GetApiDevicesByUuidUsersResponse, GetApiDevicesByUuidUsersActiveData, GetApiDevicesByUuidUsersActiveResponse, GetApiDevicesByUuidVideoInModeData, GetApiDevicesByUuidVideoInModeResponse, PostApiDevicesByUuidVideoInModeSyncData, PostApiDevicesByUuidVideoInModeSyncResponse, GetApiEmailEndpointsResponse, PutApiEmailEndpointsData, PutApiEmailEndpointsResponse, PostApiEmailEndpointsCreateData, PostApiEmailEndpointsCreateResponse, DeleteApiEndpointsByUuidData, DeleteApiEndpointsByUuidResponse, GetApiEndpointsByUuidData, GetApiEndpointsByUuidResponse, PostApiEndpointsByUuidData, PostApiEndpointsByUuidResponse, GetApiEventActionsResponse, GetApiEventCodesResponse, GetApiEventsData, GetApiEventsResponse, GetApiPagesHomeResponse, DeleteApiSettingsResponse, GetApiSettingsResponse, PatchApiSettingsData, PatchApiSettingsResponse, PutApiSettingsData, PutApiSettingsResponse, GetApiStorageDestinationsResponse, PutApiStorageDestinationsData, PutApiStorageDestinationsResponse } from './types.gen';

/**
 * List devices
 * @returns Device OK
 * @throws ApiError
 */
export const getApiDevices = (): CancelablePromise<GetApiDevicesResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/devices',
    errors: {
        default: 'Error'
    }
}); };

/**
 * Put devices
 * @param data The data for the request.
 * @param data.requestBody
 * @returns Device OK
 * @throws ApiError
 */
export const putApiDevices = (data: PutApiDevicesData): CancelablePromise<PutApiDevicesResponse> => { return __request(OpenAPI, {
    method: 'PUT',
    url: '/api/devices',
    body: data.requestBody,
    mediaType: 'application/json',
    errors: {
        default: 'Error'
    }
}); };

/**
 * Create device
 * @param data The data for the request.
 * @param data.requestBody
 * @returns Device OK
 * @throws ApiError
 */
export const postApiDevicesCreate = (data: PostApiDevicesCreateData): CancelablePromise<PostApiDevicesCreateResponse> => { return __request(OpenAPI, {
    method: 'POST',
    url: '/api/devices/create',
    body: data.requestBody,
    mediaType: 'application/json',
    errors: {
        default: 'Error'
    }
}); };

/**
 * Delete device
 * @param data The data for the request.
 * @param data.uuid
 * @returns void No Content
 * @throws ApiError
 */
export const deleteApiDevicesByUuid = (data: DeleteApiDevicesByUuidData): CancelablePromise<DeleteApiDevicesByUuidResponse> => { return __request(OpenAPI, {
    method: 'DELETE',
    url: '/api/devices/{uuid}',
    path: {
        uuid: data.uuid
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * Get device
 * @param data The data for the request.
 * @param data.uuid
 * @returns Device OK
 * @throws ApiError
 */
export const getApiDevicesByUuid = (data: GetApiDevicesByUuidData): CancelablePromise<GetApiDevicesByUuidResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/devices/{uuid}',
    path: {
        uuid: data.uuid
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * Update device
 * @param data The data for the request.
 * @param data.uuid
 * @param data.requestBody
 * @returns Device OK
 * @throws ApiError
 */
export const postApiDevicesByUuid = (data: PostApiDevicesByUuidData): CancelablePromise<PostApiDevicesByUuidResponse> => { return __request(OpenAPI, {
    method: 'POST',
    url: '/api/devices/{uuid}',
    path: {
        uuid: data.uuid
    },
    body: data.requestBody,
    mediaType: 'application/json',
    errors: {
        default: 'Error'
    }
}); };

/**
 * Get device coaxial caps
 * @param data The data for the request.
 * @param data.uuid
 * @param data.channel
 * @returns DeviceCoaxialCaps OK
 * @throws ApiError
 */
export const getApiDevicesByUuidCoaxialCaps = (data: GetApiDevicesByUuidCoaxialCapsData): CancelablePromise<GetApiDevicesByUuidCoaxialCapsResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/devices/{uuid}/coaxial/caps',
    path: {
        uuid: data.uuid
    },
    query: {
        channel: data.channel
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * Set device speaker state
 * @param data The data for the request.
 * @param data.uuid
 * @param data.channel
 * @param data.action
 * @returns void No Content
 * @throws ApiError
 */
export const postApiDevicesByUuidCoaxialSpeaker = (data: PostApiDevicesByUuidCoaxialSpeakerData): CancelablePromise<PostApiDevicesByUuidCoaxialSpeakerResponse> => { return __request(OpenAPI, {
    method: 'POST',
    url: '/api/devices/{uuid}/coaxial/speaker',
    path: {
        uuid: data.uuid
    },
    query: {
        channel: data.channel,
        action: data.action
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * Get device coaxial status
 * @param data The data for the request.
 * @param data.uuid
 * @param data.channel
 * @returns DeviceCoaxialStatus OK
 * @throws ApiError
 */
export const getApiDevicesByUuidCoaxialStatus = (data: GetApiDevicesByUuidCoaxialStatusData): CancelablePromise<GetApiDevicesByUuidCoaxialStatusResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/devices/{uuid}/coaxial/status',
    path: {
        uuid: data.uuid
    },
    query: {
        channel: data.channel
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * Set device white light state
 * @param data The data for the request.
 * @param data.uuid
 * @param data.channel
 * @param data.action
 * @returns void No Content
 * @throws ApiError
 */
export const postApiDevicesByUuidCoaxialWhiteLight = (data: PostApiDevicesByUuidCoaxialWhiteLightData): CancelablePromise<PostApiDevicesByUuidCoaxialWhiteLightResponse> => { return __request(OpenAPI, {
    method: 'POST',
    url: '/api/devices/{uuid}/coaxial/white-light',
    path: {
        uuid: data.uuid
    },
    query: {
        channel: data.channel,
        action: data.action
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * Get device detail
 * @param data The data for the request.
 * @param data.uuid
 * @returns DeviceDetail OK
 * @throws ApiError
 */
export const getApiDevicesByUuidDetail = (data: GetApiDevicesByUuidDetailData): CancelablePromise<GetApiDevicesByUuidDetailResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/devices/{uuid}/detail',
    path: {
        uuid: data.uuid
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * Download device file
 * @param data The data for the request.
 * @param data.uuid
 * @param data.name
 * @returns unknown File from camera
 * @throws ApiError
 */
export const getApiDevicesByUuidFile = (data: GetApiDevicesByUuidFileData): CancelablePromise<GetApiDevicesByUuidFileResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/devices/{uuid}/file',
    path: {
        uuid: data.uuid
    },
    query: {
        name: data.name
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * List device groups
 * @param data The data for the request.
 * @param data.uuid
 * @returns DeviceGroup OK
 * @throws ApiError
 */
export const getApiDevicesByUuidGroups = (data: GetApiDevicesByUuidGroupsData): CancelablePromise<GetApiDevicesByUuidGroupsResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/devices/{uuid}/groups',
    path: {
        uuid: data.uuid
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * Get device licenses
 * @param data The data for the request.
 * @param data.uuid
 * @returns DeviceLicense OK
 * @throws ApiError
 */
export const getApiDevicesByUuidLicenses = (data: GetApiDevicesByUuidLicensesData): CancelablePromise<GetApiDevicesByUuidLicensesResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/devices/{uuid}/licenses',
    path: {
        uuid: data.uuid
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * List device ptz presets
 * @param data The data for the request.
 * @param data.uuid
 * @param data.channel
 * @returns DevicePTZPreset OK
 * @throws ApiError
 */
export const getApiDevicesByUuidPtzPresets = (data: GetApiDevicesByUuidPtzPresetsData): CancelablePromise<GetApiDevicesByUuidPtzPresetsResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/devices/{uuid}/ptz/presets',
    path: {
        uuid: data.uuid
    },
    query: {
        channel: data.channel
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * Reboot device
 * @param data The data for the request.
 * @param data.uuid
 * @returns void No Content
 * @throws ApiError
 */
export const postApiDevicesByUuidReboot = (data: PostApiDevicesByUuidRebootData): CancelablePromise<PostApiDevicesByUuidRebootResponse> => { return __request(OpenAPI, {
    method: 'POST',
    url: '/api/devices/{uuid}/reboot',
    path: {
        uuid: data.uuid
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * Reset device file scan cursor
 * @param data The data for the request.
 * @param data.uuid
 * @returns FileScanCursor OK
 * @throws ApiError
 */
export const deleteApiDevicesByUuidScanCursor = (data: DeleteApiDevicesByUuidScanCursorData): CancelablePromise<DeleteApiDevicesByUuidScanCursorResponse> => { return __request(OpenAPI, {
    method: 'DELETE',
    url: '/api/devices/{uuid}/scan/cursor',
    path: {
        uuid: data.uuid
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * Get file scan cursor
 * @param data The data for the request.
 * @param data.uuid
 * @returns FileScanCursor OK
 * @throws ApiError
 */
export const getApiDevicesByUuidScanCursor = (data: GetApiDevicesByUuidScanCursorData): CancelablePromise<GetApiDevicesByUuidScanCursorResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/devices/{uuid}/scan/cursor',
    path: {
        uuid: data.uuid
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * Full scan device files
 * @param data The data for the request.
 * @param data.uuid
 * @returns ScanResult OK
 * @throws ApiError
 */
export const postApiDevicesByUuidScanFull = (data: PostApiDevicesByUuidScanFullData): CancelablePromise<PostApiDevicesByUuidScanFullResponse> => { return __request(OpenAPI, {
    method: 'POST',
    url: '/api/devices/{uuid}/scan/full',
    path: {
        uuid: data.uuid
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * Manual scan device files
 * @param data The data for the request.
 * @param data.uuid
 * @param data.requestBody
 * @returns ScanResult OK
 * @throws ApiError
 */
export const postApiDevicesByUuidScanManual = (data: PostApiDevicesByUuidScanManualData): CancelablePromise<PostApiDevicesByUuidScanManualResponse> => { return __request(OpenAPI, {
    method: 'POST',
    url: '/api/devices/{uuid}/scan/manual',
    path: {
        uuid: data.uuid
    },
    body: data.requestBody,
    mediaType: 'application/json',
    errors: {
        default: 'Error'
    }
}); };

/**
 * Quick scan device files
 * @param data The data for the request.
 * @param data.uuid
 * @returns ScanResult OK
 * @throws ApiError
 */
export const postApiDevicesByUuidScanQuick = (data: PostApiDevicesByUuidScanQuickData): CancelablePromise<PostApiDevicesByUuidScanQuickResponse> => { return __request(OpenAPI, {
    method: 'POST',
    url: '/api/devices/{uuid}/scan/quick',
    path: {
        uuid: data.uuid
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * Get device snapshot
 * @param data The data for the request.
 * @param data.uuid
 * @param data.channel
 * @param data.type
 * @returns unknown Current snapshot of camera
 * @throws ApiError
 */
export const getApiDevicesByUuidSnapshot = (data: GetApiDevicesByUuidSnapshotData): CancelablePromise<GetApiDevicesByUuidSnapshotResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/devices/{uuid}/snapshot',
    path: {
        uuid: data.uuid
    },
    query: {
        channel: data.channel,
        type: data.type
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * Get device software versions
 * @param data The data for the request.
 * @param data.uuid
 * @returns DeviceSoftwareVersion OK
 * @throws ApiError
 */
export const getApiDevicesByUuidSoftware = (data: GetApiDevicesByUuidSoftwareData): CancelablePromise<GetApiDevicesByUuidSoftwareResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/devices/{uuid}/software',
    path: {
        uuid: data.uuid
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * Get device status
 * @param data The data for the request.
 * @param data.uuid
 * @returns DeviceStatus OK
 * @throws ApiError
 */
export const getApiDevicesByUuidStatus = (data: GetApiDevicesByUuidStatusData): CancelablePromise<GetApiDevicesByUuidStatusResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/devices/{uuid}/status',
    path: {
        uuid: data.uuid
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * List device storage
 * @param data The data for the request.
 * @param data.uuid
 * @returns DeviceStorage OK
 * @throws ApiError
 */
export const getApiDevicesByUuidStorage = (data: GetApiDevicesByUuidStorageData): CancelablePromise<GetApiDevicesByUuidStorageResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/devices/{uuid}/storage',
    path: {
        uuid: data.uuid
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * Get device uptime
 * @param data The data for the request.
 * @param data.uuid
 * @returns DeviceUptime OK
 * @throws ApiError
 */
export const getApiDevicesByUuidUptime = (data: GetApiDevicesByUuidUptimeData): CancelablePromise<GetApiDevicesByUuidUptimeResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/devices/{uuid}/uptime',
    path: {
        uuid: data.uuid
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * List device users
 * @param data The data for the request.
 * @param data.uuid
 * @returns DeviceUser OK
 * @throws ApiError
 */
export const getApiDevicesByUuidUsers = (data: GetApiDevicesByUuidUsersData): CancelablePromise<GetApiDevicesByUuidUsersResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/devices/{uuid}/users',
    path: {
        uuid: data.uuid
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * List device active users
 * @param data The data for the request.
 * @param data.uuid
 * @returns DeviceActiveUser OK
 * @throws ApiError
 */
export const getApiDevicesByUuidUsersActive = (data: GetApiDevicesByUuidUsersActiveData): CancelablePromise<GetApiDevicesByUuidUsersActiveResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/devices/{uuid}/users/active',
    path: {
        uuid: data.uuid
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * Get device VideoInMode
 * @param data The data for the request.
 * @param data.uuid
 * @returns DeviceVideoInMode OK
 * @throws ApiError
 */
export const getApiDevicesByUuidVideoInMode = (data: GetApiDevicesByUuidVideoInModeData): CancelablePromise<GetApiDevicesByUuidVideoInModeResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/devices/{uuid}/video-in-mode',
    path: {
        uuid: data.uuid
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * Sync device VideoInMode
 * @param data The data for the request.
 * @param data.uuid
 * @param data.requestBody
 * @returns DeviceVideoInMode OK
 * @throws ApiError
 */
export const postApiDevicesByUuidVideoInModeSync = (data: PostApiDevicesByUuidVideoInModeSyncData): CancelablePromise<PostApiDevicesByUuidVideoInModeSyncResponse> => { return __request(OpenAPI, {
    method: 'POST',
    url: '/api/devices/{uuid}/video-in-mode/sync',
    path: {
        uuid: data.uuid
    },
    body: data.requestBody,
    mediaType: 'application/json',
    errors: {
        default: 'Error'
    }
}); };

/**
 * List email endpoints
 * @returns EmailEndpoint OK
 * @throws ApiError
 */
export const getApiEmailEndpoints = (): CancelablePromise<GetApiEmailEndpointsResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/email-endpoints',
    errors: {
        default: 'Error'
    }
}); };

/**
 * Put email endpoints
 * @param data The data for the request.
 * @param data.requestBody
 * @returns EmailEndpoint OK
 * @throws ApiError
 */
export const putApiEmailEndpoints = (data: PutApiEmailEndpointsData): CancelablePromise<PutApiEmailEndpointsResponse> => { return __request(OpenAPI, {
    method: 'PUT',
    url: '/api/email-endpoints',
    body: data.requestBody,
    mediaType: 'application/json',
    errors: {
        default: 'Error'
    }
}); };

/**
 * Create email endpoint
 * @param data The data for the request.
 * @param data.requestBody
 * @returns EmailEndpoint OK
 * @throws ApiError
 */
export const postApiEmailEndpointsCreate = (data: PostApiEmailEndpointsCreateData): CancelablePromise<PostApiEmailEndpointsCreateResponse> => { return __request(OpenAPI, {
    method: 'POST',
    url: '/api/email-endpoints/create',
    body: data.requestBody,
    mediaType: 'application/json',
    errors: {
        default: 'Error'
    }
}); };

/**
 * Delete endpoint
 * @param data The data for the request.
 * @param data.uuid
 * @returns void No Content
 * @throws ApiError
 */
export const deleteApiEndpointsByUuid = (data: DeleteApiEndpointsByUuidData): CancelablePromise<DeleteApiEndpointsByUuidResponse> => { return __request(OpenAPI, {
    method: 'DELETE',
    url: '/api/endpoints/{uuid}',
    path: {
        uuid: data.uuid
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * Get email endpoint
 * @param data The data for the request.
 * @param data.uuid
 * @returns EmailEndpoint OK
 * @throws ApiError
 */
export const getApiEndpointsByUuid = (data: GetApiEndpointsByUuidData): CancelablePromise<GetApiEndpointsByUuidResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/endpoints/{uuid}',
    path: {
        uuid: data.uuid
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * Update email endpoint
 * @param data The data for the request.
 * @param data.uuid
 * @param data.requestBody
 * @returns EmailEndpoint OK
 * @throws ApiError
 */
export const postApiEndpointsByUuid = (data: PostApiEndpointsByUuidData): CancelablePromise<PostApiEndpointsByUuidResponse> => { return __request(OpenAPI, {
    method: 'POST',
    url: '/api/endpoints/{uuid}',
    path: {
        uuid: data.uuid
    },
    body: data.requestBody,
    mediaType: 'application/json',
    errors: {
        default: 'Error'
    }
}); };

/**
 * List event actions
 * @returns string OK
 * @throws ApiError
 */
export const getApiEventActions = (): CancelablePromise<GetApiEventActionsResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/event-actions',
    errors: {
        default: 'Error'
    }
}); };

/**
 * List event codes
 * @returns string OK
 * @throws ApiError
 */
export const getApiEventCodes = (): CancelablePromise<GetApiEventCodesResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/event-codes',
    errors: {
        default: 'Error'
    }
}); };

/**
 * Listen for events
 * @param data The data for the request.
 * @param data.deviceUuids
 * @param data.codes
 * @param data.actions
 * @returns unknown OK
 * @throws ApiError
 */
export const getApiEvents = (data: GetApiEventsData = {}): CancelablePromise<GetApiEventsResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/events',
    query: {
        'device-uuids': data.deviceUuids,
        codes: data.codes,
        actions: data.actions
    },
    errors: {
        default: 'Error'
    }
}); };

/**
 * Get home page
 * @returns GetHomePage OK
 * @throws ApiError
 */
export const getApiPagesHome = (): CancelablePromise<GetApiPagesHomeResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/pages/home',
    errors: {
        default: 'Error'
    }
}); };

/**
 * Default settings
 * @returns Settings OK
 * @throws ApiError
 */
export const deleteApiSettings = (): CancelablePromise<DeleteApiSettingsResponse> => { return __request(OpenAPI, {
    method: 'DELETE',
    url: '/api/settings',
    errors: {
        default: 'Error'
    }
}); };

/**
 * Get settings
 * @returns Settings OK
 * @throws ApiError
 */
export const getApiSettings = (): CancelablePromise<GetApiSettingsResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/settings',
    errors: {
        default: 'Error'
    }
}); };

/**
 * Patch settings
 * @param data The data for the request.
 * @param data.requestBody
 * @returns Settings OK
 * @throws ApiError
 */
export const patchApiSettings = (data: PatchApiSettingsData): CancelablePromise<PatchApiSettingsResponse> => { return __request(OpenAPI, {
    method: 'PATCH',
    url: '/api/settings',
    body: data.requestBody,
    mediaType: 'application/json',
    errors: {
        default: 'Error'
    }
}); };

/**
 * Update settings
 * @param data The data for the request.
 * @param data.requestBody
 * @returns Settings OK
 * @throws ApiError
 */
export const putApiSettings = (data: PutApiSettingsData): CancelablePromise<PutApiSettingsResponse> => { return __request(OpenAPI, {
    method: 'PUT',
    url: '/api/settings',
    body: data.requestBody,
    mediaType: 'application/json',
    errors: {
        default: 'Error'
    }
}); };

/**
 * List storage destinations
 * @returns StorageDestination OK
 * @throws ApiError
 */
export const getApiStorageDestinations = (): CancelablePromise<GetApiStorageDestinationsResponse> => { return __request(OpenAPI, {
    method: 'GET',
    url: '/api/storage-destinations',
    errors: {
        default: 'Error'
    }
}); };

/**
 * Put storage destinations
 * @param data The data for the request.
 * @param data.requestBody
 * @returns StorageDestination OK
 * @throws ApiError
 */
export const putApiStorageDestinations = (data: PutApiStorageDestinationsData): CancelablePromise<PutApiStorageDestinationsResponse> => { return __request(OpenAPI, {
    method: 'PUT',
    url: '/api/storage-destinations',
    body: data.requestBody,
    mediaType: 'application/json',
    errors: {
        default: 'Error'
    }
}); };