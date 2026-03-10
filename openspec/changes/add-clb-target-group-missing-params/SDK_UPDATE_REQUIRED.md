# SDK Update Required for SnatEnable Parameter

## Issue

The Tencent Cloud API documentation for `CreateTargetGroup` (updated 2026-01-26) includes a `SnatEnable` parameter, but the current Go SDK version does not include this field in the `CreateTargetGroupRequest` structure.

## API Documentation

**Source**: https://cloud.tencent.com/document/product/214/40559

**Parameter Details**:
- **Name**: SnatEnable
- **Type**: Boolean
- **Required**: No
- **Description**: 是否开启SNAT(源IP替换)。True(开启)、False(关闭)。默认为关闭。注意:SnatEnable开启时会替换客户端源IP,此时`透传客户端源IP`选项关闭,反之亦然。

## Current SDK Status

**Package**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb`
**Current Version**: v1.3.52
**Upgraded From**: v1.3.8 → v1.3.52 (on 2026-03-10)

### Verified SDK Fields

Checked `CreateTargetGroupRequest` structure in `v1.3.52`:
- ✅ TargetGroupName
- ✅ VpcId
- ✅ Port
- ✅ TargetGroupInstances
- ✅ Type
- ✅ Protocol
- ✅ HealthCheck
- ✅ ScheduleAlgorithm
- ✅ Tags
- ✅ Weight
- ✅ FullListenSwitch
- ✅ KeepaliveEnable
- ✅ SessionExpireTime
- ✅ IpVersion
- ❌ **SnatEnable** - NOT PRESENT

### Investigation

Checked GitHub repository (master branch):
```bash
https://github.com/TencentCloud/tencentcloud-sdk-go/blob/master/tencentcloud/clb/v20180317/models.go
```

According to web fetch results, the latest GitHub code **DOES** include `SnatEnable` field:
```go
// 是否开启SNAT(源IP替换),True(开启)、False(关闭)。默认为关闭。注意:SnatEnable开启时会替换客户端源IP,此时`透传客户端源IP`选项关闭,反之亦然。
SnatEnable *bool `json:"SnatEnable,omitnil,omitempty" name:"SnatEnable"`
```

## Root Cause

The SDK Go module releases lag behind the GitHub master branch. The `SnatEnable` field has been added to the SDK source code but has not been released in a tagged version yet.

## Required Actions

### Option 1: Wait for Official SDK Release (Recommended)

1. Monitor SDK releases: https://github.com/TencentCloud/tencentcloud-sdk-go/releases
2. When a new version (e.g., v1.3.53+) is released with `SnatEnable` support:
   ```bash
   go get github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb@latest
   go mod tidy
   go mod vendor
   ```
3. Complete the implementation of this change proposal

**Timeline**: Could be days to weeks depending on SDK release schedule.

### Option 2: Contact Tencent Cloud SDK Team

File an issue or contact the SDK maintainers to request:
- Expedited release of the `SnatEnable` field
- Provide context that the API documentation has been updated but SDK is lagging

**Contact Methods**:
- GitHub Issues: https://github.com/TencentCloud/tencentcloud-sdk-go/issues
- Tencent Cloud Support: 工单系统

### Option 3: Defer SnatEnable Implementation

Implement all other 8 parameters first:
1. health_check
2. schedule_algorithm
3. tags
4. weight
5. full_listen_switch
6. keepalive_enable
7. session_expire_time
8. ip_version

Add `snat_enable` in a follow-up change when SDK is updated.

## Implementation Notes

When SDK is updated, add to `resource_tc_clb_target_group.go`:

### Schema Definition
```go
"snat_enable": {
    Type:        schema.TypeBool,
    Optional:    true,
    Default:     false,
    Description: "Enable SNAT (source IP replacement). When enabled, client source IPs are replaced with the load balancer's IP. Default is false. Note: When SnatEnable is enabled, the 'transparent client IP' option is disabled, and vice versa.",
},
```

### Create Logic
```go
if v, ok := d.GetOkExists("snat_enable"); ok {
    request.SnatEnable = helper.Bool(v.(bool))
}
```

### Read Logic
```go
_ = d.Set("snat_enable", targetGroup.SnatEnable)
```

### Update Logic
```go
if d.HasChange("snat_enable") {
    request.SnatEnable = helper.Bool(d.Get("snat_enable").(bool))
}
```

## Verification

To verify SDK version includes `SnatEnable`:

```bash
# Check structure definition
grep -A 50 "type CreateTargetGroupRequest struct" \
  vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317/models.go \
  | grep SnatEnable

# Expected output:
# SnatEnable *bool `json:"SnatEnable,omitnil,omitempty" name:"SnatEnable"`
```

## Status

- **Proposal Updated**: ✅ SnatEnable added to spec.md and proposal.md
- **SDK Updated**: ❌ Waiting for SDK release with SnatEnable field
- **Implementation**: ⏸️ Blocked pending SDK update
- **Next Action**: Monitor SDK releases or contact SDK team

## References

- API Documentation: https://cloud.tencent.com/document/product/214/40559
- SDK Repository: https://github.com/TencentCloud/tencentcloud-sdk-go
- Current SDK Version: v1.3.52
- Provider go.mod: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb v1.3.52`
