## 1. 添加 name 拆分辅助函数

- [x] 1.1 在 `tencentcloud/services/teo/resource_tc_teo_function_extension.go` 中添加 `parseTeoFunctionOriginalName` 函数，通过查找 `-zone-` 子串的位置，从拼接后的 name 中提取原始 name 部分。若未找到 `-zone-` 则返回原始值。

## 2. 修改 Read 方法

- [x] 2.1 在 `tencentcloud/services/teo/resource_tc_teo_function.go` 的 `resourceTencentCloudTeoFunctionRead` 方法中，将 `respData.Name` 的设置改为先调用 `parseTeoFunctionOriginalName` 函数进行拆分，再将原始 name 设置到 state

## 3. 补充单元测试

- [x] 3.1 在 `tencentcloud/services/teo/resource_tc_teo_function_test.go` 中添加 `parseTeoFunctionOriginalName` 函数的单元测试，覆盖以下场景：标准拼接 name 拆分、原始 name 不包含连字符、原始 name 包含多个连字符、不含 -zone- 子串、空字符串

## 4. 更新文档示例

- [x] 4.1 更新 `tencentcloud/services/teo/resource_tc_teo_function.md` 中的示例，将 name 值从拼接后的完整函数名改为原始 name 值
