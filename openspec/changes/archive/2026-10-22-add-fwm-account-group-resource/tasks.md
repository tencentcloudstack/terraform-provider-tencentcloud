## 1. FWM 连接层接入

- [x] 1.1 在 `tencentcloud/connectivity/client.go` import 区新增 `fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`（字母序位于 es 与 ga2 之间的 fwm 段）
- [x] 1.2 在 `TencentCloudClient` struct 新增字段 `fwmv20250611Conn *fwmv20250611.Client`
- [x] 1.3 新增方法 `UseFwmV20250611Client() *fwmv20250611.Client`（参照 UseIgtmV20231024Client：缓存、NewClientProfile(300)、WithHttpTransport(&LogRoundTripper{})）

## 2. FWM 服务层

- [x] 2.1 创建 `tencentcloud/services/fwm/service_tencentcloud_fwm.go`，定义 `FwmService` struct（含 `client *connectivity.TencentCloudClient`）与构造函数
- [x] 2.2 实现 `DescribeFwmAccountGroupById(ctx, groupId) (*fwmv20250611.AccountGroup, error)`：构建 Filters=[{Name:"GroupId",Values:[groupId],OperatorType:1}]、Limit=20、Offset=0，Retry(ReadRetryTimeout) 调 DescribeAccountGroupsWithContext；AccountGroups 为空返回 nil,nil（视为删除）；否则返回首元素
- [x] 2.3 实现 `DeleteFwmAccountGroupById(ctx, groupId) error`：构建 DeleteAccountGroupRequest(GroupId)，Retry(WriteRetryTimeout) 调 DeleteAccountGroupWithContext；InvalidParameter 业务拒绝原样透传

## 3. 资源文件

- [x] 3.1 创建 `tencentcloud/services/fwm/resource_tc_fwm_account_group.go`，包名 `fwm`，import fwmv20250611 SDK
- [x] 3.2 定义 `ResourceTencentCloudFwmAccountGroup() *schema.Resource`，含 Create/Read/Update/Delete + Importer(ImportStatePassthrough)，Schema：group_name(Required) / description(Optional,ForceNew) / group_id,member_count,can_delete,create_time(Computed)
- [x] 3.3 实现 `resourceTencentCloudFwmAccountGroupCreate`：取 group_name/description（Description 用 d.GetOk 判断），Retry(WriteRetryTimeout) 调 CreateAccountGroupWithContext，校验 Response.GroupId 非空，d.SetId(groupId)，调 Read
- [x] 3.4 实现 `resourceTencentCloudFwmAccountGroupRead`：d.Id() 为 GroupId，服务层 DescribeFwmAccountGroupById；nil 时 `log.Printf("[CRUD] fwm_account_group id=%s", d.Id())` 后 d.SetId("")；回填各字段前判 nil；description 空值保护
- [x] 3.5 实现 `resourceTencentCloudFwmAccountGroupUpdate`：mutableArgs=["group_name"]；HasChange 时构建 ModifyAccountGroupRequest(GroupId+GroupName)，Retry(WriteRetryTimeout) 调 ModifyAccountGroupWithContext；description 变更因 ForceNew 不走 Update；末尾调 Read
- [x] 3.6 实现 `resourceTencentCloudFwmAccountGroupDelete`：d.Id() 为 GroupId，服务层 DeleteFwmAccountGroupById；错误透传

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/fwm/resource_tc_fwm_account_group_test.go`，package `fwm_test`，使用 gomonkey mock（参考 resource_tc_cfw_cluster_fw_bypass_config_test.go 的 mockMeta 模式）
- [x] 4.2 mock `UseFwmV20250611Client` 及四个 WithContext 方法，覆盖 create/read/update/delete 业务逻辑（含 read 空列表=删除场景、create 空 GroupId 返回错误场景）

## 5. Provider 注册

- [x] 5.1 在 `tencentcloud/provider.go` import 区新增 `"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"`
- [x] 5.2 在 ResourcesMap 新增 `"tencentcloud_fwm_account_group": fwm.ResourceTencentCloudFwmAccountGroup(),`

## 6. 文档

- [x] 6.1 创建 `tencentcloud/services/fwm/resource_tc_fwm_account_group.md`：一句话描述（带 FWM 产品名）、Example Usage（group_name=test-account-group、description=test account group）、Import（直接写 GroupId acg-xxxxxxxx）；不添加 Argument/Attribute Reference

## 7. 验证（收尾阶段执行）

- [x] 7.1 gofmt 格式化新生成的 Go 文件（由 tfpacer-finalize skill 统一执行）
- [x] 7.2 通过 `make doc` 生成 website/docs 文档（由 tfpacer-finalize skill 统一执行）
- [x] 7.3 代码正确性检查：确认 Create 参数在 CreateAccountGroup 中存在、Update 参数在 ModifyAccountGroup 中存在、各 CRUD nil 检查完备