# Design: Add Deletion Options for CLS Cloud Product Log Task

## Overview

本设计文档描述如何为 `tencentcloud_cls_cloud_product_log_task_v2` 资源添加 `is_delete_topic` 和 `is_delete_logset` 字段，以支持更精细的删除控制。

## Architecture

### 当前架构

```
resourceTencentCloudClsCloudProductLogTaskV2Delete
  ├── DeleteCloudProductLogCollection (删除投递任务)
  └── if force_delete == true:
      ├── DeleteTopic (手动删除 Topic)
      └── DeleteLogset (手动删除 Logset)
```

### 新架构

```
resourceTencentCloudClsCloudProductLogTaskV2Delete
  └── DeleteCloudProductLogCollection (删除投递任务 + 同时删除 Topic/Logset)
      ├── IsDeleteTopic: 由 force_delete 或 is_delete_topic 控制
      └── IsDeleteLogset: 由 force_delete 或 is_delete_logset 控制
```

## Component Design

### 1. Schema 定义

**文件**: `tencentcloud/services/cls/resource_tc_cls_cloud_product_log_task_v2.go`

**位置**: 在现有 `force_delete` 字段后添加

```go
"is_delete_topic": {
    Type:        schema.TypeBool,
    Optional:    true,
    Default:     false,
    Description: "Whether to delete the associated Topic when deleting the log collection task. This field only takes effect when `force_delete` is false. Default is false.",
},

"is_delete_logset": {
    Type:        schema.TypeBool,
    Optional:    true,
    Default:     false,
    Description: "Whether to delete the associated Logset when deleting the log collection task. This field only takes effect when `force_delete` is false. If the Logset has other Topics, it will not be deleted. Default is false.",
},
```

### 2. Read 模块更新

**文件**: `tencentcloud/services/cls/resource_tc_cls_cloud_product_log_task_v2.go`

**函数**: `resourceTencentCloudClsCloudProductLogTaskV2Read`

**修改**: 在读取 `force_delete` 后，添加读取新字段的逻辑

```go
// 现有代码
if v, ok := d.GetOkExists("force_delete"); ok {
    deleteForce = v.(bool)
}
_ = d.Set("force_delete", deleteForce)

// 新增代码
var (
    isDeleteTopic  bool
    isDeleteLogset bool
)

if v, ok := d.GetOkExists("is_delete_topic"); ok {
    isDeleteTopic = v.(bool)
}
_ = d.Set("is_delete_topic", isDeleteTopic)

if v, ok := d.GetOkExists("is_delete_logset"); ok {
    isDeleteLogset = v.(bool)
}
_ = d.Set("is_delete_logset", isDeleteLogset)
```

### 3. Delete 模块重构

**文件**: `tencentcloud/services/cls/resource_tc_cls_cloud_product_log_task_v2.go`

**函数**: `resourceTencentCloudClsCloudProductLogTaskV2Delete`

#### 3.1 变量声明

```go
var (
    logId          = tccommon.GetLogId(tccommon.ContextNil)
    ctx            = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
    request        = clsv20201016.NewDeleteCloudProductLogCollectionRequest()
    deleteForce    bool
    isDeleteTopic  bool
    isDeleteLogset bool
)
```

#### 3.2 读取配置

```go
// 读取 force_delete
if v, ok := d.GetOkExists("force_delete"); ok {
    deleteForce = v.(bool)
}

// 读取 is_delete_topic (只在 force_delete 为 false 时有效)
if !deleteForce {
    if v, ok := d.GetOkExists("is_delete_topic"); ok {
        isDeleteTopic = v.(bool)
    }
}

// 读取 is_delete_logset (只在 force_delete 为 false 时有效)
if !deleteForce {
    if v, ok := d.GetOkExists("is_delete_logset"); ok {
        isDeleteLogset = v.(bool)
    }
}
```

#### 3.3 设置 API 请求参数

```go
request.InstanceId = helper.String(instanceId)
request.AssumerName = helper.String(assumerName)
request.LogType = helper.String(logType)
request.CloudProductRegion = helper.String(cloudProductRegion)

// 设置删除选项
if deleteForce {
    // force_delete 为 true 时，强制删除 Topic 和 Logset
    request.IsDeleteTopic = helper.Bool(true)
    request.IsDeleteLogset = helper.Bool(true)
} else {
    // force_delete 为 false 时，根据用户配置决定
    if isDeleteTopic {
        request.IsDeleteTopic = helper.Bool(true)
    }
    if isDeleteLogset {
        request.IsDeleteLogset = helper.Bool(true)
    }
}
```

#### 3.4 移除手动删除代码

**删除**以下代码块（第 399-444 行）：

```go
if deleteForce {
    var (
        request1 = clsv20201016.NewDeleteTopicRequest()
        request2 = clsv20201016.NewDeleteLogsetRequest()
    )
    
    // ... DeleteTopic 调用
    // ... DeleteLogset 调用
}
```

因为 API 会自动处理删除，不需要手动调用。

## Data Flow

### 删除流程

