## 1. 客户端基础设施

- [x] 1.1 在 `tencentcloud/connectivity/client.go` 新增 FWM SDK import：`fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`
- [x] 1.2 在 `tencentcloud/connectivity/client.go` 的 `TencentCloudClient` struct 新增连接字段：`fwmv20250611Conn *fwmv20250611.Client`
- [x] 1.3 在 `tencentcloud/connectivity/client.go` 新增方法 `UseFwmV20250611Client()`（参照 `UseIgtmV20231024Client()` line 2248 模式：懒加载 + NewClient + WithHttpTransport）

## 2. 资源实现文件

- [x] 2.1 创建 `tencentcloud/services/fwm/resource_tc_fwm_strategy.go`，定义包 `fwm`、import（context/fmt/log、resource.Retry、schema、fwmv20250611、tccommon、helper）、`ResourceTencentCloudFwmStrategy()` 函数骨架（Create/Read/Update/Delete/Importer）
- [x] 2.2 定义 Schema：`product`（string, Required, ForceNew）、`sub_account_uin`（string, Required, ForceNew）、`exec_area`（string, Required, ForceNew）、`group_id`（string, Optional）、`sequence`（int, Optional）、`receive_account`（list of string, Optional）、`receive_group`（list of string, Optional）+ 只读 Computed 字段（`strategy_id`/`group_name`/`rule_count`/`rule_status`/`err_msg`/`error_type`/`create_time`/`update_time`/`create_by`/`create_name`/`update_by`/`update_name`）
- [x] 2.3 实现 Create：填充 `CreateStrategyRequest`（按 `exec_area` 决定传 `PreStrategy` 还是 `PostStrategy`，单元素 `{GroupId,Sequence}`；`Product`/`ReceiveAccount`/`ReceiveGroup`）；retry 包装；校验 `StrategyIds` 非空/长度>0/`StrategyIds[0]` 非空（打印 logId 与 d.Id()），空则 NonRetryableError；retry 块外 `d.SetId(StrategyIds[0])`；调 Read
- [x] 2.4 实现 Read：调用 `DescribeStrategy`（retry + ReadRetryTimeout）；遇 ResourceNotFound 或 `Strategy==nil` 时先 `log.Printf("[CRUD] fwm_strategy id=%s", d.Id())` 再 `d.SetId("")` 返回 nil；回填可读字段（先判 nil）；按 ReceiverType 拆分回填 `receive_account`(0)/`receive_group`(1)，只比 Uin；**不 set** `product`/`sub_account_uin`/`exec_area` 保留 state
- [x] 2.5 实现 Update：`mutableArgs=[group_id,sequence,receive_account,receive_group]` 判 HasChange；仅把变更字段写入 `ModifyStrategyRequest`（增量语义，未变字段指针保持 nil）；`StrategyId` 从 `d.Id()` 取；retry 包装；调 Read
- [x] 2.6 实现 Delete：`DeleteStrategyRequest.StrategyIds` 传单元素 `[d.Id()]`；retry 包装；不存在的 ID 静默成功（RetryError 处理 ResourceNotFound 不报错）

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` import 块新增：`"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"`
- [x] 3.2 在 `tencentcloud/provider.go` 的 ResourcesMap 中新增：`"tencentcloud_fwm_strategy": fwm.ResourceTencentCloudFwmStrategy(),`

## 4. 文档与测试

- [x] 4.1 创建 `tencentcloud/services/fwm/resource_tc_fwm_strategy.md`：一句话描述（带上 FWM 产品名）+ Example Usage（product=enterprise_sg、sub_account_uin、group_id、sequence、receive_account）+ Import 部分（说明用 StrategyId）
- [x] 4.2 创建 `tencentcloud/services/fwm/resource_tc_fwm_strategy_test.go`：使用 gomonkey mock 云 API（Mock `UseFwmV20250611Client` 的 CreateStrategy/DescribeStrategy/ModifyStrategy/DeleteStrategy），覆盖 Create/Read/Update/Delete 业务逻辑；禁止使用 terraform 测试套件、禁止 go test 执行
- [ ] 4.3 收尾阶段通过 `make doc` 生成 `website/docs/r/tencentcloud_fwm_strategy.html.markdown` 与 provider.md 条目（禁止手动编辑 website/）

## 5. 验证（收尾阶段执行）

- [ ] 5.1 执行 `gofmt -w` 格式化新增/修改的 .go 文件（仅收尾阶段）
- [ ] 5.2 通过 `make doc` 生成文档（仅收尾阶段）
- [x] 5.3 代码正确性检查：确认 Create 入参字段在 CreateStrategy 中存在、Update 入参字段在 ModifyStrategy 中存在、Read 读取字段在 DescribeStrategy 响应中存在；确认所有 error 被检查