// This file is auto-generated by @hey-api/openapi-ts

export const $Build = {
    additionalProperties: false,
    properties: {
        commit: {
            type: 'string'
        },
        commit_url: {
            type: 'string'
        },
        date: {
            format: 'date-time',
            type: 'string'
        },
        license_url: {
            type: 'string'
        },
        release_url: {
            type: 'string'
        },
        repo_url: {
            type: 'string'
        },
        version: {
            type: 'string'
        }
    },
    type: 'object'
} as const;

export const $CreateDevice = {
    additionalProperties: false,
    properties: {
        '$schema': {
            description: 'A URL to the JSON Schema for this object.',
            examples: ['https://example.com/schemas/CreateDevice.json'],
            format: 'uri',
            readOnly: true,
            type: 'string'
        },
        email: {
            type: 'string'
        },
        features: {
            items: {
                type: 'string'
            },
            type: 'array'
        },
        ip: {
            format: 'ipv4',
            type: 'string'
        },
        latitude: {
            format: 'double',
            type: 'number'
        },
        location: {
            examples: ['UTC'],
            type: 'string'
        },
        longitude: {
            format: 'double',
            type: 'number'
        },
        name: {
            type: 'string'
        },
        password: {
            type: 'string'
        },
        sunrise_offset: {
            examples: ['0s'],
            type: 'string'
        },
        sunset_offset: {
            examples: ['0s'],
            type: 'string'
        },
        sync_video_in_mode: {
            type: 'boolean'
        },
        username: {
            type: 'string'
        }
    },
    type: 'object'
} as const;

export const $CreateEndpointInput = {
    additionalProperties: false,
    properties: {
        '$schema': {
            description: 'A URL to the JSON Schema for this object.',
            examples: ['https://example.com/schemas/CreateEndpointInput.json'],
            format: 'uri',
            readOnly: true,
            type: 'string'
        },
        gorise_url: {
            type: 'string'
        }
    },
    required: ['gorise_url'],
    type: 'object'
} as const;

export const $Device = {
    additionalProperties: false,
    properties: {
        '$schema': {
            description: 'A URL to the JSON Schema for this object.',
            examples: ['https://example.com/schemas/Device.json'],
            format: 'uri',
            readOnly: true,
            type: 'string'
        },
        created_at: {
            format: 'date-time',
            type: 'string'
        },
        email: {
            type: 'string'
        },
        features: {
            items: {
                type: 'string'
            },
            type: 'array'
        },
        ip: {
            format: 'ipv4',
            type: 'string'
        },
        latitude: {
            format: 'double',
            type: ['number', 'null']
        },
        location: {
            examples: ['UTC'],
            type: 'string'
        },
        longitude: {
            format: 'double',
            type: ['number', 'null']
        },
        name: {
            type: 'string'
        },
        sunrise_offset: {
            examples: ['0s'],
            type: 'string'
        },
        sunset_offset: {
            examples: ['0s'],
            type: 'string'
        },
        sync_video_in_mode: {
            type: ['boolean', 'null']
        },
        updated_at: {
            format: 'date-time',
            type: 'string'
        },
        username: {
            type: 'string'
        },
        uuid: {
            type: 'string'
        }
    },
    required: ['uuid', 'name', 'ip', 'username', 'location', 'features', 'email', 'created_at', 'updated_at', 'latitude', 'longitude', 'sunrise_offset', 'sunset_offset', 'sync_video_in_mode'],
    type: 'object'
} as const;

export const $DeviceActiveUser = {
    additionalProperties: false,
    properties: {
        client_address: {
            type: 'string'
        },
        client_type: {
            type: 'string'
        },
        group: {
            type: 'string'
        },
        id: {
            format: 'int64',
            type: 'integer'
        },
        login_time: {
            format: 'date-time',
            type: 'string'
        },
        name: {
            type: 'string'
        }
    },
    required: ['client_address', 'client_type', 'group', 'id', 'login_time', 'name'],
    type: 'object'
} as const;

export const $DeviceCoaxialCaps = {
    additionalProperties: false,
    properties: {
        '$schema': {
            description: 'A URL to the JSON Schema for this object.',
            examples: ['https://example.com/schemas/DeviceCoaxialCaps.json'],
            format: 'uri',
            readOnly: true,
            type: 'string'
        },
        support_control_fullcolor_light: {
            type: 'boolean'
        },
        support_control_light: {
            type: 'boolean'
        },
        support_control_speaker: {
            type: 'boolean'
        }
    },
    required: ['support_control_fullcolor_light', 'support_control_light', 'support_control_speaker'],
    type: 'object'
} as const;

