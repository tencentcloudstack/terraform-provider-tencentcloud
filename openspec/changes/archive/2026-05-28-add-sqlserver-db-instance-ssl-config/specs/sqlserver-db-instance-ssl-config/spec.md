## ADDED Requirements

### Requirement: encryption is declarative desired state
- `encryption` (Required): 用户声明 SSL 期望状态 `enable`/`disable`
- Update 根据 encryption 值决定 API 的 Type 参数
- Read 从 SSLConfig.Encryption 映射回 enable/disable

### Requirement: Polling uses DescribeFlowStatus
- ModifyDBInstanceSSL 返回 FlowId
- 通过 DescribeFlowStatus 轮询 FlowId 直到 Status==0（SQLSERVER_TASK_SUCCESS）

### Requirement: No renew support
- 暂不支持 type=renew 操作

### Requirement: WaitSwitch hardcoded to 0
- ModifyDBInstanceSSL 调用时 WaitSwitch 固定为 0
