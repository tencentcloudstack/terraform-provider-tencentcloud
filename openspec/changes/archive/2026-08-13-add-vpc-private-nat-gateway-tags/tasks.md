## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/vpc/resource_tc_vpc_private_nat_gateway.go` 的 `ResourceTencentCloudVpcPrivateNatGateway()` schema map 中新增 `tags` 参数，类型为 `schema.TypeMap`，`Optional: true`，Description 为 "Tag description of the instance."

## 2. CRUD 函数修改

- [x] 2.1 修改 `resourceTencentCloudVpcPrivateNatGatewayCreate` 函数：在创建请求构建部分，使用 `d.GetOk("tags")` 获取 tags map，遍历 map 构建 `[]*vpc.Tag` 列表（每个 Tag 设置 Key 和 Value），并赋值到 `request.Tags`
- [x] 2.2 修改 `resourceTencentCloudVpcPrivateNatGatewayRead` 函数：在读取并 set 字段的部分，判断 `privateNatGateway.TagSet` 不为 nil 且长度大于 0 时，遍历 TagSet 构建 `map[string]string`，并执行 `d.Set("tags", tagsMap)`
- [x] 2.3 修改 `resourceTencentCloudVpcPrivateNatGatewayUpdate` 函数：在 `immutableArgs` 数组中添加 `"tags"`，使 tags 变更时报错提示不可修改

## 3. 测试用例补充

- [x] 3.1 在 `tencentcloud/services/vpc/resource_tc_vpc_private_nat_gateway_test.go` 中补充 tags 参数的单元测试用例，使用 gomonkey 对云 API 进行 mock 处理，测试 Create 方法中 tags map 到 `[]*vpc.Tag` 的转换逻辑

## 4. 文档更新

- [x] 4.1 更新 `tencentcloud/services/vpc/resource_tc_vpc_private_nat_gateway.md` 文件，在 Example Usage 中补充 tags 参数的使用示例
