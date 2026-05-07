## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/tcr/resource_tc_tcr_instance.go` 的 Schema 中添加 `enable_cos_maz` 字段：Type=TypeBool, Optional=true, Computed=true, Description="Whether to enable COS bucket multi-AZ feature. Default is `false`."
- [x] 1.2 在同一文件的 Schema 中添加 `enable_cos_versioning` 字段：Type=TypeBool, Optional=true, Computed=true, Description="Whether to enable COS bucket versioning. Default is `false`."
- [x] 1.3 运行 `gofmt -w tencentcloud/services/tcr/resource_tc_tcr_instance.go` 格式化代码

## 2. Create 函数实现

- [x] 2.1 在 `resourceTencentCloudTcrInstanceCreate` 函数中，在现有 params 构建逻辑之后（约 225 行附近），添加读取 `enable_cos_maz` 字段的逻辑：使用 `d.GetOkExists("enable_cos_maz")` 判断用户是否显式设置，如果设置则添加到 params 中
- [x] 2.2 在同一位置添加读取 `enable_cos_versioning` 字段的逻辑：使用 `d.GetOkExists("enable_cos_versioning")` 判断用户是否显式设置，如果设置则添加到 params 中
- [x] 2.3 在调用 `tcrService.CreateTCRInstance` 时，确保 params 中的这两个字段会被传递给 SDK（检查 service 层是否需要修改）
- [x] 2.4 运行 `gofmt -w tencentcloud/services/tcr/resource_tc_tcr_instance.go` 格式化代码

## 3. Service 层检查与修改

- [x] 3.1 检查 `tencentcloud/services/tcr/service_tencentcloud_tcr.go` 中的 `CreateTCRInstance` 函数，确认是否需要添加对 `enable_cos_maz` 和 `enable_cos_versioning` 参数的处理
- [x] 3.2 如果需要修改 service 层，添加对这两个字段的处理逻辑，将 params 中的值设置到 `request.EnableCosMAZ` 和 `request.EnableCosVersioning`
- [x] 3.3 如果修改了 service 层，运行 `gofmt -w tencentcloud/services/tcr/service_tencentcloud_tcr.go` 格式化代码

## 4. Read 函数实现

- [x] 4.1 在 `resourceTencentCloudTcrInstanceRead` 函数中，找到设置其他字段的位置（约 410-440 行附近），添加读取 `EnableCosMAZ` 字段的逻辑：`_ = d.Set("enable_cos_maz", instance.EnableCosMAZ)`
- [x] 4.2 在同一位置添加读取 `EnableCosVersioning` 字段的逻辑：`_ = d.Set("enable_cos_versioning", instance.EnableCosVersioning)`
- [x] 4.3 运行 `gofmt -w tencentcloud/services/tcr/resource_tc_tcr_instance.go` 格式化代码

## 5. 文档更新

- [x] 5.1 在 `tencentcloud/services/tcr/resource_tc_tcr_instance.md` 的 Argument Reference 部分，添加 `enable_cos_maz` 字段说明：类型、可选性、默认值、描述
- [x] 5.2 在同一文件的 Argument Reference 部分，添加 `enable_cos_versioning` 字段说明：类型、可选性、默认值、描述
- [x] 5.3 运行 `make doc` 命令生成 `website/docs/r/tcr_instance.html.markdown` 文档

## 6. 代码验证

- [x] 6.1 运行 `go build ./tencentcloud/services/tcr/` 确保编译通过
- [x] 6.2 运行 `go vet ./tencentcloud/services/tcr/` 检查代码质量
- [x] 6.3 检查 linter 是否报告新的错误或警告
