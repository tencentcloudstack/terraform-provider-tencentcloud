# Add tencentcloud_config_recorder_config Resource

## What

Add a new Terraform resource `tencentcloud_config_recorder_config` for managing the Tencent Cloud Config recorder configuration. This is a global singleton resource that controls:

1. Whether resource monitoring is enabled/disabled (`status` bool — maps to `OpenConfigRecorder` / `CloseConfigRecorder`)
2. Which resource types are monitored (`resource_types` list — maps to `UpdateConfigRecorder`)

## Why

Users need to configure the Config recorder (监控范围) as infrastructure code. Currently there is no Terraform resource for this global config. It allows enabling/disabling monitoring and specifying which resource types to audit.

## APIs Used

| Operation | API | Notes |
|---|---|---|
| Create (initial) | `OpenConfigRecorder` / `CloseConfigRecorder` + `UpdateConfigRecorder` | No dedicated Create; uses Update path |
| Read | `DescribeConfigRecorder` | Returns `Status` (0=off/1=on) and `Items` (monitored resource types) |
| Update (switch) | `OpenConfigRecorder` or `CloseConfigRecorder` | Called when `status` changes |
| Update (types) | `UpdateConfigRecorder` | Called when `resource_types` changes |
| Delete | No-op | Global config cannot be deleted |

## Resource ID

Uses `helper.BuildToken()` (singleton — no natural unique key).
