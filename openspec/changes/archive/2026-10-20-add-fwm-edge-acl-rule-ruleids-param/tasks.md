## 1. 资源代码实现

- [x] 1.1 创建 `tencentcloud/services/fwm/resource_tc_fwm_edge_acl_rule.go`，参考 `tencentcloud/services/igtm/resource_tc_igtm_strategy.go` 风格，按 RESOURCE_KIND_GENERAL 实现完整 CRUD（Create/Read/Update/Delete），调用 fwm SDK：
  - Create：`CreateEdgeAclRule`（入参 `GroupId` + `Rules []*EdgeAclRuleInfo`）
  - Read：`DescribeEdgeAclRules`（入参 `GroupId`，分页 `Limit` 取云 API 注释最大值 1000）
  - Delete：`DeleteEdgeAclRule`（入参 `GroupId` + `RuleIds []*string`）
  - Update：fwm 无独立 Modify 接口，按先删后建实现
- [x] 1.2 在 schema 中定义 `RuleIds` 出参字段：`Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}, Computed: true`（不设 Optional/Required/ForceNew）
- [x] 1.3 在 Create 方法中：调用 `CreateEdgeAclRule` 后检查 `response`/`response.Response`/`RuleIds` 是否为空，为空则返回 `NonRetryableError`；非空时先 `d.SetId(...)` 再 `d.Set("RuleIds", RuleIds)`
- [x] 1.4 在 Read 方法中：从 `DescribeEdgeAclRules` 返回的规则列表提取 RuleId 回填 `RuleIds`；response 为空时先 `log.Printf("[CRUD] fwm_edge_acl_rule id=%s", d.Id())` 再 `d.SetId("")`
- [x] 1.5 在 CRUD 调用云 API 处使用 `tccommon.ReadRetryTimeout` 超时 + `helper.Retry` 包装，失败用 `tccommon.RetryError()` 包装返回；retry 块内只做接口调用，SetId 等Success 操作放到 retry 块外

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 中注册 `tencentcloud_fwm_edge_acl_rule` 资源
- [x] 2.2 在 `tencentcloud/provider.md` 中添加该资源说明（import 示例需说明使用联合 id）

## 3. 测试

- [x] 3.1 创建 `tencentcloud/services/fwm/resource_tc_fwm_edge_acl_rule_test.go`，使用 gomonkey mock 云 API 进行业务逻辑单元测试，覆盖 Create 回填 `RuleIds`、Read 刷新 `RuleIds`、空返回报错等场景（禁止使用 terraform 测试套件，禁止执行 go test）

## 4. 文档

- [ ] 4.1 通过收尾阶段 `make doc` 命令生成 `website/` 下文档（禁止手动编写 website/ 文件）