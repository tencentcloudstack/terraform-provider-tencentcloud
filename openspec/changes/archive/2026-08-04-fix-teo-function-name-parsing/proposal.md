## Why

调用CreateFunction接口后，DescribeFunctions接口返回的name字段值是拼接后的完整函数名（格式为：原始name + "-" + zone_id + "-" + app_id），而非用户原始的入参name值。这导致terraform在后续执行plan/apply时，检测到name字段不一致，错误地提示需要修改name，而实际上云API接口并不支持修改name。需要在Read方法中对DescribeFunctions返回的name进行拆分，仅保留原始name部分作为资源状态值。

## What Changes

- 修改 `resourceTencentCloudTeoFunctionRead` 方法，在从DescribeFunctions返回的name设置到state之前，按照"原始name + '-' + zone_id + '-' + app_id"的固定格式进行拆分，仅将原始name部分设置到state中
- 添加辅助函数 `parseTeoFunctionOriginalName`，根据zone_id格式（"zone-xxxxxxx"，xxxxxxx为数字加字母）和app_id格式（纯数字字符串）从拼接后的name中提取原始name
- 修改 `resource_tc_teo_function_test.go` 中的单元测试，验证name拆分逻辑的正确性
- 修改 `resource_tc_teo_function.md` 文档中的示例，使用原始name而非拼接后的完整函数名

## Capabilities

### New Capabilities

- `teo-function-name-parsing`: 对DescribeFunctions返回的拼接name进行拆分，提取原始name，避免terraform误判资源变更

### Modified Capabilities

## Impact

- `tencentcloud/services/teo/resource_tc_teo_function.go`: Read方法中name字段的处理逻辑
- `tencentcloud/services/teo/resource_tc_teo_function_extension.go`: 可能需要在此文件中添加name拆分辅助函数
- `tencentcloud/services/teo/resource_tc_teo_function_test.go`: 补充name拆分逻辑的单元测试
- `tencentcloud/services/teo/resource_tc_teo_function.md`: 更新示例中的name值
