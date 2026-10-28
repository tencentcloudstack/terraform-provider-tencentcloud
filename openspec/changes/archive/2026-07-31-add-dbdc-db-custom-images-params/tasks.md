## 1. Schema 与 Read 函数修改

- [x] 1.1 在 `data_source_tc_dbdc_db_custom_images.go` 的 `Schema` 中新增顶层可选入参 `name`（TypeString, Optional）和 `values`（TypeList, Elem TypeString, Optional），用于构造过滤条件。
- [x] 1.2 在 `data_source_tc_dbdc_db_custom_images.go` 的 `image_set.Elem.Schema` 中新增 Computed 出参字段 `os_type`（TypeString）。
- [x] 1.3 在 `dataSourceTencentCloudDbdcDbCustomImagesRead` 方法中，读取 `name`/`values` 并写入 `paramMap`（仅当 `name` 非空时构造），传给 `DescribeDBCustomImagesByFilter`。
- [x] 1.4 在 Read 方法遍历 `respData` 设置 `image_set` 时，新增对 `image.OsType != nil` 的判断并设置 `imageMap["os_type"]`。

## 2. Service 层修改

- [x] 2.1 在 `service_tencentcloud_dbdc.go` 的 `DescribeDBCustomImagesByFilter` 方法中，从 `paramMap` 读取 `name` 与 `values`，当 `name` 非空时构造 `[]*dbdcv20201029.Filter`（`Name`=name，`Values`=[]*string）并赋值给 `request.Filters`。

## 3. 单元测试补充

- [x] 3.1 在 `data_source_tc_dbdc_db_custom_images_test.go` 的 schema 测试中新增对 `name`、`values`、`os_type` 字段的断言。
- [x] 3.2 在 `data_source_tc_dbdc_db_custom_images_test.go` 的 Read 基础测试中，为 mock 返回的 `DBCustomImage` 增加 `OsType` 字段，并断言 `os_type` 正确写入。
- [x] 3.3 新增测试用例：传入 `name`/`values` 过滤参数时，断言请求 `Filters` 被正确构造。

## 4. 文档更新

- [x] 4.1 更新 `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_images.md`，在示例中补充 `name`、`values` 入参与 `os_type` 出参说明（由 `make doc` 生成 `website/docs/` 文档）。