export const $DeviceCoaxialStatus = {
    additionalProperties: false,
    properties: {
        '$schema': {
            description: 'A URL to the JSON Schema for this object.',
            examples: ['https://example.com/schemas/DeviceCoaxialStatus.json'],
            format: 'uri',
            readOnly: true,
            type: 'string'
        },
        speaker: {
            type: 'boolean'
        },
        white_light: {
            type: 'boolean'
        }
    },
    required: ['white_light', 'speaker'],
    type: 'object'
} as const;

export const $DeviceDetail = {
    additionalProperties: false,
    properties: {
        '$schema': {
            description: 'A URL to the JSON Schema for this object.',
            examples: ['https://example.com/schemas/DeviceDetail.json'],
            format: 'uri',
            readOnly: true,
            type: 'string'
        },
        algorithm_version: {
            type: 'string'
        },
        device_class: {
            type: 'string'
        },
        device_type: {
            type: 'string'
        },
        hardware_version: {
            type: 'string'
        },
        market_area: {
            type: 'string'
        },
        onvif_version: {
            type: 'string'
        },
        process_info: {
            type: 'string'
        },
        sn: {
            type: 'string'
        },
        vendor: {
            type: 'string'
        }
    },
    required: ['sn', 'device_class', 'device_type', 'hardware_version', 'market_area', 'process_info', 'vendor', 'onvif_version', 'algorithm_version'],
    type: 'object'
} as const;

export const $DeviceEventsOutput = {
    additionalProperties: false,
    properties: {
        action: {
            type: 'string'
        },
        code: {
            type: 'string'
        },
        created_at: {
            format: 'date-time',
            type: 'string'
        },
        data: {},
        device_uuid: {
            type: 'string'
        },
        id: {
            type: 'string'
        },
        index: {
            format: 'int64',
            type: 'integer'
        }
    },
    required: ['id', 'device_uuid', 'code', 'action', 'index', 'data', 'created_at'],
    type: 'object'
} as const;

export const $DeviceGroup = {
    additionalProperties: false,
    properties: {
        authority_list: {
            items: {
                type: 'string'
            },
            type: 'array'
        },
        id: {
            format: 'int64',
            type: 'integer'
        },
        memo: {
            type: 'string'
        },
        name: {
            type: 'string'
        }
    },
    required: ['authority_list', 'id', 'memo', 'name'],
    type: 'object'
} as const;

export const $DeviceLicense = {
    additionalProperties: false,
    properties: {
        abroad_info: {
            type: 'string'
        },
        all_type: {
            type: 'boolean'
        },
        digit_channel: {
            format: 'int64',
            type: 'integer'
        },
        effective_days: {
            format: 'int64',
            type: 'integer'
        },
        effective_time: {
            format: 'date-time',
            type: 'string'
        },
        license_id: {
            format: 'int64',
            type: 'integer'
        },
        product_type: {
            type: 'string'
        },
        status: {
            format: 'int64',
            type: 'integer'
        },
        username: {
            type: 'string'
        }
    },
    required: ['abroad_info', 'all_type', 'digit_channel', 'effective_days', 'effective_time', 'license_id', 'product_type', 'status', 'username'],
    type: 'object'
} as const;

export const $DevicePTZPreset = {
    additionalProperties: false,
    properties: {
        index: {
            format: 'int64',
            type: 'integer'
        },
        name: {
            type: 'string'
        }
    },
    required: ['index', 'name'],
    type: 'object'
} as const;

export const $DeviceSoftwareVersion = {
    additionalProperties: false,
    properties: {
        '$schema': {
            description: 'A URL to the JSON Schema for this object.',
            examples: ['https://example.com/schemas/DeviceSoftwareVersion.json'],
            format: 'uri',
            readOnly: true,
            type: 'string'
        },
        build: {
            type: 'string'
        },
        build_date: {
            type: 'string'
        },
        security_base_line_version: {
            type: 'string'
        },
        version: {
            type: 'string'
        },
        web_version: {
            type: 'string'
        }
    },
    required: ['build', 'build_date', 'security_base_line_version', 'version', 'web_version'],
    type: 'object'
} as const;

export const $DeviceStatus = {
    additionalProperties: false,
    properties: {
        '$schema': {
            description: 'A URL to the JSON Schema for this object.',
            examples: ['https://example.com/schemas/DeviceStatus.json'],
            format: 'uri',
            readOnly: true,
            type: 'string'
        },
        error: {
            type: 'string'
        },
        last_login: {
            format: 'date-time',
            type: 'string'
        },
        state: {
            type: 'string'
        }
    },
    required: ['error', 'state', 'last_login'],
    type: 'object'
} as const;

