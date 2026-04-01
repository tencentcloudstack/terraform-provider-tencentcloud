## 1. Schema 定义修改

- [x] 1.1 在 resource_tencentcloud_vpc_end_point.go 中添加 SecurityGroupId 字段定义（TypeString, Optional）
- [x] 1.2 在 resource_tencentcloud_vpc_end_point.go 中添加 Tags 字段定义（TypeList + TypeMap, Optional, Key 为必填，Value 为可选）
- [x] 1.3 在 resource_tencentcloud_vpc_end_point.go 中添加 IpAddressType 字段定义（TypeString, Optional, Default: "Ipv4"）
- [x] 1.4 为 IpAddressType 字段添加验证逻辑，限制值为 "Ipv4" 或 "Ipv6"

## 2. Create 函数更新

- [x] 2.1 在 resourceTencentCloudVpcEndPointCreate 函数中从 d.Get() 获取 SecurityGroupId
- [x] 2.2 在 Create 函数中将 SecurityGroupId 传递到 CreateVpcEndPoint API 请求
- [x] 2.3 在 resourceTencentCloudVpcEndPointCreate 函数中从 d.Get() 获取 Tags
- [x] 2.4 在 Create 函数中将 Tags 转换为 API 要求格式并传递到 CreateVpcEndPoint API 请求
- [x] 2.5 在 resourceTencentCloudVpcEndPointCreate 函数中从 d.Get() 获取 IpAddressType
- [x] 2.6 在 Create 函数中将 IpAddressType（或默认值 "Ipv4"）传递到 CreateVpcEndPoint API 请求

## 3. Read 函数更新

- [x] 3.1 在 resourceTencentCloudVpcEndPointRead 函数中从 DescribeVpcEndPoints API 响应获取 SecurityGroupId
- [x] 3.2 在 Read 函数中将 SecurityGroupId 设置到 Terraform state
- [x] 3.3 在 resourceTencentCloudVpcEndPointRead 函数中从 DescribeVpcEndPoints API 响应获取 Tags
- [x] 3.4 在 Read 函数中将 Tags 转换为 Terraform 格式并设置到 Terraform state
- [x] 3.5 在 resourceTencentCloudVpcEndPointRead 函数中从 DescribeVpcEndPoints API 响应获取 IpAddressType
- [x] 3.6 在 Read 函数中将 IpAddressType 设置到 Terraform state

## 4. Update 函数更新

- [x] 4.1 在 resourceTencentCloudVpcEndPointUpdate 函数中检测 SecurityGroupId 变更
- [x] 4.2 确定是否通过 ModifyVpcEndPointAttribute API 更新 SecurityGroupId 或标记为 ForceNew
- [x] 4.3 实现更新 SecurityGroupId 的逻辑（API 调用或 ForceNew 触发重建）
- [x] 4.4 在 Update 函数中将更新后的 SecurityGroupId 设置到 Terraform state
- [x] 4.5 在 resourceTencentCloudVpcEndPointUpdate 函数中检测 Tags 变更
- [x] 4.6 确定是否通过 ModifyVpcEndPointAttribute API 或标签管理 API 更新 Tags
- [x] 4.7 实现更新 Tags 的逻辑（API 调用）
- [x] 4.8 在 Update 函数中将更新后的 Tags 设置到 Terraform state
- [x] 4.9 在 resourceTencentCloudVpcEndPointUpdate 函数中检测 IpAddressType 变更
- [x] 4.10 确定是否通过 ModifyVpcEndPointAttribute API 更新 IpAddressType 或标记为 ForceNew
- [x] 4.11 实现更新 IpAddressType 的逻辑（API 调用或 ForceNew 触发重建）
- [x] 4.12 在 Update 函数中将更新后的 IpAddressType 设置到 Terraform state

## 5. 单元测试编写

- [x] 5.1 添加 TestAccTencentCloudVpcEndPoint_SecurityGroupId_Create 测试用例（创建时指定 SecurityGroupId）
- [x] 5.2 添加 TestAccTencentCloudVpcEndPoint_SecurityGroupId_Read 测试用例（读取 SecurityGroupId）
- [x] 5.3 添加 TestAccTencentCloudVpcEndPoint_SecurityGroupId_Update 测试用例（更新 SecurityGroupId）
- [x] 5.4 添加 TestAccTencentCloudVpcEndPoint_Tags_Create 测试用例（创建时指定 Tags）
- [x] 5.5 添加 TestAccTencentCloudVpcEndPoint_Tags_Read 测试用例（读取 Tags）
- [x] 5.6 添加 TestAccTencentCloudVpcEndPoint_Tags_Update 测试用例（更新 Tags）
- [x] 5.7 添加 TestAccTencentCloudVpcEndPoint_Tags_ValidateKey 测试用例（验证 Key 必填）
- [x] 5.8 添加 TestAccTencentCloudVpcEndPoint_IpAddressType_Create 测试用例（创建时指定 IpAddressType）
- [x] 5.9 添加 TestAccTencentCloudVpcEndPoint_IpAddressType_Default 测试用例（默认值 Ipv4）
- [x] 5.10 添加 TestAccTencentCloudVpcEndPoint_IpAddressType_Read 测试用例（读取 IpAddressType）
- [x] 5.11 添加 TestAccTencentCloudVpcEndPoint_IpAddressType_Validate 测试用例（验证值为 Ipv4 或 Ipv6）

## 6. 验收测试编写

- [x] 6.1 创建完整的 TF 配置文件，包含所有三个新字段的定义
- [x] 6.2 运行验收测试（TF_ACC=1）验证资源的创建和读取
- [x] 6.3 运行验收测试验证资源的更新功能
- [x] 6.4 运行验收测试验证资源删除功能

## 7. 文档和示例更新

- [x] 7.1 更新 resource_tc_vpc_end_point.md 示例文件，添加 SecurityGroupId 字段的使用示例
- [x] 7.2 更新 resource_tc_vpc_end_point.md 示例文件，添加 Tags 字段的使用示例
- [x] 7.3 更新 resource_tc_vpc_end_point.md 示例文件，添加 IpAddressType 字段的使用示例
- [x] 7.4 运行 make doc 命令自动生成 website/docs/ 下的 markdown 文档
- [x] 7.5 验证生成的文档内容是否正确

## 8. 代码验证

- [x] 8.1 运行 go build 编译检查，确保代码无语法错误
- [x] 8.2 运行 go vet 静态检查，确保代码无潜在问题
- [x] 8.3 运行 go fmt 格式化代码，确保代码风格一致
- [x] 8.4 运行单元测试，确保所有测试用例通过
- [x] 8.5 运行验收测试（TF_ACC=1），确保与云 API 集成正常
