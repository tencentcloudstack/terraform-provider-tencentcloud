## 1. connectivity 层接入 fwm 客户端

- [x] 1.1 在 `tencentcloud/connectivity/client.go` import 区（约第 46 行 cfw 导入旁）新增 `fwm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`
- [x] 1.2 在 `TencentCloudClient` 结构体字段区（约第 230 行 cfwConn 旁）新增 `fwmConn *fwm.Client` 字段
- [x] 1.3 在方法区（约第 1771 行 UseCfwClient 旁）新增 `UseFwmClient() *fwm.Client` 方法（参考 UseCfwClient 实现）

## 2. 服务层实现

- [x] 2.1 创建 `tencentcloud/services/fwm/service_tencentcloud_fwm.go`，定义 `FwmService` 结构体与 `NewFwmService` 构造函数（参考 CfwService 模式）
- [x] 2.2 实现 `DescribeFwmCfwManagedByAppId(ctx, appId, memberId)` 方法：用 `Filters=[{Name:"AppId",Values:[appId],OperatorType:1}]`、Limit 较大值（如 100）调用 DescribeOrganMembers，从返回 Members 按 MemberId 匹配返回目标成员对象（*OrganMemberItem）；用 `resource.Retry(tccommon.ReadRetryTimeout,...)` 包裹，失败用 `tccommon.RetryError` 包装
- [x] 2.3 实现 `ModifyFwmCfwManaged(ctx, memberId, cfwManaged)` 方法：构建 ModifyCfwManagedRequest 调用云 API，用 `resource.Retry(tccommon.WriteRetryTimeout,...)` 包裹，失败用 `tccommon.RetryError` 包装

## 3. 资源核心实现

- [x] 3.1 创建 `tencentcloud/services/fwm/resource_tc_fwm_cfw_managed_config.go`，定义 `ResourceTencentCloudFwmCfwManaged()` 返回 `*schema.Resource`，包含 Create/Read/Update/Delete 与 Importer
- [x] 3.2 定义 Schema：`member_id`(TypeString/Required)、`app_id`(TypeString/Required/ForceNew)、`cfw_managed`(TypeInt/Required)，description 英文
- [x] 3.3 实现 Create（empty）：读取 member_id → `d.SetId(member_id)` → 复用 `resourceTencentCloudFwmCfwManagedUpdate(d, meta)`
- [x] 3.4 实现 Read：从 `d.Id()` split FILED_SP 解析 app_id/member_id；调 service.DescribeFwmCfwManagedByAppId；nil/匹配不到 → 打印日志后 `d.SetId("")` 不报错；set 字段前判 nil
- [x] 3.5 实现 Update：比较 cfw_managed 与当前值，相同跳过，不同调 service.ModifyFwmCfwManaged；用 retry 包裹；不生成轮询；错误透传
- [x] 3.6 实现 Delete（empty）：`return nil`，不调云 API
- [x] 3.7 实现 Importer：`schema.ImportStatePassthrough`（ID 为 app_id#member_id）

## 4. Provider 注册与文档

- [x] 4.1 在 `tencentcloud/provider.go` ResourcesMap 中注册 `"tencentcloud_fwm_cfw_managed": fwm.ResourceTencentCloudFwmCfwManaged()`，并添加 fwm 包 import（若 provider.go 未见 fwm 导入则新增 `fwm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"`，参考其他 services import 风格）
- [x] 4.2 在 `tencentcloud/provider.md` 中新增 fwm 产品章节与资源条目（参考 cfw 章节格式）
- [x] 4.3 创建 `tencentcloud/services/fwm/resource_tc_fwm_cfw_managed_config.md` 文档样例：一句话描述（带 fwm）、Example Usage、Import 部分（说明使用 app_id#member_id 两段式 ID）；不含 Argument Reference / Attribute Reference

## 5. 单元测试

- [x] 5.1 创建 `tencentcloud/services/fwm/resource_tc_fwm_cfw_managed_config_test.go`，使用 mock（gomonkey）方法，不使用 terraform 验收测试套件
- [x] 5.2 mock UseFwmClient / DescribeOrganMembers / ModifyCfwManaged，覆盖 Create、Read（成员存在/不存在）、Update（cfw_managed 变更/未变更）、Delete 业务逻辑
- [x] 5.3 确保所有函数返回的 error 被检查，必定不出错的用 `_ = func()` 赋值给 _，代码可正确构建执行

## 6. 代码正确性检查

- [x] 6.1 核对 DescribeOrganMembers 入参（Filters/Limit/Offset/By/Order）与 ModifyCfwManaged 入参（MemberId/CfwManaged）在 vendor 中的定义一致
- [x] 6.2 核对 cfw_managed 为 int（*int64），非 bool；MemberId JSON tag 为 "MemberId"
- [x] 6.3 核对 Read 中不传 Filters[{Name:"MemberId"}]，只用 AppId 过滤
- [x] 6.4 核对 Update 同步接口无轮询，Create 在 SetId 后再复用 Update