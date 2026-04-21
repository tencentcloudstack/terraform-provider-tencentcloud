# Design: EMR Boot Script Config Resource

## File Layout

| File | Action |
|---|---|
| `tencentcloud/services/emr/resource_tc_emr_boot_script_config.go` | New |
| `tencentcloud/services/emr/resource_tc_emr_boot_script_config.md` | New |
| `tencentcloud/services/emr/resource_tc_emr_boot_script_config_test.go` | New |
| `tencentcloud/provider.go` | Modified — register resource |

## Code Style Reference

Strictly follow `resource_tc_config_deliver_config.go`:
- `Create`: `d.SetId(...)` then `return resourceXxxUpdate(d, meta)`
- `Update`: build request, call API wrapped in `resource.Retry(WriteRetryTimeout)`, then call Read
- `Read`: call Describe, populate fields from response
- `Delete`: no-op (just defer log/check)
- Use `tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)` for ctx
- Direct client call: `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().ModifyBootScriptWithContext(ctx, request)`

## SDK Types

- Client method: `UseEmrClient().ModifyBootScriptWithContext` / `DescribeBootScriptWithContext`
- SDK package alias: `emrv20190103 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"`
- `PreExecuteFileSetting` struct fields: `Path`, `Args`, `Bucket`, `Region`, `Domain`, `RunOrder` (*int64), `WhenRun`, `CosFileName`, `CosFileURI`, `CosSecretId`, `CosSecretKey`, `AppId`, `Remark`

## Schema

### Input / ForceNew
```
instance_id   Required, ForceNew, String   — EMR instance ID
boot_type     Required, ForceNew, String   — Valid: resourceAfter, clusterAfter, clusterBefore
```

### Writable (pre_executed_file_settings list)
Each item maps directly to `PreExecuteFileSetting`:
```
path          Optional, String
args          Optional, String
bucket        Optional, String
region        Optional, String
domain        Optional, String
run_order     Optional, Int
when_run      Optional, String   (resourceAfter | clusterAfter)
cos_file_name Optional, String
cos_file_uri  Optional, String
cos_secret_id Optional, String
cos_secret_key Optional, String
app_id        Optional, String
remark        Optional, String
```

## SetId

`strings.Join([]string{instanceId, bootType}, "#")`

## Read Handler

Call `DescribeBootScript` with both `InstanceId` and `BootType`. The response `Detail` contains three slices (`ResourceAfter`, `ClusterBefore`, `ClusterAfter`) — read the slice corresponding to `boot_type` and set it as `pre_executed_file_settings`.

Mapping:
- `resourceAfter` → `Detail.ResourceAfter`
- `clusterBefore` → `Detail.ClusterBefore`
- `clusterAfter` → `Detail.ClusterAfter`
