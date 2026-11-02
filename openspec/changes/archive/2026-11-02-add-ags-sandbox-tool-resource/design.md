## Context

TencentCloud AGS (Agent Sandbox) provides sandbox tools for running code, browser, and other environments. The `tencentcloud_ags_sandbox_tool` Terraform resource manages these tools' full lifecycle. The AGS SDK v20250920 is already vendored at `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920/`.

**Code patterns** follow `resource_tc_igtm_strategy.go` in `tencentcloud/services/igtm/`.

## Goals / Non-Goals

**Goals:**
- Provide `tencentcloud_ags_sandbox_tool` resource with full CRUD (CreateSandboxTool, DescribeSandboxToolList, UpdateSandboxTool, DeleteSandboxTool)
- Support all sandbox tool parameters: tool_name, tool_type, network_configuration, description, default_timeout, tags, client_token, role_arn, storage_mounts, custom_configuration, computer_configuration, log_configuration, persistent
- Support import via `tool_id`
- Native synchronous APIs (no async polling needed for CRUD)

**Non-Goals:**
- No data source for sandbox tool list in this change
- No support for sandbox instance management
- No support for sandbox snapshot management

## Decisions

### 1. Resource ID: Use single `tool_id`
The AGS sandbox tool is identified by a single `ToolId`. Unlike `igtm_strategy` (which uses `instanceId#strategyId`), the resource ID is simply the `tool_id` string. No composite key needed.

### 2. Schema nesting approach
Follow the SDK's nested struct pattern: `network_configuration` is a `TypeList` with `MaxItems: 1`, containing `network_mode` (string) and `vpc_config` (nested `TypeList`). All nested objects use `TypeList` with `MaxItems: 1` for single-object containers (following Terraform convention for nested blocks).

### 3. Tags handling
Tags are passed as `[]*Tag` array to CreateSandboxTool and UpdateSandboxTool. In the schema, tags use `TypeList` of `TypeMap`. Read response also returns `[]*Tag`. Use `helper.SetTags` / `helper.GetTags` pattern for tag serialization.

### 4. Update immutability
Fields that cannot be changed after creation (ForceNew): `tool_name`, `tool_type`, `client_token`, `role_arn`, `storage_mounts`, `persistent`. Other fields can be updated via UpdateSandboxTool.

### 5. Schema design for nested lists
- **storage_mounts**: `TypeList` of objects, each containing `name`, `storage_source` (nested), `mount_path`, `read_only`, `volume` (nested), `volume_template` (nested)
- **custom_configuration**: `TypeList` with `MaxItems: 1`, containing `image`, `image_registry_type`, `command`, `args`, `env`, `ports`, `resources`, `probe`, `dns_config`, `security_context`
- **computer_configuration**: `TypeList` with `MaxItems: 1`, containing `waa_configuration` and `os_world_configuration`
- **log_configuration**: `TypeList` with `MaxItems: 1`, containing `cls_config` and `log_sources`

### 6. Read implementation
Use `DescribeSandboxToolList` with `ToolIds` set to the resource id. Check response for empty set → set `d.SetId("")`. The response returns `SandboxToolSet` array; read the first element.

## API Reference (from vendor SDK)

### CreateSandboxTool

**Request** (`models.go` lines 1386-1427):
```go
type CreateSandboxToolRequest struct {
    ToolName             *string                `json:"ToolName"`
    ToolType             *string                `json:"ToolType"`
    NetworkConfiguration *NetworkConfiguration  `json:"NetworkConfiguration"`
    Description          *string                `json:"Description"`
    DefaultTimeout       *string                `json:"DefaultTimeout"`
    Tags                 []*Tag                 `json:"Tags"`
    ClientToken          *string                `json:"ClientToken"`
    RoleArn              *string                `json:"RoleArn"`
    StorageMounts        []*StorageMount        `json:"StorageMounts"`
    CustomConfiguration  *CustomConfiguration   `json:"CustomConfiguration"`
    ComputerConfiguration *ComputerConfiguration `json:"ComputerConfiguration"`
    LogConfiguration     *LogConfiguration      `json:"LogConfiguration"`
    Persistent           *bool                  `json:"Persistent"`
}
```

**Response** (`models.go` lines 1461-1472):
```go
type CreateSandboxToolResponseParams struct {
    ToolId    *string `json:"ToolId"`
    RequestId *string `json:"RequestId"`
}
```

### DescribeSandboxToolList

**Request** (`models.go` lines 3677-3691):
```go
type DescribeSandboxToolListRequest struct {
    ToolIds []*string `json:"ToolIds"`
    Offset  *int64    `json:"Offset"`
    Limit   *int64    `json:"Limit"`
    Filters []*Filter `json:"Filters"`
}
```

