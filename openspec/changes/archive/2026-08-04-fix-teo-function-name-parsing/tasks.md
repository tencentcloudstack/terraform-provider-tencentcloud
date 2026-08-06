## 1. 修改 name 拆分辅助函数

- [x] 1.1 修改 `tencentcloud/services/teo/resource_tc_teo_function.go` 中的 `ParseTeoFunctionOriginalName` 函数，接收 `name` 和 `zoneId` 两个参数，通过查找 `-zoneId` 后缀（`-` + zone_id，如 `-zone-2qtuhspy7cr6`）在拼接 name 中的位置来提取原始 name。相比查找 `-zone-` 子串，此方法能正确处理原始 name 中包含 `-zone-` 字符串的场景（如 `my-zone-func`）。若未找到匹配后缀则返回原始值。

## 2. 修改 Read 方法

- [x] 2.1 在 `tencentcloud/services/teo/resource_tc_teo_function.go` 的 `resourceTencentCloudTeoFunctionRead` 方法中，将 `respData.Name` 的设置改为先调用 `ParseTeoFunctionOriginalName(*respData.Name, zoneId)` 函数进行拆分，再将原始 name 设置到 state

## 3. 补充单元测试

- [x] 3.1 在 `tencentcloud/services/teo/resource_tc_teo_function_test.go` 中添加 `ParseTeoFunctionOriginalName` 函数的单元测试，覆盖以下场景：标准拼接 name 拆分、原始 name 不包含连字符、原始 name 包含多个连字符、原始 name 包含 -zone- 子串、不含匹配后缀、空字符串

## 4. 更新文档示例

- [x] 4.1 更新 `tencentcloud/services/teo/resource_tc_teo_function.md` 中的示例，将 name 值从拼接后的完整函数名改为原始 name 值
