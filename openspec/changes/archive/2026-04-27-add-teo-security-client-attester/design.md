## Context

TEO (EdgeOne) 是腾讯云的边缘安全加速平台。客户端认证选项（Security Client Attester）是 TEO 安全功能的一部分，允许用户配置不同类型的客户端认证方式（TC-RCE 风险识别、TC-CAPTCHA 天御验证码、TC-EO-CAPTCHA EdgeOne 人机校验）来保护站点。

当前 terraform-provider-tencentcloud 中已有 TEO 安全相关资源如 `tencentcloud_teo_security_ip_group` 和 `tencentcloud_teo_security_policy_config`，但缺少客户端认证选项的管理能力。新增资源 `tencentcloud_teo_security_client_attester` 将补齐此能力。

云API接口情况：
- `CreateSecurityClientAttester`: 创建认证选项，入参 ZoneId + ClientAttesters，出参 ClientAttesterIds
- `DescribeSecurityClientAttester`: 查询认证选项，入参 ZoneId + 分页参数，出参 ClientAttesters + TotalCount
- `ModifySecurityClientAttester`: 修改认证选项，入参 ZoneId + ClientAttesters
- `DeleteSecurityClientAttester`: 删除认证选项，入参 ZoneId + ClientAttesterIds

所有接口均为同步接口，无需异步轮询。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_security_client_attester` 资源的完整 CRUD 支持
- 支持三种认证方式（TC-RCE、TC-CAPTCHA、TC-EO-CAPTCHA）的配置
- 遵循现有 TEO 安全资源的代码风格和模式
- 在 provider.go 中注册资源并在 provider.md 中添加注释
- 生成 .md 文档文件
- 编写使用 gomonkey mock 的单元测试

**Non-Goals:**
- 不实现数据源（datasource）
- 不修改已有资源的 schema 或行为
- 不支持异步接口轮询（API均为同步）

## Decisions

### 1. 复合 ID 设计：zone_id + client_attester_ids
**决策**: 使用 `zone_id` + `client_attester_ids`（以逗号连接后再用 FILED_SP 拼接）作为复合 ID。

**理由**: 
- Create 接口返回 `ClientAttesterIds`（字符串列表），Delete 接口需要 `ZoneId` + `ClientAttesterIds`，因此需要存储这两个信息
- 遵循现有 TEO 资源（如 `tencentcloud_teo_security_ip_group`）的复合 ID 模式
- 格式: `zone_id#id1,id2,id3`

**替代方案**: 仅用 `zone_id` 作为 ID。但这样在 Delete 时需要额外查询获取 client_attester_ids，增加复杂性和潜在的一致性问题。

### 2. client_attesters 参数设计为 TypeList
**决策**: `client_attesters` 参数设计为 TypeList，每个元素是包含嵌套块的 Resource。

**理由**:
- 云 API 的 `ClientAttesters` 字段是 `[]*ClientAttester` 数组类型
- 每个认证选项包含条件性子配置（根据 `attester_source` 不同需要不同的 option 块）
- TypeList + nested Resource 可以完整映射 API 结构

### 3. 嵌套 option 块设计为 TypeList（MaxItems: 1）
**决策**: `tc_rce_option`、`tc_captcha_option`、`tc_eo_captcha_option` 均设计为 TypeList 且 MaxItems=1。

**理由**:
- 对应云 API 中的结构体指针（`*TCRCEOption`、`*TCCaptchaOption`、`*TCEOCaptchaOption`）
- MaxItems=1 确保每个认证选项只有一份配置
- 遵循现有资源中类似嵌套结构的设计模式

### 4. DescribeSecurityClientAttester 分页查询
**决策**: 查询时使用 Limit=100（API注释中标注的最大值）进行分页查询，获取所有认证选项后在客户端过滤匹配 ID 的记录。

**理由**:
- Describe 接口只有 ZoneId 作为过滤条件，没有按 ID 过滤的能力
- 需要获取全量数据后在本地筛选
- 使用 Limit=100 减少分页请求次数

### 5. Update 策略：整体替换
**决策**: Update 时将当前所有 `client_attesters` 参数整体传入 ModifySecurityClientAttester 接口。

**理由**:
- ModifySecurityClientAttester 接口接受完整的 ClientAttesters 列表，是整体替换语义
- 与 Create 接口入参结构一致
- 遵循现有 TEO 资源（如 security_ip_group）的更新模式

## Risks / Trade-offs

- **[Risk] client_attesters 列表较大时性能问题** → 缓解: Describe 接口使用 Limit=100 减少分页次数，一般场景下认证选项数量不会太多
- **[Risk] 复合 ID 中 client_attester_ids 的排序一致性** → 缓解: 在 Create 时按 API 返回顺序存储，Read 时按 API 返回顺序更新，确保一致性
- **[Risk] 条件性子配置的验证** → 缓解: 依赖云 API 端验证，不在 Terraform schema 层做复杂条件校验，避免与 API 行为不一致