```
用户调用 terraform destroy
    ↓
Read force_delete, is_delete_topic, is_delete_logset
    ↓
判断 force_delete
    ├─ true  → IsDeleteTopic=true, IsDeleteLogset=true
    └─ false → IsDeleteTopic=is_delete_topic, IsDeleteLogset=is_delete_logset
    ↓
调用 DeleteCloudProductLogCollection(IsDeleteTopic, IsDeleteLogset)
    ↓
API 内部处理：删除投递任务 + 根据参数删除 Topic/Logset
    ↓
等待删除完成 (Status = 3)
    ↓
完成
```

## API Integration

### DeleteCloudProductLogCollection 接口

**接口文档**: https://cloud.tencent.com/document/api/614/117420

**请求参数**:

| 参数名 | 类型 | 必选 | 描述 |
|--------|------|------|------|
| InstanceId | String | 是 | 实例 ID |
| AssumerName | String | 是 | 云产品标识 |
| LogType | String | 是 | 日志类型 |
| CloudProductRegion | String | 是 | 云产品地域 |
| **IsDeleteTopic** | Boolean | 否 | 是否删除关联的 Topic |
| **IsDeleteLogset** | Boolean | 否 | 是否删除关联的 Logset |

**响应参数**:

| 参数名 | 类型 | 描述 |
|--------|------|------|
| Status | Integer | 0-创建中, 1-创建完成, 2-删除中, 3-删除完成 |
| RequestId | String | 唯一请求 ID |

## Error Handling

### 场景 1: Logset 有其他 Topic

**情况**: `IsDeleteLogset=true`，但 Logset 下还有其他 Topic

**行为**: API 不会删除 Logset，但不会返回错误

**处理**: 无需特殊处理，API 内部会正确判断

### 场景 2: 资源已被删除

**情况**: Topic 或 Logset 已在 Terraform 外部被删除

**行为**: API 可能返回资源不存在错误

**处理**: 由 `tccommon.RetryError` 处理，符合 Terraform 的幂等性要求

## Validation Rules

### Schema 级别验证

- `is_delete_topic`: 布尔类型，可选，默认 `false`
- `is_delete_logset`: 布尔类型，可选，默认 `false`
- `force_delete`: 布尔类型，可选，默认 `false`（已存在）

### 逻辑级别验证

在 Delete 函数中：

```go
// force_delete 优先级最高
if deleteForce {
    // 忽略 is_delete_topic 和 is_delete_logset 的值
    request.IsDeleteTopic = helper.Bool(true)
    request.IsDeleteLogset = helper.Bool(true)
} else {
    // 只有 force_delete 为 false 时，才使用 is_delete_* 字段
    if isDeleteTopic {
        request.IsDeleteTopic = helper.Bool(true)
    }
    if isDeleteLogset {
        request.IsDeleteLogset = helper.Bool(true)
    }
}
```

## Testing Strategy

### 手动测试场景

1. **场景 1**: `force_delete=true` → 应删除 Topic 和 Logset
2. **场景 2**: `force_delete=false, is_delete_topic=true` → 只删除 Topic
3. **场景 3**: `force_delete=false, is_delete_logset=true` → 只删除 Logset
4. **场景 4**: `force_delete=false, is_delete_topic=true, is_delete_logset=true` → 删除两者
5. **场景 5**: `force_delete=false` (默认) → 不删除 Topic 和 Logset

### 验证点

- Schema 字段定义正确
- Read 模块正确读取和设置字段值
- Delete 模块逻辑正确，API 调用成功
- 代码格式化（`go fmt`）通过
- 编译无错误

## Implementation Notes

### 注意事项

1. **向后兼容**: 保留 `force_delete` 字段，确保现有配置不受影响
2. **代码格式化**: 每次修改后执行 `go fmt`
3. **API 字段映射**: 
   - `force_delete=true` → `IsDeleteTopic=true, IsDeleteLogset=true`
   - `is_delete_topic` → `IsDeleteTopic`
   - `is_delete_logset` → `IsDeleteLogset`
4. **删除手动调用**: 移除 DeleteTopic 和 DeleteLogset 的手动调用代码
5. **变量命名**: 使用 `isDeleteTopic` 和 `isDeleteLogset`（驼峰命名）

### 代码位置

- **文件**: `tencentcloud/services/cls/resource_tc_cls_cloud_product_log_task_v2.go`
- **Schema 修改**: 第 101-106 行附近
- **Read 修改**: 第 273-278 行附近
- **Delete 修改**: 第 347-447 行

## Dependencies

- SDK: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016`
- API 版本: 2020-10-16
- 最低 API 支持: DeleteCloudProductLogCollection 接口需支持 IsDeleteTopic 和 IsDeleteLogset 参数

## Rollback Plan

如果实现出现问题，回滚步骤：

1. 恢复 Delete 模块中手动调用 DeleteTopic 和 DeleteLogset 的代码
2. 移除 `is_delete_topic` 和 `is_delete_logset` 字段
3. 保持 `force_delete` 字段的原有行为

## Success Metrics

- ✅ Schema 定义正确，包含新字段
- ✅ Read 模块正确处理新字段
- ✅ Delete 模块使用 API 原生参数
- ✅ 移除手动删除代码
- ✅ 代码通过 `go fmt` 格式化
- ✅ 编译无错误
- ✅ 向后兼容性保持
