## 1. fwm 客户端接入（基座）

- [x] 1.1 在 `tencentcloud/connectivity/client.go` import `fwm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`
- [x] 1.2 在 `TencentCloudClient` 结构体新增字段 `fwmConn *fwm.Client`
- [x] 1.3 新增 `UseFwmClient()` 方法，仿 `UseCfwClient`（`NewClientProfile(300)`、`NewClient`、`WithHttpTransport`）

## 2. 服务层

- [x] 2.1 创建 `tencentcloud/services/fwm/service_tencentcloud_fwm.go`，定义 `FwmService{client}` 与辅助方法
- [x] 2.2 实现 `DescribeFwmMemberGroup(ctx, appId, memberId)`：用 `DescribeOrganMembers` + `Filters=[{Name:"AppId",Values:[appId],OperatorType:1}]` 查询，遍历 `Members` 按 `MemberId` 精确匹配，返回目标 `*OrganMemberItem`（无匹配返回 nil）

## 3. 资源实现

- [x] 3.1 创建 `tencentcloud/services/fwm/resource_tc_fwm_member_group_config.go`，定义 `ResourceTencentCloudFwmMemberGroup()` schema：`member_id`(required)、`app_id`(required+ForceNew)、`group_id`(optional，空=未分组)，description 中文、字段 description 英文，Importer=ImportStatePassthrough
- [x] 3.2 Create：解析 app_id/member_id，`d.SetId(app_id#member_id)`（`tccommon.FILED_SP`）后转 Update
- [x] 3.3 Read：从 `d.Id()` 解析 app_id、member_id，retry 包内调 `DescribeFwmMemberGroup`；成员查到则取 `AccountGroup.GroupId` 回填（null 置空字符串，不清 id）；成员查不到先打印日志再 `d.SetId("")`
- [x] 3.4 Update：`d.HasChange("group_id")` 才调 `ModifyMemberGroup`，`UserMemberId=member_id`、`GroupId=group_id`（空字符串移出）；写成功后调 Read 回填；同步无轮询
- [x] 3.5 Delete：空实现 `return nil`

## 4. Provider 注册

- [x] 4.1 在 `tencentcloud/provider.go` 资源 map 注册 `"tencentcloud_fwm_member_group": fwm.ResourceTencentCloudFwmMemberGroup()`
- [x] 4.2 在 `tencentcloud/provider.md` 新增资源条目

## 5. 文档

- [x] 5.1 创建 `tencentcloud/services/fwm/resource_tc_fwm_member_group.md`：一句话描述含 FWM；Example Usage（`depends_on` 指向 cfw_managed、`group_id` 引用组 id）；Import 段说明两段式 `app_id#member_id`；不手写 Argument/Attribute Reference
- [ ] 5.2 通过收尾阶段 `make doc` 生成 `website/docs/r/tencentcloud_fwm_member_group.html.markdown`

## 6. 测试

- [x] 6.1 创建 `tencentcloud/services/fwm/resource_tc_fwm_member_group_test.go`，使用 gomonkey mock 云 API（DescribeOrganMembers / ModifyMemberGroup），覆盖：建组、读未分组、读已分组、更新 group_id、移出（空字符串）、skip unchanged、成员不存在等业务逻辑

## 7. 验证（独立任务，本轮不执行）

- [x] 7.1 代码正确性检查：确认 ModifyMemberGroup 入参 UserMemberId/GroupId 在云 API 中存在；DescribeOrganMembers 过滤字段使用 AppId（不支持 MemberId）
- [ ] 7.2 gofmt / make doc（仅收尾阶段 tfpacer-finalize 执行）