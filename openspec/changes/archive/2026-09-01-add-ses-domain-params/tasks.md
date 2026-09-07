## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/ses/resource_tc_ses_domain.go` 的 Schema 中新增 `dkim_option` 字段（TypeInt, Optional, ForceNew, 描述 DKIM 密钥长度 0:1024 1:2048）
- [x] 1.2 在 Schema 中新增 `tag_list` 嵌套块字段（TypeList, Optional, ForceNew），其子字段为 `tag_key`（TypeString, Required, 描述标签键）和 `tag_value`（TypeString, Required, 描述标签值），支持多个标签

## 2. Service 层修改

- [x] 2.1 修改 `tencentcloud/services/ses/service_tencentcloud_ses.go` 中的 `DescribeSesDomain` 方法，使其返回 `*ses.GetEmailIdentityResponseParams` 而非 `[]*ses.DNSAttributes`，以便 Read 函数可以访问 `DKIMOption` 和 `TagList` 字段
- [x] 2.2 更新 `DescribeSesDomain` 方法内部的空响应检查逻辑，确保 response 为 nil 或 Response 为 nil 时正确返回

## 3. CRUD 函数修改

- [x] 3.1 修改 `resourceTencentCloudSesDomainCreate` 函数：从 schema 读取 `dkim_option` 并设置到 `request.DKIMOption`（将 int 转换为 uint64）
- [x] 3.2 修改 `resourceTencentCloudSesDomainCreate` 函数：从 schema 读取 `tag_list` 列表，遍历每个元素的 `tag_key` 和 `tag_value`，构造 `[]*ses.TagList` 切片并设置到 `request.TagList`
- [x] 3.3 修改 `resourceTencentCloudSesDomainRead` 函数：适配 `DescribeSesDomain` 新的返回类型，从 `GetEmailIdentityResponseParams` 中读取 `Attributes`、`DKIMOption`、`TagList`
- [x] 3.4 修改 `resourceTencentCloudSesDomainRead` 函数：当 `DKIMOption` 不为 nil 时，调用 `d.Set("dkim_option", *response.DKIMOption)` 设置状态
- [x] 3.5 修改 `resourceTencentCloudSesDomainRead` 函数：当 `TagList` 非空时，遍历每个元素读取 `TagKey` 和 `TagValue` 构造 `tag_list` 列表设置到 state；为空时不设置

## 4. 测试更新

- [x] 4.1 更新 `tencentcloud/services/ses/resource_tc_ses_domain_test.go`，在测试用例中补充 `dkim_option`、`tag_list`（含 `tag_key`、`tag_value`）参数的测试覆盖

## 5. 文档更新

- [x] 5.1 更新 `tencentcloud/services/ses/resource_tc_ses_domain.md` 的 Example Usage，补充 `dkim_option`、`tag_list` 参数示例
