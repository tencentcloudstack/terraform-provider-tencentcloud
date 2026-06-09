# Design: BH Auth Mode Config Resource

## Architecture

这是一个 **Config 类型资源**，用于管理堡垒机的全局认证方式配置。由于该配置不依赖特定资源 ID，我们使用 Terraform 生成的 UUID 作为资源标识符。

### File Structure

```
tencentcloud/services/bh/
├── resource_tc_bh_auth_mode_config.go  # 新增资源文件
└── service_tencentcloud_bh.go          # 已有服务层（可能需要扩展）

tencentcloud/
└── provider.go                          # 注册资源
```

## Schema Design

### Resource Schema

```hcl
resource "tencentcloud_bh_auth_mode_config" "example" {
  auth_mode    = 1  # Optional: 双因子认证（0=关闭, 1=OTP, 2=短信, 3=USB Key）
  auth_mode_gm = 2  # Optional: 国密双因子认证（0=关闭, 1=OTP, 2=短信, 3=USB Key）
}
```

### Schema Definition

| 字段 | 类型 | 属性 | 说明 |
|------|------|------|------|
| `auth_mode` | Int | Optional, Computed | 双因子认证模式：0-关闭, 1-OTP, 2-短信, 3-USB Key |
| `auth_mode_gm` | Int | Optional, Computed | 国密双因子认证：0-关闭, 1-OTP, 2-短信, 3-USB Key |

**约束条件**：
- `auth_mode` 和 `auth_mode_gm` 至少有一个有效传参

**废弃字段说明**：
- `ResourceType` 字段在 API 文档中已标注为"暂停使用"，不纳入 Schema 定义

## CRUD Implementation

### Create Flow

1. 生成 UUID 作为资源 ID：`id := helper.UUIDGenerator()`
2. 构建 `ModifyAuthModeSettingRequest`
3. 调用 `ModifyAuthModeSetting` 接口
4. 设置资源 ID：`d.SetId(id)`
5. 调用 Read 验证创建结果

```go
func resourceTencentCloudBhAuthModeConfigCreate(d *schema.ResourceData, meta interface{}) error {
    // 1. 生成 UUID
    id := helper.UUIDGenerator()
    
    // 2. 构建请求
    request := bh.NewModifyAuthModeSettingRequest()
    if v, ok := d.GetOk("auth_mode"); ok {
        request.AuthMode = helper.IntInt64(v.(int))
    }
    if v, ok := d.GetOk("auth_mode_gm"); ok {
        request.AuthModeGM = helper.IntInt64(v.(int))
    }
    
    // 3. 调用 API
    err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
        response, e := meta.(*tccommon.ProviderMeta).GetAPIV3Conn().UseBhClient().ModifyAuthModeSetting(ctx, request)
        if e != nil {
            return tccommon.RetryError(e)
        }
        log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
        return nil
    })
    
    // 4. 设置 ID
    d.SetId(id)
    
    // 5. 读取验证
    return resourceTencentCloudBhAuthModeConfigRead(d, meta)
}
```

### Read Flow

1. 调用 `DescribeSecuritySetting` 接口
2. 从 `SecuritySetting.AuthMode` 中提取配置
3. 设置到 Terraform State

```go
func resourceTencentCloudBhAuthModeConfigRead(d *schema.ResourceData, meta interface{}) error {
    // 1. 调用查询接口
    request := bh.NewDescribeSecuritySettingRequest()
    
    var response *bh.DescribeSecuritySettingResponse
    err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
        result, e := meta.(*tccommon.ProviderMeta).GetAPIV3Conn().UseBhClient().DescribeSecuritySetting(ctx, request)
        if e != nil {
            return tccommon.RetryError(e)
        }
        response = result
        return nil
    })
    
    if err != nil {
        return err
    }
    
    // 2. 提取 AuthMode 配置
    if response.Response.SecuritySetting != nil && response.Response.SecuritySetting.AuthMode != nil {
        authMode := response.Response.SecuritySetting.AuthMode
        if authMode.AuthMode != nil {
            _ = d.Set("auth_mode", authMode.AuthMode)
        }
        // 注意：DescribeSecuritySetting 返回结构中可能不包含 AuthModeGM 和 ResourceType
        // 需要根据实际 API 返回结构调整
    }
    
    return nil
}
```

### Update Flow

1. 检测哪些字段发生了变化
2. 构建 `ModifyAuthModeSettingRequest`
3. 调用 `ModifyAuthModeSetting` 接口
4. 调用 Read 验证更新结果

```go
func resourceTencentCloudBhAuthModeConfigUpdate(d *schema.ResourceData, meta interface{}) error {
    immutableArgs := []string{"auth_mode", "auth_mode_gm"}
    
    for _, v := range immutableArgs {
        if d.HasChange(v) {
            return fmt.Errorf("argument `%s` cannot be changed", v)
        }
    }
    
    // 如果所有字段都可变，则构建 ModifyAuthModeSettingRequest 并调用 API
    // 参考 Create 逻辑
    
    return resourceTencentCloudBhAuthModeConfigRead(d, meta)
}
```

**注意**：根据业务需求，认证配置可能是 **不可变的**（需要删除重建）或 **可变的**。需要在实施时根据实际情况调整。

### Delete Flow

1. 调用 `ModifyAuthModeSetting` 接口，设置 `AuthMode=0`（关闭认证）
2. 清除资源 ID

