## 1. 资源实现

- [x] 1.1 创建 `tencentcloud/services/crs/resource_tc_redis_instance_password_policy_config.go`，定义 `ResourceTencentCloudRedisInstancePasswordPolicyConfig()` schema 函数
- [x] 1.2 实现 schema：`instance_id`(String, Required, ForceNew)、`enabled`(Bool, Required)、`min_letter_count`(Int, Optional, Computed)、`min_digit_count`(Int, Optional, Computed)、`min_special_count`(Int, Optional, Computed)、`min_length`(Int, Optional, Computed)，字段平铺在顶层
- [x] 1.3 添加 Import 支持（`schema.ImportStatePassthrough`）
- [x] 1.4 实现 `resourceTencentCloudRedisInstancePasswordPolicyConfigCreate`：设置 `d.SetId(instanceId)` 后调用 Update
- [x] 1.5 实现 `resourceTencentCloudRedisInstancePasswordPolicyConfigRead`：调用 `DescribeInstancePasswordPolicy`，设置字段前检查 nil；空响应先打印 `log.Printf("[CRUD] redis_instance_password_policy_config id=%s", d.Id())` 再 `d.SetId("")`
- [x] 1.6 实现 `resourceTencentCloudRedisInstancePasswordPolicyConfigUpdate`：调用 `ModifyInstancePasswordPolicy`，用 `resource.Retry(tccommon.WriteRetryTimeout)` + `tccommon.RetryError` 包装错误，retry 块外调用 Read 刷新
- [x] 1.7 实现 `resourceTencentCloudRedisInstancePasswordPolicyConfigDelete`：no-op，返回 nil

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 的 ResourcesMap 中注册 `tencentcloud_redis_instance_password_policy_config` → `crs.ResourceTencentCloudRedisInstancePasswordPolicyConfig()`
- [x] 2.2 在 `tencentcloud/provider.md` 中添加 `tencentcloud_redis_instance_password_policy_config` 资源声明

## 3. 文档

- [x] 3.1 创建 `tencentcloud/services/crs/resource_tc_redis_instance_password_policy_config.md`，包含一句话描述（带云产品名称 redis）、Example Usage、Import 示例（说明使用 instance_id 导入）
- [x] 3.2 运行 `make doc` 生成 `website/docs/r/redis_instance_password_policy_config.html.markdown`（由收尾阶段 tfpacer-finalize skill 执行）

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/crs/resource_tc_redis_instance_password_policy_config_test.go`，使用 gomonkey mock `DescribeInstancePasswordPolicy` 和 `ModifyInstancePasswordPolicy` 接口
- [x] 4.2 编写 Create 业务逻辑测试（mock ModifyInstancePasswordPolicy 成功 + Read 回填）
- [x] 4.3 编写 Update 业务逻辑测试（mock ModifyInstancePasswordPolicy 成功）
- [x] 4.4 编写 Read 业务逻辑测试（mock DescribeInstancePasswordPolicy 返回完整策略 + 返回空的场景）
- [x] 4.5 确保 error 返回值均被检查，必不出错的函数用 `_ =` 处理

## 5. 代码正确性检查

- [x] 5.1 确认 `ModifyInstancePasswordPolicy` 请求参数（InstanceId、PasswordPolicy.Enabled/MinLetterCount/MinDigitCount/MinSpecialCount/MinLength）均在云 API 入参中存在
- [x] 5.2 确认 `DescribeInstancePasswordPolicy` 出参字段（PasswordPolicy.Enabled/MinLetterCount/MinDigitCount/MinSpecialCount/MinLength）与 schema 字段一一对应
- [x] 5.3 确认 Read 中设置字段前均检查 nil，Create 中检查返回值非空

## 6. 收尾

- [x] 6.1 由 tfpacer-finalize skill 执行 `gofmt`、`make doc`、生成 changelog 并推送
