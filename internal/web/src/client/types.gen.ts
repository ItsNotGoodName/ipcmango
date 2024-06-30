// This file is auto-generated by @hey-api/openapi-ts

export type Build = {
    commit?: string;
    commit_url?: string;
    date?: string;
    license_url?: string;
    release_url?: string;
    repo_url?: string;
    version?: string;
};

export type CreateDevice = {
    /**
     * A URL to the JSON Schema for this object.
     */
    readonly $schema?: string;
    email?: string;
    features?: Array<('camera')>;
    ip?: string;
    latitude?: number;
    location?: string;
    longitude?: number;
    name?: string;
    password?: string;
    sunrise_offset?: string;
    sunset_offset?: string;
    sync_video_in_mode?: boolean;
    username?: string;
    uuid?: string;
};

export type CreateEmailEndpoint = {
    /**
     * A URL to the JSON Schema for this object.
     */
    readonly $schema?: string;
    attachments?: boolean;
    body_template?: string;
    device_uuids?: Array<(string)>;
    disabled?: boolean;
    expression?: string;
    global?: boolean;
    title_template?: string;
    urls: Array<(string)>;
    uuid?: string;
};

export type CreateStorageDestination = {
    name: string;
    password: string;
    port: number;
    remote_directory: string;
    server_address: string;
    storage: 'sftp' | 'ftp';
    username: string;
    uuid?: string;
};

export type storage = 'sftp' | 'ftp';

export type Device = {
    /**
     * A URL to the JSON Schema for this object.
     */
    readonly $schema?: string;
    created_at: string;
    email: string;
    features: Array<(string)>;
    ip: string;
    latitude: number | null;
    location: string;
    longitude: number | null;
    name: string;
    sunrise_offset: string;
    sunset_offset: string;
    sync_video_in_mode: boolean | null;
    updated_at: string;
    username: string;
    uuid: string;
};

export type DeviceActiveUser = {
    client_address: string;
    client_type: string;
    group: string;
    id: number;
    login_time: string;
    name: string;
};

export type DeviceCoaxialCaps = {
    /**
     * A URL to the JSON Schema for this object.
     */
    readonly $schema?: string;
    support_control_fullcolor_light: boolean;
    support_control_light: boolean;
    support_control_speaker: boolean;
};

export type DeviceCoaxialStatus = {
    /**
     * A URL to the JSON Schema for this object.
     */
    readonly $schema?: string;
    speaker: boolean;
    white_light: boolean;
};

export type DeviceDetail = {
    /**
     * A URL to the JSON Schema for this object.
     */
    readonly $schema?: string;
    algorithm_version: string;
    device_class: string;
    device_type: string;
    hardware_version: string;
    market_area: string;
    onvif_version: string;
    process_info: string;
    sn: string;
    vendor: string;
};

export type DeviceEvent = {
    action: string;
    code: string;
    created_at: string;
    data: unknown;
    device_uuid: string;
    id: string;
    index: number;
};

export type DeviceGroup = {
    authority_list: Array<(string)>;
    id: number;
    memo: string;
    name: string;
};

export type DeviceLicense = {
    abroad_info: string;
    all_type: boolean;
    digit_channel: number;
    effective_days: number;
    effective_time: string;
    license_id: number;
    product_type: string;
    status: number;
    username: string;
};

export type DevicePTZPreset = {
    index: number;
    name: string;
};

export type DeviceSoftwareVersion = {
    /**
     * A URL to the JSON Schema for this object.
     */
    readonly $schema?: string;
    build: string;
    build_date: string;
    security_base_line_version: string;
    version: string;
    web_version: string;
};

export type DeviceStatus = {
    /**
     * A URL to the JSON Schema for this object.
     */
    readonly $schema?: string;
    error: string;
    last_login: string;
    state: string;
};

export type DeviceStorage = {
    is_error: boolean;
    name: string;
    path: string;
    state: string;
    total_bytes: number;
    type: string;
    used_bytes: number;
};

export type DeviceUptime = {
    /**
     * A URL to the JSON Schema for this object.
     */
    readonly $schema?: string;
    last: string;
    supported: boolean;
    total: string;
};

export type DeviceUser = {
    anonymous: boolean;
    authority_list: Array<(string)>;
    group: string;
    id: number;
    memo: string;
    name: string;
    password: string;
    password_modified_time: string;
    pwd_score: number;
    reserved: boolean;
    sharable: boolean;
};

