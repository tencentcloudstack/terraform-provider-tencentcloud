# Tasks: Add BH Auth Mode Config Resource

## Task 1: 创建资源文件 resource_tc_bh_auth_mode_config.go

**文件**: `tencentcloud/services/bh/resource_tc_bh_auth_mode_config.go`

### 实施步骤

#### 1.1 创建文件并定义 Schema

创建新文件，参考 `resource_tc_igtm_strategy.go` 的代码结构：

**包声明和导入**：
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

**Schema 定义**：
```go
func ResourceTencentCloudBhAuthModeConfig() *schema.Resource {
    return &schema.Resource{
        Create: resourceTencentCloudBhAuthModeConfigCreate,
        Read:   resourceTencentCloudBhAuthModeConfigRead,
        Update: resourceTencentCloudBhAuthModeConfigUpdate,
        Delete: resourceTencentCloudBhAuthModeConfigDelete,
        Importer: &schema.ResourceImporter{
            State: schema.ImportStatePassthrough,
        },
        Schema: map[string]*schema.Schema{
            "auth_mode": {
                Type:        schema.TypeInt,
                Optional:    true,
                Computed:    true,
                Description: "Double factor authentication mode. 0: Disabled, 1: OTP, 2: SMS, 3: USB Key. Note: At least one of auth_mode or auth_mode_gm must be set.",
            },
            "auth_mode_gm": {
                Type:        schema.TypeInt,
                Optional:    true,
                Computed:    true,
                Description: "National secret double factor authentication. 0: Disabled, 1: OTP, 2: SMS, 3: USB Key.",
            },
        },
    }
}
```

#### 1.2 实现 Create 函数

```go
func resourceTencentCloudBhAuthModeConfigCreate(d *schema.ResourceData, meta interface{}) error {
    defer tccommon.LogElapsed("resource.tencentcloud_bh_auth_mode_config.create")()
    defer tccommon.InconsistentCheck(d, meta)()

    var (
        logId   = tccommon.GetLogId(tccommon.ContextNil)
        ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
        request = bh.NewModifyAuthModeSettingRequest()
        id      string
    )

    // 验证至少有一个字段被设置
    authModeSet := false
    if _, ok := d.GetOk("auth_mode"); ok {
        authModeSet = true
    }
    if _, ok := d.GetOk("auth_mode_gm"); ok {
        authModeSet = true
    }
    if !authModeSet {
        return fmt.Errorf("at least one of auth_mode or auth_mode_gm must be set")
    }

    // 构建请求
    if v, ok := d.GetOk("auth_mode"); ok {
        request.AuthMode = helper.IntInt64(v.(int))
    }

    if v, ok := d.GetOk("auth_mode_gm"); ok {
        request.AuthModeGM = helper.IntInt64(v.(int))
    }

    // 调用 API
    err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
        result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhClient().ModifyAuthModeSettingWithContext(ctx, request)
        if e != nil {
            return tccommon.RetryError(e)
        }
        log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
        return nil
    })

    if err != nil {
        log.Printf("[CRITAL]%s create bh auth mode config failed, reason:%+v", logId, err)
        return err
    }

    // 生成 UUID 作为资源 ID
    id = helper.UUIDGenerator()
    d.SetId(id)

    return resourceTencentCloudBhAuthModeConfigRead(d, meta)
}
```

#### 1.3 实现 Read 函数

