# Tasks: Add Isolation Wait Mechanism for SQL Server General Instance Deletion

> **Note**: 本变更同时应用于两个资源：
> - `tencentcloud_sqlserver_general_cloud_instance` ✅ 已完成
> - `tencentcloud_sqlserver_instance` ✅ 已完成

## Task 1: 修改 Delete 函数 - 添加等待隔离逻辑

**文件**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_general_cloud_instance.go`

**函数**: `resourceTencentCloudSqlserverGeneralCloudInstanceDelete`（第 652-672 行）

### 实施步骤

在 `TerminateSqlserverInstanceById` 调用之后，`DeleteSqlserverInstanceById` 调用之前，插入等待逻辑：

#### 1.1 定位插入位置

找到第 663-665 行：

```go
if err := service.TerminateSqlserverInstanceById(ctx, instanceId); err != nil {
    return err
}

// 在这里插入等待逻辑 ← 插入位置

if err := service.DeleteSqlserverInstanceById(ctx, instanceId); err != nil {
    return err
}
```

#### 1.2 插入等待代码

在上述位置插入以下代码：

```go
	// Wait for instance to be isolated (status = 4)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		// Query instance status using DescribeDBInstances API
		instances, err := service.DescribeSqlserverInstances(ctx, instanceId, "", -1, "", "", 0)
		if err != nil {
			return tccommon.RetryError(err)
		}

		// Check if instance exists
		if len(instances) == 0 {
			return resource.NonRetryableError(fmt.Errorf("instance %s not found", instanceId))
		}

		instance := instances[0]

		// Check instance status
		if instance.Status != nil {
			status := *instance.Status
			log.Printf("[DEBUG]%s instance %s current status: %d", logId, instanceId, status)

			if status == 4 {
				// Instance is isolated, ready to delete
				log.Printf("[INFO]%s instance %s is isolated (status=4), ready to delete", logId, instanceId)
				return nil
			}

			// Continue waiting for other statuses
			return resource.RetryableError(fmt.Errorf("waiting for instance %s to be isolated, current status: %d", instanceId, status))
		}

		return resource.RetryableError(fmt.Errorf("instance %s status is nil", instanceId))
	})

	if err != nil {
		log.Printf("[CRITAL]%s wait for instance %s isolation failed, reason: %+v", logId, instanceId, err)
		return err
	}
```

### 完整的 Delete 函数（修改后）

```go
func resourceTencentCloudSqlserverGeneralCloudInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_general_cloud_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	if err := service.TerminateSqlserverInstanceById(ctx, instanceId); err != nil {
		return err
	}

	// Wait for instance to be isolated (status = 4)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		// Query instance status using DescribeDBInstances API
		instances, err := service.DescribeSqlserverInstances(ctx, instanceId, "", -1, "", "", 0)
		if err != nil {
			return tccommon.RetryError(err)
		}

		// Check if instance exists
		if len(instances) == 0 {
			return resource.NonRetryableError(fmt.Errorf("instance %s not found", instanceId))
		}

		instance := instances[0]

		// Check instance status
		if instance.Status != nil {
			status := *instance.Status
			log.Printf("[DEBUG]%s instance %s current status: %d", logId, instanceId, status)

			if status == 4 {
				// Instance is isolated, ready to delete
				log.Printf("[INFO]%s instance %s is isolated (status=4), ready to delete", logId, instanceId)
				return nil
			}

			// Continue waiting for other statuses
			return resource.RetryableError(fmt.Errorf("waiting for instance %s to be isolated, current status: %d", instanceId, status))
		}

		return resource.RetryableError(fmt.Errorf("instance %s status is nil", instanceId))
	})

	if err != nil {
		log.Printf("[CRITAL]%s wait for instance %s isolation failed, reason: %+v", logId, instanceId, err)
		return err
	}

	if err := service.DeleteSqlserverInstanceById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
