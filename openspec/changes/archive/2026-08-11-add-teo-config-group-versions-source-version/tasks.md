## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/teo/data_source_tc_teo_config_group_versions.go` 的 `config_group_version_infos` 嵌套 schema 中新增 `source_version` 字段（`schema.TypeString`，`Optional: true`），用于映射云 API `DescribeConfigGroupVersions` 返回的 `ConfigGroupVersionInfos.SourceVersion`，并添加对应描述（来源版本 ID，格式 `ver-xxxxxxxx`）。

## 2. Read 函数实现

- [x] 2.1 在 `dataSourceTencentCloudTeoConfigGroupVersionsRead` 的 `respData` 遍历循环中，增加对 `configGroupVersionInfos.SourceVersion` 的 nil 判断，并在非 nil 时将其写入 `configGroupVersionInfosMap["source_version"]`，位置与现有字段保持一致。

## 3. 单元测试

- [x] 3.1 在 `tencentcloud/services/eo/data_source_tc_teo_config_group_versions_test.go` 中，使用 gomonkey 对 `DescribeTeoConfigGroupVersionsByFilter` 进行 mock，新增单元测试用例，断言当 API 返回非空 `SourceVersion` 时 `config_group_version_infos.*.source_version` 被正确填充。
- [x] 3.2 在同一测试文件中补充用例，断言当 API 返回的 `SourceVersion` 为 nil 时 `source_version` 字段被安全省略且不报错。

## 4. 文档更新

- [x] 4.1 更新 `tencentcloud/services/teo/data_source_tc_teo_config_group_versions.md`，在 `config_group_version_infos` 块的字段说明中补充 `source_version` 字段描述。
- [x] 4.2 在收尾阶段通过 `make doc` 命令重新生成 `website/docs/d/teo_config_group_versions.html.markdown`（不在本阶段手动编辑 `website/` 目录）。
