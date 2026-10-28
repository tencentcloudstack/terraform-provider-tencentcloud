## 1. Schema 定义与 CRUD 代码修改

- [x] 1.1 在 `tencentcloud/services/cam/resource_tc_cam_role_sso.go` 的 `ResourceTencentCloudCamRoleSSO()` schema map 中新增 `auto_rotate_key` 字段（`schema.TypeInt`，Optional，非 ForceNew，Description 说明 0=关闭/1=开启/默认0）
- [x] 1.2 在 `resourceTencentCloudCamRoleSSOCreate` 中读取 `auto_rotate_key` 并设置到 `request.AutoRotateKey`（使用 helper 将 int 转为 `*uint64`）
- [x] 1.3 在 `resourceTencentCloudCamRoleSSORead` 中增加 nil 判断：当 `response.Response.AutoRotateKey != nil` 时回填 `auto_rotate_key` 到 state
- [x] 1.4 在 `resourceTencentCloudCamRoleSSOUpdate` 的 `d.HasChange(...)` 检测条件中加入 `auto_rotate_key`，并在变更时将 `auto_rotate_key` 设置到 `request.AutoRotateKey`

## 2. 文档更新

- [x] 2.1 在 `tencentcloud/services/cam/resource_tc_cam_role_sso.md` 的 Example Usage 中补充 `auto_rotate_key` 参数示例

## 3. 测试补充

- [x] 3.1 在 `tencentcloud/services/cam/resource_tc_cam_role_sso_test.go` 中使用 gomonkey mock `CreateOIDCConfig` / `UpdateOIDCConfig` / `DescribeOIDCConfig` / `DeleteOIDCConfig`，新增覆盖 `auto_rotate_key` 创建场景的单元测试用例
- [x] 3.2 新增覆盖 `auto_rotate_key` 更新场景的单元测试用例（验证变更触发 UpdateOIDCConfig 且参数正确传递）
- [x] 3.3 新增覆盖 `auto_rotate_key` 读取回填场景的单元测试用例（验证 nil 判断与值回填逻辑）

## 4. 验证

- [x] 4.1 使用 `go test -gcflags=all=-l` 运行 `resource_tc_cam_role_sso_test.go` 涉及的单元测试并确保通过
- [x] 4.2 检查所有函数返回的 error 已正确处理（无未使用变量）

## 5. 收尾（由 tfpacer-finalize skill 统一执行，非本阶段任务）

- [ ] 5.1 执行 `gofmt` 格式化变更的 Go 代码
- [ ] 5.2 执行 `make doc` 生成 `website/docs/r/cam_role_sso.html.markdown` 文档
- [ ] 5.3 在 `.changelog/` 目录下创建 changelog 文件
