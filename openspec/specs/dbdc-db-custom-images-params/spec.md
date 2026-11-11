# dbdc-db-custom-images-params

## Requirements

### Requirement: dbdc_db_custom_images 数据源支持过滤入参 Name 和 Values
`dbdc_db_custom_images` 数据源 SHALL 提供可选入参 `Name`（string）和 `Values`（字符串列表），用于构造 `DescribeDBCustomImages` 接口的 `Filters` 过滤条件。

#### Scenario: 同时传入 Name 和 Values 过滤镜像
- **WHEN** 用户在 `tencentcloud_dbdc_db_custom_images` 数据源配置中同时设置 `Name` 与 `Values`
- **THEN** 系统 SHALL 将其组装为一个 `Filter`（`Name` + `Values`）并传入 `DescribeDBCustomImages` 请求的 `Filters` 字段，返回符合过滤条件的镜像列表

#### Scenario: 不传入过滤参数时返回全部镜像
- **WHEN** 用户未设置 `Name` 与 `Values`（或 `Name` 为空）
- **THEN** 系统 SHALL 不构造 `Filters`，`DescribeDBCustomImages` 请求不携带过滤条件，返回全部镜像，行为与变更前一致

### Requirement: dbdc_db_custom_images 数据源 image_set 输出 OsType 字段
`dbdc_db_custom_images` 数据源的 `image_set` 出参 SHALL 包含 `os_type`（string，Computed）字段，对应云 API 返回的 `response.ImageSet.OsType`。

#### Scenario: 云 API 返回的镜像包含 OsType
- **WHEN** `DescribeDBCustomImages` 返回的 `ImageSet` 元素中 `OsType` 非空
- **THEN** 系统 SHALL 将该值写入对应 `image_set` 元素的 `os_type` 字段

#### Scenario: 云 API 返回的镜像 OsType 为空
- **WHEN** `DescribeDBCustomImages` 返回的 `ImageSet` 元素中 `OsType` 为 nil
- **THEN** 系统 SHALL 跳过设置 `os_type`，不引发错误