**Response** (`models.go` lines 3716-3730):
```go
type DescribeSandboxToolListResponseParams struct {
    SandboxToolSet []*SandboxTool `json:"SandboxToolSet"`
    TotalCount     *int64         `json:"TotalCount"`
    RequestId      *string        `json:"RequestId"`
}
```

### UpdateSandboxTool

**Request** (`models.go` lines 5881-5901):
```go
type UpdateSandboxToolRequest struct {
    ToolId               *string                `json:"ToolId"`
    Description          *string                `json:"Description"`
    NetworkConfiguration *NetworkConfiguration  `json:"NetworkConfiguration"`
    Tags                 []*Tag                 `json:"Tags"`
    CustomConfiguration  *CustomConfiguration   `json:"CustomConfiguration"`
    ComputerConfiguration *ComputerConfiguration `json:"ComputerConfiguration"`
}
```

**Response**: Only `RequestId` (no data fields).

### DeleteSandboxTool

**Request** (`models.go` lines 2107-2112):
```go
type DeleteSandboxToolRequest struct {
    ToolId *string `json:"ToolId"`
}
```

**Response**: Only `RequestId` (no data fields).

### Key Sub-structs

**NetworkConfiguration** (line 4647):
```go
type NetworkConfiguration struct {
    NetworkMode *string     `json:"NetworkMode"`
    VpcConfig   *VPCConfig  `json:"VpcConfig"`
}
```

**VPCConfig** (line 5949):
```go
type VPCConfig struct {
    SubnetIds        []*string `json:"SubnetIds"`
    SecurityGroupIds []*string `json:"SecurityGroupIds"`
}
```

**StorageMount** (line 5440):
```go
type StorageMount struct {
    Name            *string              `json:"Name"`
    StorageSource   *StorageSource       `json:"StorageSource"`
    MountPath       *string              `json:"MountPath"`
    ReadOnly        *bool                `json:"ReadOnly"`
    Volume          *VolumeRef           `json:"Volume"`
    VolumeTemplate  *VolumeTemplateRef   `json:"VolumeTemplate"`
}
```

**StorageSource** (line 5463):
```go
type StorageSource struct {
    Cos         *CosStorageSource         `json:"Cos"`
    Image       *ImageStorageSource       `json:"Image"`
    Cfs         *CfsStorageSource         `json:"Cfs"`
    AgentBucket *AgentBucketStorageSource `json:"AgentBucket"`
}
```

**CosStorageSource** (line 807): `Endpoint`, `BucketName`, `BucketPath` (all *string)
**CfsStorageSource** (line 446): `FileSystemId`, `Path` (all *string)
**AgentBucketStorageSource** (line 179): `LibraryId`, `SpaceId`, `AccessDomain` (all *string)
**ImageStorageSource** (line 4345): `Reference`, `ImageRegistryType`, `SubPath`, `Digest` (all *string)

**VolumeRef** (line 5957):
```go
type VolumeRef struct {
    VolumeName *string `json:"VolumeName"`
    VolumeId   *string `json:"VolumeId"`
}
```

**VolumeTemplateRef** (line 5967):
```go
type VolumeTemplateRef struct {
    VolumeTemplateName *string `json:"VolumeTemplateName"`
    ReuseKey           *string `json:"ReuseKey"`
}
```

**CustomConfiguration** (line 1657):
```go
type CustomConfiguration struct {
    Image             *string                `json:"Image"`
    ImageRegistryType *string                `json:"ImageRegistryType"`
    Command           []*string              `json:"Command"`
    Args              []*string              `json:"Args"`
    Env               []*EnvVar              `json:"Env"`
    Ports             []*PortConfiguration   `json:"Ports"`
    Resources         *ResourceConfiguration `json:"Resources"`
    Probe             *ProbeConfiguration    `json:"Probe"`
    DNSConfig         *DNSConfig             `json:"DNSConfig"`
    SecurityContext   *SecurityContext       `json:"SecurityContext"`
}
```

**CustomConfigurationDetail** (response, line 1689): Same as above but adds `ImageDigest *string`.

**EnvVar** (line 4072): `Name`, `Value` (all *string)
**PortConfiguration** (line 4724): `Name *string`, `Port *int64`, `Protocol *string`
**ResourceConfiguration** (line 4980): `CPU *string`, `Memory *string`, `Storage *string`
**ProbeConfiguration** (line 4845): `HttpGet *HttpGetAction`, `ReadyTimeoutMs`, `ProbeTimeoutMs`, `ProbePeriodMs`, `SuccessThreshold`, `FailureThreshold` (all *int64)
**HttpGetAction** (line 4334): `Path *string`, `Port *int64`, `Scheme *string`
**DNSConfig** (line 1724): `Servers []*string`, `Searches []*string`, `Options []*string`
**SecurityContext** (line 5196): `Privileged *bool`, `Capabilities *Capabilities`
**Capabilities** (line 438): `Add []*string`, `Drop []*string`

