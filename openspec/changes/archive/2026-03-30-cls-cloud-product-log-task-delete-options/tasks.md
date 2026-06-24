# Tasks: Add Deletion Options for CLS Cloud Product Log Task

## Task 1: 更新 Schema 定义

**文件**: `tencentcloud/services/cls/resource_tc_cls_cloud_product_log_task_v2.go`

在现有 `force_delete` 字段（第 101-105 行）之后添加两个新字段：

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

**验收标准**:
- [x] 两个新字段添加在 `force_delete` 之后
- [x] 字段类型为 `TypeBool`，`Optional: true`，`Default: false`
- [x] Description 清晰说明字段用途和约束条件
- [x] 执行 `go fmt tencentcloud/services/cls/resource_tc_cls_cloud_product_log_task_v2.go`

---

## Task 2: 更新 Read 模块

**文件**: `tencentcloud/services/cls/resource_tc_cls_cloud_product_log_task_v2.go`

**函数**: `resourceTencentCloudClsCloudProductLogTaskV2Read`（第 200-280 行）

在现有的 `force_delete` 读取逻辑（第 273-277 行）之后添加新字段的读取：

```go
// 现有代码保持不变
if v, ok := d.GetOkExists("force_delete"); ok {
	deleteForce = v.(bool)
}
_ = d.Set("force_delete", deleteForce)

// 新增：读取 is_delete_topic
var (
	isDeleteTopic  bool
	isDeleteLogset bool
)

if v, ok := d.GetOkExists("is_delete_topic"); ok {
	isDeleteTopic = v.(bool)
}
_ = d.Set("is_delete_topic", isDeleteTopic)

// 新增：读取 is_delete_logset
if v, ok := d.GetOkExists("is_delete_logset"); ok {
	isDeleteLogset = v.(bool)
}
_ = d.Set("is_delete_logset", isDeleteLogset)
```

**注意**:
- 在函数开头的变量声明部分（第 208 行）也需要添加 `isDeleteTopic` 和 `isDeleteLogset` 变量声明（如果需要）
- 或者直接在使用处声明（如上面代码所示）

**验收标准**:
- [x] 新字段的读取逻辑添加在 `force_delete` 之后
- [x] 使用 `GetOkExists` 读取字段值
- [x] 使用 `d.Set` 设置字段值
- [x] 执行 `go fmt tencentcloud/services/cls/resource_tc_cls_cloud_product_log_task_v2.go`

---

## Task 3: 重构 Delete 模块

**文件**: `tencentcloud/services/cls/resource_tc_cls_cloud_product_log_task_v2.go`

**函数**: `resourceTencentCloudClsCloudProductLogTaskV2Delete`（第 347-447 行）

### 3.1 更新变量声明

在函数开头的变量声明（第 351-356 行）中添加新变量：

```go
var (
	logId          = tccommon.GetLogId(tccommon.ContextNil)
	ctx            = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	request        = clsv20201016.NewDeleteCloudProductLogCollectionRequest()
	deleteForce    bool
	isDeleteTopic  bool  // 新增
	isDeleteLogset bool  // 新增
)
```

### 3.2 添加新字段读取逻辑

在 ID 解析代码之后（第 358-371 行），API 调用之前，添加配置读取：

```go
// 解析 ID（现有代码保持不变）
idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
if len(idSplit) != 4 {
	return fmt.Errorf("id is broken,%s", d.Id())
}

instanceId := idSplit[0]
assumerName := idSplit[1]
logType := idSplit[2]
cloudProductRegion := idSplit[3]

// 新增：读取删除选项配置
if v, ok := d.GetOkExists("force_delete"); ok {
	deleteForce = v.(bool)
}

if !deleteForce {
	// 只有 force_delete 为 false 时，才读取 is_delete_* 字段
	if v, ok := d.GetOkExists("is_delete_topic"); ok {
		isDeleteTopic = v.(bool)
	}
	
	if v, ok := d.GetOkExists("is_delete_logset"); ok {
		isDeleteLogset = v.(bool)
	}
}

// 构建 API 请求（现有代码）
request.InstanceId = helper.String(instanceId)
request.AssumerName = helper.String(assumerName)
request.LogType = helper.String(logType)
request.CloudProductRegion = helper.String(cloudProductRegion)

// 新增：设置删除选项
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

### 3.3 删除手动删除代码

**删除**第 395-444 行的代码块：

```go
// 删除以下整个代码块
if v, ok := d.GetOkExists("force_delete"); ok {
	deleteForce = v.(bool)
}

