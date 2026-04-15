## ADDED Requirements

### Requirement: Create DNS Record
用户能够通过 Terraform 创建一个新的 TEO DNS 记录，系统应调用云 API `CreateDnsRecord` 创建记录，并返回记录 ID。

#### Scenario: 成功创建 A 类型 DNS 记录
- **WHEN** 用户在 Terraform 配置中指定 zone_id、name、type 为 "A"、content 为 IP 地址
- **THEN** 系统调用 CreateDnsRecord API 创建记录
- **THEN** 系统将返回的 record_id 保存到 Terraform 状态中
- **THEN** 资源 ID 格式为 "zone_id#record_id"

#### Scenario: 成功创建 CNAME 类型 DNS 记录
- **WHEN** 用户在 Terraform 配置中指定 zone_id、name、type 为 "CNAME"、content 为域名
- **THEN** 系统调用 CreateDnsRecord API 创建记录
- **THEN** 系统将返回的 record_id 保存到 Terraform 状态中

#### Scenario: 创建记录时指定 TTL
- **WHEN** 用户在 Terraform 配置中指定 ttl 参数（范围 60-86400）
- **THEN** 系统在创建记录时包含 TTL 字段
- **THEN** 云 API 接受 TTL 参数并生效

#### Scenario: 创建记录时指定权重
- **WHEN** 用户在 Terraform 配置中指定 weight 参数（范围 -1~100）
- **AND** 记录类型为 A、AAAA 或 CNAME
- **THEN** 系统在创建记录时包含 Weight 字段
- **THEN** 云 API 接受 Weight 参数并生效

#### Scenario: 创建记录时指定优先级
- **WHEN** 用户在 Terraform 配置中指定 priority 参数（范围 0~50）
- **AND** 记录类型为 MX
- **THEN** 系统在创建记录时包含 Priority 字段
- **THEN** 云 API 接受 Priority 参数并生效

#### Scenario: 创建记录失败
- **WHEN** 云 API 返回错误（如参数错误、权限不足）
- **THEN** 系统返回清晰的错误信息给用户
- **THEN** 资源创建失败，不修改 Terraform 状态

### Requirement: Read DNS Record
用户能够通过 Terraform 读取已存在的 TEO DNS 记录信息，系统应调用云 API `DescribeDnsRecords` 查询记录并返回完整的记录信息。

#### Scenario: 成功读取 DNS 记录
- **WHEN** Terraform 执行 refresh 操作
- **AND** 资源已存在于云 API 中
- **THEN** 系统调用 DescribeDnsRecords API，使用 record_id 过滤条件
- **THEN** 系统更新 Terraform 状态中的所有字段
- **THEN** 包含 computed 字段如 record_id、zone_id

#### Scenario: 读取不存在的记录
- **WHEN** Terraform 执行 refresh 操作
- **AND** 资源在云 API 中不存在
- **THEN** 系统返回 nil 或标记资源为已删除
- **THEN** Terraform 状态显示资源不存在

#### Scenario: 读取记录时网络错误
- **WHEN** 调用 DescribeDnsRecords API 时发生网络错误
- **THEN** 系统返回错误信息
- **THEN** Terraform 操作失败，但状态保持不变

### Requirement: Update DNS Record
用户能够通过 Terraform 更新已存在的 TEO DNS 记录，系统应调用云 API `ModifyDnsRecords` 更新记录的指定字段。

#### Scenario: 成功更新 DNS 记录内容
- **WHEN** 用户在 Terraform 配置中修改 content 字段
- **THEN** 系统从状态中读取当前记录的完整数据
- **THEN** 系统调用 ModifyDnsRecords API，包含更新后的 content 和其他未修改字段
- **THEN** 云 API 接受修改并生效
- **THEN** Terraform 状态更新为新的内容

#### Scenario: 成功更新 DNS 记录的 TTL
- **WHEN** 用户在 Terraform 配置中修改 ttl 字段
- **THEN** 系统调用 ModifyDnsRecords API，包含更新后的 ttl
- **THEN** 云 API 接受修改并生效

