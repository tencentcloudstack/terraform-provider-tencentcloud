## Context

Terraform Provider for TencentCloud 目前不支持管理 EdgeOne (TEO) 多通道安全加速网关的接入密钥配置。该配置属于 RESOURCE_KIND_CONFIG 类型，即站点存在则配置就存在，主要管理配置的读取和更新。

云 API 接口情况：
- `DescribeMultiPathGatewaySecretKey`：查询多通道安全加速网关接入密钥，入参为 ZoneId，出参为 SecretKey
- `CreateMultiPathGatewaySecretKey`：创建多通道安全加速网关接入密钥，入参为 ZoneId 和 SecretKey（SecretKey 非必填，不填系统自动生成）
- `ModifyMultiPathGatewaySecretKey`：修改多通道安全加速网关接入密钥，入参为 ZoneId 和 SecretKey
- 另有 `RefreshMultiPathGatewaySecretKey` 接口，但本次需求不涉及

参考资源：`tencentcloud_teo_ddos_protection_config`（同属 TEO 产品下的 CONFIG 类型资源）

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_multi_path_gateway_secret_key` 资源，支持管理 TEO 多通道安全加速网关接入密钥
- 实现 Create（设置 ID 后调用 Update）、Read、Update、Delete 四个标准 CRUD 方法
- 支持 Import 导入已有配置
- 提供 .md 文档用于 make doc 生成
- 编写单元测试覆盖核心逻辑

**Non-Goals:**
- 不实现 `RefreshMultiPathGatewaySecretKey` 接口的调用
- 不实现数据源（data source）
- 不实现异步轮询机制（Describe、Create 和 Modify 均为同步接口）

## Decisions

1. **资源 ID 设计**：使用 `zone_id` 作为资源 ID。因为每个站点只有一个密钥配置，zone_id 可以唯一标识一个密钥配置。zone_id 设置为 ForceNew，变更 zone_id 时需要重建资源。

2. **Create/Update 方法设计**：采用 RESOURCE_KIND_CONFIG 的标准模式，Create 方法设置 zone_id 为资源 ID，然后直接调用 Update 方法。在 Update 方法中，首先调用 `DescribeMultiPathGatewaySecretKey` 查询当前密钥是否存在：
   - 如果密钥**存在**（Describe 返回非 nil SecretKey），调用 `CreateMultiPathGatewaySecretKey` API 替换密钥
   - 如果密钥**不存在**（Describe 返回 nil SecretKey 或 ResourceNotFound 错误），调用 `ModifyMultiPathGatewaySecretKey` API 设置密钥
   这种基于查询结果的判断方式比 `d.IsNewResource()` 更可靠，能够正确处理外部状态变更（如密钥被手动删除后重新创建的场景）。

3. **Delete 方法设计**：由于密钥配置是站点级别的，站点存在则密钥配置就存在，无法真正删除。Delete 方法仅从 Terraform state 中移除资源记录，不做实际的 API 调用。

4. **Schema 设计**：
   - `zone_id`：Required, ForceNew, TypeString - 站点 ID
   - `secret_key`：Required, TypeString, Sensitive - 多通道安全加速网关接入密钥（base64 字符串，编码前长度 32-48 个字符）

5. **密钥敏感性**：`secret_key` 为敏感字段，需设置 `Sensitive: true`，避免在 Terraform 日志和输出中暴露密钥内容。

6. **重试机制**：Read 操作使用 `tccommon.ReadRetryTimeout` 进行重试，Create 和 Update 操作使用 `tccommon.WriteRetryTimeout` 进行重试，失败时使用 `tccommon.RetryError()` 包装错误返回。

7. **代码文件命名**：`resource_tc_teo_multi_path_gateway_secret_key_config.go`，遵循 RESOURCE_KIND_CONFIG 的命名规范 `resource_tc_<Product>_<资源名>_config.go`。

## Risks / Trade-offs

- [密钥删除不可逆] → Delete 方法仅清理 state，不做实际 API 调用。用户通过 Terraform 销毁资源后再 apply 会重新设置密钥。
- [secret_key 为 Sensitive 字段] → Terraform 会在 state 中标记为 sensitive，但 state 文件中仍以明文存储，需提醒用户妥善保护 state 文件。
- [Describe 接口返回的 SecretKey 可能与 Terraform 配置不一致] → Read 方法会比较远端返回值与本地配置，若不一致则更新 state，确保 Terraform 能检测到配置漂移。