```go
func resourceTencentCloudBhAuthModeConfigRead(d *schema.ResourceData, meta interface{}) error {
    defer tccommon.LogElapsed("resource.tencentcloud_bh_auth_mode_config.read")()
    defer tccommon.InconsistentCheck(d, meta)()

    var (
        logId   = tccommon.GetLogId(tccommon.ContextNil)
        ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
        request = bh.NewDescribeSecuritySettingRequest()
    )

    var response *bh.DescribeSecuritySettingResponse
    err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
        result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhClient().DescribeSecuritySettingWithContext(ctx, request)
        if e != nil {
            return tccommon.RetryError(e)
        }
        log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
        response = result
        return nil
    })

    if err != nil {
        log.Printf("[CRITAL]%s read bh auth mode config failed, reason:%+v", logId, err)
        return err
    }

    if response == nil || response.Response == nil || response.Response.SecuritySetting == nil {
        log.Printf("[WARN]%s resource `tencentcloud_bh_auth_mode_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
        d.SetId("")
        return nil
    }

    // 设置字段值
    securitySetting := response.Response.SecuritySetting
    if securitySetting.AuthMode != nil && securitySetting.AuthMode.AuthMode != nil {
        _ = d.Set("auth_mode", securitySetting.AuthMode.AuthMode)
    }

    // 注意：根据 API 实际返回调整以下字段
    // 如果 API 不返回这些字段，则需要保持 State 中的值不变

    return nil
}
```

#### 1.4 实现 Update 函数

```go
func resourceTencentCloudBhAuthModeConfigUpdate(d *schema.ResourceData, meta interface{}) error {
    defer tccommon.LogElapsed("resource.tencentcloud_bh_auth_mode_config.update")()
    defer tccommon.InconsistentCheck(d, meta)()

    var (
        logId   = tccommon.GetLogId(tccommon.ContextNil)
        ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
        request = bh.NewModifyAuthModeSettingRequest()
    )

    // 检测变更
    if d.HasChange("auth_mode") || d.HasChange("auth_mode_gm") {
        // 构建请求
        if v, ok := d.GetOk("auth_mode"); ok {
            request.AuthMode = helper.IntInt64(v.(int))
        }

        if v, ok := d.GetOk("auth_mode_gm"); ok {
            request.AuthModeGM = helper.IntInt64(v.(int))
        }

        // 调用 API
        err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
            result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhClient().ModifyAuthModeSettingWithContext(ctx, request)
            if e != nil {
                return tccommon.RetryError(e)
            }
            log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
            return nil
        })

        if err != nil {
            log.Printf("[CRITAL]%s update bh auth mode config failed, reason:%+v", logId, err)
            return err
        }
    }

    return resourceTencentCloudBhAuthModeConfigRead(d, meta)
}
```

#### 1.5 实现 Delete 函数

```go
func resourceTencentCloudBhAuthModeConfigDelete(d *schema.ResourceData, meta interface{}) error {
    defer tccommon.LogElapsed("resource.tencentcloud_bh_auth_mode_config.delete")()
    defer tccommon.InconsistentCheck(d, meta)()

    var (
        logId   = tccommon.GetLogId(tccommon.ContextNil)
        ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
        request = bh.NewModifyAuthModeSettingRequest()
    )

    // 重置为默认配置（关闭认证）
    request.AuthMode = helper.Int64(0)

    err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
        result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhClient().ModifyAuthModeSettingWithContext(ctx, request)
        if e != nil {
            return tccommon.RetryError(e)
        }
        log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
        return nil
    })

    if err != nil {
        log.Printf("[CRITAL]%s delete bh auth mode config failed, reason:%+v", logId, err)
        return err
    }

    return nil
}
```

### 验收标准

- [ ] 文件路径正确：`tencentcloud/services/bh/resource_tc_bh_auth_mode_config.go`
- [ ] 包名为 `bh`
- [ ] Import 语句完整
- [ ] Schema 定义包含所有字段（auth_mode, auth_mode_gm）
- [ ] Create 函数使用 UUID 作为资源 ID
- [ ] Create 函数调用 ModifyAuthModeSetting
- [ ] Read 函数调用 DescribeSecuritySetting 并正确设置 state
- [ ] Update 函数检测变更并调用 ModifyAuthModeSetting
- [ ] Delete 函数重置为默认配置（AuthMode=0）
- [ ] 支持 Import（使用 ImportStatePassthrough）
- [ ] 代码风格与 resource_tc_igtm_strategy.go 一致
- [ ] 使用 `resource.Retry` 包装 API 调用
- [ ] 完整的错误处理和日志记录
- [ ] 正确的 nil 检查
- [ ] 使用 `UseBhClient()` 获取客户端
- [ ] 日志格式正确（DEBUG/CRITAL/WARN）

---

## Task 2: 注册资源到 Provider

**文件**: `tencentcloud/provider.go`

### 实施步骤

#### 2.1 找到 BH 资源注册位置

在 `provider.go` 中搜索已有的 BH 资源（如 `tencentcloud_bh_user`），在同一区域添加新资源。

#### 2.2 添加资源注册

在 `ResourcesMap` 中添加：
```go
"tencentcloud_bh_auth_mode_config": bh.ResourceTencentCloudBhAuthModeConfig(),
```

**注意**：
- 按字母顺序插入
- 确保与其他 BH 资源在同一区域
- 保持缩进和格式一致

### 验收标准

- [x] 在 provider.go 的 ResourcesMap 中正确注册
- [x] 资源名称为 `tencentcloud_bh_auth_mode_config`
- [x] 引用函数名为 `bh.ResourceTencentCloudBhAuthModeConfig()`
- [x] 注册位置与其他 BH 资源相邻
- [x] 按字母顺序排列

---

## Task 3: 代码编译验证

### 实施步骤

#### 3.1 格式化代码

运行以下命令格式化新创建的资源文件：

```bash
cd /Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud
gofmt -w tencentcloud/services/bh/resource_tc_bh_auth_mode_config.go
```

#### 3.2 编译验证

运行以下命令编译 BH 服务包：

```bash
cd /Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud
go build ./tencentcloud/services/bh/
```

#### 3.3 检查 Linter 错误

使用 IDE 或命令行工具检查是否有 linter 错误（可接受 deprecated 警告）。

### 验收标准

- [ ] `gofmt` 格式化成功
- [ ] `go build` 编译成功，无语法错误
- [ ] 无 linter 错误（可接受已有的 deprecated 警告）

---

## 验收清单

### 功能完整性

- [ ] 资源文件包含完整的 CRUD 四个函数
- [ ] Create 调用 ModifyAuthModeSetting 并使用 UUID 作为 ID
- [ ] Read 调用 DescribeSecuritySetting
- [ ] Update 调用 ModifyAuthModeSetting（或返回不支持错误）
- [ ] Delete 调用 ModifyAuthModeSetting 重置配置
- [ ] 支持 Import（使用 ImportStatePassthrough）

### Schema 完整性

- [ ] auth_mode: Optional, Computed
- [ ] auth_mode_gm: Optional, Computed
- [ ] Description 描述清晰准确

### 代码质量

- [ ] 严格参考 igtm_strategy 代码风格
- [ ] 使用 resource.Retry 包装 API 调用
- [ ] 完整的错误处理和日志记录
- [ ] 正确的 nil 检查
- [ ] 参数校验逻辑正确（至少设置一个字段）

### API 集成

- [ ] 使用 `UseBhClient()` 获取客户端
- [ ] API 版本正确：`v20230418`
- [ ] 请求构建正确
- [ ] 响应解析正确

---

## 实施顺序

1. **Task 1**: 创建资源文件 `resource_tc_bh_auth_mode_config.go`（核心任务）
2. **Task 2**: 注册资源到 Provider（集成）
3. **Task 3**: 代码编译验证（质量保证）

---

## 参考信息

### API 文档

- **查询接口**: https://cloud.tencent.com/document/api/1025/125050
  - 接口名称: DescribeSecuritySetting
  - 返回字段: SecuritySetting.AuthMode.AuthMode
  
- **修改接口**: https://cloud.tencent.com/document/api/1025/125048
  - 接口名称: ModifyAuthModeSetting
  - 请求参数: AuthMode, AuthModeGM, ResourceType

### 参考代码

- **代码模板**: `/Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud/tencentcloud/services/igtm/resource_tc_igtm_strategy.go`
- **服务目录**: `/Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud/tencentcloud/services/bh/`

### SDK 信息

- **包路径**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bh/v20230418`
- **客户端**: `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhClient()`

### 字段约束

| 字段 | 可选值 | 约束 |
|------|--------|------|
| auth_mode | 0, 1, 2, 3 | 双因子认证模式 |
| auth_mode_gm | 0, 1, 2, 3 | 国密双因子认证 |

**约束规则**：
- `auth_mode` 和 `auth_mode_gm` 至少有一个有效传参

**废弃字段**：
- `ResourceType` 字段在 API 文档中已标注为"暂停使用"，不纳入 Schema
