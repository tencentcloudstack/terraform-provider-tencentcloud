## 1. 资源代码实现

- [x] 1.1 创建 `tencentcloud/services/ssl/resource_tc_ssl_upload_certificate_instance_config.go`，定义 `ResourceTencentCloudSslUploadCertificateInstance()`，schema 包含入参 `old_certificate_id`(Required)、`resource_types`(Required, TypeSet)、`certificate_public_key`(Required)、`certificate_private_key`(Required, Sensitive)、`resource_types_regions`(Optional, TypeList)、`certificate_id`(Optional)，及 Computed 出参 `deploy_record_id`、`deploy_status`、`update_sync_progress`(TypeList，含子 list `update_sync_progress_regions`)、`deploy_record_list`(TypeList，展平不套第二层)；声明 Importer(State passthrough)
- [x] 1.2 实现 `resourceTencentCloudSslUploadCertificateInstanceCreate`/`Update`：映射入参到 `UploadUpdateCertificateInstanceRequest`，外包 `resource.Retry(tccommon.WriteRetryTimeout)`，retry 块内只调接口并判空 `RetryableError`；成功后于 retry 外取 `DeployRecordId`，校验非空打印 logId/Id，先 `d.SetId()`，保留 `deploy_status`/`update_sync_progress`，再调 Read 轮询生效
- [x] 1.3 实现 `resourceTencentCloudSslUploadCertificateInstanceRead`：调用 `DescribeHostUploadUpdateRecord`（`OldCertificateId` 过滤，内部 `Limit=200` 分页），外包 `resource.Retry(tccommon.ReadRetryTimeout)`；空结果先 `log.Printf("[CRUD] ssl_upload_certificate_instance id=%s", d.Id())` 再 `d.SetId("")`；从 `DeployRecordList` 按 id 匹配记录回填字段（nil 判断后 setXx），并增加调用 `DescribeHostUploadUpdateRecordDetail`（`DeployRecordId`=id，`Limit=200`）获取详情
- [x] 1.4 实现 `resourceTencentCloudSslUploadCertificateInstanceDelete`：解析 `d.Id()` 为 `DeployRecordId`(int64)，调用 `UploadUpdateCertificateRecordRollback`，外包 `resource.Retry(tccommon.WriteRetryTimeout)`，失败用 `RetryError` 包装

## 2. Provider 注册与文档

- [x] 2.1 在 `tencentcloud/provider.go` 的 resources map 中新增 `"tencentcloud_ssl_upload_certificate_instance": ssl.ResourceTencentCloudSslUploadCertificateInstance()`
- [x] 2.2 创建 `tencentcloud/services/ssl/resource_tc_ssl_upload_certificate_instance_config.md` 文档（含一句话描述带产品名、Example Usage、Import 部分）
- [x] 2.3 在 `tencentcloud/provider.md` 中补齐资源文档条目（由 make doc 生成，禁止手编 website/）

## 3. 单元测试

- [x] 3.1 创建 `tencentcloud/services/ssl/resource_tc_ssl_upload_certificate_instance_config_test.go`，使用 gomonkey mock 云 API（`UploadUpdateCertificateInstance`/`DescribeHostUploadUpdateRecord`/`DescribeHostUploadUpdateRecordDetail`/`UploadUpdateCertificateRecordRollback`），对 Create/Read/Update/Delete 业务逻辑进行单元测试

## 4. 验证

- [ ] 4.1 执行 make doc 生成 website 文档，验证文档生成
- [x] 4.2 验证编排：go 代码可正确构建执行