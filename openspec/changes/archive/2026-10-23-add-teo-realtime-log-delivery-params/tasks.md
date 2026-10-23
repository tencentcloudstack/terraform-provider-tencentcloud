## 1. Schema 调整

- [x] 1.1 在 `tencentcloud/services/teo/resource_tc_teo_realtime_log_delivery.go` 的 `ResourceTencentCloudTeoRealtimeLogDelivery` schema 中新增顶层块 `aws_s3`（TypeList, Optional, Computed, MaxItems:1），Elem.Schema 包含子字段：`endpoint`(Required, TypeString)、`bucket_name`(Required, TypeString)、`path`(Required, TypeString)、`role_arn`(Required, TypeString)、`compress_type`(Optional, TypeString)，参照现有 `s3` 块结构
- [x] 1.2 在 schema 中新增顶层块 `custom_expression_fields`（TypeList, Optional, Computed），Elem.Schema（schema.Resource）包含子字段：`name`(Required, TypeString)、`expression`(Required, TypeString)、`enabled`(Optional, TypeBool)，参照现有 `custom_fields` 块结构
- [x] 1.3 在已有 `log_format` 块的 Elem.Schema 中追加子字段 `record_template`（Optional, TypeString），Description 注明仅当 `format_type = template` 时生效

## 2. Create 函数扩展

- [x] 2.1 在 `resourceTencentCloudTeoRealtimeLogDeliveryCreate` 中新增 `custom_expression_fields` 填充逻辑：遍历 `d.GetOk("custom_expression_fields")`，构造 `[]*teo.CustomExpressionField`，填充 Name/Expression/Enabled，赋值给 `request.CustomExpressionFields`
- [x] 2.2 在 Create 中新增 `aws_s3` 填充逻辑：使用 `helper.InterfacesHeadMap(d, "aws_s3")`，构造 `teo.AWSS3{}`，填充 Endpoint/BucketName/Path/RoleARN/CompressType，赋值给 `request.AWSS3`
- [x] 2.3 在 Create 的 `log_format` 填充分支内追加：`if v, ok := logFormatMap["record_template"]; ok { logFormat.RecordTemplate = helper.String(v.(string)) }`
- [x] 2.4 复查 Create 在调用 API 后检查 `response.Response.TaskId` 非空再取值并 `d.SetId(zoneId#taskId)`，并在拿到 TaskId 后、轮询等待前完成 SetId

## 3. Read 函数扩展

- [x] 3.1 在 `resourceTencentCloudTeoRealtimeLogDeliveryRead` 中新增 `custom_expression_fields` 回填逻辑：遍历 `respData.CustomExpressionFields`（nil 检查），逐元素构造 map（Name/Expression/Enabled），组装 list 后 `_ = d.Set("custom_expression_fields", customExpressionFieldsList)`
- [x] 3.2 在 Read 中新增 `aws_s3` 回填逻辑：当 `respData.AWSS3 != nil` 时构造 map（Endpoint/BucketName/Path/RoleARN/CompressType 均做 nil 检查），`_ = d.Set("aws_s3", []interface{}{awsS3Map})`
- [x] 3.3 在 Read 的 `log_format` 回填分支内追加：`if respData.LogFormat.RecordTemplate != nil { logFormatMap["record_template"] = respData.LogFormat.RecordTemplate }`

## 4. Update 函数扩展

- [x] 4.1 在 `resourceTencentCloudTeoRealtimeLogDeliveryUpdate` 的 `mutableArgs` 数组中追加 `"custom_expression_fields"` 和 `"aws_s3"`（`log_format` 已在其中，record_template 子字段变更自动触发）
- [x] 4.2 在 Update 构造 `ModifyRealtimeLogDeliveryTaskRequest` 的逻辑中，新增 `custom_expression_fields` 填充（与 Create 一致，构造 `[]*teo.CustomExpressionField` 赋给 `request.CustomExpressionFields`）
- [x] 4.3 在 Update 中新增 `aws_s3` 填充（与 Create 一致，构造 `teo.AWSS3{}` 赋给 `request.AWSS3`）
- [x] 4.4 在 Update 的 `log_format` 填充分支内追加 `record_template` 填充（与 Create 一致）
- [x] 4.5 复查 `immutableArgs` 保持 `["task_type", "log_type", "area", "cls"]` 不变，Update 末尾仍调用 Read 回写最新状态

## 5. 单元测试

- [x] 5.1 在 `tencentcloud/services/teo/resource_tc_teo_realtime_log_delivery_test.go` 中使用 gomonkey mock 云 API，新增测试用例覆盖含 `custom_expression_fields`、`aws_s3`、`log_format.record_template` 的 Create 场景，断言 request 中对应字段被正确填充
- [x] 5.2 新增测试用例覆盖含新字段的 Read 场景，mock `DescribeRealtimeLogDeliveryTasks` 返回含 CustomExpressionFields/AWSS3/RecordTemplate 的响应，断言 state 正确回填
- [x] 5.3 新增测试用例覆盖含新字段的 Update 场景，mock `ModifyRealtimeLogDeliveryTask`，断言新字段变更触发对应 request 填充
- [x] 5.4 保证已有测试用例继续通过（不修改其行为）

## 6. 文档同步

- [x] 6.1 在 `tencentcloud/services/teo/resource_tc_teo_realtime_log_delivery.md` 的 Example Usage 中新增 `aws_s3` 块示例、`custom_expression_fields` 块示例，并在 `log_format` 块示例中补充 `record_template` 字段
- [x] 6.2 在收尾阶段执行 `make doc`，根据 provider 规范重新生成 `website/docs/` 下的 markdown 文档（禁止手改 website 目录）

## 7. 验证

- [x] 7.1 检查所有函数返回的 error，若函数必定不出错则用 `_ = func()` 将 err 赋值给 `_`
- [x] 7.2 确认新增参数在 Create/Update 接口中存在、Read 出参中存在（CRUD 参数一致性校验，参照 vendor 中 CreateRealtimeLogDeliveryTaskRequest/ModifyRealtimeLogDeliveryTaskRequest/RealtimeLogDeliveryTask 字段）
- [x] 7.3 代码正确性检查：确认生成的代码在当前环境下可正确构建执行（禁止在本阶段执行 go build/go vet/golint）