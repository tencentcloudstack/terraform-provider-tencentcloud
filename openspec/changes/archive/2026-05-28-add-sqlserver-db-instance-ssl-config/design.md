## Context

云 API：
- `ModifyDBInstanceSSL(InstanceId, Type, WaitSwitch, IsKMS, KeyId, KeyRegion)` → `FlowId`
- `DescribeDBInstancesAttribute` → `SSLConfig{Encryption, SSLValidityPeriod, SSLValidity, IsKMS, CMKId, CMKRegion}`
- `DescribeFlowStatus(FlowId)` → `Status`（0=成功）

## Decisions

1. **encryption 为声明式参数**: Required，用户声明期望状态 `enable`/`disable`
2. **移除 type**: Update 根据 encryption 值推导 API Type
3. **轮询使用 DescribeFlowStatus**: ModifyDBInstanceSSL 返回 FlowId，通过 DescribeFlowStatus 轮询直到 Status==0（成功）
4. **不支持 renew**: 暂不支持证书续期
5. **WaitSwitch=0**: 硬编码立即执行
6. **Read 映射**: enable/enable_doing/renew_doing → "enable"，其他 → "disable"

## Schema
- `instance_id`: Required, ForceNew
- `encryption`: Required (enable/disable)
- `is_kms`: Optional, Computed
- `cmk_id`: Optional, Computed
- `cmk_region`: Optional, Computed
- `ssl_validity_period`: Computed
- `ssl_validity`: Computed
