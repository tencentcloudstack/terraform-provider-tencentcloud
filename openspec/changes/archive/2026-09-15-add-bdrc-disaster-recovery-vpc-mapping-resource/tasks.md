# Tasks for `add-bdrc-disaster-recovery-vpc-mapping-resource`

## 1. Connectivity layer（新增 BDRC 客户端）

- [x] 1.1 在 `tencentcloud/connectivity/client.go` 中新增 import `bdrcv20260330 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330"`（放在已有 SDK import 段，按字母序紧邻其他 `tencentcloud-sdk-go/tencentcloud/...` import）。
- [x] 1.2 在 `TencentCloudClient` struct 中新增字段 `bdrcV20260330Conn *bdrcv20260330.Client`（参考 `igtmv20231024Conn` 字段位置与命名风格）。
- [x] 1.3 新增方法 `UseBdrcV20260330Client() *bdrcv20260330.Client`，实现完全参考 `UseIgtmV20231024Client`：nil 缓存检查 → `me.NewClientProfile(300)` → `bdrcv20260330.NewClient(me.Credential, me.Region, cpf)` → `WithHttpTransport(&LogRoundTripper{})` → 返回；注释行 `// UseBdrcV20260330Client return BDRC client for service`。

## 2. Resource implementation（新增资源主文件）

- [x] 2.1 创建 `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_vpc_mapping.go`，package `bdrc`；文件开头不加注释；import 段参考 `resource_tc_igtm_strategy.go`：`context`、`fmt`、`log`、`strconv`、`strings`、terraform plugin sdk v2 helpers、`bdrcv20260330`、`tccommon`、`helper`。
- [x] 2.2 声明 `ResourceTencentCloudBdrcDisasterRecoveryVpcMapping() *schema.Resource`，含 `Create/Read/Update/Delete` 回调与 `Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}`（支持 import）。
- [x] 2.3 Schema 定义（严格按 spec 的字段表）：5 个 Required+ForceNew 输入字段（`site_pair_id`、`source_vpc_id`、`source_subnet_id`、`target_vpc_id`、`target_subnet_id`，TypeString），7 个 Computed 输出字段（`id` TypeInt、`source_vpc`/`source_subnet`/`target_vpc`/`target_subnet`/`status`/`life_state` TypeString）；禁止出现 `vpc_mapping_set` 嵌套层；每个字段带 Description。
- [x] 2.4 实现 `resourceTencentCloudBdrcDisasterRecoveryVpcMappingCreate`：
  - `defer tccommon.LogElapsed("resource.tencentcloud_bdrc_disaster_recovery_vpc_mapping.create")()` 与 `defer tccommon.InconsistentCheck(d, meta)()`。
  - `var (logId, ctx, request, sitePairId, sourceVpcId, sourceSubnetId, targetVpcId, targetSubnetId)` 块。
  - 从 schema 读取 5 个输入字段填入 `CreateDisasterRecoveryVpcMappingRequest`（`helper.String(...)`）；同时保存到局部变量供回查使用。
  - `resource.Retry(tccommon.WriteRetryTimeout, ...)` 调用 `CreateDisasterRecoveryVpcMappingWithContext`；错误 `tccommon.RetryError(e)`；成功校验 `result == nil || result.Response == nil` → `resource.NonRetryableError(fmt.Errorf(...))`；日志用 `[DEBUG]` 含 logId、`request.GetAction()`、`request.ToJsonString()`、`result.ToJsonString()`。
  - retry 块外、错误处理后：再用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 调用 `DescribeVpcMappingsWithContext` 回查：
    - `request2 := bdrcv20260330.NewDescribeVpcMappingsRequest()`；`SitePairId = helper.String(sitePairId)`；`Limit = helper.IntUint64(100)`；
    - `Filters` 填两项：`{Name: "source-vpc-id", Values: []*string{&sourceVpcId}}`、`{Name: "source-subnet-id", Values: []*string{&sourceSubnetId}}`；
    - retry 内：SDK 错误 `tccommon.RetryError(e)`；`result == nil || result.Response == nil` → `resource.NonRetryableError(...)`；`len(result.Response.VpcMappingSet) == 0` → `resource.RetryableError(fmt.Errorf("bdrc_disaster_recovery_vpc_mapping not found yet, waiting for eventual consistency"))`（可重试，等待最终一致）；
    - 在 `VpcMappingSet` 中匹配 `SourceVpc==sourceVpcId && SourceSubnet==sourceSubnetId && TargetVpc==targetVpcId && TargetSubnet==targetSubnetId` 的项；匹配项 `Id == nil` → `resource.NonRetryableError(...)`；匹配项数 > 1 → `resource.NonRetryableError(...)`（人工介入）；命中唯一项取 `vpcMappingIdStr = strconv.FormatUint(*matched.Id, 10)` 后 `return nil`。
  - 第二个 retry 块外、错误处理后：`log.Printf("[CRITICAL]%s create bdrc_disaster_recovery_vpc_mapping success, vpcMappingId=%s", logId, vpcMappingIdStr)`（注意本项目沿用 `CRITAL` 拼写，与 igtm 一致）；`d.SetId(strings.Join([]string{sitePairId, vpcMappingIdStr}, tccommon.FILED_SP))`；`return resourceTencentCloudBdrcDisasterRecoveryVpcMappingRead(d, meta)`。
