## 1. 移除 zone_setting

- [x] 1.1 从 `resource_tc_teo_l7_acc_setting.go` 的 Schema 中移除 `zone_setting` computed 属性定义（原 lines 550-1076）
- [x] 1.2 从 `resourceTencentCloudTeoL7AccSettingRead` 函数中移除 `zone_setting` 的填充逻辑
- [x] 1.3 从 `resource_tc_teo_l7_acc_setting_test.go` 中移除 `zone_setting` 相关的单元测试及其依赖的 mock 辅助代码

## 2. 添加 network_error_logging

- [x] 2.1 在 `zone_config` schema 中添加 `network_error_logging` 字段（TypeList, Optional, MaxItems:1, 包含 `switch` 字段）
- [x] 2.2 在 `resourceTencentCloudTeoL7AccSettingRead` 函数中添加 `respData.ZoneConfig.NetworkErrorLogging` 的读取映射
- [x] 2.3 在 `resourceTencentCloudTeoL7AccSettingUpdate` 函数中添加 `network_error_logging` 到 `NetworkErrorLoggingParameters` 的映射

## 3. 文档

- [x] 3.1 更新 `resource_tc_teo_l7_acc_setting.md` 示例，添加 `network_error_logging` 配置块
