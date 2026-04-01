## 1. Schema 定义

- [x] 1.1 在 tencentcloud_cvm_launch_template resource schema 中新增 `launch_template_id` 字段，类型为 String，Computed 属性为 true
- [x] 1.2 验证字段名称遵循 snake_case 命名约定

## 2. Create 函数修改

- [x] 2.1 在 Create 函数中，从 CreateLaunchTemplate API 响应中读取 LaunchTemplateId 字段
- [x] 2.2 使用 d.Set("launch_template_id", value) 将 LaunchTemplateId 设置到 state
- [x] 2.3 添加错误处理逻辑，如果 API 响应中不包含 LaunchTemplateId 字段，忽略该字段不报错

## 3. Read 函数修改

- [x] 3.1 在 Read 函数中，从 DescribeLaunchTemplates API 响应中读取 LaunchTemplateId 字段
- [x] 3.2 使用 d.Set("launch_template_id", value) 更新 state 中的 launch_template_id 字段
- [x] 3.3 添加错误处理逻辑，如果 API 响应中不包含 LaunchTemplateId 字段，忽略该字段不报错

## 4. 单元测试更新

- [x] 4.1 在 resource_tencentcloud_cvm_launch_template_test.go 中添加单元测试用例，验证 Create 函数正确设置 launch_template_id
- [x] 4.2 添加单元测试用例，验证 Read 函数正确更新 launch_template_id
- [x] 4.3 添加测试用例验证向后兼容性（现有配置不包含 launch_template_id 的情况下仍然正常工作）

## 5. 验收测试更新

- [x] 5.1 更新验收测试，添加测试用例验证 launch_template_id 字段在创建时正确返回
- [x] 5.2 更新验收测试，添加测试用例验证 launch_template_id 字段在读取时正确返回

## 6. 文档更新

- [x] 6.1 更新 resource_tc_cvm_launch_template.md 样例文件，添加 launch_template_id 字段说明
- [ ] 6.2 运行 make doc 命令自动生成 website/docs/ 下的 markdown 文档

## 7. 验证任务

- [ ] 7.1 运行单元测试验证所有测试通过：go test ./tencentcloud/services/cvm -run TestAccTencentCloudCvmLaunchTemplate -v
- [ ] 7.2 运行验收测试验证实际 API 调用：TF_ACC=1 go test ./tencentcloud/services/cvm -run TestAccTencentCloudCvmLaunchTemplate -v
- [ ] 7.3 运行 make lint 检查代码规范
- [ ] 7.4 运行 make build 验证构建成功
- [ ] 7.5 验证新增字段不影响现有功能和用户配置
