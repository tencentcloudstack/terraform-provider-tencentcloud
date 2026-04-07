## 1. Schema 定义

- [x] 1.1 在 `resource_tc_vpn_ssl_client.go` 的 Schema 中添加 `tags` 字段(TypeMap, Optional)
- [x] 1.2 添加 tags 字段的 description 说明
- [x] 1.3 **移除 tags 字段的 `ForceNew: true` 属性,允许原地更新**

## 2. 资源定义更新

- [x] 2.1 在资源定义中添加 `Update: resourceTencentCloudVpnSslClientUpdate` 回调函数

## 3. Create 函数实现

- [x] 3.1 在 Create 函数中读取用户配置的 tags
- [x] 3.2 将 tags map 转换为 API 所需的 `[]*Tag` 格式
- [x] 3.3 在 `CreateVpnGatewaySslClient` 请求中设置 Tags 参数
- [x] 3.4 验证 tags 在创建资源时正确传递给 API

## 4. Read 函数实现

- [x] 4.1 导入 Tag Service: `svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"`
- [x] 4.2 在 Read 函数中创建 Tag Service 实例
- [x] 4.3 使用 `tagService.DescribeResourceTags()` 查询资源标签(service: "vpc", resourceType: "vpnx")
- [x] 4.4 处理查询到的标签并使用 `d.Set("tags", tags)` 设置到 state
- [x] 4.5 处理空标签情况(graceful handling)

## 5. Update 函数实现(新增)

- [x] 5.1 创建 `resourceTencentCloudVpnSslClientUpdate` 函数
- [x] 5.2 添加标准的日志和一致性检查(defer LogElapsed, defer InconsistentCheck)
- [x] 5.3 使用 `d.HasChange("tags")` 检测 tags 是否变更
- [x] 5.4 使用 `d.GetChange("tags")` 获取旧值和新值
- [x] 5.5 使用 `svctag.DiffTags()` 计算需要添加/修改和删除的标签
- [x] 5.6 创建 Tag Service 实例
- [x] 5.7 使用 `tccommon.BuildTagResourceName("vpc", "vpnx", region, sslClientId)` 构建资源名称
- [x] 5.8 调用 `tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)` 更新标签
- [x] 5.9 处理错误情况并返回
- [x] 5.10 如果没有 tags 变更,直接返回成功(Update 函数可能被其他字段触发,但当前只有 tags 可更新)

## 6. 示例文档更新

- [x] 6.1 在 `tencentcloud/services/vpn/resource_tc_vpn_ssl_client.md` 中添加 tags 使用示例
- [x] 6.2 在示例中展示如何设置标签(至少包含 2 个标签)
- [x] 6.3 添加 tags 参数到 Arguments Reference 部分
- [x] 6.4 **更新文档说明 tags 可以原地更新,移除 ForceNew 的说明**

## 7. 测试用例

- [x] 7.1 创建 `tencentcloud/services/vpn/resource_tc_vpn_ssl_client_test.go` (如果不存在)
- [x] 7.2 添加基本测试用例 `TestAccTencentCloudVpnSslClientResource_basic` (不含 tags,验证向后兼容)
- [x] 7.3 添加 tags 测试用例 `TestAccTencentCloudVpnSslClientResource_withTags`
- [x] 7.4 在 tags 测试中验证创建时设置的标签可以正确读取
- [x] 7.5 添加测试辅助函数 `testAccCheckVpnSslClientExists` 和 `testAccCheckVpnSslClientDestroy`
- [x] 7.6 **添加 tags 更新测试用例 `TestAccTencentCloudVpnSslClientResource_updateTags`**
- [x] 7.7 **在更新测试中验证添加、删除、修改标签的场景**

## 8. 文档生成

- [x] 8.1 运行 `make doc` 生成 `website/docs/r/vpn_ssl_client.html.markdown` 文档
- [x] 8.2 验证生成的文档包含 tags 参数说明和示例
- [x] 8.3 **验证生成的文档说明 tags 可更新**

## 9. 代码验证和测试

- [x] 9.1 运行 `go build ./tencentcloud/services/vpn/...` 验证代码编译通过
- [x] 9.2 运行 `gofmt -l tencentcloud/services/vpn/resource_tc_vpn_ssl_client.go` 检查代码格式化
- [x] 9.3 运行基本测试用例验证向后兼容性(无 tags 的配置应继续工作)
- [x] 9.4 运行 tags 测试用例验证新功能
- [x] 9.5 验证 Read 函数正确处理 API 返回的标签(包括空标签情况)
- [x] 9.6 **运行 tags 更新测试用例验证更新功能**

## 10. 示例项目更新(可选)

- [x] 10.1 检查 `examples/tencentcloud-vpn/` 目录是否有 SSL Client 示例
- [x] 10.2 如果存在示例,添加 tags 参数展示用法

**注意事项**:
- tags 字段必须是 Optional,保持向后兼容
- **tags 字段移除 ForceNew,支持原地更新**
- 使用 Tag Service 查询和更新标签时,resource type 为 "vpnx"
- 参考 `resource_tc_vpn_ssl_server.go` 的 tags 更新实现方式保持一致
- Update 函数的实现应参考 vpn_ssl_server 的第 445-460 行
