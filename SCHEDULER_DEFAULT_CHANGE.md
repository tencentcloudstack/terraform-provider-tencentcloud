# Scheduler Default Value Change - Compatibility Guide

## 变更概述

修改 `tencentcloud_clb_listener_rule` 资源中 `scheduler` 字段的默认值逻辑,使其仅在 `target_type` 不为 `TARGETGROUP-V2` 时应用默认值 `WRR`。

## 变更详情

### Schema 变更
- **移除**: 静态 `Default: CLB_LISTENER_SCHEDULER_WRR`
- **添加**: `Computed: true` 标记
- **更新**: Description 说明默认值的条件

### 代码逻辑变更

#### Create 函数 (resourceTencentCloudClbListenerRuleCreate)
```go
// 之前: 始终有默认值 WRR (在 schema 级别)
// 现在: 条件性应用默认值
scheduler := ""
targetType := d.Get("target_type").(string)
if v, ok := d.GetOk("scheduler"); ok {
    // 用户显式设置了 scheduler
    scheduler = v.(string)
    if targetType != CLB_TARGET_TYPE_TARGETGROUP_V2 {
        rule.Scheduler = helper.String(scheduler)
    }
} else {
    // 用户未设置 scheduler,应用条件默认值
    if targetType != CLB_TARGET_TYPE_TARGETGROUP_V2 {
        scheduler = CLB_LISTENER_SCHEDULER_WRR
        rule.Scheduler = helper.String(scheduler)
    }
    // 当 target_type 为 TARGETGROUP-V2 时,不设置 scheduler
}
```

#### Update 函数 (resourceTencentCloudClbListenerRuleUpdate)
```go
// 只有当 target_type 不是 TARGETGROUP-V2 时才允许修改 scheduler
if d.HasChange("scheduler") {
    targetType := d.Get("target_type").(string)
    scheduler = d.Get("scheduler").(string)
    if targetType != CLB_TARGET_TYPE_TARGETGROUP_V2 {
        request.Scheduler = helper.String(scheduler)
    }
}
```

## 兼容性分析

### 1. 现有资源 (已经创建的 listener rules)

#### 场景 1.1: target_type = NODE/TARGETGROUP, 未显式设置 scheduler
- **状态变更前**: scheduler = "WRR" (schema default)
- **状态变更后**: scheduler = "WRR" (computed from API)
- **兼容性**: ✅ 完全兼容,无实际变化

#### 场景 1.2: target_type = NODE/TARGETGROUP, 显式设置 scheduler = "WRR"
- **状态变更前**: scheduler = "WRR" (用户设置)
- **状态变更后**: scheduler = "WRR" (保持不变)
- **兼容性**: ✅ 完全兼容

#### 场景 1.3: target_type = NODE/TARGETGROUP, 显式设置 scheduler = "LEAST_CONN"
- **状态变更前**: scheduler = "LEAST_CONN" (用户设置)
- **状态变更后**: scheduler = "LEAST_CONN" (保持不变)
- **兼容性**: ✅ 完全兼容

#### 场景 1.4: target_type = TARGETGROUP-V2, 未显式设置 scheduler
- **状态变更前**: scheduler = "WRR" (schema default,但实际 API 创建时被忽略)
- **状态变更后**: scheduler = "" 或从 API 读取的实际值
- **兼容性**: ⚠️ 可能显示 diff,但这是修正行为
- **说明**: 之前的行为是错误的(为 TARGETGROUP-V2 设置了默认值),现在修正为正确行为

#### 场景 1.5: target_type = TARGETGROUP-V2, 显式设置了 scheduler
- **状态变更前**: scheduler 值被忽略(API 不接受)
- **状态变更后**: scheduler 值仍被忽略
- **兼容性**: ✅ 行为一致

### 2. 新创建的资源

#### 场景 2.1: target_type = NODE, 不设置 scheduler
```hcl
resource "tencentcloud_clb_listener_rule" "example" {
  clb_id      = "lb-xxx"
  listener_id = "lbl-xxx"
  domain      = "example.com"
  url         = "/"
  target_type = "NODE"
  # scheduler 未设置
}
```
- **行为**: scheduler 自动设置为 "WRR"
- **兼容性**: ✅ 与之前行为完全一致

