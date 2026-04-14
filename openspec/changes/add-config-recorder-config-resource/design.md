# Design: tencentcloud_config_recorder_config Resource

## Architecture

Follows the `tencentcloud_bh_reconnection_setting_config` pattern (Create delegates to Update):

```
provider.go
    └─ resource_tc_config_recorder_config.go (Create=setId+Update, Read, Update, Delete=no-op)
           └─ service_tencentcloud_config.go (DescribeConfigRecorder)
                  └─ config SDK v20220802
```

## File Layout

| File | Action |
|---|---|
| `resource_tc_config_recorder_config.go` | New |
| `resource_tc_config_recorder_config.md` | New |
| `resource_tc_config_recorder_config_test.go` | New |
| `service_tencentcloud_config.go` | Modified — append `DescribeConfigRecorder` |
| `provider.go` | Modified — register resource |

## Schema

| Field | Type | Required | Description |
|---|---|---|---|
| `status` | Bool | Required | `true`: enable monitoring (OpenConfigRecorder), `false`: disable (CloseConfigRecorder) |
| `resource_types` | List of String | Optional | Resource type list to monitor (e.g. `QCS::CAM::Group`) |

### Computed

| Field | Type | Description |
|---|---|---|
| `create_time` | String | Recorder creation time |
| `trigger_count` | Int | Number of snapshots taken today |
| `open_count` | Int | Number of times monitoring was opened today |
| `update_count` | Int | Number of monitoring range updates today |

## Update Logic

```
if status changed:
    if status == true  → call OpenConfigRecorder
    if status == false → call CloseConfigRecorder
if resource_types changed:
    call UpdateConfigRecorder
```

## Read Strategy

`DescribeConfigRecorder` returns `Status` (`0`=off, `1`=on) and `Items` (`[]*UserConfigResource`). Map `Status` to `bool` for the `status` field.
