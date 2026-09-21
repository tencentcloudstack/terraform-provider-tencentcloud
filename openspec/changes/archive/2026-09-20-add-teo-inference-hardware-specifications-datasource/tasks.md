## 1. Service 层实现

- [x] 1.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增 `DescribeTeoInferenceHardwareSpecificationsByFilter` 方法，接收 `paramMap`，构造 `DescribeInferenceHardwareSpecificationsRequest` 并设置 `ZoneId`，调用 `DescribeInferenceHardwareSpecificationsWithContext`，返回 `[]*teov20220901.InferenceHardwareSpecification`

## 2. 数据源核心实现

- [x] 2.1 创建 `tencentcloud/services/teo/data_source_tc_teo_inference_hardware_specifications.go` 文件
- [x] 2.2 实现 Schema 定义：必填参数 `zone_id`（TypeString）；Computed 列表参数 `hardware_specifications`（TypeList），列表内每个元素平铺 `spec`（TypeString，已废弃）、`hardware_spec_id`（TypeString）、`name`（TypeString）、`gpu_num`（TypeFloat）、`cpu_num`（TypeFloat）、`mem_size`（TypeInt）、`gpu_mem_size`（TypeInt）、`disk_size`（TypeInt）、`allowed_gpu_nums`（TypeList+TypeFloat）；可选参数 `result_output_file`（TypeString）
- [x] 2.3 实现 `dataSourceTencentCloudTeoInferenceHardwareSpecificationsRead` 函数，使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装 service 调用
- [x] 2.4 在 retry 块内对返回空（`response == nil` 或 `len(HardwareSpecifications) == 0`）返回 `NonRetryableError`，并在 retry 失败路径打印 `log.Printf("[DATASOURCE] read empty, skip SetId")`
- [x] 2.5 在 set 各字段前判断对应 Response 字段是否为 nil，为 nil 则不调用 set；将列表展开写入 `hardware_specifications`，使用 `helper.DataResourceIdsHash(ids)` 设置数据源 ID（ids 取各 HardwareSpecId）

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 中以名称 `tencentcloud_teo_inference_hardware_specifications` 注册新数据源
- [x] 3.2 在 `tencentcloud/provider.md` 中添加数据源文档条目

## 4. 文档实现

- [x] 4.1 创建 `tencentcloud/services/teo/data_source_tc_teo_inference_hardware_specifications.md` 文档文件，包含一句话描述（带云产品名称 TEO）、Example Usage、不带 Argument Reference/Attribute Reference 部分

## 5. 单元测试

- [x] 5.1 创建 `tencentcloud/services/teo/data_source_tc_teo_inference_hardware_specifications_test.go` 文件，使用 gomonkey mock `DescribeTeoInferenceHardwareSpecificationsByFilter` 方法，校验 flatten 后各字段值正确，不使用 Terraform 测试套件

## 6. 验证与收尾

- [x] 6.1 检查代码正确性：确认 CRUD 参数与云 API 接口参数一致、函数返回 error 均已检查、空 id 校验完整
- [ ] 6.2 通过收尾阶段 tfpacer-finalize skill 执行 `gofmt`、`make doc`、changelog 生成与 PR 创建