if deleteForce {
	var (
		request1 = clsv20201016.NewDeleteTopicRequest()
		request2 = clsv20201016.NewDeleteLogsetRequest()
	)

	if v, ok := d.GetOk("topic_id"); ok {
		request1.TopicId = helper.String(v.(string))
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsV20201016Client().DeleteTopicWithContext(ctx, request1)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request1.GetAction(), request1.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete cls cloud product log task topic failed, reason:%+v", logId, err)
		return err
	}

	if v, ok := d.GetOk("logset_id"); ok {
		request2.LogsetId = helper.String(v.(string))
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsV20201016Client().DeleteLogsetWithContext(ctx, request2)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request2.GetAction(), request2.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete cls cloud product log task logset failed, reason:%+v", logId, err)
		return err
	}
}
```

**原因**: API 的 `IsDeleteTopic` 和 `IsDeleteLogset` 参数会自动处理删除，不需要手动调用。

**验收标准**:
- [x] 变量声明中添加了 `isDeleteTopic` 和 `isDeleteLogset`
- [x] 在 DeleteCloudProductLogCollection 调用前添加了配置读取逻辑
- [x] 正确设置 `request.IsDeleteTopic` 和 `request.IsDeleteLogset`
- [x] 删除了手动调用 DeleteTopic 和 DeleteLogset 的代码块
- [x] 逻辑正确：`force_delete=true` 时忽略 `is_delete_*` 字段
- [x] 执行 `go fmt tencentcloud/services/cls/resource_tc_cls_cloud_product_log_task_v2.go`

---

## Task 4: 代码格式化和验证

### 4.1 格式化代码

```bash
cd /Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud
go fmt tencentcloud/services/cls/resource_tc_cls_cloud_product_log_task_v2.go
```

### 4.2 编译验证

```bash
cd /Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud
go build ./tencentcloud/services/cls/
```

**验收标准**:
- [x] `go fmt` 执行成功，代码格式规范
- [x] `go build` 编译成功，无语法错误
- [x] 无 linter 错误（可以有 deprecated 警告）

---

## Task 5: 文档更新（可选）

**文件**: `tencentcloud/services/cls/resource_tc_cls_cloud_product_log_task_v2.md`（如果存在）

或者更新相关的示例文档，添加新字段的使用说明：

```hcl
resource "tencentcloud_cls_cloud_product_log_task_v2" "example" {
  instance_id         = "cdb-xxxxx"
  assumer_name        = "CDB"
  log_type            = "CDB-AUDIT"
  cloud_product_region = "gz"
  cls_region          = "ap-guangzhou"
  
  # 选项 1: 使用 force_delete（向后兼容）
  # force_delete = true  # 强制删除 Topic 和 Logset
  
  # 选项 2: 精细控制（推荐）
  force_delete     = false
  is_delete_topic  = true   # 只删除 Topic
  is_delete_logset = false  # 保留 Logset
}
```

**验收标准**:
- [ ] 如果文档文件存在，添加新字段的使用示例
- [ ] 说明 `force_delete` 和 `is_delete_*` 字段的关系

---

## 验收清单

### Schema
- [ ] `is_delete_topic` 字段定义正确
- [ ] `is_delete_logset` 字段定义正确
- [ ] 字段类型、属性、描述符合要求

### Read 模块
- [ ] 正确读取 `is_delete_topic` 字段
- [ ] 正确读取 `is_delete_logset` 字段
- [ ] 使用 `GetOkExists` 和 `d.Set`

### Delete 模块
- [ ] 变量声明包含新字段
- [ ] 配置读取逻辑正确（force_delete 优先级处理）
- [ ] API 请求参数正确设置（IsDeleteTopic, IsDeleteLogset）
- [ ] 删除了手动调用 DeleteTopic/DeleteLogset 的代码
- [ ] 逻辑清晰，符合设计文档

### 代码质量
- [ ] 代码通过 `go fmt` 格式化
- [ ] 代码编译通过
- [ ] 无语法错误
- [ ] 变量命名规范（驼峰命名）

### 功能验证
- [ ] `force_delete=true` 行为正确（删除 Topic 和 Logset）
- [ ] `force_delete=false, is_delete_topic=true` 行为正确（只删除 Topic）
- [ ] `force_delete=false, is_delete_logset=true` 行为正确（只删除 Logset）
- [ ] 默认行为正确（不删除任何资源）

---

## 注意事项

1. **每次修改后都要执行 `go fmt`**，确保代码格式规范
2. **保持向后兼容**：`force_delete` 的行为不应改变
3. **逻辑优先级**：`force_delete=true` 时，忽略 `is_delete_*` 字段的值
4. **API 字段命名**：Terraform 使用下划线 `is_delete_topic`，API 使用驼峰 `IsDeleteTopic`
5. **删除旧代码**：移除手动调用 DeleteTopic 和 DeleteLogset 的逻辑，因为 API 会自动处理
6. **测试建议**：实际测试时注意 Logset 可能因为有其他 Topic 而无法删除，这是正常行为
