## 1. Connectivity Client Registration

- [x] 1.1 在 `tencentcloud/connectivity/client.go` 添加 import 别名 `fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`。
- [x] 1.2 在 `TencentCloudClient` 结构体添加私有缓存字段 `fwmv20250611Conn *fwmv20250611.Client`。
- [x] 1.3 添加 `UseFwmV20250611Client()` 方法，lazy-init 模式引用 `UseIgtmV20231024Client`（line 2247）：`me.NewClientProfile(300)` + `WithHttpTransport(&LogRoundTripper{})`。

## 2. Service Layer

- [x] 2.1 创建 `tencentcloud/services/fwm/service_tencentcloud_fwm.go`，声明 `FwmService` 结构体与 `NewFwmService(client)` 构造函数（引用 `IgtmService` 模式）。
- [x] 2.2 实现 `DescribeFwmEdgeAclRuleById(ctx, groupId, ruleId, direction int64)` 方法：构造 `DescribeEdgeAclRulesRequest`，设 `GroupId`、`Direction`、`Filters=[{Name:"RuleId",Values:[ruleId],OperatorType:1}]`、`Limit=1000`，包裹 `resource.Retry(tccommon.ReadRetryTimeout)`；取 `Rules` 首元素返回，未找到返回 `(nil, nil)`。

## 3. Resource Schema

- [x] 3.1 创建 `tencentcloud/services/fwm/resource_tc_fwm_edge_acl_rule.go`，定义 `ResourceTencentCloudFwmEdgeAclRule()`，含 Create/Read/Update/Delete + Importer(`schema.ImportStatePassthrough`)。
- [x] 3.2 定义 schema 顶层可写字段：`group_id`(Required)、`order_index`(Required,可变)、`direction`(Required+ForceNew,int)、`source_content`、`source_type`(ForceNew)、`target_content`、`target_type`(ForceNew)、`port`、`protocol`、`rule_action`、`description`、`scope`、`belong_member_id`、`param_template_id`。
- [x] 3.3 定义 Computed 只读字段：`rule_id`、`sequence`、`create_time`、`update_time`、`source_name`、`target_name`、`param_template_name`、`country_name`、`city_name`、`instance_name`。

## 4. CRUD Implementation

- [x] 4.1 实现 `resourceTencentCloudFwmEdgeAclRuleCreate`：从 schema 构造单元素 `Rules []*EdgeAclRuleInfo`（写入字段名: OrderIndex/Direction/SourceContent/SourceType/TargetContent/TargetType/Port/Protocol/RuleAction/Description/Scope/BelongMemberId/ParamTemplateId），`group_id`→`GroupId`；包裹 `WriteRetryTimeout` retry；retry 内检查 Response 空值 NonRetryableError；retry 外打印 logId/d.Id 后检查 `RuleIds[0]` 非空，取首元素为 ruleId，`d.SetId` 设复合 ID；调用 Read。不生成异步轮询。
- [x] 4.2 实现 `resourceTencentCloudFwmEdgeAclRuleRead`：拆分 ID 校验 2 段得 groupId/ruleId；调用 `FwmService.DescribeFwmEdgeAclRuleById(ctx, groupId, ruleId, direction)`；nil 时打印 `[CRUD] fwm_edge_acl_rule id=%s` 后 `d.SetId("")` 返回 nil；非 nil 时按错误位映射回填：`order_index`←Sequence、`belong_member_id`←BelongMember.MemberId(非空才 set)、scope/rule_action 直透传、其余字段 nil 防御后 `_ = d.Set(...)`。
- [x] 4.3 实现 `resourceTencentCloudFwmEdgeAclRuleUpdate`：拆分 ID；检测可变字段 HasChange(`order_index`/`source_content`/`target_content`/`port`/`protocol`/`rule_action`/`description`/`scope`/`belong_member_id`/`param_template_id`)；构造 `ModifyEdgeAclRuleRequest`，`Rule` 为含 RuleId + 全部可写字段的 `EdgeAclRuleInfo`（全量覆盖），`GroupId` 透传；包裹 `WriteRetryTimeout` retry；成功后调用 Read。无变更不下发。
- [x] 4.4 实现 `resourceTencentCloudFwmEdgeAclRuleDelete`：拆分 ID；构造 `DeleteEdgeAclRuleRequest`，`RuleIds` 传单元素；包裹 `WriteRetryTimeout` retry；retry 内错误码含 `ResourceNotFound` 返回 nil（幂等），否则 `tccommon.RetryError(e)`。

## 5. Provider Registration

- [x] 5.1 在 `tencentcloud/provider.go` 添加 import `"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"`。
- [x] 5.2 在 resources map 注册 `"tencentcloud_fwm_edge_acl_rule": fwm.ResourceTencentCloudFwmEdgeAclRule()`。
- [x] 5.3 在 `tencentcloud/provider.md` 添加 `tencentcloud_fwm_edge_acl_rule` 资源条目。

## 6. Documentation

- [x] 6.1 创建 `tencentcloud/services/fwm/resource_tc_fwm_edge_acl_rule.md`：开头中文一句话描述带上云产品名（云防火墙）；Example Usage（testValue 保证可创建成功：direction=0、order_index=1、source_type=ip、source_content=1.1.1.1、target_type=ip、target_content=2.2.2.2、protocol=TCP、port=80、rule_action=accept、scope=all）；Import 部分（RESOURCE_KIND_GENERAL 有 Import，说明用 `GroupId:RuleId` 两段式）。不写 Argument Reference / Attribute Reference。

## 7. Tests (gomonkey, 不用 terraform 测试套件)

- [x] 7.1 创建 `tencentcloud/services/fwm/resource_tc_fwm_edge_acl_rule_test.go`（`package fwm_test`）。定义 `mockMeta`/`newMockMeta`/`ptrString`（引用 teo 风格）。
- [x] 7.2 `Test..._Create`：mock `UseFwmV20250611Client` 与 `CreateEdgeAclRule` 返回 RuleIds，验证 `d.SetId` 为 `groupId#ruleId`、字段正确构造。
- [x] 7.3 `Test..._Read`：mock `DescribeEdgeAclRules` 返回首元素（含 Sequence/BelongMember.MemberId），验证 order_index/belong_member_id 错位映射与其他字段回填；测试 Rules 为空时 `d.SetId("")`。
- [x] 7.4 `Test..._Update`：mock `ModifyEdgeAclRule` 成功，验证请求携带全量字段。
- [x] 7.5 `Test..._Delete`：mock `DeleteEdgeAclRule` 成功；mock 返回 `ResourceNotFound` 验证幂等（无错）。
- [x] 7.6 `Test..._Schema`：校验关键字段 Required/ForceNew/Computed 配置。

## 8. Verification (由其他流程执行，不在本阶段运行)

- [x] 8.1 确认 `tencentcloud/services/fwm/` 包可编译。
- [x] 8.2 确认 `tencentcloud/connectivity/` 新增 client 方法可编译。
- [x] 8.3 确认 `tencentcloud/` 全 provider 可编译。
- [x] 8.4 单元测试编译通过（`-gcflags="all=-l"`）。