# 变更提案:修复 tencentcloud_postgresql_instance 资源 auto_renew_flag 字段的 Read 实现

## 变更类型

**Bug 修复** - 此变更修复了现有功能的缺陷,不引入新能力,因此不需要 spec deltas。

## Why

### 问题描述

当前 `tencentcloud_postgresql_instance` 资源存在以下问题:

1. **Schema 定义不完整**: `auto_renew_flag` 字段在 schema 中只标记为 `Optional: true`,但未标记为 `Computed: true`
2. **Read 方法未设置值**: `resourceTencentCloudPostgresqlInstanceRead` 函数中没有调用 `d.Set("auto_renew_flag", ...)` 来读取并设置该字段的值
3. **状态不一致**: 导致 Terraform state 与实际云上资源的自动续费状态不一致,可能引发配置漂移检测问题

### 根本原因

通过检查代码发现:

1. **API 支持**: PostgreSQL SDK 中 `DBInstance` 结构体包含 `AutoRenew *uint64` 字段 (见 `models.go:2169`)
   ```go
   // 是否自动续费:
   // <li>0:手动续费</li>
   // <li>1:自动续费</li>
   // 默认值:0
   AutoRenew *uint64 `json:"AutoRenew,omitnil,omitempty" name:"AutoRenew"`
   ```

2. **数据源已实现**: `data_source_tc_postgresql_instances.go` 中已经正确读取并设置了该字段:
   ```go
   listItem["auto_renew_flag"] = v.AutoRenew  // line 546
   ```

3. **只读实例已实现**: `resource_tc_postgresql_readonly_instance.go` 的 Read 函数中正确设置了该字段:
   ```go
   _ = d.Set("auto_renew_flag", instance.AutoRenew)  // line 407
   ```

4. **主实例缺失**: 唯独 `resource_tc_postgresql_instance.go` 的 Read 函数中遗漏了这个字段的设置

### 影响

- **配置漂移**: 用户通过控制台或其他方式修改自动续费设置后,Terraform 无法检测到变化
- **状态不准确**: Terraform state 中 `auto_renew_flag` 的值可能与实际资源状态不符
- **用户困惑**: 用户设置的 `auto_renew_flag` 在 `terraform plan` 时可能显示为未知变化
- **导入问题**: 通过 `terraform import` 导入实例时,无法获取正确的自动续费标志

## What Changes

### 代码变更

在 `tencentcloud/services/postgresql/resource_tc_postgresql_instance.go` 中进行以下修改:

#### 1. 更新 Schema 定义 (line 57-61)

**修改前:**
```go
"auto_renew_flag": {
    Type:        schema.TypeInt,
    Optional:    true,
    Description: "Auto renew flag, `1` for enabled. NOTES: Only support prepaid instance.",
},
```

**修改后:**
```go
"auto_renew_flag": {
    Type:        schema.TypeInt,
    Optional:    true,
    Computed:    true,  // 新增:支持从 API 读取
    Description: "Auto renew flag, `1` for enabled. NOTES: Only support prepaid instance.",
},
```

**理由**: 
- 设置为 `Computed: true` 表示该字段可以从 API 读取,而不仅仅由用户输入
- 与其他类似字段(如 `engine_version`, `cpu` 等)保持一致
- 支持 Terraform import 和 refresh 操作

#### 2. 在 Read 函数中设置 auto_renew_flag (在 line 903 之后插入)

**位置**: 在设置 `charge_type` 之后,设置 `security_groups` 之前

```go
if *instance.PayType == POSTGRESQL_PAYTYPE_PREPAID || *instance.PayType == COMMON_PAYTYPE_PREPAID {
    _ = d.Set("charge_type", COMMON_PAYTYPE_PREPAID)
} else {
    _ = d.Set("charge_type", COMMON_PAYTYPE_POSTPAID)
}

// 新增:读取自动续费标志
if instance.AutoRenew != nil {
    _ = d.Set("auto_renew_flag", int(*instance.AutoRenew))
}

// security groups
sg, err := postgresqlService.DescribeDBInstanceSecurityGroupsById(ctx, d.Id())
```

**实现细节**:
- 检查 `instance.AutoRenew` 是否为 `nil` 以避免空指针异常
- 将 `*uint64` 类型转换为 `int` 以匹配 schema 中的 `TypeInt`
- 对于按量计费实例,如果 `AutoRenew` 为 `nil`,则不设置该字段(保持默认行为)

### 影响范围

- **影响的规范**: 无 - 这是一个 **bug 修复**,补充了已存在字段的状态读取逻辑,不涉及新能力或行为变更,因此不需要创建或修改 spec deltas
- **影响的文件**:
  - `tencentcloud/services/postgresql/resource_tc_postgresql_instance.go` - Schema 定义和 Read 函数
- **破坏性变更**: 无
  - 现有用户配置完全兼容
  - 只是修复了状态读取的缺失,不改变任何行为
- **迁移需求**: 不适用

### 测试建议

1. **手动测试场景**:
   - 创建预付费实例并设置 `auto_renew_flag = 1`
   - 执行 `terraform refresh` 验证 state 中 `auto_renew_flag` 值正确
   - 在控制台修改自动续费设置
   - 执行 `terraform plan` 验证能检测到配置漂移
   - 测试 `terraform import` 能正确读取自动续费标志

2. **验证步骤**:
   ```bash
   # 创建实例
   terraform apply
   
   # 刷新状态
   terraform refresh
   
   # 查看状态
   terraform show | grep auto_renew_flag
   
   # 导入测试
   terraform import tencentcloud_postgresql_instance.test postgres-xxxxx
   ```

### 向后兼容性

✅ **完全向后兼容**:
- 现有的 Terraform 配置无需修改
- 不影响 Create/Update/Delete 操作
- 只增强了 Read 操作的完整性
- 用户不指定 `auto_renew_flag` 时行为不变

### 参考

- **类似实现**: 
  - `resource_tc_postgresql_readonly_instance.go:407` - 已正确实现
  - `data_source_tc_postgresql_instances.go:546` - 数据源中的实现
- **SDK 文档**: `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312/models.go:2169`