export type DeviceVideoInMode = {
    /**
     * A URL to the JSON Schema for this object.
     */
    readonly $schema?: string;
    switch_mode: string;
    time_section: string;
};

export type DeviceVideoInModeSync = {
    /**
     * A URL to the JSON Schema for this object.
     */
    readonly $schema?: string;
    latitude?: number;
    location?: string;
    longitude?: number;
    sunrise_offset?: string;
    sunset_offset?: string;
};

export type EmailEndpoint = {
    /**
     * A URL to the JSON Schema for this object.
     */
    readonly $schema?: string;
    attachments: boolean;
    body_template: string;
    created_at: string;
    device_uuids: Array<(string)>;
    disabled: boolean;
    expression: string;
    global: boolean;
    title_template: string;
    updated_at: string;
    urls: Array<(string)>;
    uuid: string;
};

export type ErrorDetail = {
    /**
     * Where the error occurred, e.g. 'body.items[3].tags' or 'path.thing-id'
     */
    location?: string;
    /**
     * Error message text
     */
    message?: string;
    /**
     * The value at the given location
     */
    value?: unknown;
};

export type ErrorModel = {
    /**
     * A URL to the JSON Schema for this object.
     */
    readonly $schema?: string;
    /**
     * A human-readable explanation specific to this occurrence of the problem.
     */
    detail?: string;
    /**
     * Optional list of individual error details
     */
    errors?: Array<ErrorDetail>;
    /**
     * A URI reference that identifies the specific occurrence of the problem.
     */
    instance?: string;
    /**
     * HTTP status code
     */
    status?: number;
    /**
     * A short, human-readable summary of the problem type. This value should not change between occurrences of the error.
     */
    title?: string;
    /**
     * A URI reference to human-readable documentation for the error.
     */
    type?: string;
};

export type EventRule = {
    code: string;
    deletable: boolean;
    ignore_db: boolean;
    ignore_live: boolean;
    ignore_mqtt: boolean;
};

export type FileScanCursor = {
    /**
     * A URL to the JSON Schema for this object.
     */
    readonly $schema?: string;
    full_complete: boolean;
    full_cursor: string;
    full_epoch: string;
    quick_cursor: string;
    updated_at: string;
};

export type GetHomePage = {
    /**
     * A URL to the JSON Schema for this object.
     */
    readonly $schema?: string;
    build: Build;
    db_usage: number;
    device_count: number;
    email_count: number;
    event_count: number;
    file_count: number;
    file_usage: number;
};

export type ListEvents = {
    /**
     * A URL to the JSON Schema for this object.
     */
    readonly $schema?: string;
    data: Array<DeviceEvent>;
    pagination: PagePagination;
};

export type ManualFileScan = {
    /**
     * A URL to the JSON Schema for this object.
     */
    readonly $schema?: string;
    end_time?: string;
    start_time?: string;
};

export type PagePagination = {
    next_page: number;
    page: number;
    per_page: number;
    previous_page: number;
    seen_items: number;
    total_items: number;
    total_pages: number;
};

export type PatchSettings = {
    /**
     * A URL to the JSON Schema for this object.
     */
    readonly $schema?: string;
    latitude?: number;
    location?: string;
    longitude?: number;
    sunrise_offset?: string;
    sunset_offset?: string;
    sync_video_in_mode?: boolean;
};

export type ScanResult = {
    /**
     * A URL to the JSON Schema for this object.
     */
    readonly $schema?: string;
    created_count: number;
    deleted_count: number;
    updated_count: number;
};

export type Settings = {
    /**
     * A URL to the JSON Schema for this object.
     */
    readonly $schema?: string;
    latitude: number;
    location: string;
    longitude: number;
    sunrise_offset: string;
    sunset_offset: string;
    sync_video_in_mode: boolean;
};

export type StorageDestination = {
    created_at: string;
    name: string;
    password: string;
    port: number;
    remote_directory: string;
    server_address: string;
    storage: string;
    updated_at: string;
    username: string;
    uuid: string;
};

export type UpdateDevice = {
    /**
     * A URL to the JSON Schema for this object.
     */
    readonly $schema?: string;
    email?: string;
    features?: Array<('camera')>;
    ip: string;
    latitude: number | null;
    location?: string;
    longitude: number | null;
    name: string;
    password?: string;
    sunrise_offset?: string;
    sunset_offset?: string;
    sync_video_in_mode?: boolean;
    username: string;
};

