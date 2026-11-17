## 1. FWM 客户端注册

- [ ] 1.1 在 `tencentcloud/connectivity/client.go` import 段新增 `fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`（参照 igtmv20231024 导入位置）。
- [ ] 1.2 在 `TencentCloudClient` struct 新增字段 `fwmv20250611Conn *fwmv20250611.Client`（参照 `igtmv20231024Conn` 字段位置）。
- [ ] 1.3 新增访问器 `func (me *TencentCloudClient) UseFwmV20250611Client() *fwmv20250611.Client`，懒初始化（`NewClient(me.Credential, me.Region, cpf)` + `WithHttpTransport(&LogRoundTripper{})`），完全参照 `UseIgtmV20231024Client` 写法。
- [ ] 1.4 确认 `vendor/modules.txt` 与 `go.mod` 已纳入 `tencentcloud-sdk-go/tencentcloud/fwm/v20250611`（工作区已有 vendor，必要时按 go-module-vendor skill 补全 vendor 内容；不执行 go build）。

## 2. FWM service 层

- [ ] 2.1 创建 `tencentcloud/services/fwm/service_tencentcloud_fwm.go`，package `fwm`，import 参照 `service_tencentcloud_igtm.go`（context、fmt、log、terraform-plugin-sdk resource、fwmv20250611、tccommon、connectivity、helper、ratelimit）。
- [ ] 2.2 实现 `NewFwmService` 构造器与 `FwmService` struct（`client *connectivity.TencentCloudClient`）。
- [ ] 2.3 实现 `DescribeFwmSecurityGroupRuleById(ctx, groupId, ruleId string) (ret *fwmv20250611.SecGroupRuleResp, errRet error)`：构造 `DescribeSecurityGroupRulesRequest`，`GroupId`、`Filters=[{Name:"RuleId",Values:[ruleId],OperatorType:helper.Int64(1)}]`、`Limit=helper.Int64(100)`；包 `resource.Retry(tccommon.ReadRetryTimeout)`，retry 内 `ratelimit.Check` + 调 `UseFwmV20250611Client().DescribeSecurityGroupRulesWithContext(ctx, request)` + `tccommon.RetryError(e)`；retry 内空响应（`result==nil||result.Response==nil`）返回 `NonRetryableError`；retry 成功后若 `len(result.Response.Rules)==0` 返回 `ret=nil, errRet=nil`（让资源层判空后 SetId）；否则 `ret = result.Response.Rules[0]`。
- [ ] 2.4 实现 `DeleteFwmSecurityGroupRuleById(ctx, groupId, ruleId string) error`：构造 `DeleteSecurityGroupRuleRequest`，`GroupId`、`RuleIds=[ruleId]`（1 元素）；包 `resource.Retry(tccommon.WriteRetryTimeout)`，遇 `ResourceNotFound` 错误码返回 nil（幂等），其余 `tccommon.RetryError(e)`。

## 3. 资源 CRUD 实现

- [ ] 3.1 创建 `tencentcloud/services/fwm/resource_tc_fwm_security_group_rule.go`，package `fwm`，import 参照 `resource_tc_igtm_strategy.go`（context、fmt、log、strings、terraform-plugin-sdk schema/resource、fwmv20250611、tccommon、helper）。
- [ ] 3.2 实现 `ResourceTencentCloudFwmSecurityGroupRule() *schema.Resource`：声明 `Create/Read/Update/Delete` + `Importer{State: schema.ImportStatePassthrough}`；按 design D1 表定义全部顶层 schema 字段（writable 输入字段 + Computed 只读展示字段）。`group_id` Required；`ip_version` Required+ForceNew；`order_index` Required(TypeInt)；`source_type`/`dest_type` Required+ForceNew；`source_content`/`dest_content`/`protocol`/`port`/`rule_action` Required；`scope` Optional+ForceNew；`service_template_id` Optional+Computed；`description` Optional；`protocol_port_type` Optional(TypeInt)；`belong_member_id` Optional。`rule_id` 及其余只读字段全部 Computed。description 中文一句话，字段 description 英文。
- [ ] 3.3 实现 `resourceTencentCloudFwmSecurityGroupRuleCreate`：`defer LogElapsed/InconsistentCheck`；`var (logId, ctx)`；构造 `CreateSecurityGroupRuleRequest`，`GroupId` 从 `d.GetOk("group_id")`；构造 1 元素 `Rules []*SecurityGroupRule`，填入全部入参字段（string 用 `helper.String`，`order_index`/`protocol_port_type` 用 `helper.Int64`，bool/int 可选字段用 `d.GetOkExists`）；包 `resource.Retry(tccommon.WriteRetryTimeout)` 调 `CreateSecurityGroupRuleWithContext`，`tccommon.RetryError(e)`；retry 内检查 `result==nil||result.Response==nil||len(result.Response.RuleIds)==0||result.Response.RuleIds[0]==nil||*result.Response.RuleIds[0]==""` → `NonRetryableError`（先打印 logId）；retry 外打印 logId+d.Id()，`ruleId := *response.Response.RuleIds[0]`，`d.SetId(strings.Join([]string{groupId, ruleId}, tccommon.FILED_SP))`；末尾 `return resourceTencentCloudFwmSecurityGroupRuleRead(d, meta)`。不生成轮询。
- [ ] 3.4 实现 `resourceTencentCloudFwmSecurityGroupRuleRead`：`defer LogElapsed/InconsistentCheck`；`var (logId, ctx, service)`；拆 `d.Id()` by `tccommon.FILED_SP` 得 `groupId`/`ruleId`，长度校验；调 `service.DescribeFwmSecurityGroupRuleById(ctx, groupId, ruleId)`；`ret==nil` 时 `log.Printf("[CRUD] fwm_security_group_rule id=%s", d.Id())` 后 `d.SetId("")` 返回 nil；取 `ret`（`*SecGroupRuleResp`）回写：所有 `d.Set` 前 nil 判断；int 枚举反查字符串（`SourceType`→`source_type`、`TargetType`→`dest_type`、`Strategy`→`rule_action`，未知 int 保留十进制字符串 + warn）；`Detail`→`description`、`SourceId`→`source_content`+`source_id`、`TargetId`→`dest_content`+`target_id`、`BelongMember.MemberId` 非空才 `d.Set("belong_member_id")`；只读字段 `rule_id`/`region`/`protocol_port_name`/`sou_*`/`instance_name`/`public_ip`/`private_ip`/`cidr`/`parameter_name`/`strategy`/`detail`/`id`/`is_new`/`all_total_count`/`dns_parse_count`(嵌套) 等逐一回写。
- [ ] 3.5 实现 `resourceTencentCloudFwmSecurityGroupRuleUpdate`：`defer LogElapsed/InconsistentCheck`；拆 id 得 `groupId`/`ruleId`；`mutableArgs := []string{"order_index","source_content","dest_content","protocol","port","service_template_id","rule_action","description","protocol_port_type","belong_member_id"}`（注意 `dest_type`/`source_type`/`ip_version`/`scope` 是 ForceNew 不在此）；`d.HasChange` 任一命中 → 构造 `ModifySecurityGroupRuleRequest`，`GroupId`、`Rule *SecurityGroupRule`（`RuleId=ruleId` + 全部可写字段重传，覆盖语义，bool/int 用 GetOkExists），不传服务端未使用的 `ModifyType`；包 `resource.Retry(tccommon.WriteRetryTimeout)` 调 `ModifySecurityGroupRuleWithContext`，`tccommon.RetryError(e)`，空响应 `NonRetryableError`；末尾 `return resourceTencentCloudFwmSecurityGroupRuleRead(d, meta)`。`group_id` 不在 mutableArgs（不可改 GroupId）。
- [ ] 3.6 实现 `resourceTencentCloudFwmSecurityGroupRuleDelete`：`defer LogElapsed/InconsistentCheck`；拆 id 得 `groupId`/`ruleId`；调 `service.DeleteFwmSecurityGroupRuleById(ctx, groupId, ruleId)`（幂等 ResourceNotFound 已在 service 层处理）；错误返回前 `log.Printf("[CRITAL]%s ...", logId, err)`；返回 nil。

