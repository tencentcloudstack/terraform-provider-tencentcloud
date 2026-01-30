## Context
当前 `tencentcloud_cls_data_transform` 资源仅支持基础的数据加工任务创建功能，缺少腾讯云 CLS CreateDataTransform API 提供的高级功能字段。用户需要这些字段来实现完整的数据加工场景，包括：
- 预览测试数据处理
- 动态创建场景下的备份策略
- 服务日志投递
- 前置加工任务
- 失败日志保留
- 时间范围处理
- 外部数据库关联
- 环境变量配置

## Goals / Non-Goals

### Goals
- 补齐所有腾讯云 CLS CreateDataTransform API 支持的字段
- 保持与现有代码风格和架构模式的一致性
- 确保向后兼容，所有新字段均为可选字段
- 支持完整的 CRUD 操作（Create, Read, Update, Delete）

### Non-Goals
- 不修改现有字段的行为或验证逻辑
- 不重构现有代码结构
- 不添加自定义验证逻辑（依赖 API 层验证）
- 不创建新的资源或数据源

## Decisions

### Decision 1: Schema Field Types and Structure
**选择**: 严格按照腾讯云 SDK 模型定义字段类型和结构

**理由**:
- `preview_log_statistics`: 使用 `TypeList` 嵌套 `Resource`，包含 `log_content` (TypeString), `line_num` (TypeInt), `dst_topic_id` (TypeString)
- `backup_give_up_data`: 使用 `TypeBool`，符合 Go SDK 中的 `*bool` 类型
- `has_services_log`, `data_transform_type`, `keep_failure_log`: 使用 `TypeInt`，符合 SDK 中的枚举整数类型
- `failure_log_key`: 使用 `TypeString`
- `process_from_timestamp`, `process_to_timestamp`: 使用 `TypeInt`，表示 Unix 时间戳
- `data_transform_sql_data_sources`: 使用 `TypeList` 嵌套 `Resource`，包含 6 个字段（data_source, region, instance_id, user, alias_name, password）
- `env_infos`: 使用 `TypeList` 嵌套 `Resource`，包含 key-value 键值对

**替代方案**:
- 考虑过使用自定义类型或验证函数，但决定遵循项目现有模式，依赖 API 层验证

### Decision 2: Optional vs Required Fields
**选择**: 所有新增字段均设置为 `Optional: true`

**理由**:
- API 文档中这些字段均为非必填
- 确保向后兼容，不影响现有 Terraform 配置
- 用户可根据实际需求选择性配置

### Decision 3: Sensitive Field Handling
**选择**: 将 `data_transform_sql_data_sources` 中的 `password` 字段标记为 `Sensitive: true`

**理由**:
- 遵循安全最佳实践，避免敏感信息在日志和状态文件中明文显示
- 符合项目现有的敏感字段处理模式

### Decision 4: Immutable Fields
**选择**: 将 `preview_log_statistics` 添加到 `immutableArgs` 列表

**理由**:
- 根据 API 文档和现有代码，`preview_log_statistics` 仅用于创建时的预览，不支持更新
- 如果用户尝试修改，应触发资源重建而非更新

### Decision 5: Read Function Implementation
**选择**: 在 Read 函数中读取所有新增字段（除了 `preview_log_statistics`）

**理由**:
- `preview_log_statistics` 是创建时的输入参数，API 不返回该字段
- 其他字段均可从 DescribeDataTransform API 响应中获取
- 确保 Terraform 状态与实际资源状态同步

### Decision 6: Update Function Implementation
**选择**: 在 Update 函数中支持所有可变字段的更新

**理由**:
- ModifyDataTransform API 支持这些字段的更新
- 使用 `d.HasChange()` 检测变更，仅在字段变更时发送更新请求
- 遵循项目现有的更新模式

### Decision 7: Helper Function Usage
**选择**: 使用项目提供的 `helper` 包函数进行类型转换

**理由**:
- `helper.String()`, `helper.IntInt64()`, `helper.IntUint64()`, `helper.Bool()` 等函数处理指针类型转换
- 保持代码一致性
- 简化空值处理逻辑

### Decision 8: Code Organization
**选择**: 在现有文件中按照字段声明顺序组织新增代码

**理由**:
- Schema 定义: 在现有字段后按照 API 文档顺序添加
- Create 函数: 在现有参数处理逻辑后按顺序添加
- Read 函数: 在现有字段读取逻辑后添加
- Update 函数: 在现有更新逻辑中插入 `d.HasChange()` 检查

## Implementation Details