export type UpdateEmailEndpoint = {
    /**
     * A URL to the JSON Schema for this object.
     */
    readonly $schema?: string;
    attachments?: boolean;
    body_template?: string;
    device_uuids?: Array<(string)>;
    disabled?: boolean;
    expression?: string;
    global?: boolean;
    title_template?: string;
    urls: Array<(string)>;
};

export type UpdateEventRule = {
    code: string;
    ignore_db: boolean;
    ignore_live: boolean;
    ignore_mqtt: boolean;
    uuid: string;
};

export type UpdateSettings = {
    /**
     * A URL to the JSON Schema for this object.
     */
    readonly $schema?: string;
    latitude: number;
    location: string;
    longitude: number;
    sunrise_offset: string;
    sunset_offset: string;
    sync_video_in_mode: boolean;
};

export type GetApiDevicesResponse = Array<Device>;

export type PutApiDevicesData = {
    requestBody: Array<CreateDevice>;
};

export type PutApiDevicesResponse = Array<Device>;

export type PostApiDevicesCreateData = {
    requestBody: CreateDevice;
};

export type PostApiDevicesCreateResponse = Device;

export type DeleteApiDevicesByUuidData = {
    uuid: string;
};

export type DeleteApiDevicesByUuidResponse = void;

export type GetApiDevicesByUuidData = {
    uuid: string;
};

export type GetApiDevicesByUuidResponse = Device;

export type PostApiDevicesByUuidData = {
    requestBody: UpdateDevice;
    uuid: string;
};

export type PostApiDevicesByUuidResponse = Device;

export type GetApiDevicesByUuidCoaxialCapsData = {
    channel?: number;
    uuid: string;
};

export type GetApiDevicesByUuidCoaxialCapsResponse = DeviceCoaxialCaps;

export type PostApiDevicesByUuidCoaxialSpeakerData = {
    action?: 'on' | 'off' | 'toggle';
    channel?: number;
    uuid: string;
};

export type PostApiDevicesByUuidCoaxialSpeakerResponse = void;

export type GetApiDevicesByUuidCoaxialStatusData = {
    channel?: number;
    uuid: string;
};

export type GetApiDevicesByUuidCoaxialStatusResponse = DeviceCoaxialStatus;

export type PostApiDevicesByUuidCoaxialWhiteLightData = {
    action?: 'on' | 'off' | 'toggle';
    channel?: number;
    uuid: string;
};

export type PostApiDevicesByUuidCoaxialWhiteLightResponse = void;

export type GetApiDevicesByUuidDetailData = {
    uuid: string;
};

export type GetApiDevicesByUuidDetailResponse = DeviceDetail;

export type GetApiDevicesByUuidFileData = {
    name: string;
    uuid: string;
};

export type GetApiDevicesByUuidFileResponse = unknown;

export type GetApiDevicesByUuidGroupsData = {
    uuid: string;
};

export type GetApiDevicesByUuidGroupsResponse = Array<DeviceGroup>;

export type GetApiDevicesByUuidLicensesData = {
    uuid: string;
};

export type GetApiDevicesByUuidLicensesResponse = Array<DeviceLicense>;

export type GetApiDevicesByUuidPtzPresetsData = {
    channel?: number;
    uuid: string;
};

export type GetApiDevicesByUuidPtzPresetsResponse = Array<DevicePTZPreset>;

export type PostApiDevicesByUuidRebootData = {
    uuid: string;
};

export type PostApiDevicesByUuidRebootResponse = void;

export type DeleteApiDevicesByUuidScanCursorData = {
    uuid: string;
};

export type DeleteApiDevicesByUuidScanCursorResponse = FileScanCursor;

export type GetApiDevicesByUuidScanCursorData = {
    uuid: string;
};

export type GetApiDevicesByUuidScanCursorResponse = FileScanCursor;

export type PostApiDevicesByUuidScanFullData = {
    uuid: string;
};

export type PostApiDevicesByUuidScanFullResponse = ScanResult;

export type PostApiDevicesByUuidScanManualData = {
    requestBody: ManualFileScan;
    uuid: string;
};

export type PostApiDevicesByUuidScanManualResponse = ScanResult;

