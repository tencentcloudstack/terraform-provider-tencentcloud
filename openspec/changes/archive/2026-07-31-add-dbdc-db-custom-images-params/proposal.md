## Why

`dbdc_db_custom_images` 数据源当前不支持按过滤条件查询操作系统镜像列表，也不输出镜像的操作系统类型（OsType）。用户无法在 Terraform 中按镜像 ID、操作系统类型、镜像类型、架构等条件筛选镜像，也无法在查询结果中获取 OsType 字段，导致无法精确选择所需镜像。

## What Changes

- 在数据源 `tencentcloud_dbdc_db_custom_images` 中新增可选入参 `Name`（string）和 `Values`（list），用于构造 `DescribeDBCustomImages` 接口的 `Filters` 过滤条件。
- 在数据源 `tencentcloud_dbdc_db_custom_images` 的 `image_set` 出参中新增 `OsType`（string）字段，对应云 API 返回的 `response.ImageSet.OsType`。
- 同步更新 `data_source_tc_dbdc_db_custom_images_test.go` 中的单元测试用例，补充对新字段的测试。

## Capabilities

### New Capabilities
- `dbdc-db-custom-images-params`: 为 `dbdc_db_custom_images` 数据源新增过滤入参 `Name`/`Values` 及出参 `OsType`

### Modified Capabilities

## Impact

- 受影响代码文件：
  - `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_images.go`：新增 schema 字段、Read 方法中构造 Filters 请求参数、设置 OsType 出参。
  - `tencentcloud/services/dbdc/service_tencentcloud_dbdc.go`：`DescribeDBCustomImagesByFilter` 方法支持从 paramMap 读取并构造 `Filters` 请求参数。
  - `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_images_test.go`：补充单元测试用例。
- 依赖的云 API：`DescribeDBCustomImages`（dbdc v20201029），该接口已支持 `Filters` 入参与 `ImageSet.OsType` 出参。
- 文档：对应 `.md` 数据源文档需同步更新（由 `make doc` 生成）。