export const $DeviceStorage = {
    additionalProperties: false,
    properties: {
        is_error: {
            type: 'boolean'
        },
        name: {
            type: 'string'
        },
        path: {
            type: 'string'
        },
        state: {
            type: 'string'
        },
        total_bytes: {
            format: 'int64',
            type: 'integer'
        },
        type: {
            type: 'string'
        },
        used_bytes: {
            format: 'int64',
            type: 'integer'
        }
    },
    required: ['name', 'state', 'path', 'type', 'total_bytes', 'used_bytes', 'is_error'],
    type: 'object'
} as const;

export const $DeviceUptime = {
    additionalProperties: false,
    properties: {
        '$schema': {
            description: 'A URL to the JSON Schema for this object.',
            examples: ['https://example.com/schemas/DeviceUptime.json'],
            format: 'uri',
            readOnly: true,
            type: 'string'
        },
        last: {
            format: 'date-time',
            type: 'string'
        },
        supported: {
            type: 'boolean'
        },
        total: {
            format: 'date-time',
            type: 'string'
        }
    },
    required: ['last', 'total', 'supported'],
    type: 'object'
} as const;

export const $DeviceUser = {
    additionalProperties: false,
    properties: {
        anonymous: {
            type: 'boolean'
        },
        authority_list: {
            items: {
                type: 'string'
            },
            type: 'array'
        },
        group: {
            type: 'string'
        },
        id: {
            format: 'int64',
            type: 'integer'
        },
        memo: {
            type: 'string'
        },
        name: {
            type: 'string'
        },
        password: {
            type: 'string'
        },
        password_modified_time: {
            format: 'date-time',
            type: 'string'
        },
        pwd_score: {
            format: 'int64',
            type: 'integer'
        },
        reserved: {
            type: 'boolean'
        },
        sharable: {
            type: 'boolean'
        }
    },
    required: ['anonymous', 'authority_list', 'group', 'id', 'memo', 'name', 'password', 'password_modified_time', 'pwd_score', 'reserved', 'sharable'],
    type: 'object'
} as const;

export const $DeviceVideoInMode = {
    additionalProperties: false,
    properties: {
        '$schema': {
            description: 'A URL to the JSON Schema for this object.',
            examples: ['https://example.com/schemas/DeviceVideoInMode.json'],
            format: 'uri',
            readOnly: true,
            type: 'string'
        },
        switch_mode: {
            type: 'string'
        },
        time_section: {
            type: 'string'
        }
    },
    required: ['switch_mode', 'time_section'],
    type: 'object'
} as const;

export const $DeviceVideoInModeSchedule = {
    additionalProperties: false,
    properties: {
        '$schema': {
            description: 'A URL to the JSON Schema for this object.',
            examples: ['https://example.com/schemas/DeviceVideoInModeSchedule.json'],
            format: 'uri',
            readOnly: true,
            type: 'string'
        },
        latitude: {
            format: 'double',
            type: 'number'
        },
        location: {
            examples: ['UTC'],
            type: 'string'
        },
        longitude: {
            format: 'double',
            type: 'number'
        },
        sunrise_offset: {
            examples: ['0s'],
            type: 'string'
        },
        sunset_offset: {
            examples: ['0s'],
            type: 'string'
        }
    },
    type: 'object'
} as const;

export const $Endpoint = {
    additionalProperties: false,
    properties: {
        '$schema': {
            description: 'A URL to the JSON Schema for this object.',
            examples: ['https://example.com/schemas/Endpoint.json'],
            format: 'uri',
            readOnly: true,
            type: 'string'
        },
        CreatedAt: {
            format: 'date-time',
            type: 'string'
        },
        GoriseURL: {
            type: 'string'
        },
        UUID: {
            type: 'string'
        },
        UpdatedAt: {
            format: 'date-time',
            type: 'string'
        }
    },
    required: ['UUID', 'GoriseURL', 'CreatedAt', 'UpdatedAt'],
    type: 'object'
} as const;

export const $ErrorDetail = {
    additionalProperties: false,
    properties: {
        location: {
            description: "Where the error occurred, e.g. 'body.items[3].tags' or 'path.thing-id'",
            type: 'string'
        },
        message: {
            description: 'Error message text',
            type: 'string'
        },
        value: {
            description: 'The value at the given location'
        }
    },
    type: 'object'
} as const;

