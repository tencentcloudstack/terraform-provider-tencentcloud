## 1. 客户端访问器与服务目录初始化

- [x] 1.1 在 `tencentcloud/connectivity/client.go` 中新增 bdrc SDK import（别名 `bdrcv20260330`，路径 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330`）
- [x] 1.2 在 `tencentcloud/connectivity/client.go` 的 `TencentCloudClient` 结构体中新增 `bdrcv20260330Conn *bdrcv20260330.Client` 字段
- [x] 1.3 在 `tencentcloud/connectivity/client.go` 中新增 `UseBdrcV20260330Client()` 方法（懒加载，`NewClientProfile(300)`，`WithHttpTransport(&LogRoundTripper{})`）
- [x] 1.4 创建 `tencentcloud/services/bdrc/` 目录及 `service_tencentcloud_bdrc.go` 服务层文件骨架（定义 `BdrcService{client}` 结构体）

## 2. 服务层实现

- [x] 2.1 在 `service_tencentcloud_bdrc.go` 中实现 `DescribeSecurityGroupMappingById(ctx, sitePairId, securityGroupMappingId)` 方法：调用 `DescribeSecurityGroupMappings`（`SitePairId` 入参，`Limit=500`），在返回 `SecurityGroupMappingSet` 中按 `SecurityGroupMappingId` 精确匹配，返回 `*bdrcv20260330.SecurityGroupMapping`；使用 `resource.Retry(tccommon.ReadRetryTimeout)` 包装，失败用 `tccommon.RetryError` 包装
- [x] 2.2 在 `service_tencentcloud_bdrc.go` 中实现 `DescribeSecurityGroupMappingByFilter(ctx, sitePairId, srcSecurityGroupId, targetSecurityGroupId)` 方法：调用 `DescribeSecurityGroupMappings`，`Filters` 设置 `src-security-group-id` 与 `target-security-group-id`，`Limit=500`，返回命中记录（用于 Create 后轮询定位新建映射）；使用 `resource.Retry(tccommon.ReadRetryTimeout)` 包装

## 3. 资源 Schema 与 CRUD 实现

- [x] 3.1 创建 `tencentcloud/services/bdrc/resource_tc_bdrc_security_group_mapping.go`，定义 `ResourceTencentCloudBdrcSecurityGroupMapping()` 返回 `*schema.Resource`，包含 Create/Read/Update/Delete/Importer
- [x] 3.2 定义 Schema 字段：`site_pair_id`(Required,ForceNew)、`src_security_group_id`(Required,ForceNew)、`target_security_group_id`(Required,ForceNew)、`security_group_mapping_id`(Computed)、`source_security_group_id`(Computed)、`life_state`(Computed)
- [x] 3.3 实现 `resourceTencentCloudBdrcSecurityGroupMappingCreate`：构建 `CreateSecurityGroupMappingRequest`（映射 `src_security_group_id`→`SrcSecurityGroupId` 等），用 `resource.Retry(tccommon.WriteRetryTimeout)` 调用；调用后用 `helper.Retry()` 轮询 `DescribeSecurityGroupMappingByFilter` 直到命中新建映射，取 `SecurityGroupMappingId`；打印 `logId` 与 `d.Id()`；设置联合 ID `site_pair_id#security_group_mapping_id`；调用 Read 回填
- [x] 3.4 实现 `resourceTencentCloudBdrcSecurityGroupMappingRead`：从 `d.Id()` 解析 `site_pair_id` 与 `security_group_mapping_id`（长度不为 2 返回 `id is broken` 错误）；调用 `DescribeSecurityGroupMappingById`；空返回时先 `log.Printf("[CRUD] bdrc security_group_mapping id=%s", d.Id())` 再 `d.SetId("")`；非空时按 nil 检查后 `d.Set` 各字段
- [x] 3.5 实现 `resourceTencentCloudBdrcSecurityGroupMappingUpdate`：`immutableArgs = []string{"src_security_group_id", "target_security_group_id", "site_pair_id"}`，检测 `d.HasChange(v)` 返回 error；无变更则调用 Read
- [x] 3.6 实现 `resourceTencentCloudBdrcSecurityGroupMappingDelete`：从 `d.Id()` 解析 `site_pair_id` 与 `security_group_mapping_id`；构建 `DeleteSecurityGroupMappingRequest`（`SitePairId` + `SecurityGroupMappingIds=[mappingId]`）；用 `resource.Retry(tccommon.WriteRetryTimeout)` 调用，失败用 `tccommon.RetryError` 包装

## 4. Provider 注册

- [x] 4.1 在 `tencentcloud/provider.go` 中新增 bdrc 服务 import（`github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/bdrc`）
- [x] 4.2 在 `tencentcloud/provider.go` 的 `ResourcesMap` 中注册 `"tencentcloud_bdrc_security_group_mapping": bdrc.ResourceTencentCloudBdrcSecurityGroupMapping()`

## 5. 资源文档

- [x] 5.1 创建 `tencentcloud/services/bdrc/resource_tc_bdrc_security_group_mapping.md`：一句话描述（带上 BDRC 云产品名称）、Example Usage、Import 部分（说明需使用联合 ID `sitePairId#securityGroupMappingId`）；不添加 Argument Reference / Attribute Reference

## 6. 单元测试

- [x] 6.1 创建 `tencentcloud/services/bdrc/resource_tc_bdrc_security_group_mapping_test.go`，使用 gomonkey mock 云 API（不使用 terraform 测试套件）
- [x] 6.2 补充 Create 测试用例：mock `CreateSecurityGroupMappingWithContext` 与 `DescribeSecurityGroupMappingsWithContext` 返回命中映射，验证联合 ID 设置与字段回填
- [x] 6.3 补充 Read 测试用例：mock 命中与未命中（空返回）场景
- [x] 6.4 补充 Update 测试用例：验证不可变字段变更返回 error
- [x] 6.5 补充 Delete 测试用例：mock `DeleteSecurityGroupMappingWithContext` 成功路径
- [x] 6.6 确保 `go vet`/`go build` 之外不执行测试命令，保证生成的测试代码在当前环境下可正确构建执行；所有函数返回的 error 均被检查或赋值给 `_`

## 7. 验证

- [x] 7.1 检查所有 CRUD 函数中云 API 参数与 vendor SDK 结构体字段一致（Create 入参在 `CreateSecurityGroupMappingRequest`、Read/Delete 入参在对应 Request、Read 出参在 `SecurityGroupMapping` 结构体）
- [x] 7.2 检查 retry 块内仅执行云 API 调用，设置 id 等成功操作置于 retry 块外
- [x] 7.3 检查无 `_extension.go` 文件生成、资源 go 文件开头无注释、日志/错误描述统一使用资源名 `bdrc security_group_mapping`
