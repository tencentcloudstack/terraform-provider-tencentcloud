## Why

腾讯云防火墙（FWM，云防火墙 CFW 的新一代产品族，SDK 包 `tencentcloud-sdk-go/tencentcloud/fwm/v20250611`）已支持「企业安全组规则」能力：在企业安全组规则组（Product=enterprise_sg）内创建、查询、修改、删除单条规则。该能力在云 API 已发布（`CreateSecurityGroupRule` / `DescribeSecurityGroupRules` / `ModifySecurityGroupRule` / `DeleteSecurityGroupRule` 四个接口，全同步），但 Terraform Provider 尚未暴露，导致用户无法用 HCL 声明式管理企业安全组规则，必须在控制台或脚本中维护，无法做基础设施代码化（IaC）。

本 change 新增通用型资源 `tencentcloud_fwm_security_group_rule`，让用户在 HCL 里声明一条企业安全组规则的完整生命周期，由 Provider 调用上述四个同步 CRUD 接口完成管理。

## What Changes

- 新增通用型（RESOURCE_KIND_GENERAL）资源 `tencentcloud_fwm_security_group_rule`，覆盖完整 CRUD 生命周期，资源代码文件 `tencentcloud/services/fwm/resource_tc_fwm_security_group_rule.go`。
- 新增 FWM 产品在 Provider 的客户端注册：在 `tencentcloud/connectivity/client.go` 增加 `fwmv20250611` 包 import、`fwmv20250611Conn` 字段与 `UseFwmV20250611Client()` 访问器（参照 `UseIgtmV20231024Client` 与 `ga2v20250115` 模式）。在 `tencentcloud/provider.go` 注册资源 `tencentcloud_fwm_security_group_rule`，在 `tencentcloud/provider.md` 追加资源名。
- Schema 采用平铺结构：入参字段全部定义在资源顶层（不嵌套 `rules`/`rule` 一层数组），Create 时构造 1 个元素的 `Rules` 数组，Modify 时构造 `Rule` 对象。
- 资源 ID 为复合结构 `GroupId<tccommon.FILED_SP>RuleId`；import 支持 `GroupId:RuleId`。
- 字段约束（依据用户自定义要求第 13、18 条与云 API 实际定义）：
  - `group_id`：Required（四接口均必填，父资源引用）。
  - `ip_version`：Required + ForceNew。
  - `source_type` / `target_type` / `scope`：ForceNew（变更视为重建）。
  - `order_index`：Required，可变（原地更新）。
  - 其余可写字段（`source_content` / `dest_content` / `dest_type` / `protocol` / `port` / `service_template_id` / `rule_action` / `description` / `protocol_port_type` / `belong_member_id`）为 Optional，可变。
- 创建入参（字符串）与查询响应（int）双向映射，避免 plan 永久 diff（用户自定义要求第 8 条）：
  - `source_type` / `dest_type`：入参字符串（net/vpc/subnet/cvm/clb/eni/cdb/template/tag/region/resourcegroup/addressbook/instance），响应 `SourceType`/`TargetType` 为 int（net=0,vpc=1,subnet=2,cvm=3,clb=4,eni=5,cdb=6,template=7,tag=8,region=9,resourcegroup=100）。
  - `rule_action`：入参字符串（accept/drop/log），响应 `Strategy` 为 int（LOG=0,DROP=1,ACCEPT=2）。
  - `description`（入参）↔ `Detail`（响应）；`port`↔`Port`；`source_content`↔`SourceId`；`dest_content`↔`TargetId`；`belong_member_id`↔`BelongMember.MemberId`。
- 只读展示字段（computed，不得映射进可比字段）：`rule_id`、`source_id`、`target_id`、`strategy`、`detail`、`region`、`service_template_id`（展示）、`protocol_port_name`、`sou_instance_name`、`sou_public_ip`、`sou_private_ip`、`sou_cidr`、`sou_parameter_name`、`instance_name`、`public_ip`、`private_ip`、`cidr`、`parameter_name`、`dns_parse_count`、`all_total_count`、`id` 等。
- Read 使用描述：用 `DescribeSecurityGroupRules`（列表查询无单实例 Get）+ `Filters=[{Name:"RuleId",Values:[ruleId],OperatorType:1}]` 精准查单条，取 `Rules` 首元素。Read 遇空（规则被外部删除）打印日志后 `d.SetId("")` 移除 state 不报错。
- 四接口全同步（无任务 ID、无中间态）：不生成 waitReady/async 轮询。
- 新增 `.md` 资源文档（`resource_tc_fwm_security_group_rule.md`，含 Import 段说明使用 `GroupId:RuleId` 复合 id）+ 单元测试 `_test.go`（使用 gomonkey mock 云 API，做业务逻辑单测，不使用 TF 测试套件）。

## Capabilities

### New Capabilities

- `fwm-security-group-rule-resource`: 新增 `tencentcloud_fwm_security_group_rule` 资源的 schema、CRUD（Create / Read / Update / Delete）行为、复合 ID 约定、入参/响应双向类型映射、只读展示字段、import 支持、文档与测试规范，作为企业安全组规则管理的 IaC 入口。

### Modified Capabilities

<!-- 不修改任何已有 capability 的 requirement。仅在 provider.go / provider.md / connectivity/client.go 增加注册与客户端条目。 -->

## Impact

- 新文件：
  - `tencentcloud/services/fwm/service_tencentcloud_fwm.go`（FWM service 层，含 `DescribeFwmSecurityGroupRuleById`、`DeleteFwmSecurityGroupRuleById` 等方法）。
  - `tencentcloud/services/fwm/resource_tc_fwm_security_group_rule.go`（资源 CRUD）。
  - `tencentcloud/services/fwm/resource_tc_fwm_security_group_rule.md`（资源文档）。
  - `tencentcloud/services/fwm/resource_tc_fwm_security_group_rule_test.go`（单元测试，gomonkey mock）。
  - `website/docs/r/fwm_security_group_rule.html.markdown`（由 `make doc` 生成）。
- 既有文件改动：
  - `tencentcloud/connectivity/client.go`：新增 `fwmv20250611` import、`fwmv20250611Conn` 字段、`UseFwmV20250611Client()` 访问器。
  - `tencentcloud/provider.go`：新增 fwm 服务 import 与 `"tencentcloud_fwm_security_group_rule": fwm.ResourceTencentCloudFwmSecurityGroupRule()` 注册。
  - `tencentcloud/provider.md`：新增 FWM（Firewall）段及资源名一行，供 gendoc 索引识别。
- 依赖：vendor 下 SDK `tencentcloud-sdk-go/tencentcloud/fwm/v20250611`（含 client.go / models.go）已在工作区（`go.mod` / `vendor/modules.txt` 已纳入 fwm 模块），无需升级 SDK。
- 不修改任何既有资源/数据源 schema，对现有用户零影响。