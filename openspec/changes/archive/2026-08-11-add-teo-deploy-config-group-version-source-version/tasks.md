## 1. Schema 修改

- [x] 1.1 在 `tencentcloud/services/teo/resource_tc_teo_deploy_config_group_version.go` 的 `config_group_version_infos` 嵌套 schema 中新增 `source_version` 字段：`Type: schema.TypeString`、`Computed: true`，并补充 Description（来源版本 ID，只读）

## 2. Read 函数修改

- [x] 2.1 在 `resourceTencentCloudTeoDeployConfigGroupVersionRead` 函数遍历 `deployRecord.ConfigGroupVersionInfos` 的循环中，新增 nil 检查并回填 `source_version`：
```go
if configGroupVersionInfo.SourceVersion != nil {
    configGroupVersionInfoMap["source_version"] = configGroupVersionInfo.SourceVersion
}
```

## 3. 文档更新

- [x] 3.1 在 `tencentcloud/services/teo/resource_tc_teo_deploy_config_group_version.md` 的 `config_group_version_infos` 嵌套块字段说明中补充 `source_version` 字段描述
- [x] 3.2 在收尾阶段通过 `make doc` 重新生成 `website/docs/r/teo_deploy_config_group_version.html.markdown` 文档

## 4. 单元测试补充

- [x] 4.1 在 `tencentcloud/services/teo/resource_tc_teo_deploy_config_group_version_test.go` 中使用 gomonkey mock 的方式补充单元测试，覆盖 Read 回填 `source_version` 的场景（云 API 返回非空 / 返回 nil 两种情况）

## 5. 代码正确性检查

- [x] 5.1 核对新增字段仅在 Read 中回填，Create/未传入 `DeployConfigGroupVersion` 请求（保持现有行为）
- [x] 5.2 核对云 API `ConfigGroupVersionInfo.SourceVersion` 字段在 vendor 包中存在（`teo/v20220901`，已确认）

## 6. 收尾阶段（tfpacer-finalize）

- [x] 6.1 执行 `gofmt` 格式化变更的 Go 代码
- [x] 6.2 执行 `make doc` 更新 website 文档
- [x] 6.3 在 `.changelog/` 下创建 changelog 文件
