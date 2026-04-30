## 1. Schema 调整

- [x] 1.1 在 `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway.go` 中，为 `status` 字段新增 `Optional: true` 属性（保留 `Computed: true`），并更新 Description 说明合法取值 `online`/`offline`
- [x] 1.2 检查资源 schema 是否需要在 `Timeouts` 块中补充 `Update` 字段（若已有则复用，未声明则追加 `schema.ResourceTimeout{ Update: schema.DefaultTimeout(10 * time.Minute) }`）

## 2. Update 函数扩展

- [x] 2.1 在 `resourceTencentCloudTeoMultiPathGatewayUpdate` 中，新增独立分支 `if d.HasChange("status")`，仅在 `d.GetOk("status")` 为 true 时触发
- [x] 2.2 构造 `teov20220901.NewModifyMultiPathGatewayStatusRequest()`，填充 `ZoneId`、`GatewayId`、`GatewayStatus`
- [x] 2.3 使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装 `ModifyMultiPathGatewayStatusWithContext` 调用，按 provider 现有模式处理错误（`tccommon.RetryError` / `log.Printf("[DEBUG]...")`）
- [x] 2.4 调用成功后，使用 `service.DescribeTeoMultiPathGatewayById` 在 `d.Timeout(schema.TimeoutUpdate)` 内轮询，等待 `respData.Status` 达到目标值（`online`/`offline`）或中间态（`creating` 等）结束
- [x] 2.5 复查 Update 函数末尾仍调用 `resourceTencentCloudTeoMultiPathGatewayRead(d, meta)` 以回写最新状态

## 3. 单元测试

- [x] 3.1 在 `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_test.go` 中新增测试用例 `TestAccTeoMultiPathGateway_UpdateStatusToOffline`，使用 gomonkey mock `ModifyMultiPathGatewayStatusWithContext` 和 `DescribeMultiPathGateways`
- [x] 3.2 新增测试用例 `TestAccTeoMultiPathGateway_UpdateStatusToOnline`，覆盖启用分支
- [x] 3.3 新增测试用例 `TestAccTeoMultiPathGateway_UpdateStatusNotSet`，验证未配置 status 时不调用 `ModifyMultiPathGatewayStatus`
- [x] 3.4 保证已有 Create/Read/Update(name)/Delete 测试用例继续通过（不修改其行为）

## 4. 文档同步

- [x] 4.1 在 `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway.md` 的 Example Usage 中新增 `status = "online"` 字段示例
- [x] 4.2 在示例描述中标注 status 的可选取值（`online`/`offline`）
- [x] 4.3 执行 `make doc`，根据 provider 规范重新生成 `website/docs/` 下的 markdown 文档（禁止手改）

## 5. 验证

- [x] 5.1 运行 `go build ./...` 确保编译通过
- [x] 5.2 运行 `go vet ./tencentcloud/services/teo/...`
- [x] 5.3 运行 `go test ./tencentcloud/services/teo/ -run TestAccTeoMultiPathGateway -v -count=1 -gcflags="all=-l"` 确认所有单测通过（包括新增 status 变更用例）
- [x] 5.4 在本地 `terraform plan` 已有配置（未声明 status）验证无 plan drift（向后兼容检查）
