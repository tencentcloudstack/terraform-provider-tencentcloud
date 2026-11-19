## Context

The `tencentcloud_instance` resource (file `tencentcloud/services/cvm/resource_tc_instance.go`, resource func `ResourceTencentCloudInstance` at line 32) manages CVM instances via the TencentCloud CVM SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312`). Today it exposes many parameters but does NOT expose jumbo frame support.

The CVM cloud API supports an `EnableJumboFrame` boolean parameter on three relevant operations:
- Create: `RunInstances` (`RunInstancesRequest.EnableJumboFrame`)
- Update: `ModifyInstancesAttribute` (`ModifyInstancesAttributeRequest.EnableJumboFrame`)
- Read: returned on the `Instance` struct by `DescribeInstances` (`Instance.EnableJumboFrame`)

The vendored SDK already includes this field on all three structs, so no SDK upgrade is needed.

## Goals / Non-Goals

**Goals:**
- Expose `enable_jumbo_frame` (Optional + Computed bool) on `tencentcloud_instance` so users can enable/disable jumbo frames via Terraform across the full lifecycle (create, update, read/import).

**Non-Goals:**
- Not adding jumbo frame support to other CVM resources (e.g. instance set).
- Not changing instance type validation logic; the cloud API rejects `EnableJumboFrame=true` for unsupported instance types and that error surfaces to the user.
- Not introducing a separate dedicated Modify API; reuse the existing `ModifyInstancesAttribute` endpoint (consistent with how `disable_api_termination`, `cam_role_name`, etc. are updated via per-field service wrappers).

## Decisions

### Decision 1: Schema field definition — `enable_jumbo_frame` (Optional, Computed, TypeBool)

Add the field next to the similar boolean field `disable_api_termination` (currently at schema line ~524). Mirror its definition style:

```go
"enable_jumbo_frame": {
    Type:        schema.TypeBool,
    Optional:    true,
    Computed:    true,
    Description: "Whether to enable jumbo frames for the instance. Only supported instance types accept `true`. See https://cloud.tencent.com/document/product/213/11518 for supported instance types.",
},
```

Rationale: `Computed: true` so that resources created/imported without the field get populated from the API Read, and so the cloud-side default is reflected without forcing the user to set it. Optional (not ForceNew) because `ModifyInstancesAttribute` supports changing it after creation.

### Decision 2: Create flow — set `RunInstancesRequest.EnableJumboFrame` when provided

In `resourceTencentCloudInstanceCreate`, near the existing `disable_api_termination` handling (currently at line ~1066), add:

```go
if v, ok := d.GetOkExists("enable_jumbo_frame"); ok {
    request.EnableJumboFrame = helper.Bool(v.(bool))
}
```

Rationale: follows the exact existing pattern for `disable_api_termination` (`d.GetOkExists` + `helper.Bool`). Using `GetOkExists` (rather than `GetOk`) is consistent with the sibling bool field and correctly distinguishes "explicitly set false" from "unset".

### Decision 3: Update flow — new service method `ModifyEnableJumboFrame`

Add a new service method in `tencentcloud/services/cvm/service_tencentcloud_cvm.go`, mirroring the existing `ModifyDisableApiTermination` (currently at line ~297):

```go
func (me *CvmService) ModifyEnableJumboFrame(ctx context.Context, instanceId string, enableJumboFrame bool) error {
    logId := tccommon.GetLogId(ctx)
    request := cvm.NewModifyInstancesAttributeRequest()
    request.InstanceIds = []*string{&instanceId}
    request.EnableJumboFrame = &enableJumboFrame

    ratelimit.Check(request.GetAction())
    response, err := me.client.UseCvmClient().ModifyInstancesAttribute(request)
    if err != nil {
        log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
            logId, request.GetAction(), request.ToJsonString(), err.Error())
        return err
    }
    log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
        logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

    return nil
}
```

In `resourceTencentCloudInstanceUpdate`, near the existing `disable_api_termination` change handling (currently at line ~2088), add:

```go
if d.HasChange("enable_jumbo_frame") {
    err := cvmService.ModifyEnableJumboFrame(ctx, instanceId, d.Get("enable_jumbo_frame").(bool))
    if err != nil {
        return err
    }
}
```

Rationale: every other mutable attribute on this resource uses a dedicated `ModifyXxx` service wrapper that internally calls `ModifyInstancesAttribute`. Following that convention keeps the codebase consistent and the Update flow readable. The field is NOT added to any `immutableArgs`/ForceNew path because it is mutable.

### Decision 4: Read flow — read `Instance.EnableJumboFrame` into state

In `resourceTencentCloudInstanceRead`, near the existing `disable_api_termination` set (currently at line ~1383), add:

```go
_ = d.Set("enable_jumbo_frame", instance.EnableJumboFrame)
```

Rationale: `DescribeInstances` already returns `EnableJumboFrame` on the `Instance` struct (vendor line 6811). Direct `d.Set` matches the existing pattern for sibling scalar fields (`disable_api_termination`, `cam_role_name`, etc.). No nil-guard is strictly required because `d.Set` of a nil `*bool` is a no-op that leaves state untouched (consistent with the neighboring `disable_api_termination` set).

## Risks / Trade-offs

- [Unsupported instance type sets `EnableJumboFrame=true`] → The cloud API returns an error which surfaces to the user; no provider-side mitigation needed. The schema Description points users to the supported-instance-type doc.
- [Read returns nil for older instances] → `Computed: true` + nil-safe `d.Set` means the field stays unset in state rather than erroring. Acceptable and consistent with sibling fields.
- [Field ordering in Create/Read/Update blocks] → Placing the new code adjacent to `disable_api_termination` keeps diffs localized and reviewable; no behavioral ordering dependency exists between these independent attributes.

## Reference: vendored cloud API field definitions

The following excerpts are from `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312/models.go` and confirm the SDK already supports `EnableJumboFrame`. Implementation should reference these rather than re-reading the large vendor file.

### `RunInstancesRequest.EnableJumboFrame` (models.go line 9995-9996)

```go
// <p>实例是否开启巨型帧，取值范围：<br>&lt;li/&gt; true：表示实例开启巨型帧，只有支持巨型帧的机型可设置为true。<br>&lt;li/&gt;false：表示实例关闭巨型帧，只有支持巨型帧的机型可设置为false。<br> 支持巨型帧的实例规格： <a href="https://cloud.tencent.com/document/product/213/11518">实例规格</a></p>
EnableJumboFrame *bool `json:"EnableJumboFrame,omitnil,omitempty" name:"EnableJumboFrame"`
```

(`RunInstancesRequest` struct starts at line 10002; `EnableJumboFrame` is the last field before `NetworkInterfaces` at line 9996.)

### `ModifyInstancesAttributeRequest.EnableJumboFrame` (models.go line 7973-7974)

```go
// <p>实例是否开启巨型帧，取值范围：<br></p><ul><li>true：表示实例开启巨型帧，只有支持巨型帧的机型可设置为true。</li><li>false：表示实例关闭巨型帧，只有支持巨型帧的机型可设置为false。<br><br>支持巨型帧的实例规格： <a href="https://cloud.tencent.com/document/product/213/11518">实例规格</a></li></ul>
EnableJumboFrame *bool `json:"EnableJumboFrame,omitnil,omitempty" name:"EnableJumboFrame"`
```

(`ModifyInstancesAttributeRequest` struct starts at line 7943; `EnableJumboFrame` is the last field at line 7974, followed by the struct close at 7975.)

### `Instance.EnableJumboFrame` (models.go line 6810-6811) — returned by `DescribeInstances`

```go
// <p>实例是否开启巨型帧</p>
EnableJumboFrame *bool `json:"EnableJumboFrame,omitnil,omitempty" name:"EnableJumboFrame"`
```

(`Instance` struct starts at line 6670; `EnableJumboFrame` is at line 6811, between `LatestOperationErrorMsg` and `Metadata`.)

### Client methods (vendor client.go)

- `func (c *Client) RunInstances(request *RunInstancesRequest) (response *RunInstancesResponse, err error)` — client.go line 9291
- `func (c *Client) ModifyInstancesAttribute(request *ModifyInstancesAttributeRequest) (response *ModifyInstancesAttributeResponse, err error)` — client.go line 6729

Both already wired in the existing resource code paths (`RunInstances` used by Create; `ModifyInstancesAttribute` used by the `ModifyXxx` service wrappers).

## Reference: existing resource code anchors

From `tencentcloud/services/cvm/resource_tc_instance.go` (3276 lines):

- `ResourceTencentCloudInstance()` schema func — line 32; `disable_api_termination` schema block at line ~524 (reference style for the new bool field).
- `resourceTencentCloudInstanceCreate` — line 660; `disable_api_termination` Create handling at line ~1066 (reference pattern: `d.GetOkExists` + `helper.Bool`).
- `resourceTencentCloudInstanceRead` — line 1270; `disable_api_termination` Read set at line ~1383 (reference pattern: `_ = d.Set(...)`).
- `resourceTencentCloudInstanceUpdate` — line 1964; `disable_api_termination` Update handling at line ~2088 (reference pattern: `d.HasChange` + `cvmService.ModifyXxx`).

From `tencentcloud/services/cvm/service_tencentcloud_cvm.go`:

- `ModifyDisableApiTermination(ctx, instanceId, disableApiTermination bool) error` — line 297 (reference pattern for the new `ModifyEnableJumboFrame` wrapper).
