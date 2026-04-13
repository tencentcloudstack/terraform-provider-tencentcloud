# Design: tencentcloud_config_deliver_config Resource

## Architecture

Follows the igtm_strategy code style, singleton global config pattern (no Create/Delete API):

```
provider.go (registration)
    └─ resource_tc_config_deliver_config.go   (Create=Update, Read, Update, Delete=no-op)
           └─ service_tencentcloud_config.go  (DescribeConfigDeliver method)
                  └─ config SDK v20220802
```

## File Layout

| File | Action |
|---|---|
| `tencentcloud/services/config/resource_tc_config_deliver_config.go` | New — resource |
| `tencentcloud/services/config/resource_tc_config_deliver_config.md` | New — doc |
| `tencentcloud/services/config/resource_tc_config_deliver_config_test.go` | New — test |
| `tencentcloud/services/config/service_tencentcloud_config.go` | Modified — append `DescribeConfigDeliver` |
| `tencentcloud/provider.go` | Modified — register resource |

## Schema

### Required

| Field | Type | Description |
|---|---|---|
| `status` | Int | Delivery switch. `0`: disabled, `1`: enabled |

### Optional

| Field | Type | Description |
|---|---|---|
| `deliver_name` | String | Delivery service name |
| `target_arn` | String | Resource ARN (COS: `qcs::cos:$region:$account:prefix/$appid/$BucketName`; CLS: `qcs::cls:$region:$account:cls/topicId`) |
| `deliver_prefix` | String | Log prefix for stored delivery content |
| `deliver_type` | String | Delivery type. Valid values: `COS`, `CLS` |
| `deliver_content_type` | Int | Content type. `1`: config change, `2`: resource list, `3`: all |

### Computed

| Field | Type | Description |
|---|---|---|
| `create_time` | String | Creation time |

## Lifecycle

| Handler | Behavior |
|---|---|
| Create | Call `UpdateConfigDeliver` with all provided fields; `d.SetId(helper.BuildToken())` |
| Read | Call `DescribeConfigDeliver`; populate state |
| Update | Call `UpdateConfigDeliver` with changed fields |
| Delete | No-op (set `Status=0` externally if needed) |

No `Importer` — singleton cannot be meaningfully imported by ID.