- [x] 2.5 实现 `resourceTencentCloudBdrcDisasterRecoveryVpcMappingRead`：
  - `defer tccommon.LogElapsed(...)` 与 `defer tccommon.InconsistentCheck(d, meta)()`；`var (logId, ctx)` 块。
  - `idSplit := strings.Split(d.Id(), tccommon.FILED_SP)`；`len(idSplit) != 2` → `return fmt.Errorf("id is broken,%s", d.Id())`；`sitePairId := idSplit[0]`、`vpcMappingId := idSplit[1]`。
  - `resource.Retry(tccommon.ReadRetryTimeout, ...)` 调用 `DescribeVpcMappingsWithContext`：`SitePairId = helper.String(sitePairId)`、`Limit = helper.IntUint64(100)`（不使用 Filters，全量后按 Id 匹配）；
    - retry 内：SDK 错误 `tccommon.RetryError(e)`；`result == nil || result.Response == nil` → `resource.NonRetryableError(...)`；
    - 遍历 `VpcMappingSet`，`strconv.FormatUint(*item.Id, 10) == vpcMappingId` 命中即赋值给外层 `matched` 变量并 `return nil`；遍历完未命中 `return nil`（让外层按"未找到"处理，而非在 retry 内立即 NonRetryable——因为列表可能分页，但 Limit=100 已覆盖单站点对常见规模；若未来超 100 条需改分页，此处保持简单）。
  - retry 块外：若 `matched == nil`：**先** `log.Printf("[CRUD] bdrc_disaster_recovery_vpc_mapping id=%s", d.Id())`，**再** `d.SetId("")`，`return nil`；
  - 若 `matched != nil`：逐字段判 nil 后 `d.Set(...)`：`id`→`*matched.Id`、`site_pair_id`→`*matched.SitePairId`、`source_vpc`→`*matched.SourceVpc`、`source_subnet`→`*matched.SourceSubnet`、`target_vpc`→`*matched.TargetVpc`、`target_subnet`→`*matched.TargetSubnet`、`status`→`*matched.Status`、`life_state`→`*matched.LifeState`；每个 set 用 `_ = d.Set(...)` 忽略 error（避免未使用变量）。
- [x] 2.6 实现 `resourceTencentCloudBdrcDisasterRecoveryVpcMappingUpdate`：
  - `defer tccommon.LogElapsed(...)` 与 `defer tccommon.InconsistentCheck(d, meta)()`。
  - `immutableArgs := []string{"site_pair_id", "source_vpc_id", "source_subnet_id", "target_vpc_id", "target_subnet_id"}`；遍历 `d.HasChange(v)`，命中即 `return fmt.Errorf("Update bdrc_disaster_recovery_vpc_mapping is not supported, all business arguments are immutable (CRD-only API), please recreate the resource.")`。
  - 未命中则 `return resourceTencentCloudBdrcDisasterRecoveryVpcMappingRead(d, meta)`；**不调用任何 BDRC SDK 方法**。
- [x] 2.7 实现 `resourceTencentCloudBdrcDisasterRecoveryVpcMappingDelete`：
  - `defer tccommon.LogElapsed(...)` 与 `defer tccommon.InconsistentCheck(d, meta)()`；`var (logId, ctx, request)` 块。
  - 拆分联合 ID 得到 `sitePairId`、`vpcMappingId`；`vpcMappingIdPoint := helper.StrToUint64Point(vpcMappingId)`。
  - `request := bdrcv20260330.NewDeleteDisasterRecoveryVpcMappingRequest()`；`request.VpcMappingIds = []*uint64{vpcMappingIdPoint}`（严格单元素）。
  - `resource.Retry(tccommon.WriteRetryTimeout, ...)` 调用 `DeleteDisasterRecoveryVpcMappingWithContext`；错误 `tccommon.RetryError(e)`；成功校验 `result == nil || result.Response == nil` → `resource.NonRetryableError(...)`；`[DEBUG]` 日志。
  - retry 块外、错误处理后 `return nil`（不做删除后回查）。