**ComputerConfiguration** (line 799):
```go
type ComputerConfiguration struct {
    WAAConfiguration     *WAAConfiguration     `json:"WAAConfiguration"`
    OSWorldConfiguration *OSWorldConfiguration `json:"OSWorldConfiguration"`
}
```

**WAAConfiguration** (line 5977): `ImageId *string`
**OSWorldConfiguration** (line 4655): `Version *string`

**LogConfiguration** (line 4377):
```go
type LogConfiguration struct {
    CLSConfig  *CLSConfig  `json:"CLSConfig"`
    LogSources *LogSources `json:"LogSources"`
}
```

**CLSConfig** (line 355): `TopicId *string`
**LogSources** (line 4385): `Files []*string`

**Filter** (line 4159): `Name *string`, `Values []*string`
**Tag** (line 5585): `Key *string`, `Value *string`

**SandboxTool** (response, line 5132):
```go
type SandboxTool struct {
    ToolId               *string                    `json:"ToolId"`
    ToolName             *string                    `json:"ToolName"`
    ToolType             *string                    `json:"ToolType"`
    Status               *string                    `json:"Status"`
    Description          *string                    `json:"Description"`
    Persistent           *bool                      `json:"Persistent"`
    DefaultTimeoutSeconds *uint64                   `json:"DefaultTimeoutSeconds"`
    NetworkConfiguration *NetworkConfiguration      `json:"NetworkConfiguration"`
    Tags                 []*Tag                     `json:"Tags"`
    CreateTime           *string                    `json:"CreateTime"`
    UpdateTime           *string                    `json:"UpdateTime"`
    RoleArn              *string                    `json:"RoleArn"`
    StorageMounts        []*StorageMount            `json:"StorageMounts"`
    CustomConfiguration  *CustomConfigurationDetail `json:"CustomConfiguration"`
    LogConfiguration     *LogConfiguration          `json:"LogConfiguration"`
    ComputerConfiguration *ComputerConfiguration    `json:"ComputerConfiguration"`
    StatusReason         *string                    `json:"StatusReason"`
}
```

## Reference Code Pattern (resource_tc_igtm_strategy.go)

**Resource function signature:**
```go
func ResourceTencentCloudIgtmStrategy() *schema.Resource {
    return &schema.Resource{
        Create: resourceTencentCloudIgtmStrategyCreate,
        Read:   resourceTencentCloudIgtmStrategyRead,
        Update: resourceTencentCloudIgtmStrategyUpdate,
        Delete: resourceTencentCloudIgtmStrategyDelete,
        Importer: &schema.ResourceImporter{
            State: schema.ImportStatePassthrough,
        },
        Schema: map[string]*schema.Schema{...},
    }
}
```

**Create flow** (lines 176-322):
1. Build request from schema.ResourceData using `d.GetOk()`
2. Use `resource.Retry(tccommon.WriteRetryTimeout, ...)` for API call
3. Check `result == nil || result.Response == nil` → `NonRetryableError`
4. Check `response.Response.XxxId == nil` → return error
5. `d.SetId(...)` with the ID
6. Return `resourceXxxRead(d, meta)`

**Read flow** (lines 324-472):
1. Parse `d.Id()`
2. Call describe API via service layer (or directly)
3. If response empty → `log.Printf("[WARN]...")` + `d.SetId("")` + `return nil`
4. Set each field via `d.Set()` after nil-check

**Update flow** (lines 474-635):
1. Parse `d.Id()`
2. Check `d.HasChange()` for each field
3. Build update request with changed fields
4. Call update API with retry
5. Return `resourceXxxRead(d, meta)`

**Delete flow** (lines 635-680):
1. Parse `d.Id()`
2. Build delete request
3. Call delete API with retry
4. Return nil

**Key imports:**
```go
import (
    tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
    "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)
```

## Risks / Trade-offs

- **API is synchronous**: No async polling needed for CRUD; simplifies implementation
- **UpdateSandboxTool has limited fields**: Only `Description`, `NetworkConfiguration`, `Tags`, `CustomConfiguration`, `ComputerConfiguration` are updatable. Other fields require ForceNew (destroy + recreate)
- **Schema complexity**: Deeply nested structures (e.g., storage_mounts → storage_source → cos/image/cfs/agent_bucket) require careful flatten/expand logic
- **Read response uses CustomConfigurationDetail** (different from CustomConfiguration used in Create/Update): `CustomConfigurationDetail` adds `ImageDigest` field not present in input `CustomConfiguration`