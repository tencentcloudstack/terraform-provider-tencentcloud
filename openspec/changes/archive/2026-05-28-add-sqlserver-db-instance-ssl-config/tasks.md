## 1. Schema 与 CRUD

- [x] 1.1 Schema: instance_id(Required+ForceNew), encryption(Required, enable/disable), is_kms/cmk_id/cmk_region(Optional+Computed), ssl_validity_period/ssl_validity(Computed)
- [x] 1.2 Create: 设置 ID，调用 Update
- [x] 1.3 Read: 从 SSLConfig.Encryption 映射为 enable/disable 设置到 encryption；读取 is_kms/cmk_id/cmk_region
- [x] 1.4 Update: 根据 encryption 期望值设置 Type(enable/disable)，调用 ModifyDBInstanceSSL，轮询 DescribeSqlserverInstanceSslById 直到状态达标
- [x] 1.5 Delete: 空操作

## 2. 单元测试

- [x] 2.1 覆盖 Read、Create(enable)、Update(disable)、Delete、Schema 验证