```

### 验收标准

- [x] 在 `TerminateSqlserverInstanceById` 之后添加了等待逻辑
- [x] 使用 `resource.Retry` 和 `tccommon.ReadRetryTimeout`
- [x] 调用 `DescribeSqlserverInstances` 查询实例状态
- [x] 判断 `Status == 4` 后继续执行删除
- [x] 添加了详细的日志记录（DEBUG、INFO、CRITAL 级别）
- [x] 正确处理错误情况（实例不存在、Status 为 nil）
- [x] 代码逻辑清晰，注释完整

### 1.3 应用到 `tencentcloud_sqlserver_instance` 资源

**文件**: `tencentcloud/services/sqlserver/resource_tc_sqlserver_instance.go`

**函数**: `resourceTencentCLoudSqlserverInstanceDelete`（第 733-795 行）

**插入位置**: 在第 774-776 行之后（`TerminateSqlserverInstance` 调用成功后）

**插入代码**: 与上述相同的等待隔离逻辑（第 778-814 行）

**验收标准**:
- [x] 在 `TerminateSqlserverInstance` 之后添加了等待逻辑
- [x] 使用 `resource.Retry` 和 `tccommon.ReadRetryTimeout*10`
- [x] 调用 `DescribeSqlserverRestartDBInstanceById` 查询实例状态
- [x] 判断 `Status == 4` 后继续执行删除
- [x] 添加了详细的日志记录
- [x] 正确处理错误情况
- [x] 代码逻辑与 `tencentcloud_sqlserver_general_cloud_instance` 一致

---

## Task 2: 代码格式化和编译验证

### 2.1 格式化代码

```bash
cd /Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud
go fmt tencentcloud/services/sqlserver/resource_tc_sqlserver_general_cloud_instance.go
go fmt tencentcloud/services/sqlserver/resource_tc_sqlserver_instance.go
```

### 2.2 编译验证

```bash
cd /Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud
go build ./tencentcloud/services/sqlserver/
```

### 验收标准

- [x] 两个文件都通过 `go fmt` 格式化，代码格式规范
- [x] `go build` 编译成功，无语法错误
- [x] 无 linter 错误或警告（可接受已有的 deprecated 警告）

---

## Task 3: 代码审查和验证

### 3.1 代码审查清单

检查以下内容：

- [ ] Import 语句正确（已有 `resource` 和 `fmt` 包）
- [ ] 变量作用域正确（`logId`、`ctx`、`instanceId` 在闭包中可访问）
- [ ] 错误处理正确：
  - [ ] 使用 `tccommon.RetryError()` 包装 API 错误
  - [ ] 使用 `resource.NonRetryableError()` 处理不可重试错误
  - [ ] 使用 `resource.RetryableError()` 处理需要重试的情况
- [ ] 日志记录完整：
  - [ ] DEBUG 日志记录当前状态
  - [ ] INFO 日志记录隔离成功
  - [ ] CRITAL 日志记录失败原因
- [ ] 注释清晰，解释关键逻辑

### 3.2 逻辑验证

确认逻辑正确：

- [ ] Status = 4 时返回 `nil`（成功）
- [ ] Status != 4 时返回 `RetryableError`（继续等待）
- [ ] 实例不存在时返回 `NonRetryableError`（不可重试）
- [ ] Status 为 nil 时返回 `RetryableError`（继续等待）
- [ ] 等待超时时返回错误并记录日志

### 3.3 与现有代码风格一致

- [ ] 变量命名符合 Go 规范
- [ ] 缩进和换行符合项目风格
- [ ] 错误处理模式与其他资源一致
- [ ] 日志格式与项目一致

---

## 验收清单

### 功能完整性

- [x] 删除流程中正确添加了等待隔离的逻辑
- [x] 使用 `DescribeDBInstances` API 查询实例状态
- [x] 等待 Status = 4 后继续执行删除
- [x] 超时机制正常工作

### 代码质量

- [x] 代码通过 `go fmt` 格式化
- [x] 代码编译通过，无语法错误
- [x] 日志记录完整清晰
- [x] 错误处理健壮
- [x] 注释完整，易于理解

### 兼容性

- [x] 不修改 Service 层方法
- [x] 不修改 API 调用方式
- [x] 不影响现有功能
- [x] 向后兼容

### 文档

- [x] 代码注释清晰
- [x] 日志信息详细
- [x] 符合项目规范

---

## 实施顺序

1. **Task 1**: 修改 Delete 函数（核心任务）
2. **Task 2**: 代码格式化和编译验证（质量保证）
3. **Task 3**: 代码审查和验证（最终检查）

---

## 注意事项

1. **每次修改后立即执行 `go fmt`**，确保代码格式规范
2. **仔细检查变量作用域**，确保在闭包中正确使用 `logId`、`ctx`、`instanceId`
3. **Status 值含义**：
   - 4 = 已隔离（目标状态）
   - 其他值需要继续等待
4. **超时时间**：使用 `tccommon.ReadRetryTimeout`（默认 5 分钟）
5. **日志级别**：
   - DEBUG：用于调试信息（当前状态）
   - INFO：用于重要信息（隔离成功）
   - CRITAL：用于错误信息（等待失败）
6. **错误处理**：
   - API 错误：使用 `tccommon.RetryError()` 包装
   - 不可重试错误：使用 `resource.NonRetryableError()`
   - 可重试错误：使用 `resource.RetryableError()`

---

## 测试建议

### 手动测试步骤

1. 创建 SQL Server General Instance
   ```hcl
   resource "tencentcloud_sqlserver_general_cloud_instance" "test" {
     name         = "test-instance"
     zone         = "ap-guangzhou-6"
     memory       = 4
     storage      = 100
     cpu          = 2
     machine_type = "CLOUD_HSSD"
     # ... 其他配置
   }
   ```

2. 执行 `terraform destroy`

3. 观察日志输出，确认：
   ```
   [INFO] instance mssql-xxxxx is isolated (status=4), ready to delete
   [DEBUG] api[DeleteDBInstance] success
   ```

### 预期行为

- ✅ `TerminateSqlserverInstanceById` 成功
- ✅ 等待实例状态变为 4（已隔离）
- ✅ `DeleteSqlserverInstanceById` 成功
- ✅ 删除操作完成

### 异常场景测试

1. **超时场景**：如果实例一直未隔离（理论场景）
   - 预期：5 分钟后返回超时错误
   
2. **实例不存在**：如果隔离后实例意外消失
   - 预期：返回 "instance not found" 错误

3. **API 错误**：如果查询 API 临时失败
   - 预期：重试直到成功或超时

---

## 参考信息

### SQL Server 实例状态枚举

| Status | 状态描述 | 说明 |
|--------|---------|------|
| 1 | 申请中 | 实例创建中 |
| 2 | 运行中 | 正常运行 |
| 3 | 受限运行中 | 主备切换中 |
| **4** | **已隔离** | **目标状态** ✅ |
| 5 | 回收中 | 准备物理销毁 |
| 6 | 已回收 | 已物理销毁 |
| 7 | 任务执行中 | 备份、回档等操作 |
| 8 | 已下线 | 实例已下线 |
| 9 | 实例扩容中 | 扩容操作中 |
| 10 | 实例迁移中 | 迁移操作中 |
| 11 | 只读 | 只读实例 |
| 12 | 重启中 | 重启操作中 |

### 相关 API

- **TerminateDBInstance**: 隔离实例（将实例放入回收站）
- **DescribeDBInstances**: 查询实例列表和状态
- **DeleteDBInstance**: 删除实例（从回收站彻底删除）

### 相关代码位置

- Delete 函数：`resource_tc_sqlserver_general_cloud_instance.go` 第 652-672 行
- Service 方法：`service_tencentcloud_sqlserver.go` 第 302-358 行（DescribeSqlserverInstances）
- Service 方法：`service_tencentcloud_sqlserver.go` 第 2197-2218 行（TerminateSqlserverInstanceById）
- Service 方法：`service_tencentcloud_sqlserver.go` 第 2220-2241 行（DeleteSqlserverInstanceById）