export const $ErrorModel = {
    additionalProperties: false,
    properties: {
        '$schema': {
            description: 'A URL to the JSON Schema for this object.',
            examples: ['https://example.com/schemas/ErrorModel.json'],
            format: 'uri',
            readOnly: true,
            type: 'string'
        },
        detail: {
            description: 'A human-readable explanation specific to this occurrence of the problem.',
            examples: ['Property foo is required but is missing.'],
            type: 'string'
        },
        errors: {
            description: 'Optional list of individual error details',
            items: {
                '$ref': '#/components/schemas/ErrorDetail'
            },
            type: 'array'
        },
        instance: {
            description: 'A URI reference that identifies the specific occurrence of the problem.',
            examples: ['https://example.com/error-log/abc123'],
            format: 'uri',
            type: 'string'
        },
        status: {
            description: 'HTTP status code',
            examples: [400],
            format: 'int64',
            type: 'integer'
        },
        title: {
            description: 'A short, human-readable summary of the problem type. This value should not change between occurrences of the error.',
            examples: ['Bad Request'],
            type: 'string'
        },
        type: {
            default: 'about:blank',
            description: 'A URI reference to human-readable documentation for the error.',
            examples: ['https://example.com/errors/example'],
            format: 'uri',
            type: 'string'
        }
    },
    type: 'object'
} as const;

export const $GetHomePage = {
    additionalProperties: false,
    properties: {
        '$schema': {
            description: 'A URL to the JSON Schema for this object.',
            examples: ['https://example.com/schemas/GetHomePage.json'],
            format: 'uri',
            readOnly: true,
            type: 'string'
        },
        build: {
            '$ref': '#/components/schemas/Build'
        },
        db_usage: {
            format: 'int64',
            type: 'integer'
        },
        device_count: {
            format: 'int64',
            type: 'integer'
        },
        email_count: {
            format: 'int64',
            type: 'integer'
        },
        event_count: {
            format: 'int64',
            type: 'integer'
        },
        file_count: {
            format: 'int64',
            type: 'integer'
        },
        file_usage: {
            format: 'int64',
            type: 'integer'
        }
    },
    required: ['device_count', 'event_count', 'email_count', 'file_count', 'db_usage', 'file_usage', 'build'],
    type: 'object'
} as const;

export const $Settings = {
    additionalProperties: false,
    properties: {
        '$schema': {
            description: 'A URL to the JSON Schema for this object.',
            examples: ['https://example.com/schemas/Settings.json'],
            format: 'uri',
            readOnly: true,
            type: 'string'
        },
        latitude: {
            format: 'double',
            type: 'number'
        },
        location: {
            examples: ['UTC'],
            type: 'string'
        },
        longitude: {
            format: 'double',
            type: 'number'
        },
        sunrise_offset: {
            examples: ['0s'],
            type: 'string'
        },
        sunset_offset: {
            examples: ['0s'],
            type: 'string'
        },
        sync_video_in_mode: {
            type: 'boolean'
        },
        updated_at: {
            format: 'date-time',
            type: 'string'
        }
    },
    required: ['location', 'latitude', 'longitude', 'sunrise_offset', 'sunset_offset', 'updated_at', 'sync_video_in_mode'],
    type: 'object'
} as const;

export const $UpdateDevice = {
    additionalProperties: false,
    properties: {
        '$schema': {
            description: 'A URL to the JSON Schema for this object.',
            examples: ['https://example.com/schemas/UpdateDevice.json'],
            format: 'uri',
            readOnly: true,
            type: 'string'
        },
        email: {
            type: 'string'
        },
        features: {
            items: {
                type: 'string'
            },
            type: 'array'
        },
        ip: {
            format: 'ipv4',
            type: 'string'
        },
        latitude: {
            format: 'double',
            type: ['number', 'null']
        },
        location: {
            examples: ['UTC'],
            type: 'string'
        },
        longitude: {
            format: 'double',
            type: ['number', 'null']
        },
        name: {
            type: 'string'
        },
        password: {
            type: 'string'
        },
        sunrise_offset: {
            examples: ['0s'],
            type: 'string'
        },
        sunset_offset: {
            examples: ['0s'],
            type: 'string'
        },
        sync_video_in_mode: {
            type: ['boolean', 'null']
        },
        username: {
            type: 'string'
        }
    },
    required: ['name', 'ip', 'username', 'latitude', 'longitude', 'sync_video_in_mode'],
    type: 'object'
} as const;

export const $UpdateSettingsInput = {
    additionalProperties: false,
    properties: {
        '$schema': {
            description: 'A URL to the JSON Schema for this object.',
            examples: ['https://example.com/schemas/UpdateSettingsInput.json'],
            format: 'uri',
            readOnly: true,
            type: 'string'
        },
        latitude: {
            format: 'double',
            type: 'number'
        },
        location: {
            examples: ['UTC'],
            type: 'string'
        },
        longitude: {
            format: 'double',
            type: 'number'
        },
        sunrise_offset: {
            examples: ['0s'],
            type: 'string'
        },
        sunset_offset: {
            examples: ['0s'],
            type: 'string'
        },
        sync_video_in_mode: {
            type: 'boolean'
        }
    },
    required: ['location', 'latitude', 'longitude', 'sunrise_offset', 'sunset_offset', 'sync_video_in_mode'],
    type: 'object'
} as const;