## 3. Provider registration

- [x] 3.1 在 `tencentcloud/provider.go` 的 import 段新增 `"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/bdrc"`（按字母序紧邻 `as`/`bh` 等服务 import）。
- [x] 3.2 在 `tencentcloud/provider.go` 的资源注册 map 中追加 `"tencentcloud_bdrc_disaster_recovery_vpc_mapping": bdrc.ResourceTencentCloudBdrcDisasterRecoveryVpcMapping()`（紧邻其他 `bdrc.Resource...` 或在字母序合适位置；当前无其他 bdrc 资源，放在合理位置即可）。
- [x] 3.3 在 `tencentcloud/provider.md` 中新增 BDRC 产品段并登记 `tencentcloud_bdrc_disaster_recovery_vpc_mapping` 资源名一行，使 gendoc 能扫描到（参考其他产品段格式，如 Bastion Host(BH)）。

## 4. Documentation

- [x] 4.1 创建 `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_vpc_mapping.md`，内容包含：
  - 一句话描述（带 BDRC 产品名）："Provides a resource to create a disaster recovery VPC mapping for site pair of BDRC."
  - `Example Usage` HCL 块：包含 5 个输入字段（`site_pair_id`、`source_vpc_id`、`source_subnet_id`、`target_vpc_id`、`target_subnet_id`）。
  - `Import` 部分（RESOURCE_KIND_GENERAL 有 Import）：说明使用联合 ID `sitePairId#vpcMappingId`，给出 `terraform import` 示例命令。
  - 不手写 `Argument Reference` / `Attribute Reference`（由工具自动生成）。
- [x] 4.2 `website/docs/r/bdrc_disaster_recovery_vpc_mapping.html.markdown` 由收尾阶段 `make doc` 生成，本阶段不手动创建。

## 5. Unit test

- [x] 5.1 创建 `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_vpc_mapping_test.go`，package `bdrc_test`；使用 `gomonkey` mock 三个 BDRC SDK 方法（`CreateDisasterRecoveryVpcMappingWithContext`、`DescribeVpcMappingsWithContext`、`DeleteDisasterRecoveryVpcMappingWithContext`），不使用 terraform 测试套件。
- [x] 5.2 测试用例覆盖：
  - Create 成功路径：mock Create 返回仅 RequestId 的响应，mock Describe 返回含匹配 `VpcMapping`（`Id` 非空、四字段匹配）的响应，断言 `d.Id()` 为 `sitePairId#vpcMappingId` 联合格式、computed 字段已 set。
  - Read 未找到路径：mock Describe 返回空 `VpcMappingSet`，断言 `d.Id()` 被清空（`""`）。
  - Delete 路径：mock Delete 返回成功响应，断言 `VpcMappingIds` 为单元素切片。
  - Update 拒绝路径：构造 `d.HasChange("target_vpc_id")` 场景，断言返回 error 且未调用任何 SDK 方法。
- [x] 5.3 保证生成的测试代码在当前环境下可正确构建（import 路径、mock 目标函数签名正确），但不通过 `go test` 执行。

## 6. Validation（代码正确性检查，禁止 go build/vet/lint）

- [x] 6.1 人工核对：Create 入参 5 个字段均在 `CreateDisasterRecoveryVpcMappingRequest` 中存在；Describe 入参 `SitePairId`/`Filters`/`Limit` 均在 `DescribeVpcMappingsRequest` 中存在；Delete 入参 `VpcMappingIds` 在 `DeleteDisasterRecoveryVpcMappingRequest` 中存在。
- [x] 6.2 人工核对：所有函数返回的 error 均被检查；必定不出错的函数用 `_ = func()` 赋值给 `_`。
- [x] 6.3 人工核对：资源 Read 中 set 前判 nil；云 API 返回空时先 `log.Printf` 保留 id 再 `d.SetId("")`；Create 完成后检查返回值是否为空、Id 是否为空，空则返回 `NonRetryableError`。
- [x] 6.4 人工核对：retry 块内只做接口调用与空值校验，`d.SetId` 等成功操作在 retry 块外；错误用 `tccommon.RetryError()` 包装。
- [x] 6.5 人工核对：`.md` 文档格式符合 gendoc/README.md 与其他资源 `.md` 样式；provider.go/provider.md 注册条目完整。
