# Tasks for `add-fwm-vpc-acl-rule-resource`

## 1. Connectivity client 接入

- [x] 1.1 在 `tencentcloud/connectivity/client.go` 新增 import `fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`（按字母序插入）。
- [x] 1.2 在 `TencentCloudClient` struct 新增字段 `fwmv20250611Conn *fwmv20250611.Client`（line 268 附近，与其他 Conn 字段并列）。
- [x] 1.3 新增访问器方法 `UseFwmV20250611Client() *fwmv20250611.Client`，参照 `UseIgtmV20231024Client` 模式（NewClientProfile(300)、NewClient、WithHttpTransport(&LogRoundTripper{})）。

## 2. Service layer

- [x] 2.1 创建 `tencentcloud/services/fwm/service_tencentcloud_fwm.go`，package `fwm`。包含 `NewFwmService(client *connectivity.TencentCloudClient) FwmService` 与 `type FwmService struct { client *connectivity.TencentCloudClient }`。import 风格对齐 `service_tencentcloud_igtm.go`（context、fmt、log、resource、fwmv20250611 SDK、tccommon、connectivity、helper、ratelimit）。
- [x] 2.2 实现 `DescribeFwmVpcAclRuleById(ctx context.Context, groupId, ruleId string) (ret *fwmv20250611.VpcAclRuleResp, allTotalCount *int64, errRet error)`：构建 `DescribeVpcAclRulesRequest`，设置 `GroupId = groupId`、`Filters = [{Name: "RuleId", Values: [ruleId], OperatorType: 1}]`、`Limit = 0`、`Offset = 0`，用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` + `ratelimit.Check` + `tccommon.RetryError(e)` 包装调用 `UseFwmV20250611Client().DescribeVpcAclRules`。retry 块内校验 `result == nil || result.Response == nil` → NonRetryableError。取 `Rules` 首元素；若 `len(Rules) == 0` 返回 `(nil, nil, nil)`。返回 `(Rules[0], Response.AllTotalCount, nil)`。
- [x] 2.3 实现 `DeleteFwmVpcAclRuleById(ctx context.Context, groupId, ruleId string) error`：构建 `DeleteVpcAclRuleRequest`，设置 `GroupId = groupId`、`RuleIds = [ruleId]`（1 元素），用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装调用。捕获 ResourceNotFound 错误码，视为删除成功返回 nil（幂等）。

## 3. Resource implementation

- [x] 3.1 创建 `tencentcloud/services/fwm/resource_tc_fwm_vpc_acl_rule.go`，package `fwm`。文件头无注释，import 风格对齐 `resource_tc_igtm_strategy.go`（context、fmt、log、strings、resource、schema、fwmv20250611、tccommon、helper）。
- [x] 3.2 声明 `ResourceTencentCloudFwmVpcAclRule() *schema.Resource`，含 `Create/Read/Update/Delete` 回调 + `Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}`。
- [x] 3.3 Schema 字段定义，按 design.md D1 表所有字段，Required/ForceNew/Computed/Optional 标志精确匹配。forceNew 集合：`group_id`、`edge_id`、`fw_group_id`、`ip_version`、`source_type`、`source_content`、`dest_type`、`dest_content`。description 字段英文。
- [x] 3.4 实现 `resourceTencentCloudFwmVpcAclRuleCreate`：
  - `defer LogElapsed("resource.tencentcloud_fwm_vpc_acl_rule.create")()` + `defer InconsistentCheck(d, meta)`。
  - 从 d.Get 读取 `group_id` 及所有规则字段，组装 `fwmv20250611.VpcAclRule{}`（注意：不设 IpVersion，入参结构体无此字段）。
  - `request.Rules = append(request.Rules, &rule)`（1 元素）。
  - `request.GroupId = helper.String(groupId)`。
  - retry 块（WriteRetryTimeout）调用 `UseFwmV20250611Client().CreateVpcAclRuleWithContext(ctx, request)`，校验 `result == nil || result.Response == nil` → NonRetryableError。
  - retry 外打印 `logId` 与检查 `RuleIds`：`if response.Response.RuleIds == nil || len(response.Response.RuleIds) == 0` → 返回 `NonRetryableError(fmt.Errorf("create fwm vpc_acl_rule failed, RuleIds is empty"))`。
  - `ruleId = *response.Response.RuleIds[0]`。
  - `d.SetId(strings.Join([]string{groupId, ruleId}, tccommon.FILED_SP))`。
  - `return resourceTencentCloudFwmVpcAclRuleRead(d, meta)`。
- [x] 3.5 实现 `resourceTencentCloudFwmVpcAclRuleRead`：
  - 拆分 ID：`idSplit := strings.Split(d.Id(), tccommon.FILED_SP)`，`len(idSplit) != 2` → error。
  - `groupId, ruleId := idSplit[0], idSplit[1]`。
  - 调用 `service.DescribeFwmVpcAclRuleById(ctx, groupId, ruleId)`。
  - 若返回 nil：`log.Printf("[CRUD] fwm_vpc_acl_rule id=%s", d.Id())` 保留现场，再 `d.SetId("")`，return nil。
  - 逐一回读字段，每个字段 `if respData.X != nil { _ = d.Set("x", respData.X) }`。
  - `order_index` ← `respData.Sequence`（非 nil 才设置）。
  - `belong_member_id` ← `respData.BelongMember.MemberId`（BelongMember 非 nil 且 MemberId 非 nil 才设置）。
  - `app_id`/`uin`/`nickname`/`member_id`/`node_name` ← `respData.BelongMember.*`（BelongMember 非 nil 且对应字段非 nil 才设置）。
  - `all_total_count` ← `allTotalCount`（非 nil 才设置）。
- [x] 3.6 实现 `resourceTencentCloudFwmVpcAclRuleUpdate`：
  - 拆分 ID，`groupId, ruleId := ...`。
  - `mutableArgs := []string{"protocol", "rule_action", "order_index", "description", "port", "param_template_id", "belong_member_id"}`，任一 `d.HasChange` 则 `needChange = true`。
  - 若 needChange：组装 `fwmv20250611.VpcAclRule{}`，设置 `RuleId = helper.String(ruleId)` 及所有可写字段（从 d.Get 完整读取，全量覆盖）。`request.Rule = &rule`（单对象）、`request.GroupId = helper.String(groupId)`。
  - retry 块（WriteRetryTimeout）调用 `UseFwmV20250611Client().ModifyVpcAclRuleWithContext(ctx, request)`，`tccommon.RetryError(e)`。
  - `return resourceTencentCloudFwmVpcAclRuleRead(d, meta)`。
- [x] 3.7 实现 `resourceTencentCloudFwmVpcAclRuleDelete`：
  - 拆分 ID，调用 `service.DeleteFwmVpcAclRuleById(ctx, groupId, ruleId)`，返回 err。

## 4. Provider registration

- [x] 4.1 在 `tencentcloud/provider.go` 新增 import `"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"`（按字母序，约 line 73 附近）。
- [x] 4.2 在资源注册映射新增 `"tencentcloud_fwm_vpc_acl_rule": fwm.ResourceTencentCloudFwmVpcAclRule(),`（按字母序，约 line 2072 f 开头区域）。

## 5. Documentation

- [x] 5.1 创建 `tencentcloud/services/fwm/resource_tc_fwm_vpc_acl_rule.md`，包含：
  - 一句话描述（带云产品名 FWM）：`Provides a resource to create a FWM VPC ACL rule`
  - Example Usage HCL 块：testValue 配置（ip_version=ipv4、order_index=1、source_type=ip、source_content=10.0.0.0/24、dest_type=ip、dest_content=192.168.1.0/24、protocol=TCP、port=80、rule_action=accept、description=test rule、edge_id=ALL、fw_group_id=ALL）
  - Import section：说明使用 `GroupId#RuleId` 联合 ID
