## 1. 资源 schema 与 Read 实现

- [x] 1.1 在 `tencentcloud/services/dnspod/resource_tc_dnspod_record.go` 的 `ResourceTencentCloudDnspodRecord().Schema` 中新增 `updated_on`：`Type: schema.TypeString, Computed: true, Description: "Last update time of the record."`
- [x] 1.2 在 `resourceTencentCloudDnspodRecordRead` 中、与 `monitor_status` / `remark` 同级位置加入 `if recordInfo.UpdatedOn != nil { _ = d.Set("updated_on", recordInfo.UpdatedOn) }`
- [x] 1.3 复核：不要将 `updated_on` 加入 Create/Update 的 request 构造路径；不要让其参与 `d.Partial`、`d.HasChange` 等更新判断

## 2. 数据源既有字段确认（无代码改动）

- [x] 2.1 复核 `tencentcloud/services/dnspod/data_source_tc_dnspod_record_list.go` 中 `record_list` 与 `instance_list` 子 schema 已存在 `updated_on`（`schema.TypeString, Computed: true`），如缺失则补齐
- [x] 2.2 复核 Read 函数中存在 `recordListItemMap["updated_on"] = recordListItem.UpdatedOn` 与 `instanceListItemMap["updated_on"] = recordListItem.UpdatedOn` 赋值，且包裹在 `if recordListItem.UpdatedOn != nil` 守护中

## 3. 文档与示例

- [x] 3.1 更新 `tencentcloud/services/dnspod/resource_tc_dnspod_record.md`，在示例 HCL 后追加说明并展示 `output "updated_on" { value = tencentcloud_dnspod_record.demo.updated_on }`
- [x] 3.2 复核 `tencentcloud/services/dnspod/data_source_tc_dnspod_record_list.md` 是否需要补充 `updated_on` 字段说明（如缺失则补充）
- [x] 3.3 运行 `make doc` 自动生成 `website/docs/r/dnspod_record.html.markdown` 与 `website/docs/d/dnspod_record_list.html.markdown`，确认两份文件均包含 `updated_on`

## 4. 验收测试

- [x] 4.1 在 `tencentcloud/services/dnspod/resource_tc_dnspod_record_test.go` 的现有 `TestAccTencentCloudDnsPodRecord*` 测试 Steps 中追加 `resource.TestCheckResourceAttrSet("tencentcloud_dnspod_record.<name>", "updated_on")` 断言（不新增独立 Test 函数，避免重复占用真实 DNSPod 配额）
- [x] 4.2 在 `tencentcloud/services/dnspod/data_source_tc_dnspod_record_list_test.go` 现有测试中追加 `resource.TestCheckResourceAttrSet("data.tencentcloud_dnspod_record_list.<name>", "record_list.0.updated_on")` 与 `instance_list.0.updated_on` 断言

## 5. Changelog

- [x] 5.1 在 `.changelog/` 下新增 `<PR_NUMBER>.txt`（占位 PR 号），使用 `enhancement` 块，内容：
  ```
  ```release-note:enhancement
  resource/tencentcloud_dnspod_record: support new attribute `updated_on`
  ```
  ```release-note:enhancement
  data-source/tencentcloud_dnspod_record_list: add `updated_on` attribute under `record_list` and `instance_list`
  ```
  ```

## 6. 验证

- [x] 6.1 运行 `make fmt`
- [x] 6.2 运行 `make lint`（针对 dnspod 模块），确认本变更未引入新增 lint 错误
- [x] 6.3 运行 `go build ./...`
- [x] 6.4 运行 `go vet ./tencentcloud/services/dnspod/...`
- [ ] 6.5 在具备 DNSPod 测试账号的环境下运行：`TF_ACC=1 go test -run "TestAccTencentCloudDnsPodRecord|TestAccTencentCloudDnspodRecordListDataSource" ./tencentcloud/services/dnspod/... -v -timeout 30m`
- [x] 6.6 运行 `openspec validate dnspod-record-updated-on --strict`