export type PostApiDevicesByUuidScanQuickData = {
    uuid: string;
};

export type PostApiDevicesByUuidScanQuickResponse = ScanResult;

export type GetApiDevicesByUuidSnapshotData = {
    channel?: number;
    type?: number;
    uuid: string;
};

export type GetApiDevicesByUuidSnapshotResponse = unknown;

export type GetApiDevicesByUuidSoftwareData = {
    uuid: string;
};

export type GetApiDevicesByUuidSoftwareResponse = DeviceSoftwareVersion;

export type GetApiDevicesByUuidStatusData = {
    uuid: string;
};

export type GetApiDevicesByUuidStatusResponse = DeviceStatus;

export type GetApiDevicesByUuidStorageData = {
    uuid: string;
};

export type GetApiDevicesByUuidStorageResponse = Array<DeviceStorage>;

export type GetApiDevicesByUuidUptimeData = {
    uuid: string;
};

export type GetApiDevicesByUuidUptimeResponse = DeviceUptime;

export type GetApiDevicesByUuidUsersData = {
    uuid: string;
};

export type GetApiDevicesByUuidUsersResponse = Array<DeviceUser>;

export type GetApiDevicesByUuidUsersActiveData = {
    uuid: string;
};

export type GetApiDevicesByUuidUsersActiveResponse = Array<DeviceActiveUser>;

export type GetApiDevicesByUuidVideoInModeData = {
    uuid: string;
};

export type GetApiDevicesByUuidVideoInModeResponse = DeviceVideoInMode;

export type PostApiDevicesByUuidVideoInModeSyncData = {
    requestBody: DeviceVideoInModeSync;
    uuid: string;
};

export type PostApiDevicesByUuidVideoInModeSyncResponse = DeviceVideoInMode;

export type GetApiEmailEndpointsResponse = Array<EmailEndpoint>;

export type PutApiEmailEndpointsData = {
    requestBody: Array<CreateEmailEndpoint>;
};

export type PutApiEmailEndpointsResponse = Array<EmailEndpoint>;

export type PostApiEmailEndpointsCreateData = {
    requestBody: CreateEmailEndpoint;
};

export type PostApiEmailEndpointsCreateResponse = EmailEndpoint;

export type DeleteApiEndpointsByUuidData = {
    uuid: string;
};

export type DeleteApiEndpointsByUuidResponse = void;

export type GetApiEndpointsByUuidData = {
    uuid: string;
};

export type GetApiEndpointsByUuidResponse = EmailEndpoint;

export type PostApiEndpointsByUuidData = {
    requestBody: UpdateEmailEndpoint;
    uuid: string;
};

export type PostApiEndpointsByUuidResponse = EmailEndpoint;

export type GetApiEventActionsResponse = Array<(string)>;

export type GetApiEventCodesResponse = Array<(string)>;

export type GetApiEventRulesResponse = Array<EventRule>;

export type PostApiEventRulesData = {
    requestBody: Array<UpdateEventRule>;
};

export type PostApiEventRulesResponse = Array<EventRule>;

export type DeleteApiEventsResponse = void;

export type GetApiEventsData = {
    actions?: Array<(string)>;
    codes?: Array<(string)>;
    device?: Array<(string)>;
    order?: 'ascending' | 'descending';
    page?: number;
    perPage?: number;
};

export type GetApiEventsResponse = ListEvents;

export type GetApiEventsSseData = {
    actions?: Array<(string)>;
    codes?: Array<(string)>;
    device?: Array<(string)>;
};

export type GetApiEventsSseResponse = Array<({
    data: DeviceEvent;
    /**
     * The event name.
     */
    event?: "message";
    /**
     * The event ID.
     */
    id?: number;
    /**
     * The retry time in milliseconds.
     */
    retry?: number;
})>;

export type GetApiLocationsResponse = Array<(string)>;

export type GetApiPagesHomeResponse = GetHomePage;

export type DeleteApiSettingsResponse = Settings;

export type GetApiSettingsResponse = Settings;

export type PatchApiSettingsData = {
    requestBody: PatchSettings;
};

export type PatchApiSettingsResponse = Settings;

export type PutApiSettingsData = {
    requestBody: UpdateSettings;
};

export type PutApiSettingsResponse = Settings;

export type GetApiStorageDestinationsResponse = Array<StorageDestination>;