## 4. Provider 注册

- [ ] 4.1 在 `tencentcloud/provider.go` import 段新增 `"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"`。
- [ ] 4.2 在 `tencentcloud/provider.go` 资源注册 map 追加 `"tencentcloud_fwm_security_group_rule": fwm.ResourceTencentCloudFwmSecurityGroupRule(),`（按字母序放置，参照 igtm 段邻近）。
- [ ] 4.3 在 `tencentcloud/provider.md` 新增 Firewall(FWM) Resource 段，追加 `tencentcloud_fwm_security_group_rule` 一行，确保 gendoc 索引识别。

## 5. 文档

- [ ] 5.1 创建 `tencentcloud/services/fwm/resource_tc_fwm_security_group_rule.md`：一句话描述（带 FWM 产品名，格式 "Provides a resource to ..."）；`Example Usage` HCL 块（group_id 引用、ip_version=ipv4、order_index=1、source_type=net、source_content=1.1.1.1、dest_type=net、dest_content=2.2.2.2、protocol=TCP、port=80、rule_action=accept、scope=SG）；`Import` 段说明使用复合 id `terraform import tencentcloud_fwm_security_group_rule.example <group-id>#<rule-id>`；不含 Argument/Attribute Reference。

## 6. 单元测试

- [ ] 6.1 创建 `tencentcloud/services/fwm/resource_tc_fwm_security_group_rule_test.go`，package `fwm_test`，import gomonkey、fwmv20250611、schema、tccommon 等。
- [ ] 6.2 用 gomonkey `ApplyMethod`（或 `ApplyMethodFunc`）mock `fwmv20250611.Client` 的 `CreateSecurityGroupRuleWithContext`、`DescribeSecurityGroupRulesWithContext`、`ModifySecurityGroupRuleWithContext`、`DeleteSecurityGroupRuleWithContext` 四个方法，返回构造好的 response（含 RuleIds、Rules 首元素 SecGroupRuleResp 等）。
- [ ] 6.3 编写测试用例覆盖：Create（断言 d.Id() 复合格式、rule_id 写入）、Read（断言 int 枚举反查为字符串、空响应 SetId）、Update（断言 Modify 请求含全部可写字段、RuleId）、Delete（断言 RuleIds 1 元素、ResourceNotFound 走幂等）。不使用 TF 测试套件，直接调业务函数验证逻辑；确保代码可正确构建执行。

## 7. 验证（由后续流程执行，不在本阶段手动跑）

- [ ] 7.1 由后续流程执行 `gofmt`（收尾阶段 tfpacer-finalize skill 统一执行）。
- [ ] 7.2 由后续流程执行 `make doc` 生成 `website/docs/r/fwm_security_group_rule.html.markdown`（收尾阶段 tfpacer-finalize skill 统一执行）。
- [ ] 7.3 由后续流程生成 `.changelog/` 文件（收尾阶段 tfpacer-finalize skill 统一执行）。
- [ ] 7.4 代码正确性自检：确认 Create 入参字段均在云 API CreateSecurityGroupRule 入参中存在；Modify 入参字段均在 ModifySecurityGroupRule 入参中存在；Describe 出参字段映射与 `SecGroupRuleResp` 定义一致；所有函数返回 error 均检查，不会出错的用 `_ =`。