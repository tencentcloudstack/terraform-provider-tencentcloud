## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/ckafka/resource_tc_ckafka_instance.go` 的 `ResourceTencentCloudCkafkaInstance()` schema map 中新增 `delete_protection_enable` 字段：`Type: schema.TypeInt, Optional: true, Computed: true`，Description 说明"实例删除保护开关: 1 开启 0 关闭"
- [x] 1.2 确认 `delete_protection_enable` 不加入 `immutableArgs` 列表（其为可变属性，允许原地更新）

## 2. Create 函数扩展

- [x] 2.1 在 `resourceTencentCloudCkafkaInstanceCreate` 末尾构造 `modifyRequest`（`ModifyInstanceAttributesRequest`）的现有逻辑块中，新增 `if v, ok := d.GetOkExists("delete_protection_enable"); ok { needModify = true; modifyRequest.DeleteProtectionEnable = helper.Int64(int64(v.(int))) }`，确保用户显式配置 `0` 时也能传入
- [x] 2.2 复用已有的 `service.ModifyCkafkaInstanceAttributes(ctx, modifyRequest)` 调用，不新增额外 API 调用

## 3. Update 函数扩展

- [x] 3.1 在 `resourceTencentCloudCkafkaInstanceUpdate` 的 `modifyInstanceAttributesFlag` 逻辑块中，新增 `if d.HasChange("delete_protection_enable") { if v, ok := d.GetOkExists("delete_protection_enable"); ok { request.DeleteProtectionEnable = helper.Int64(int64(v.(int))); modifyInstanceAttributesFlag = true } }`
- [x] 3.2 复用已有的 `service.ModifyCkafkaInstanceAttributes(ctx, request)` 调用
- [x] 3.3 确认 Update 函数末尾仍调用 `resourceTencentCloudCkafkaInstanceRead(d, meta)` 回写最新状态

## 4. Read 函数扩展

- [x] 4.1 在 `resourceTencentCloudCkafkaInstanceRead` 调用 `DescribeInstanceAttributes` 的 retry 块内，处理 `attr := response.Response.Result` 后，新增 `if attr.DeleteProtectionEnable != nil { _ = d.Set("delete_protection_enable", attr.DeleteProtectionEnable) }`（set 前判 nil）

## 5. 单元测试

- [x] 5.1 在 `tencentcloud/services/ckafka/resource_tc_ckafka_instance_test.go` 中使用 gomonkey mock 云 API，新增 Create 时 `delete_protection_enable = 1` 的测试用例，断言 `ModifyInstanceAttributes` 被调用且 `DeleteProtectionEnable == 1`
- [x] 5.2 新增 Create 时 `delete_protection_enable = 0` 的测试用例（验证 GetOkExists 能传入显式 0）
- [x] 5.3 新增 Update 时 `delete_protection_enable` 从 `0` 变为 `1` 的测试用例，断言 `ModifyInstanceAttributes` 被调用且 `DeleteProtectionEnable == 1`
- [x] 5.4 新增 Update 时 `delete_protection_enable` 从 `1` 变为 `0` 的测试用例
- [x] 5.5 新增 Read 时 `DescribeInstanceAttributes` 返回 `DeleteProtectionEnable` 非 nil 的测试用例，断言 state 被正确回填
- [x] 5.6 新增 Read 时 `DeleteProtectionEnable` 为 nil 的测试用例，断言不触发 panic 且跳过 set
- [x] 5.7 使用 `go test ./tencentcloud/services/ckafka/ -run <TestFunc> -v -count=1 -gcflags=all=-l` 跑通涉及的单元测试文件

## 6. 文档同步

- [x] 6.1 在 `tencentcloud/services/ckafka/resource_tc_ckafka_instance.md` 的 Example Usage 中新增 `delete_protection_enable` 字段示例（取值 `1` 开启 / `0` 关闭）
- [ ] 6.2 在收尾阶段执行 `make doc`，根据 provider 规范重新生成 `website/docs/` 下的 markdown 文档（禁止手改 website/ 目录）

## 7. 验证（收尾阶段执行）

- [ ] 7.1 在收尾阶段执行 `gofmt` 格式化变更的 Go 代码
- [ ] 7.2 确认所有涉及的单元测试文件通过 `go test -gcflags=all=-l`
