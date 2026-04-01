# StorageType 字段添加 - 任务清单

## 任务概览

**总任务数**: 17  
**预计工时**: 4-6 小时

---

## 阶段 1: 服务层修改 (service_tencentcloud_postgresql.go)

### Task 1.1: 修改 CreatePostgresqlInstance 方法
- [ ] 在函数签名中添加 `storageType string` 参数
- [ ] 在 SDK 请求中设置 `request.StorageType = &storageType`
- [ ] 更新所有调用此方法的位置传入新参数

**预计时间**: 30 分钟

### Task 1.2: 修改 DescribePostgresqlDbInstanceVersionsByFilter 方法
- [ ] 检查方法是否接受 `paramMap` 参数
- [ ] 如果 `paramMap` 中包含 `StorageType`,设置到请求中
- [ ] 验证 API 调用正确性

**预计时间**: 20 分钟

### Task 1.3: 修改 DescribePostgresqlDbInstanceClassesByFilter 方法
- [ ] 添加对 `StorageType` 参数的支持
- [ ] 从 `paramMap` 中读取并设置到 SDK 请求

**预计时间**: 20 分钟

### Task 1.4: 修改 DescribeSpecinfos 方法
- [ ] 添加 `storageType` 参数或通过其他方式传递
- [ ] 设置到 SDK 请求中

**预计时间**: 20 分钟

---

## 阶段 2: Resource 修改 (resource_tc_postgresql_instance.go)

### Task 2.1: 添加 Schema 定义
- [ ] 在 Schema map 中添加 `storage_type` 字段
- [ ] 设置类型、可选性、默认值、验证函数
- [ ] 添加详细的 Description

**位置**: 约第 133 行之后  
**预计时间**: 15 分钟

### Task 2.2: Create 阶段 - 获取参数
- [ ] 在 `resourceTencentCloudPostgresqlInstanceCreate` 函数中获取 `storage_type`
- [ ] 将其传递给 `CreatePostgresqlInstance` 方法

**位置**: 约第 560-585 行  
**预计时间**: 15 分钟

### Task 2.3: Read 阶段 - 读取属性
- [ ] 在 `DescribeDBInstanceAttribute` 调用后读取 `DBInstanceStorageType`
- [ ] 使用 `d.Set("storage_type", ...)` 设置到 state

**位置**: 约第 928-933 行  
**预计时间**: 10 分钟

### Task 2.4: 格式化代码
- [ ] 执行 `go fmt` 格式化文件

**预计时间**: 2 分钟

---

## 阶段 3: DataSource 修改 - db_instance_versions

### Task 3.1: 添加 Schema 定义
- [ ] 在 `DataSourceTencentCloudPostgresqlDbInstanceVersions` 的 Schema 中添加 `storage_type`
- [ ] 设置为 Optional 字段

**文件**: `data_source_tc_postgresql_db_instance_versions.go`  
**位置**: 约第 75 行之后  
**预计时间**: 10 分钟

### Task 3.2: 传递参数到服务层
- [ ] 在 Read 函数中获取 `storage_type`
- [ ] 传递给 `DescribePostgresqlDbInstanceVersionsByFilter` 方法

**位置**: 约第 79-98 行  
**预计时间**: 10 分钟

### Task 3.3: 格式化代码
- [ ] 执行 `go fmt`

**预计时间**: 2 分钟

---

## 阶段 4: DataSource 修改 - db_instance_classes

### Task 4.1: 添加 Schema 定义
- [ ] 添加 `storage_type` 字段到 Schema

**文件**: `data_source_tc_postgresql_db_instance_classes.go`  
**位置**: 约第 36 行之后  
**预计时间**: 10 分钟

### Task 4.2: 传递参数
- [ ] 在 `paramMap` 中添加 `StorageType`

**位置**: 约第 94-106 行  
**预计时间**: 10 分钟

### Task 4.3: 格式化代码
- [ ] 执行 `go fmt`

**预计时间**: 2 分钟