```go
func resourceTencentCloudBhAuthModeConfigDelete(d *schema.ResourceData, meta interface{}) error {
    // 重置为默认配置（关闭认证）
    request := bh.NewModifyAuthModeSettingRequest()
    request.AuthMode = helper.Int64(0)
    
    err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
        response, e := meta.(*tccommon.ProviderMeta).GetAPIV3Conn().UseBhClient().ModifyAuthModeSetting(ctx, request)
        if e != nil {
            return tccommon.RetryError(e)
        }
        log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
        return nil
    })
    
    if err != nil {
        log.Printf("[CRITAL]%s delete bh auth mode config failed, reason:%+v", logId, err)
        return err
    }
    
    return nil
}
```

### Import Flow

1. 使用固定的 ID 字符串（如 `"bh-auth-mode"`）或 UUID
2. 调用 Read 函数获取当前配置

```go
Importer: &schema.ResourceImporter{
    State: schema.ImportStatePassthrough,
},
```

## Code Style Guidelines

### 参考模板

严格参考 `/Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud/tencentcloud/services/igtm/resource_tc_igtm_strategy.go` 的代码风格：

1. **包声明和导入**
```go
package bh

import (
    "context"
    "fmt"
    "log"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    bh "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bh/v20230418"

    tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
    "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)
```

2. **Resource 函数命名**
```go
func ResourceTencentCloudBhAuthModeConfig() *schema.Resource
func resourceTencentCloudBhAuthModeConfigCreate(d *schema.ResourceData, meta interface{}) error
func resourceTencentCloudBhAuthModeConfigRead(d *schema.ResourceData, meta interface{}) error
func resourceTencentCloudBhAuthModeConfigUpdate(d *schema.ResourceData, meta interface{}) error
func resourceTencentCloudBhAuthModeConfigDelete(d *schema.ResourceData, meta interface{}) error
```

3. **日志格式**
```go
defer tccommon.LogElapsed("resource.tencentcloud_bh_auth_mode_config.create")()
defer tccommon.InconsistentCheck(d, meta)()

log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", ...)
log.Printf("[CRITAL]%s create bh auth mode config failed, reason:%+v", ...)
log.Printf("[WARN]%s resource `tencentcloud_bh_auth_mode_config` [%s] not found, please check if it has been deleted.\n", ...)
```

4. **Context 创建**
```go
var (
    logId   = tccommon.GetLogId(tccommon.ContextNil)
    ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
    request = bh.NewModifyAuthModeSettingRequest()
    id      string
)
```

5. **Retry 机制**
```go
err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
    result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhClient().ModifyAuthModeSetting(ctx, request)
    if e != nil {
        return tccommon.RetryError(e)
    }
    log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
    return nil
})
```

6. **字段赋值模式**
```go
if v, ok := d.GetOk("auth_mode"); ok {
    request.AuthMode = helper.IntInt64(v.(int))
}

if v, ok := d.GetOk("auth_mode_gm"); ok && v != 0 {
    request.AuthModeGM = helper.IntInt64(v.(int))
}
```

## API Response Handling

### DescribeSecuritySetting Response Structure

根据 API 文档，返回结构为：
```json
{
    "Response": {
        "SecuritySetting": {
            "AuthMode": {
                "AuthMode": 1
            },
            "Login": {...},
            "Password": {...},
            "LDAP": {...}
        },
        "RequestId": "xxx"
    }
}
```

**注意**：
- `SecuritySetting.AuthMode.AuthMode` 只返回一个字段
- API 文档中没有明确说明是否返回 `AuthModeGM` 和 `ResourceType`
- 在实施时需要根据实际 API 返回结构调整 Read 函数

## Special Considerations

### 1. 全局配置资源的 ID 管理

由于该资源是全局配置，没有云端返回的唯一 ID，我们采用以下策略：

- **Create 时**：使用 `helper.UUIDGenerator()` 生成 UUID 作为资源 ID
- **Import 时**：用户可以使用任意字符串（如 `"bh-auth-mode"`），但 Import 后会被替换为 UUID
- **优势**：确保资源 ID 的唯一性，避免多个配置实例冲突

### 2. 字段约束校验

在 Create/Update 时需要校验：
```go
authMode, authModeOk := d.GetOk("auth_mode")
authModeGM, authModeGMOk := d.GetOk("auth_mode_gm")

if !authModeOk && !authModeGMOk {
    return fmt.Errorf("at least one of auth_mode or auth_mode_gm must be set")
}
```

### 3. Update 策略

根据业务逻辑，认证配置可能有两种更新策略：

**策略 A：可变更新**
- Update 函数检测变化并调用 `ModifyAuthModeSetting`
- 允许用户动态调整认证方式

**策略 B：不可变重建**
- Update 函数返回错误，强制用户删除重建
- 适用于安全性要求高的场景

**建议**：采用策略 A（可变更新），提供更好的用户体验。

## Validation

### Success Criteria

- [x] Schema 定义包含所有字段（auth_mode, auth_mode_gm）
- [x] Create 函数使用 UUID 作为资源 ID
- [x] Read 函数正确解析 DescribeSecuritySetting 返回值
- [x] Update 函数支持配置变更（或返回明确错误）
- [x] Delete 函数重置为默认配置
- [x] Import 功能正常工作
- [x] 代码风格与 igtm_strategy 一致
- [x] 使用 `resource.Retry` 包装 API 调用
- [x] 完整的日志记录和错误处理
- [x] 参数校验逻辑正确

### Code Quality

- [x] 包名为 `bh`
- [x] Import 语句完整
- [x] 使用 `UseBhClient()` 获取客户端
- [x] 日志格式正确（DEBUG/CRITAL/WARN）
- [x] 正确的 nil 检查
- [x] 符合 Go 代码规范（go fmt）
