## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/teo/resource_tc_teo_config_group_version.go` 的 Schema map 中新增 `source_version` 字段，放置在 `content` 字段之后、`version_id` 字段之前
- [x] 1.2 设置 `source_version` 字段类型为 `schema.TypeString`，属性为 `Optional: true, ForceNew: true`
- [x] 1.3 为 `source_version` 添加 Description 说明："Source version ID. The new version will be derived from the configuration of this source version. If not specified, the currently active production version is used as the source version by default."

## 2. Create 函数扩展

- [x] 2.1 在 `resourceTencentCloudTeoConfigGroupVersionCreate` 中，使用 `d.GetOk("source_version")` 读取用户配置的来源版本 ID
- [x] 2.2 当 `d.GetOk("source_version")` 返回 ok=true 时，设置 `request.SourceVersion = helper.String(v.(string))`；未配置时保持 nil（沿用云平台默认行为）
- [x] 2.3 确认 Create 方法在 retry 块外、调用 `CreateConfigGroupVersionWithContext` 成功后，检查 `response.Response.VersionId` 非空（已有逻辑），若为空返回 `NonRetryableError`
- [x] 2.4 确认 Create 方法末尾仍调用 `resourceTencentCloudTeoConfigGroupVersionRead(d, meta)` 回写状态

## 3. Read 函数扩展

- [x] 3.1 在 `resourceTencentCloudTeoConfigGroupVersionRead` 的 `respData.ConfigGroupVersionInfo != nil` 分支内，新增 `SourceVersion` 的 nil 检查与 `d.Set("source_version", respData.ConfigGroupVersionInfo.SourceVersion)` 回填
- [x] 3.2 确认 Read 方法中 `respData == nil` 时先打印 `log.Printf("[CRUD] teo_config_group_version id=%s", d.Id())` 再 `d.SetId("")`（已有逻辑，保持不变）

## 4. 单元测试

- [x] 4.1 在 `tencentcloud/services/teo/resource_tc_teo_config_group_version_test.go` 中新增使用 gomonkey mock 的单元测试用例，覆盖 Create 时设置 `source_version` 的场景（mock `CreateConfigGroupVersionWithContext` 返回有效 VersionId，mock `DescribeConfigGroupVersionDetail` 返回 SourceVersion）
- [x] 4.2 新增单元测试用例，覆盖 Create 时未设置 `source_version` 的场景（验证 `request.SourceVersion` 为 nil）
- [x] 4.3 新增单元测试用例，覆盖 Read 时 `SourceVersion` 为 nil 的场景（验证不报错且其余字段正常回填）
- [x] 4.4 确保新增测试代码可正确编译，遵循 gomonkey mock 模式（`-gcflags=all=-l`），不使用 terraform 测试套件运行新增用例

## 5. 文档同步

- [x] 5.1 在 `tencentcloud/services/teo/resource_tc_teo_config_group_version.md` 的 Example Usage 中新增 `source_version` 字段示例
- [x] 5.2 确认文档不包含 `Argument Reference` 和 `Attribute Reference` 部分（由工具自动生成）
- [x] 5.3 确认文档一句话描述中包含云产品名称 TEO，且包含 Import 部分说明联合 id（`zone_id#group_id#version_id`）

## 6. 代码正确性检查

- [x] 6.1 检查 `source_version` 新增参数在 Create 接口 `CreateConfigGroupVersionRequest` 中存在（vendor SDK 确认）
- [x] 6.2 检查 `source_version` 读取参数在 `ConfigGroupVersionInfo.SourceVersion` 中存在（vendor SDK 确认）
- [x] 6.3 确认所有函数返回的 error 已被检查，必定不出错的函数用 `_ = func()` 赋值
- [x] 6.4 确认未修改 `provider.go` 和 `provider.md`（资源已注册，仅新增参数）
