## 1. 服务客户端注册

- [x] 1.1 在 `tencentcloud/connectivity/client.go` 的 import 块新增 `fwm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`
- [x] 1.2 在 `TencentCloudClient` 结构体新增字段 `fwmConn *fwm.Client`
- [x] 1.3 在 `tencentcloud/connectivity/client.go` 新增 `UseFwmClient()` 方法（仿 `UseCfwClient()`，使用 `NewClientProfile(300)` 惰性创建并缓存）

## 2. 服务层实现

- [x] 2.1 新建 `tencentcloud/services/fwm/service_tencentcloud_fwm.go`，定义 `FwmService{client}` 结构体
- [x] 2.2 实现 `DescribeFwmOrganMemberById(ctx, appId, memberId)` 方法：构造 `DescribeOrganMembersRequest`（`Filters=[{Name:"AppId", Values:[appId], OperatorType:1}]`，Limit 取接口注释提示的大值），retry（ReadRetryTimeout）调用 `DescribeOrganMembersWithContext`，从 `Members` 列表按 MemberId 匹配目标成员返回 `*OrganMemberItem`（无匹配返回 nil 与 nil 错误）

## 3. 资源配置文件实现

- [x] 3.1 新建 `tencentcloud/services/fwm/resource_tc_fwm_cfw_share_role_config.go`，定义 `ResourceTencentCloudFwmCfwShareRole()`：schema 含 `member_id`(required)、`app_id`(required+ForceNew)、`cfw_share_role`(required)、`cfw_sharer_app_id`(computed)，Importer 用 ImportStatePassthrough，资源描述中文带 fwm（参考 cfw_edge_firewall_switch 风格）
- [x] 3.2 实现 `resourceTencentCloudFwmCfwShareRoleConfigCreate`：仅 `d.SetId(app_id + FILED_SP + member_id)` 后转调 Update
- [x] 3.3 实现 `resourceTencentCloudFwmCfwShareRoleConfigRead`：解析 `d.Id()` 得 app_id/member_id，调 service.DescribeFwmOrganMemberById；成员不存在则 `log.Printf("[CRUD] fwm_cfw_share_role id=%s", d.Id())` 后 `d.SetId("")` 返回 nil；存在则回读 cfw_share_role/cfw_sharer_app_id/member_id/app_id（非 nil 才 set）
- [x] 3.4 实现 Update：解析 ID；先 Read 取 current；target=d.Get("cfw_share_role")；current==target 跳过；sharer↔user 直转时先调 ModifyCfwShareRole(none) 再调 ModifyCfwShareRole(target)；其余直接调一次 target。请求仅传 CfwShareRole + UserMemberId（注意字段名 UserMemberId），写调用用 WriteRetryTimeout + RetryError 包装；完成后转 Read。中转步失败返回错误保留原 state。Delete 为 empty 实现 return nil

## 4. provider 注册

- [x] 4.1 在 `tencentcloud/provider.go` 注册 `"tencentcloud_fwm_cfw_share_role": fwm.ResourceTencentCloudFwmCfwShareRole()`（并在 import 块引入 fwm 服务包，别名按现有 cfw 模式）
- [x] 4.2 在 `tencentcloud/provider.md` 新增 `tencentcloud_fwm_cfw_share_role` 资源条目

## 5. 文档样例

- [x] 5.1 新建 `tencentcloud/services/fwm/resource_tc_fwm_cfw_share_role_config.md`：一句话描述（带 fwm）、Example Usage（含 member_id/app_id/cfw_share_role=none）、Import 部分（说明使用 `app_id:member_id` 联合 ID）。不写 Argument Reference / Attribute Reference。

## 6. 单元测试

- [x] 6.1 新建 `tencentcloud/services/fwm/resource_tc_fwm_cfw_share_role_config_test.go`，使用 gomonkey mock 云 API（不使用 terraform 测试套件）
- [x] 6.2 编写测试用例覆盖：Create empty 实现、Read 成员存在/不存在、Update 直接变更、Update 经 none 中转（sharer→user、user→sharer）、Update 目标与当前一致跳过写、Delete empty 实现

## 7. 验证

- [x] 7.1 检查新增/修改的 .go 文件可编译正确（确保 import、字段名、方法签名与 vendor SDK 一致）
- [x] 7.2 校验资源 schema 与云 API 参数映射一致（ModifyCfwShareRole 用 UserMemberId 而非 MemberId；DescribeOrganMembers 过滤用 AppId 而非 MemberId）