- [x] 5.2 `make doc` 生成 `website/docs/r/fwm_vpc_acl_rule.html.markdown`（收尾阶段由 tfpacer-finalize 处理）。

## 6. Unit test

- [x] 6.1 创建 `tencentcloud/services/fwm/resource_tc_fwm_vpc_acl_rule_test.go`，package `fwm_test`，使用 gomonkey mock 云 API（参照 `resource_tc_cfw_cluster_nat_fw_switch_test.go` 模式）：
  - mockMeta 结构体实现 `tccommon.ProviderMeta`，`GetAPIV3Conn()` 返回 mock client。
  - `TestFwmVpcAclRule_Create`：mock `UseFwmV20250611Client` + `CreateVpcAclRuleWithContext`（返回 RuleIds 1元素）+ `DescribeVpcAclRules`（返回 Rules 1元素），构造 ResourceData 调 Create，断言 id 为 `groupId#ruleId`。
  - `TestFwmVpcAclRule_Read`：mock `UseFwmV20250611Client` + `DescribeVpcAclRules`（返回 Rules 1元素含各字段），SetId 后调 Read，断言各字段回读值。
  - `TestFwmVpcAclRule_Update`：mock `UseFwmV20250611Client` + `ModifyVpcAclRuleWithContext` + `DescribeVpcAclRules`，调 Update 断言无错。
  - `TestFwmVpcAclRule_Delete`：mock `UseFwmV20250611Client` + `DeleteVpcAclRuleWithContext`，调 Delete 断言无错。
  - `TestFwmVpcAclRule_Schema`：断言 schema 包含所有字段、forceNew 字段标志正确。
  - **禁止通过 `go test` 执行测试用例**，仅保证代码可编译。

## 7. Validation

- [x] 7.1 `gofmt` 格式化（收尾阶段由 tfpacer-finalize 处理）。
- [x] 7.2 `openspec validate add-fwm-vpc-acl-rule-resource --strict` 通过。