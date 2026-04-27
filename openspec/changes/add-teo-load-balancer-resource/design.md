# Design: tencentcloud_teo_load_balancer Resource

## Architecture

Follows `tencentcloud_igtm_strategy` style:

```
provider.go
    └─ tencentcloud/services/teo/resource_tc_teo_load_balancer.go  (CRUD handlers)
           └─ tencentcloud/services/teo/service_tencentcloud_teo.go (DescribeTeoLoadBalancerById)
                  └─ teo SDK v20220901
```

## Resource ID

Composite: `<zone_id>#<instance_id>` (using `tccommon.FILED_SP` as separator), e.g. `zone-2ju9lrnpaxol#lb-2my56s2a4fw7`.

## Schema

### Required (Create)

| Field | Type | ForceNew | Description |
|---|---|---|---|
| `zone_id` | String | Yes | Site ID |
| `name` | String | No | Load balancer instance name (1-200 chars, `a-z A-Z 0-9 _ -`) |
| `type` | String | Yes | Instance type: `HTTP` or `GENERAL` |
| `origin_groups` | List (object) | No | Source origin groups with priority; each item: `priority` (String), `origin_group_id` (String) |

### Optional

| Field | Type | Description |
|---|---|---|
| `health_checker` | List (MaxItems:1) | Health check policy; sub-fields: `type`, `port`, `interval`, `timeout`, `health_threshold`, `critical_threshold`, `path`, `method`, `expected_codes`, `follow_redirect`, `headers` (list: `key`,`value`), `send_context`, `recv_context` |
| `steering_policy` | String | Traffic steering policy between origin groups. Valid value: `Pritory`. Default: `Pritory` |
| `failover_policy` | String | Retry policy on request failure. Valid values: `OtherOriginGroup`, `OtherRecordInOriginGroup`. Default: `OtherRecordInOriginGroup` |

### Computed

| Field | Type | Description |
|---|---|---|
| `instance_id` | String | Load balancer instance ID (same as resource ID suffix) |

## Read Logic

Call `DescribeLoadBalancerList` with `ZoneId` and `Filters=[{Name:"InstanceId",Values:[<id>]}]`.
If `TotalCount == 0` or list is empty → resource deleted → `d.SetId("")`.

## Update Logic

All non-ForceNew fields are updatable. On any change, call `ModifyLoadBalancer` with the full set of updated fields (name, origin_groups, health_checker, steering_policy, failover_policy).

## Delete Logic

Call `DeleteLoadBalancer` with `ZoneId` and `InstanceId`. If error contains `LoadBalancerUsedIn*`, surface a clear message to the user to remove references first.

## Key SDK Types

```go
// teo v20220901
CreateLoadBalancerRequest {
    ZoneId          *string
    Name            *string
    Type            *string
    OriginGroups    []*OriginGroupInLoadBalancer
    HealthChecker   *HealthChecker
    SteeringPolicy  *string
    FailoverPolicy  *string
}

OriginGroupInLoadBalancer {
    Priority        *string   // e.g. "priority_1"
    OriginGroupId   *string
}

HealthChecker {
    Type              *string
    Port              *int64
    Interval          *int64
    Timeout           *int64
    HealthThreshold   *int64
    CriticalThreshold *int64
    Path              *string
    Method            *string
    ExpectedCodes     []*string
    FollowRedirect    *string
    Headers           []*CustomizedHeader
    SendContext       *string
    RecvContext       *string
}
```