### Schema Definition Pattern
```go
"field_name": {
    Optional:    true,  // 或 Required: true
    Type:        schema.TypeInt,  // 或其他类型
    Description: "字段描述，包括可选值说明",
},
```

### Nested Object Pattern
```go
"nested_field": {
    Optional:    true,
    Type:        schema.TypeList,
    Description: "嵌套对象描述",
    Elem: &schema.Resource{
        Schema: map[string]*schema.Schema{
            "sub_field": {
                Type:        schema.TypeString,
                Required:    true,
                Description: "子字段描述",
            },
        },
    },
},
```

### Create Function Pattern
```go
if v, ok := d.GetOk("field_name"); ok {
    request.FieldName = helper.String(v.(string))
}

// For list of objects
if v, ok := d.GetOk("list_field"); ok {
    for _, item := range v.([]interface{}) {
        dMap := item.(map[string]interface{})
        obj := cls.ObjectType{}
        if v, ok := dMap["sub_field"]; ok {
            obj.SubField = helper.String(v.(string))
        }
        request.ListField = append(request.ListField, &obj)
    }
}
```

### Read Function Pattern
```go
if response.FieldName != nil {
    _ = d.Set("field_name", response.FieldName)
}

// For list of objects
if response.ListField != nil {
    var list []interface{}
    for _, item := range response.ListField {
        itemMap := map[string]interface{}{}
        if item.SubField != nil {
            itemMap["sub_field"] = item.SubField
        }
        list = append(list, itemMap)
    }
    _ = d.Set("list_field", list)
}
```

### Update Function Pattern
```go
if d.HasChange("field_name") {
    if v, ok := d.GetOk("field_name"); ok {
        request.FieldName = helper.String(v.(string))
    }
}
```

## Risks / Trade-offs

### Risk 1: API Compatibility
**风险**: 腾讯云 API 可能在未来版本中修改字段行为或添加新的验证规则

**缓解措施**:
- 依赖 SDK 版本锁定，避免意外升级
- 所有字段设置为 Optional，降低破坏性变更影响
- 在项目文档中说明各字段的使用场景和限制

### Risk 2: Complex Nested Structures
**风险**: `data_transform_sql_data_sources` 和 `preview_log_statistics` 等嵌套结构可能导致用户配置复杂

**缓解措施**:
- 提供清晰的字段描述和示例
- 在文档中提供完整的使用示例
- 遵循现有项目模式，保持一致性

### Risk 3: Sensitive Data Handling
**风险**: `password` 字段可能在日志或错误信息中泄露

**缓解措施**:
- 标记为 `Sensitive: true`
- 依赖 Terraform 的敏感数据保护机制
- 在文档中提醒用户使用 Terraform variables 或 secrets management

### Risk 4: State Drift
**风险**: 某些字段（如 `preview_log_statistics`）仅用于创建，可能导致状态漂移

**缓解措施**:
- 不在 Read 函数中读取 `preview_log_statistics`
- 标记为 immutable，变更时触发重建
- 在文档中说明该字段的特殊性

## Migration Plan

### Phase 1: Implementation (Current)
1. 添加所有新字段到 Schema 定义
2. 在 Create 函数中实现参数处理
3. 在 Read 函数中实现状态读取
4. 在 Update 函数中实现字段更新
5. 更新 immutableArgs 列表

### Phase 2: Testing
1. 本地测试所有新字段的 Create 操作
2. 测试 Read 操作能正确读取字段值
3. 测试 Update 操作能正确更新字段值
4. 测试 immutable 字段的重建行为

### Phase 3: Documentation
1. 创建 `.changelog` 文件记录变更
2. 更新资源文档（如果项目有文档生成流程）

### Rollback Strategy
- 所有字段均为可选，如果出现问题，用户可以简单地移除这些字段
- 不影响现有配置的正常运行
- 可以通过版本回退恢复到变更前状态

## Open Questions
1. **Q**: `preview_log_statistics` 在 Update API 中是否支持？
   **A**: 根据现有代码和 API 文档，该字段仅用于创建时预览，Update API 不支持。已添加到 immutableArgs。

2. **Q**: `data_transform_sql_data_sources` 中的 `password` 是否需要加密存储？
   **A**: Terraform 通过 `Sensitive: true` 标记提供基本保护，实际加密存储由 Terraform 状态管理机制处理。

3. **Q**: 是否需要为新字段添加自定义验证函数？
   **A**: 不需要。遵循项目现有模式，依赖 API 层进行验证，在错误发生时返回 API 错误信息。
