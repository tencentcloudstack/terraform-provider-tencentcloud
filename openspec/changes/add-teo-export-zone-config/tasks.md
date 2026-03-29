## 1. Setup

- [x] 1.1 创建数据源文件 `tencentcloud/services/teo/data_source_tc_teo_export_zone_config.go`
- [x] 1.2 创建数据源测试文件 `tencentcloud/services/teo/data_source_tc_teo_export_zone_config_test.go`
- [x] 1.3 创建数据源示例文件 `tencentcloud/services/teo/data_source_tc_teo_export_zone_config.md`

## 2. Implementation

- [x] 2.1 定义数据源 schema，包括 zone_id 和 zone_name 两个可选查询参数
- [x] 2.2 定义完整的配置输出 schema，包括基础配置、加速设置、安全规则、源站设置等嵌套结构
- [x] 2.3 实现 dataSourceTcTeoExportZoneConfig 函数，配置 schema 和读取逻辑
- [x] 2.4 实现 dataSourceTcTeoExportZoneConfigRead 函数，处理数据读取逻辑
- [x] 2.5 实现查询参数验证逻辑，确保至少提供 zone_id 或 zone_name 其中之一
- [x] 2.6 在 service 层实现调用 TEO API 的函数，获取站点配置信息
- [x] 2.7 实现配置数据的映射和转换逻辑，将 API 返回的数据转换为 schema 格式
- [x] 2.8 实现错误处理和重试机制，使用 helper.Retry() 处理最终一致性问题
- [x] 2.9 实现 zone_id 优先级逻辑，当同时提供 zone_id 和 zone_name 时优先使用 zone_id

## 3. Testing

- [x] 3.1 编写单元测试，验证 schema 定义和基本功能
- [x] 3.2 编写验收测试，测试通过 zone_id 查询站点配置的场景
- [x] 3.3 编写验收测试，测试通过 zone_name 查询站点配置的场景
- [x] 3.4 编写验收测试，测试查询不存在的站点时的错误处理
- [x] 3.5 编写验收测试，验证导出的配置包含基础信息
- [x] 3.6 编写验收测试，验证导出的配置包含加速设置
- [x] 3.7 编写验收测试，验证导出的配置包含安全规则
- [x] 3.8 编写验收测试，验证导出的配置包含源站设置
- [x] 3.9 编写验收测试，测试同时提供 zone_id 和 zone_name 时的参数优先级

## 4. Documentation

- [x] 4.1 更新数据源示例文件 `data_source_tc_teo_export_zone_config.md`，包含完整的示例代码
- [x] 4.2 运行 `make doc` 命令生成 website 文档
- [x] 4.3 验证生成的文档文件 `website/docs/d/teo_export_zone_config.html.markdown` 内容正确

## 5. Verification

- [x] 5.1 运行 `go build` 验证代码编译通过
- [x] 5.2 运行 `go fmt` 格式化代码
- [x] 5.3 运行 `go vet` 进行静态检查
- [x] 5.4 运行单元测试验证代码逻辑正确
- [x] 5.5 设置环境变量 TENCENTCLOUD_SECRET_ID 和 TENCENTCLOUD_SECRET_KEY
- [x] 5.6 运行 `TF_ACC=1 go test` 执行完整的验收测试
