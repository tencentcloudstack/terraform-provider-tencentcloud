## 1. Schema 调整

- [x] 1.1 在 `tencentcloud/services/teo/resource_tc_teo_load_balancer.go` 的 `health_checker` 嵌套块 schema 中，在 `recv_context` 字段之后新增 `probe_cluster` 字段定义：`Type: schema.TypeString, Optional: true`，Description 说明合法取值 `global`/`mainland_china`/`overseas`（默认 `global`），仅 `HTTP_V2` 类型实例支持设置
- [x] 1.2 更新顶层 `type` 字段的 Description，改为包含 `HTTP_V2`、`HTTP`、`GENERAL` 三个取值的完整释义：`HTTP_V2`（HTTP 专用型 V2，支持选择探测区域，探测请求量低但灵敏度更低）；`HTTP`（HTTP 专用型 V1，不支持选择探测区域，探测请求量大但灵敏度更高）；`GENERAL`（通用型，仅支持通用型源站组，可被站点加速和四层代理引用，不支持选择探测区域，灵敏度更高）

## 2. CRUD 转换函数扩展

- [x] 2.1 在 `buildTeoHealthChecker` 函数中（`recv_context` 处理之后），新增 `probe_cluster` 字段的写入逻辑：`if val, ok := m["probe_cluster"].(string); ok && val != "" { hc.ProbeCluster = helper.String(val) }`，该函数被 Create 和 Update 共用，无需额外修改 Create/Update 主函数
- [x] 2.2 在 `flattenTeoHealthChecker` 函数中（`recv_context` 处理之后），新增 `probe_cluster` 字段的读取逻辑：`if hc.ProbeCluster != nil { m["probe_cluster"] = hc.ProbeCluster }`，该函数被 Read 使用，无需额外修改 Read 主函数

## 3. 单元测试

- [x] 3.1 在 `tencentcloud/services/teo/resource_tc_teo_load_balancer_test.go` 中，将现有测试配置 `testAccTeoLoadBalancer` 和 `testAccTeoLoadBalancerUpdate` 的 `health_checker` 块中新增 `probe_cluster` 字段示例（如 `probe_cluster = "mainland_china"`），并补充 `resource.TestCheckResourceAttr` 断言验证该字段
- [x] 3.2 新增一个测试用例覆盖 `probe_cluster` 未设置时的行为（不触发 diff），验证向后兼容

## 4. 文档同步

- [x] 4.1 在 `tencentcloud/services/eo/resource_tc_teo_load_balancer.md` 的 Example Usage 健康检查块中新增 `probe_cluster` 字段示例
- [ ] 4.2 在收尾阶段执行 `make doc`，根据 provider 规范重新生成 `website/docs/` 下的 markdown 文档（禁止手改 website/ 目录）

## 5. 验证

- [x] 5.1 检查 `buildTeoHealthChecker` 中 `probe_cluster` 仅在非空时设置 `hc.ProbeCluster`
- [x] 5.2 检查 `flattenTeoHealthChecker` 中 `probe_cluster` 仅在 `hc.ProbeCluster != nil` 时回写
- [x] 5.3 检查 `type` 字段 Description 包含 `HTTP_V2` 取值
- [x] 5.4 检查 `probe_cluster` 字段为 Optional 且未设置 ForceNew（支持 Update）