export type PutApiStorageDestinationsData = {
    requestBody: Array<CreateStorageDestination>;
};

export type PutApiStorageDestinationsResponse = Array<StorageDestination>;

export type $OpenApiTs = {
    '/api/devices': {
        get: {
            res: {
                /**
                 * OK
                 */
                200: Array<Device>;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
        put: {
            req: PutApiDevicesData;
            res: {
                /**
                 * OK
                 */
                200: Array<Device>;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/devices/create': {
        post: {
            req: PostApiDevicesCreateData;
            res: {
                /**
                 * OK
                 */
                200: Device;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/devices/{uuid}': {
        delete: {
            req: DeleteApiDevicesByUuidData;
            res: {
                /**
                 * No Content
                 */
                204: void;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
        get: {
            req: GetApiDevicesByUuidData;
            res: {
                /**
                 * OK
                 */
                200: Device;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
        post: {
            req: PostApiDevicesByUuidData;
            res: {
                /**
                 * OK
                 */
                200: Device;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/devices/{uuid}/coaxial/caps': {
        get: {
            req: GetApiDevicesByUuidCoaxialCapsData;
            res: {
                /**
                 * OK
                 */
                200: DeviceCoaxialCaps;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/devices/{uuid}/coaxial/speaker': {
        post: {
            req: PostApiDevicesByUuidCoaxialSpeakerData;
            res: {
                /**
                 * No Content
                 */
                204: void;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/devices/{uuid}/coaxial/status': {
        get: {
            req: GetApiDevicesByUuidCoaxialStatusData;
            res: {
                /**
                 * OK
                 */
                200: DeviceCoaxialStatus;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/devices/{uuid}/coaxial/white-light': {
        post: {
            req: PostApiDevicesByUuidCoaxialWhiteLightData;
            res: {
                /**
                 * No Content
                 */
                204: void;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/devices/{uuid}/detail': {
        get: {
            req: GetApiDevicesByUuidDetailData;
            res: {
                /**
                 * OK
                 */
                200: DeviceDetail;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/devices/{uuid}/file': {
        get: {
            req: GetApiDevicesByUuidFileData;
            res: {
                /**
                 * File from camera
                 */
                200: unknown;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/devices/{uuid}/groups': {
        get: {
            req: GetApiDevicesByUuidGroupsData;
            res: {
                /**
                 * OK
                 */
                200: Array<DeviceGroup>;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/devices/{uuid}/licenses': {
        get: {
            req: GetApiDevicesByUuidLicensesData;
            res: {
                /**
                 * OK
                 */
                200: Array<DeviceLicense>;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/devices/{uuid}/ptz/presets': {
        get: {
            req: GetApiDevicesByUuidPtzPresetsData;
            res: {
                /**
                 * OK
                 */
                200: Array<DevicePTZPreset>;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/devices/{uuid}/reboot': {
        post: {
            req: PostApiDevicesByUuidRebootData;
            res: {
                /**
                 * No Content
                 */
                204: void;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/devices/{uuid}/scan/cursor': {
        delete: {
            req: DeleteApiDevicesByUuidScanCursorData;
            res: {
                /**
                 * OK
                 */
                200: FileScanCursor;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
        get: {
            req: GetApiDevicesByUuidScanCursorData;
            res: {
                /**
                 * OK
                 */
                200: FileScanCursor;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/devices/{uuid}/scan/full': {
        post: {
            req: PostApiDevicesByUuidScanFullData;
            res: {
                /**
                 * OK
                 */
                200: ScanResult;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/devices/{uuid}/scan/manual': {
        post: {
            req: PostApiDevicesByUuidScanManualData;
            res: {
                /**
                 * OK
                 */
                200: ScanResult;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/devices/{uuid}/scan/quick': {
        post: {
            req: PostApiDevicesByUuidScanQuickData;
            res: {
                /**
                 * OK
                 */
                200: ScanResult;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/devices/{uuid}/snapshot': {
        get: {
            req: GetApiDevicesByUuidSnapshotData;
            res: {
                /**
                 * Current snapshot of camera
                 */
                200: unknown;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/devices/{uuid}/software': {
        get: {
            req: GetApiDevicesByUuidSoftwareData;
            res: {
                /**
                 * OK
                 */
                200: DeviceSoftwareVersion;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/devices/{uuid}/status': {
        get: {
            req: GetApiDevicesByUuidStatusData;
            res: {
                /**
                 * OK
                 */
                200: DeviceStatus;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/devices/{uuid}/storage': {
        get: {
            req: GetApiDevicesByUuidStorageData;
            res: {
                /**
                 * OK
                 */
                200: Array<DeviceStorage>;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/devices/{uuid}/uptime': {
        get: {
            req: GetApiDevicesByUuidUptimeData;
            res: {
                /**
                 * OK
                 */
                200: DeviceUptime;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/devices/{uuid}/users': {
        get: {
            req: GetApiDevicesByUuidUsersData;
            res: {
                /**
                 * OK
                 */
                200: Array<DeviceUser>;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/devices/{uuid}/users/active': {
        get: {
            req: GetApiDevicesByUuidUsersActiveData;
            res: {
                /**
                 * OK
                 */
                200: Array<DeviceActiveUser>;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/devices/{uuid}/video-in-mode': {
        get: {
            req: GetApiDevicesByUuidVideoInModeData;
            res: {
                /**
                 * OK
                 */
                200: DeviceVideoInMode;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/devices/{uuid}/video-in-mode/sync': {
        post: {
            req: PostApiDevicesByUuidVideoInModeSyncData;
            res: {
                /**
                 * OK
                 */
                200: DeviceVideoInMode;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/email-endpoints': {
        get: {
            res: {
                /**
                 * OK
                 */
                200: Array<EmailEndpoint>;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
        put: {
            req: PutApiEmailEndpointsData;
            res: {
                /**
                 * OK
                 */
                200: Array<EmailEndpoint>;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/email-endpoints/create': {
        post: {
            req: PostApiEmailEndpointsCreateData;
            res: {
                /**
                 * OK
                 */
                200: EmailEndpoint;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/endpoints/{uuid}': {
        delete: {
            req: DeleteApiEndpointsByUuidData;
            res: {
                /**
                 * No Content
                 */
                204: void;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
        get: {
            req: GetApiEndpointsByUuidData;
            res: {
                /**
                 * OK
                 */
                200: EmailEndpoint;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
        post: {
            req: PostApiEndpointsByUuidData;
            res: {
                /**
                 * OK
                 */
                200: EmailEndpoint;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/event-actions': {
        get: {
            res: {
                /**
                 * OK
                 */
                200: Array<(string)>;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/event-codes': {
        get: {
            res: {
                /**
                 * OK
                 */
                200: Array<(string)>;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/event-rules': {
        get: {
            res: {
                /**
                 * OK
                 */
                200: Array<EventRule>;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
        post: {
            req: PostApiEventRulesData;
            res: {
                /**
                 * OK
                 */
                200: Array<EventRule>;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/events': {
        delete: {
            res: {
                /**
                 * No Content
                 */
                204: void;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
        get: {
            req: GetApiEventsData;
            res: {
                /**
                 * OK
                 */
                200: ListEvents;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/events/sse': {
        get: {
            req: GetApiEventsSseData;
            res: {
                /**
                 * OK
                 */
                200: Array<({
    data: DeviceEvent;
    /**
     * The event name.
     */
    event?: "message";
    /**
     * The event ID.
     */
    id?: number;
    /**
     * The retry time in milliseconds.
     */
    retry?: number;
})>;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/locations': {
        get: {
            res: {
                /**
                 * OK
                 */
                200: Array<(string)>;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/pages/home': {
        get: {
            res: {
                /**
                 * OK
                 */
                200: GetHomePage;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/settings': {
        delete: {
            res: {
                /**
                 * OK
                 */
                200: Settings;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
        get: {
            res: {
                /**
                 * OK
                 */
                200: Settings;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
        patch: {
            req: PatchApiSettingsData;
            res: {
                /**
                 * OK
                 */
                200: Settings;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
        put: {
            req: PutApiSettingsData;
            res: {
                /**
                 * OK
                 */
                200: Settings;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
    '/api/storage-destinations': {
        get: {
            res: {
                /**
                 * OK
                 */
                200: Array<StorageDestination>;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
        put: {
            req: PutApiStorageDestinationsData;
            res: {
                /**
                 * OK
                 */
                200: Array<StorageDestination>;
                /**
                 * Error
                 */
                default: ErrorModel;
            };
        };
    };
};