## 1. Schema 与 CRUD 逻辑

- [x] 1.1 在 `tencentcloud/services/mqtt/resource_tc_mqtt_instance.go` 的 Schema 中新增 `x509_mode` 参数（Optional, Computed, TypeString）
- [x] 1.2 在 Create 方法中，等待实例 RUNNING 后通过 ModifyInstance 设置 `X509Mode`
- [x] 1.3 在 Update 方法中，当 `x509_mode` 变更时通过 ModifyInstance 更新
- [x] 1.4 在 Read 方法中，从 DescribeInstance 响应读取 `X509Mode` 并设置到 state

## 2. 单元测试

- [x] 2.1 在 `tencentcloud/services/mqtt/resource_tc_mqtt_instance_test.go` 中新增 `x509_mode` 相关的单元测试（Read 正常值、Read nil、Schema 定义验证）

## 3. 文档与变更记录

- [x] 3.1 更新 `.changelog/4136.txt` 为 `x509_mode` 参数说明
