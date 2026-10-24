## Why

SSL 证书的云资源更新（证书 ID 不变）能力当前在 Terraform Provider 中缺失。用户需要通过 `UploadUpdateCertificateInstance` 接口在保持旧证书 ID 不变的前提下更新证书内容并联动更新已关联的云资源（如 clb），但现有 provider 仅提供了数据源和操作类资源用于部署/回滚，无法以"配置即资源"的方式托管该上传更新行为及查询其更新记录。新增一个 `RESOURCE_KIND_CONFIG` 资源 `tencentcloud_ssl_upload_certificate_instance`，可以让用户以 terraform 声明式地管理证书内容更新及关联云资源同步，并支持对更新记录的读取。

## What Changes

- 新增 Terraform 资源 `tencentcloud_ssl_upload_certificate_instance`，类型为 `RESOURCE_KIND_CONFIG`，文件名为 `resource_tc_ssl_upload_certificate_instance_config.go`（注册函数 `ResourceTencentCloudSslUploadCertificateInstance`）。
- Create/Update：调用 `UploadUpdateCertificateInstance`，上传证书公私钥并指定资源类型/地域，更新关联云资源（证书 ID 不变）。接口为异步，需在拿到 `DeployRecordId` 后先 `d.SetId()`，再轮询 `DescribeHostUploadUpdateRecord` 等待任务生效。
- Read：调用 `DescribeHostUploadUpdateRecord` 按 `OldCertificateId` 查询证书上传更新记录列表，回填 `deploy_record_list` 等出参；查询接口增加 `DescribeHostUploadUpdateRecordDetail` 用于按 `DeployRecordId` 查询更新记录详情。
- Delete：该资源为配置型资源，删除时调用 `UploadUpdateCertificateRecordRollback` 进行回滚，使配置可被销毁清理。
- 在 `tencentcloud/provider.go` 和 `tencentcloud/provider.md` 中注册新资源。
- 新增 `.md` 文档（make doc 生成），及单元测试文件 `resource_tc_ssl_upload_certificate_instance_config_test.go`（使用 gomonkey mock 云 API）。

## Capabilities

### New Capabilities
- `ssl-upload-certificate-instance`: 管理 SSL 证书内容上传更新（证书 ID 不变）及其关联云资源同步更新，支持读取更新记录与详情。

### Modified Capabilities
<!-- 无 -->

## Impact

- 新增代码文件：`tencentcloud/services/ssl/resource_tc_ssl_upload_certificate_instance_config.go`、`..._test.go`、`resource_tc_ssl_upload_certificate_instance_config.md`。
- 修改 `tencentcloud/provider.go`、`tencentcloud/provider.md`：注册新资源。
- 涉及云 API（package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205`）：
  - `UploadUpdateCertificateInstance`（异步更新）
  - `DescribeHostUploadUpdateRecord`（查询更新记录列表）
  - `DescribeHostUploadUpdateRecordDetail`（查询更新记录详情）
  - `UploadUpdateCertificateRecordRollback`（删除/回滚）
- 不影响已有资源与数据源，完全为新增能力，向后兼容。