# Design: CLS Alarm Classifications 参数支持

## Context

当前 `tencentcloud_cls_alarm` 资源已支持告警的核心配置(告警目标、触发条件、通知策略等),但缺少 `Classifications` 字段的支持。腾讯云 CLS 的 `CreateAlarm`、`ModifyAlarm` 和 `DescribeAlarms` API 已支持该参数,用于告警的分类管理。

**当前状态**:
- 资源文件: `tencentcloud/services/cls/resource_tc_cls_alarm.go`
- 涉及 API: CreateAlarm, ModifyAlarm, DescribeAlarms, DeleteAlarm
- 现有 schema 已包含约 10+ 个字段,此次新增为简单字段追加

**约束**:
- 必须保持向后兼容(Optional 字段)
- 不能破坏现有 state 和配置
- 遵循 Terraform SDK v2 标准模式

## Goals / Non-Goals

**Goals:**
- 在资源 schema 中新增 `classifications` 字段(Optional + Computed)
- 在 Create/Update 操作中支持设置和修改 `Classifications` 参数
- 在 Read 操作中读取并同步 `Classifications` 到 state
- 通过验收测试验证完整的 CRUD 流程
- 更新文档说明新字段用法

**Non-Goals:**
- 不涉及数据源(`data_source_tc_cls_alarms`)的修改(可后续单独处理)
- 不修改其他现有字段
- 不涉及 API 的 Breaking Change

## Decisions

### 决策 1: 字段类型和属性

**选择**: `schema.TypeList` + `Elem: &schema.Schema{Type: schema.TypeString}`
- `Optional: true` - 不强制用户设置,保持向后兼容
- `Computed: true` - 允许云端值自动同步到 state
- `Description` 说明用途和格式

**理由**:
- Classifications 是字符串数组类型([]string)
- Optional + Computed 确保升级不破坏现有配置
- 与现有字段(如 `alarm_notice_ids`)风格一致

**备选方案**:
- ❌ `TypeSet`: 虽然可去重,但 CLS API 返回的是 List,使用 List 更直接
- ❌ `Required`: 会破坏向后兼容性

### 决策 2: API 参数映射

**Create (CreateAlarm)**:
```go
if v, ok := d.GetOk("classifications"); ok {
    request.Classifications = helper.InterfacesStringsPoint(v.([]interface{}))
}
```

**Update (ModifyAlarm)**:
```go
if d.HasChange("classifications") {
    if v, ok := d.GetOk("classifications"); ok {
        request.Classifications = helper.InterfacesStringsPoint(v.([]interface{}))
    }
}
```

**Read (DescribeAlarms)**:
```go
if alarm.Classifications != nil {
    _ = d.Set("classifications", alarm.Classifications)
}
```

**理由**:
- 使用 `helper.InterfacesStringsPoint` 标准转换方法
- Update 中使用 `d.HasChange` 检测变更
- Read 中检查 nil 避免空指针

### 决策 3: 不需要 Timeouts 块

**选择**: 不添加 Timeouts 块

**理由**:
- CLS Alarm 的 Create/Update/Delete 都是同步操作,API 调用直接返回结果
- 现有资源代码中没有使用 `helper.Retry` 等待异步完成的逻辑
- 不符合"异步操作需要等待"的触发条件

## Risks / Trade-offs

### Risk 1: API 返回值为空数组 vs nil
**风险**: 如果 API 返回空数组 `[]` 而非 `nil`,可能导致 state 不必要的更新  
**缓解**: 在 Read 中检查 `len(alarm.Classifications) > 0` 再 Set,避免空数组写入

### Risk 2: 现有告警的 Classifications 为空
**风险**: 升级 Provider 后,现有告警 Read 时 `classifications` 字段为空,可能触发 diff  
**缓解**: 使用 `Computed: true`,允许空值,不强制用户设置

### Trade-off: List vs Set
**权衡**: List 保留顺序但允许重复,Set 去重但无序  
**选择**: List,因为云端 API 返回 List,且分类顺序可能有业务意义

## Migration Plan

### 部署步骤
1. 修改 `resource_tc_cls_alarm.go` 添加 schema 字段和 CRUD 逻辑
2. 添加验收测试用例验证新字段
3. 更新 `website/docs/r/cls_alarm.html.markdown` 文档
4. 运行 `make fmt && make lint` 确保代码规范
5. 提交 PR 进行代码审查

### 向后兼容性
- ✅ 现有配置无需修改(Optional 字段)
- ✅ 现有 state 自动升级(Computed 字段)
- ✅ 用户可选择性使用新字段

### 回滚策略
- 如有问题,可直接回滚代码(字段为 Optional,不影响存量资源)
- 用户已设置的 `classifications` 在回滚后不会丢失(保留在云端)

## Open Questions

- ☑️ Classifications 的具体格式和取值范围? → 查看 API 文档确认
- ☑️ 是否需要同步更新数据源? → 本次变更不涉及,可后续独立处理
