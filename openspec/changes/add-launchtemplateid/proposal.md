## Why

CVM Launch Template 资源需要支持从 CreateLaunchTemplate API 响应中读取 LaunchTemplateId 字段。该字段表示创建成功的实例启动模板 ID，是用户在创建模板后需要获取的重要标识符。当前 provider 中缺少该字段，导致用户无法通过 Terraform 获取新创建的模板 ID。

## What Changes

在 `tencentcloud_cvm_launch_template` 资源中新增 `launch_template_id` 字段：

- 在 Resource Schema 中新增 `launch_template_id` 字段（Computed, String）
- 更新 Create 函数：从 CreateLaunchTemplate API 响应中读取并设置 `launch_template_id`
- 更新 Read 函数：从 DescribeLaunchTemplates API 响应中读取 `launch_template_id`
- 该字段为计算属性（Computed），用户不需要手动设置
- 更新相关的单元测试和验收测试，确保字段正确读取

## Capabilities

### New Capabilities

- `cvm-launch-template-id-read`: 支持 CVM Launch Template 资源读取 `launch_template_id` 字段的能力

### Modified Capabilities

无（仅新增字段，不修改现有行为）

## Impact

**受影响的文件：**
- `tencentcloud/services/cvm/resource_tencentcloud_cvm_launch_template.go` - 资源实现文件
- `tencentcloud/services/cvm/resource_tencentcloud_cvm_launch_template_test.go` - 资源测试文件

**受影响的 API：**
- CreateLaunchTemplate - 需要从响应中读取 `LaunchTemplateId` 字段
- DescribeLaunchTemplates - 需要从响应中读取 `LaunchTemplateId` 字段

**依赖变更：**
- 无新的外部依赖
- 使用现有的 tencentcloud-sdk-go API 调用