#### 场景 2.2: target_type = TARGETGROUP-V2, 不设置 scheduler
```hcl
resource "tencentcloud_clb_listener_rule" "example" {
  clb_id      = "lb-xxx"
  listener_id = "lbl-xxx"
  domain      = "example.com"
  url         = "/"
  target_type = "TARGETGROUP-V2"
  # scheduler 未设置
}
```
- **行为**: scheduler 不设置默认值,由 API 决定
- **兼容性**: ✅ 修正了之前错误的默认值行为

#### 场景 2.3: target_type = TARGETGROUP-V2, 显式设置 scheduler
```hcl
resource "tencentcloud_clb_listener_rule" "example" {
  clb_id      = "lb-xxx"
  listener_id = "lbl-xxx"
  domain      = "example.com"
  url         = "/"
  target_type = "TARGETGROUP-V2"
  scheduler   = "WRR"  # 显式设置会被忽略
}
```
- **行为**: scheduler 值在 Create/Update 时被忽略(不发送给 API)
- **兼容性**: ✅ 与之前行为一致

### 3. 状态迁移 (terraform refresh)

用户执行 `terraform refresh` 或 `terraform plan` 时:

1. **target_type != TARGETGROUP-V2**: 
   - 无变化,scheduler 值保持原样或默认为 "WRR"
   
2. **target_type = TARGETGROUP-V2**:
   - 如果之前 state 中有 scheduler = "WRR" (旧默认值),可能会显示 diff
   - 但这是正确的修正行为,因为 TARGETGROUP-V2 不应该有 scheduler 默认值

## 测试建议

### 单元测试
1. 测试 target_type = NODE 时,scheduler 默认为 WRR
2. 测试 target_type = TARGETGROUP 时,scheduler 默认为 WRR
3. 测试 target_type = TARGETGROUP-V2 时,scheduler 无默认值
4. 测试显式设置 scheduler 时的行为

### 集成测试
1. 创建 target_type = NODE 的规则,验证 scheduler 默认值
2. 创建 target_type = TARGETGROUP-V2 的规则,验证 scheduler 不被设置
3. 更新现有规则,验证兼容性
4. 导入现有规则,验证状态一致性

### 回归测试
1. 使用旧配置文件执行 `terraform plan`,应无不期望的变更
2. 升级 provider 后执行 `terraform refresh`,验证状态同步
3. 测试所有 scheduler 相关的边界情况

## 潜在影响

### 低风险场景
- target_type = NODE 或 TARGETGROUP 的现有配置
- 所有显式设置了 scheduler 的配置

### 需要注意的场景
- target_type = TARGETGROUP-V2 且依赖默认 scheduler 值的配置
  - **影响**: 之前错误地设置了默认值,现在修正后可能显示 diff
  - **处理**: 这是修正行为,如果 API 返回了 scheduler 值,会自动同步到 state

## 回滚方案

如果需要回滚此变更:
```go
// 恢复 schema 中的 Default 属性
"scheduler": {
    Type:         schema.TypeString,
    Optional:     true,
    Default:      CLB_LISTENER_SCHEDULER_WRR,  // 恢复
    // Computed:     true,  // 移除
    ValidateFunc: tccommon.ValidateAllowedStringValue(CLB_LISTENER_SCHEDULER),
    Description:  "...",
},

// 恢复 Create 函数中的原始逻辑
scheduler := ""
if v, ok := d.GetOk("scheduler"); ok {
    scheduler = v.(string)
    if *rule.TargetType != CLB_TARGET_TYPE_TARGETGROUP_V2 {
        rule.Scheduler = helper.String(scheduler)
    }
}
// 移除 else 分支

// 恢复 Update 函数中的原始逻辑
if d.HasChange("scheduler") {
    scheduler = d.Get("scheduler").(string)
    request.Scheduler = helper.String(scheduler)
}
```

## 总结

此变更是一个**兼容性较好的修正**:
- ✅ 不影响 target_type = NODE/TARGETGROUP 的现有和新建配置
- ✅ 修正了 target_type = TARGETGROUP-V2 的错误默认值行为
- ✅ 所有显式配置保持不变
- ⚠️ 对于使用 TARGETGROUP-V2 且依赖默认值的配置,可能会出现状态 diff,但这是正确的修正

建议在发布前进行充分的集成测试,特别是针对 TARGETGROUP-V2 类型的监听规则。