#### Scenario: 更新记录时字段未变化
- **WHEN** 用户执行 apply 操作
- **AND** 所有字段与当前状态相同
- **THEN** 系统不调用 ModifyDnsRecords API
- **THEN** Terraform 显示 "No changes" 消息

#### Scenario: 更新不存在的记录
- **WHEN** 用户尝试更新一个已删除的记录
- **THEN** 系统返回错误，指示记录不存在
- **THEN** Terraform 操作失败

### Requirement: Delete DNS Record
用户能够通过 Terraform 删除已存在的 TEO DNS 记录，系统应调用云 API `DeleteDnsRecords` 删除记录。

#### Scenario: 成功删除 DNS 记录
- **WHEN** 用户执行 terraform destroy 操作
- **AND** 资源存在于云 API 中
- **THEN** 系统调用 DeleteDnsRecords API，传入 record_id
- **THEN** 云 API 删除记录
- **THEN** Terraform 从状态中移除该资源

#### Scenario: 删除不存在的记录
- **WHEN** 用户执行 terraform destroy 操作
- **AND** 资源在云 API 中不存在
- **THEN** 系统忽略删除操作
- **THEN** Terraform 状态正常更新

#### Scenario: 删除记录时权限不足
- **WHEN** 调用 DeleteDnsRecords API 时返回权限错误
- **THEN** 系统返回错误信息给用户
- **THEN** 删除操作失败，资源仍保留在状态中

### Requirement: Resource ID Handling
系统应正确处理资源的唯一标识符，使用复合 ID 格式 "zone_id#record_id" 作为 Terraform 资源 ID。

#### Scenario: 解析复合 ID
- **WHEN** 系统需要从资源 ID 中提取 zone_id 和 record_id
- **THEN** 系统按 "#" 分隔符分割 ID
- **THEN** 第一部分为 zone_id，第二部分为 record_id
- **THEN** 如果 ID 格式不正确，返回错误

#### Scenario: 构造复合 ID
- **WHEN** 系统创建新资源时
- **THEN** 系统使用 zone_id 和 record_id 构造复合 ID
- **THEN** ID 格式为 "zone_id#record_id"
- **THEN** ID 保存到 Terraform 状态中

### Requirement: Parameter Validation
系统应对用户输入的参数进行基本的校验，确保参数在合理的范围内。

#### Scenario: TTL 参数校验
- **WHEN** 用户输入的 ttl 值不在 60-86400 范围内
- **THEN** 系统返回验证错误
- **THEN** 提示 TTL 的有效范围

#### Scenario: Weight 参数校验
- **WHEN** 用户输入的 weight 值不在 -1~100 范围内
- **THEN** 系统返回验证错误
- **THEN** 提示 Weight 的有效范围

#### Scenario: Priority 参数校验
- **WHEN** 用户输入的 priority 值不在 0~50 范围内
- **THEN** 系统返回验证错误
- **THEN** 提示 Priority 的有效范围

#### Scenario: Type 参数校验
- **WHEN** 用户输入的 type 值不是有效的 DNS 记录类型
- **THEN** 系统返回验证错误
- **THEN** 提示有效的类型值（A、AAAA、CNAME、MX、TXT、NS、CAA、SRV）

### Requirement: Error Handling
系统应正确处理云 API 返回的错误，并提供友好的错误信息给用户。

#### Scenario: 处理网络超时错误
- **WHEN** 调用云 API 时发生网络超时
- **THEN** 系统重试操作（使用 helper.Retry）
- **THEN** 如果重试失败，返回超时错误
- **THEN** 错误信息包含重试次数和超时时间

#### Scenario: 处理参数错误
- **WHEN** 云 API 返回参数错误（如 InvalidParameter）
- **THEN** 系统提取错误消息
- **THEN** 返回用户友好的错误信息
- **THEN** 不修改 Terraform 状态

#### Scenario: 处理权限错误
- **WHEN** 云 API 返回权限错误（如 AuthFailure）
- **THEN** 系统返回权限错误信息
- **THEN** 提示用户检查凭证配置

#### Scenario: 处理资源不存在错误
- **WHEN** 云 API 返回资源不存在错误（如 ResourceNotFound）
- **THEN** 系统标记资源为已删除
- **THEN** 在 Read 操作中返回 nil