---

## 阶段 5: DataSource 修改 - specinfos

### Task 5.1: 添加 Schema 定义
- [ ] 添加 `storage_type` 字段

**文件**: `data_source_tc_postgresql_specinfos.go`  
**位置**: 约第 20 行之后  
**预计时间**: 10 分钟

### Task 5.2: 传递参数
- [ ] 获取 `storage_type` 并传递给 `DescribeSpecinfos`

**位置**: 约第 89-90 行  
**预计时间**: 10 分钟

### Task 5.3: 格式化代码
- [ ] 执行 `go fmt`

**预计时间**: 2 分钟

---

## 阶段 6: DataSource 修改 - instances

### Task 6.1: 添加 Schema 定义 (Computed 字段)
- [ ] 在 `instance_list` 的 Elem.Schema 中添加 `storage_type` (Computed)
- [ ] 在 `db_instance_set` 的 Elem.Schema 中添加 `storage_type` (Computed)

**文件**: `data_source_tc_postgresql_instances.go`  
**位置**: 约第 153 行和 493 行  
**预计时间**: 15 分钟

### Task 6.2: 读取响应数据
- [ ] 在数据处理逻辑中读取存储类型字段
- [ ] 设置到 `listItem["storage_type"]` 和 `dBInstanceSetMap["storage_type"]`

**位置**: 约第 542-605 行和 609-896 行  
**预计时间**: 15 分钟

### Task 6.3: 确认 API 字段名
- [ ] 检查 `DescribeDBInstances` API 返回值中存储类型的字段名
- [ ] 如果不存在,考虑从 `DescribeDBInstanceAttribute` 中补充获取

**预计时间**: 20 分钟

### Task 6.4: 格式化代码
- [ ] 执行 `go fmt`

**预计时间**: 2 分钟

---

## 阶段 7: 测试与验证

### Task 7.1: 编译检查
- [ ] 执行 `go build` 确保无编译错误
- [ ] 检查所有 import 语句正确

**预计时间**: 10 分钟

### Task 7.2: Linter 检查
- [ ] 运行 linter 工具检查代码质量
- [ ] 修复所有 linter 错误和警告

**预计时间**: 15 分钟

### Task 7.3: 本地测试 (可选)
- [ ] 创建测试配置文件
- [ ] 测试实例创建 (使用不同存储类型)
- [ ] 测试数据源查询

**预计时间**: 60 分钟 (如果环境允许)

---

## 阶段 8: 文档更新

### Task 8.1: 更新 IMPLEMENTATION.md
- [ ] 记录所有修改点
- [ ] 添加 API 字段映射表
- [ ] 记录测试结果

**预计时间**: 30 分钟

### Task 8.2: 创建示例配置
- [ ] 编写 resource 使用示例
- [ ] 编写 data source 使用示例

**预计时间**: 20 分钟

---

## 优先级

### 高优先级 (必须完成)
- Task 1.1 - 1.4: 服务层修改
- Task 2.1 - 2.3: Resource Create/Read 修改
- Task 3.1 - 6.4: 所有 DataSource 修改
- Task 7.1 - 7.2: 编译和 Linter 检查

### 中优先级 (建议完成)
- Task 7.3: 本地测试
- Task 8.1: 实施文档更新

### 低优先级 (可选)
- Task 8.2: 示例配置

---

## 检查点

### ✅ 第一阶段完成标志
- 所有服务层方法已更新
- Resource 和所有 DataSource 文件已修改
- 所有文件已执行 `go fmt`

### ✅ 第二阶段完成标志
- 代码通过编译
- 无 linter 错误

### ✅ 最终完成标志
- 所有任务完成
- 文档已更新
- 代码已提交

---

## 注意事项

1. **每修改完一个文件,立即执行 `go fmt`**
2. **确认 API 字段名称与文档一致**
3. **保持向后兼容性 - 所有新字段为 Optional/Computed**
4. **检查 `DescribeDBInstances` API 是否返回存储类型字段**
