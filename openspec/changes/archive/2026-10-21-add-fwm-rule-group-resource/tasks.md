# Tasks for `add-fwm-rule-group-resource`

> 前置：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611` 已 vendor（`vendor/modules.txt` 已包含），本 change 只需新增 connectivity 客户端工厂，无需升级 SDK。

## 1. Service layer & client factory

- [x] 1.1 在 `tencentcloud/connectivity/client.go` 的 `TencentCloudClient` 结构体中（紧跟 `ga2v20250115Conn` 之后）新增字段 `fwmV20250611Conn *fwmv20250611.Client`，并在 import 块按字母序新增 `fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`。
- [x] 1.2 在 `tencentcloud/connectivity/client.go` 新增 `UseFwmV20250611Client()`（照抄 `UseIgtmV20231024Client` 2247-2257 行写法：缓存判断 + `NewClientProfile(300)` + `fwmv20250611.NewClient(me.Credential, me.Region, cpf)` + `WithHttpTransport(&LogRoundTripper{})`）。
- [x] 1.3 新建 `tencentcloud/services/fwm/service_tencentcloud_fwm.go`：`package fwm`，`NewFwmService(client *connectivity.TencentCloudClient) FwmService`、`type FwmService struct { client *connectivity.TencentCloudClient }`。
- [x] 1.4 在 service 文件实现 `DescribeFwmRuleGroupById(ctx context.Context, groupId string) (ret *fwmv20250611.RuleGroupResp, errRet error)`：
  - `defer` 中打印 `[CRITAL]%s api[%s] fail, request body [%s], reason[%s]`；
  - 请求 `DescribeRuleGroupsRequest`：`Filters = []*fwmv20250611.CommonFilter{{Name: helper.String("GroupId"), Values: []*string{helper.String(groupId)}, OperatorType: helper.Int64(9)}}`，`Offset` 从 0 开始、`Limit = 100`；
  - 循环分页（`Offset += Limit`，直到返回条数 `< Limit` 或累计 `>= TotalCount`），每轮 `ratelimit.Check(request.GetAction())` + `resource.Retry(tccommon.ReadRetryTimeout, ...)`，`result == nil || result.Response == nil` 时 `resource.NonRetryableError`；
  - 遍历 `Response.RuleGroups`，按 `*item.GroupId == groupId` 精确匹配返回，未命中返回 `(nil, nil)`。

## 2. Resource implementation

- [x] 2.1 新建 `tencentcloud/services/fwm/resource_tc_fwm_rule_group.go`（不加文件头注释），imports 按 `resource_tc_igtm_strategy.go` 风格：`context`/`fmt`/`log`、`helper/resource`、`helper/schema`、`helper/validation`、`fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`、`tccommon`、`internal/helper`。
- [x] 2.2 声明包级变量 `fwmRuleGroupProductToName = map[string]string{"1": "enterprise_sg", "2": "cfw_nat_acl", "3": "cfw_vpc_acl", "4": "cfw_edge_acl"}` 与 `fwmRuleGroupNormalizeProduct(product string) string`（未命中原样返回）。
- [x] 2.3 实现 `ResourceTencentCloudFwmRuleGroup() *schema.Resource`：`Create/Read/Update/Delete` + `Importer{schema.ImportStatePassthrough}`；Schema 严格按 design.md D1/D2 声明 `group_name`（Required）、`product`（Required + ForceNew + `validation.StringInSlice([]string{"enterprise_sg","cfw_nat_acl","cfw_vpc_acl","cfw_edge_acl"}, false)`）与 12 个 Computed 字段（`group_id`、`app_id`、`rule_count`、`create_time`、`update_time`、`create_by`、`create_name`、`update_by`、`update_name`、`receive_account`、`exec_areas`、`direction`）。
- [x] 2.4 实现 `resourceTencentCloudFwmRuleGroupCreate`：
  - 首两行 `defer tccommon.LogElapsed("resource.tencentcloud_fwm_rule_group.create")()` / `defer tccommon.InconsistentCheck(d, meta)()`；
  - 填充 `GroupName`、`Product`（`d.GetOk`）；`resource.Retry(tccommon.WriteRetryTimeout, ...)` 调 `UseFwmV20250611Client().CreateRuleGroupWithContext(ctx, request)`，错误用 `tccommon.RetryError(e)` 包装，成功打印 `[DEBUG]...api[%s] success...`；
  - 校验 `result == nil || result.Response == nil` → `resource.NonRetryableError`；校验 `result.Response.GroupId == nil || *result.Response.GroupId == ""` → `resource.NonRetryableError`；
  - `d.SetId(*result.Response.GroupId)`，随后 `return resourceTencentCloudFwmRuleGroupRead(d, meta)`。
- [x] 2.5 实现 `resourceTencentCloudFwmRuleGroupRead`：
  - `service := FwmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}`，`respData, err := service.DescribeFwmRuleGroupById(ctx, d.Id())`；
  - `respData == nil` 时先 `log.Printf("[CRUD] tencentcloud_fwm_rule_group id=%s", d.Id())` 再 `d.SetId("")` 并返回 nil；
  - 逐个字段 nil 判断后 `d.Set`：`group_id`、`group_name`、`product`（走 `fwmRuleGroupNormalizeProduct`）、`app_id`、`rule_count`、`create_time`、`update_time`、`create_by`、`create_name`、`update_by`、`update_name`、`exec_areas`、`direction`；
  - `receive_account` 与其中 `members` 转为 `[]map[string]interface{}` 后 `d.Set`，两个层级都要判 nil / len > 0。
- [x] 2.6 实现 `resourceTencentCloudFwmRuleGroupUpdate`：
  - `immutableArgs := []string{"product"}`，命中 `d.HasChange` 则 `return fmt.Errorf("tencentcloud_fwm_rule_group: argument `%s` is immutable, please recreate the resource", v)`；
  - `mutableArgs := []string{"group_name"}`，命中变更后构造 `NewModifyRuleGroupRequest()`，`GroupId = helper.String(d.Id())`、`GroupName` 取自 `d.GetOk("group_name")`，`resource.Retry(tccommon.WriteRetryTimeout, ...)` 调 `ModifyRuleGroupWithContext`；
  - 末尾 `return resourceTencentCloudFwmRuleGroupRead(d, meta)`。
- [x] 2.7 实现 `resourceTencentCloudFwmRuleGroupDelete`：构造 `NewDeleteRuleGroupRequest()`，`GroupIds = []*string{helper.String(d.Id())}`，`resource.Retry(tccommon.WriteRetryTimeout, ...)` 调 `DeleteRuleGroupWithContext`，失败打印 `[CRITAL]%s delete fwm rule group failed, reason:%+v` 后返回。
- [x] 2.8 全量自检：所有返回的 error 均被检查；无未使用变量；`d.Set` 的 error 统一 `_ =`。

## 3. Provider registration

- [x] 3.1 `tencentcloud/provider.go` import 块新增 `"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"`。
- [x] 3.2 `tencentcloud/provider.go` 的 `ResourcesMap` 中新增 `"tencentcloud_fwm_rule_group": fwm.ResourceTencentCloudFwmRuleGroup(),`（与 `tencentcloud_igtm_strategy` 同层，只加一行）。
- [x] 3.3 `tencentcloud/provider.md` 末尾追加 `Firewall Manager(FWM)` / `Resource` / `tencentcloud_fwm_rule_group` 三行。

## 4. Documentation

- [x] 4.1 新建 `tencentcloud/services/fwm/resource_tc_fwm_rule_group.md`：首行 `Provides a resource to create a FWM rule group.`；`Example Usage` 给出最小 HCL（`group_name = "tf-example"`、`product = "enterprise_sg"`）；`Import` 段给出 `terraform import tencentcloud_fwm_rule_group.example fwmrg-wsrrhuvrib`；不写 `Argument Reference` / `Attribute Reference`。

## 5. Unit tests (gomonkey, 不使用 terraform 测试套件)

- [x] 5.1 新建 `tencentcloud/services/fwm/resource_tc_fwm_rule_group_test.go`（`package fwm_test`），参照 `resource_tc_bh_bind_device_resource_test.go`：`mockMetaFwmRuleGroup` 实现 `tccommon.ProviderMeta`，`newMockMetaFwmRuleGroup()` 返回空 `&connectivity.TencentCloudClient{}`，`patches.ApplyMethodReturn(meta.client, "UseFwmV20250611Client", &fwmv20250611.Client{})`。
- [x] 5.2 `TestFwmRuleGroupResource_Create`：mock `CreateRuleGroupWithContext` 断言 `GroupName`/`Product` 并返回 `GroupId="fwmrg-wsrrhuvrib"`，同时 mock `DescribeRuleGroupsWithContext` 返回同名记录；断言 `d.Id()` 与 `group_id`。
- [x] 5.3 `TestFwmRuleGroupResource_ReadProductNormalize`：Describe 返回 `Product="2"`、`Direction=-1`、`ExecAreas`、`ReceiveAccount`（含 `Members`），断言 state `product == "cfw_nat_acl"` 及嵌套字段落盘。
- [x] 5.4 `TestFwmRuleGroupResource_ReadNotFound`：Describe 返回空 `RuleGroups`，断言 `d.Id() == ""` 且 `err == nil`。
- [x] 5.5 `TestFwmRuleGroupResource_Update`：`group_name` 变更，mock `ModifyRuleGroupWithContext` 断言 `GroupId == d.Id()` 且 `GroupName` 为新值。
- [x] 5.6 `TestFwmRuleGroupResource_Delete`：mock `DeleteRuleGroupWithContext` 断言 `len(GroupIds) == 1 && *GroupIds[0] == d.Id()`。
- [x] 5.7 自检测试文件可编译：所有 error 均检查、无未使用 import/变量（禁止执行 `go test`）。

## 6. Validation（由后续流程统一执行，不在本阶段运行构建命令）

- [x] 6.1 确认新增/修改文件均已落盘，且 `website/` 未被手动改动（文档由收尾阶段 `make doc` 生成）。
- [x] 6.2 确认 `.changelog/` 未在本阶段新增文件（由收尾阶段 `tfpacer-finalize` 生成）。
- [x] 6.3 `pnpm dlx @fission-ai/openspec@1.2.0 validate add-fwm-rule-group-resource --strict` 通过。
