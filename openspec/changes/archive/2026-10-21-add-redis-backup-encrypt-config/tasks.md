## 1. 资源实现

- [x] 1.1 创建 `tencentcloud/services/crs/resource_tc_redis_backup_encrypt.go`，定义 `ResourceTencentCloudRedisBackupEncrypt()`，包含 schema（`instance_id` Required+ForceNew、`backup_encrypt` Required bool、`encrypt_key_type` Optional int、`kms_key_id` Optional string、`kms_region` Optional string），声明 Create/Read/Update/Delete + Importer（ImportStatePassthrough）
- [x] 1.2 实现 `resourceTencentCloudRedisBackupEncryptCreate`：取 `instance_id` 后 `d.SetId(instance_id)`，委派 `resourceTencentCloudRedisBackupEncryptUpdate`
- [x] 1.3 实现 `resourceTencentCloudRedisBackupEncryptRead`：`request.InstanceId = helper.String(d.Id())`，`resource.Retry(tccommon.ReadRetryTimeout, ...)` 调用 `UseRedisClient().DescribeBackupEncryptConfigWithContext`；对实例不存在错误码或空配置先 `log.Printf("[CRUD] redis_backup_encrypt id=%s", d.Id())` 再 `d.SetId("")`；逐字段 nil 判断后 `d.Set`；其余错误用 `tccommon.RetryError` 包装
- [x] 1.4 实现 `resourceTencentCloudRedisBackupEncryptUpdate`：构造 `ModifyBackupEncryptRequest`，`request.InstanceId = helper.String(d.Id())`，`backup_encrypt` 用 `d.GetOk`，`encrypt_key_type` 用 `d.GetOkExists`，`kms_key_id`/`kms_region` 用 `d.GetOk`；`resource.Retry(tccommon.WriteRetryTimeout, ...)` 调用 `ModifyBackupEncryptWithContext`，错误用 `tccommon.RetryError` 包装；成功后 `return resourceTencentCloudRedisBackupEncryptRead(d, meta)`
- [x] 1.5 实现 `resourceTencentCloudRedisBackupEncryptDelete`：空操作直接 `return nil`（不清空 id、不调用 API）

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 的 redis 资源块新增 `"tencentcloud_redis_backup_encrypt": crs.ResourceTencentCloudRedisBackupEncrypt(),`
- [x] 2.2 在 `tencentcloud/provider.md` 的 redis Resource 列表新增 `tencentcloud_redis_backup_encrypt`

## 3. 文档

- [x] 3.1 创建 `tencentcloud/services/crs/resource_tc_redis_backup_encrypt.md`：一句话描述（带云产品名 Redis）、Example Usage、Import（说明以实例 ID 作为 id）
- [ ] 3.2 收尾阶段通过 `make doc` 生成 `website/docs/` 下文档（不在实施阶段手动编写 website 文件）

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/crs/resource_tc_redis_backup_encrypt_test.go`，使用 gomonkey mock `UseRedisClient` 与 `DescribeBackupEncryptConfigWithContext`/`ModifyBackupEncryptWithContext`，覆盖：Schema 字段校验、Read 成功、Read 实例不存在清空 state、Create 委派 Update、Update 成功刷新、Delete 空操作
- [x] 4.2 确保所有 error 返回值被检查，函数必定不出错处用 `_ = func()` 忽略 err

## 5. 验证（收尾阶段）

- [ ] 5.1 收尾阶段执行 `gofmt` 格式化变更的 Go 代码
- [ ] 5.2 收尾阶段执行 `make doc` 生成文档
- [ ] 5.3 收尾阶段在 `.changelog/` 下生成 changelog 文件