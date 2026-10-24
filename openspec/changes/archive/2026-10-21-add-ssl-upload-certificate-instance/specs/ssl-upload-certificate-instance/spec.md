## ADDED Requirements

### Requirement: Resource manages SSL certificate upload-update with id unchanged
The provider SHALL provide a resource `tencentcloud_ssl_upload_certificate_instance` (RESOURCE_KIND_CONFIG) that uploads new certificate public/private key content and updates associated cloud resources without changing the certificate id, via the cloud API `UploadUpdateCertificateInstance`.

#### Scenario: Create/upload update succeeds
- **WHEN** the user applies a configuration with `old_certificate_id`, `resource_types`, `certificate_public_key`, `certificate_private_key`
- **THEN** the provider calls `UploadUpdateCertificateInstance`, sets the resource id to the returned `DeployRecordId`, and polls Read until the update takes effect

#### Scenario: Update changes update content
- **WHEN** the user changes `certificate_public_key`/`certificate_private_key`/`resource_types`/`resource_types_regions`/`certificate_id` on an existing resource
- **THEN** the provider calls `UploadUpdateCertificateInstance` again with the new content

### Requirement: Read returns update record list and detail
The provider SHALL read the resource state by calling `DescribeHostUploadUpdateRecord` (filtered by `OldCertificateId`, paging internally with Limit=200) and SHALL additionally call `DescribeHostUploadUpdateRecordDetail` (by `DeployRecordId`) to obtain update record details.

#### Scenario: Record found
- **WHEN** the resource is read and `DescribeHostUploadUpdateRecord` returns a matching record
- **THEN** the provider populates `deploy_record_list` (id, old_cert_id, resource_types, status, create_time, update_time) and computed fields `deploy_record_id`, `deploy_status`, `update_sync_progress`

#### Scenario: Record not found
- **WHEN** the read returns an empty list
- **THEN** the provider logs the id via `log.Printf("[CRUD] ssl_upload_certificate_instance id=%s", d.Id())` and clears the id with `d.SetId("")`

### Requirement: Destroy rolls back update record
The provider SHALL support destroying the resource by calling `UploadUpdateCertificateRecordRollback` with the `DeployRecordId` parsed from the resource id.

#### Scenario: Destroy succeeds
- **WHEN** the user runs `terraform destroy` on the resource
- **THEN** the provider calls `UploadUpdateCertificateRecordRollback` and removes the resource from state

### Requirement: Async update sets id before polling
For the asynchronous `UploadUpdateCertificateInstance` call, the provider SHALL call `d.SetId(<id>)` immediately after obtaining the `DeployRecordId` and before polling for completion, so that a failed async task still leaves an id in tfstate for cleanup.

#### Scenario: Async task fails after id assignment
- **WHEN** the update task id is obtained but polling subsequently fails
- **THEN** the resource id remains in tfstate, allowing the user to retry destroy

### Requirement: Resource registered in provider
The provider SHALL register `tencentcloud_ssl_upload_certificate_instance` in `tencentcloud/provider.go` and document it in `tencentcloud/provider.md`.

#### Scenario: Provider registration
- **WHEN** the change is implemented
- **THEN** `tencentcloud_ssl_upload_certificate_instance` maps to `ssl.ResourceTencentCloudSslUploadCertificateInstance()` in the resources map