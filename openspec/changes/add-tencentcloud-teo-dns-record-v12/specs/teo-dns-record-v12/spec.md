## ADDED Requirements

### Requirement: 用户可以创建 TEO DNS 记录
系统 SHALL 支持用户通过 Terraform 配置创建 TEO DNS 记录。用户必须提供站点 ID（zone_id）、DNS 记录名（name）、记录类型（type）和记录内容（content）。系统支持设置解析线路（location）、缓存时间（ttl）、权重（weight）和 MX 优先级（priority）等可选参数。创建成功后，系统返回 DNS 记录 ID（record_id）。

#### Scenario: 成功创建 A 类型 DNS 记录
- **WHEN** 用户配置 `zone_id` 为 `zone-12345`，`name` 为 `www.example.com`，`type` 为 `A`，`content` 为 `1.2.3.4`，`ttl` 为 `300`
- **THEN** 系统调用 `CreateDnsRecord` 接口创建 DNS 记录，返回 `record_id` 为 `record-67890`，状态为 `enable`

#### Scenario: 成功创建 CNAME 类型 DNS 记录并设置权重
- **WHEN** 用户配置 `zone_id` 为 `zone-12345`，`name` 为 `alias.example.com`，`type` 为 `CNAME`，`content` 为 `target.example.com`，`weight` 为 `50`，`location` 为 `overseas`
- **THEN** 系统调用 `CreateDnsRecord` 接口创建 DNS 记录，返回 `record_id`，权重设置为 50，解析线路为 `overseas`

#### Scenario: 成功创建 MX 类型 DNS 记录并设置优先级
- **WHEN** 用户配置 `zone_id` 为 `zone-12345`，`name` 为 `@`，`type` 为 `MX`，`content` 为 `mail.example.com`，`priority` 为 `10`
- **THEN** 系统调用 `CreateDnsRecord` 接口创建 DNS 记录，返回 `record_id`，优先级设置为 10

#### Scenario: 创建 DNS 记录时缺少必需参数
- **WHEN** 用户配置缺少 `zone_id`、`name`、`type` 或 `content` 中的任意一个
- **THEN** 系统返回验证错误，提示缺少必需参数

### Requirement: 用户可以读取 TEO DNS 记录状态
系统 SHALL 支持用户读取已创建的 TEO DNS 记录的完整状态信息，包括记录 ID、记录名、记录类型、记录内容、解析线路、缓存时间、权重、优先级、状态和创建时间。系统通过站点 ID 和记录 ID 定位到具体的 DNS 记录。

#### Scenario: 成功读取已存在的 DNS 记录
- **WHEN** 用户读取 `zone_id` 为 `zone-12345`，`record_id` 为 `record-67890` 的 DNS 记录
- **THEN** 系统调用 `DescribeDnsRecords` 接口查询，返回该记录的所有字段，包括 `name`、`type`、`content`、`location`、`ttl`、`weight`、`priority`、`status` 和 `created_on`

#### Scenario: 读取不存在的 DNS 记录
- **WHEN** 用户读取 `zone_id` 为 `zone-12345`，`record_id` 为 `record-nonexist` 的 DNS 记录
- **THEN** 系统返回记录不存在的错误

#### Scenario: 读取 DNS 记录时站点 ID 不匹配
- **WHEN** 用户读取 `zone_id` 为 `zone-wrong`，`record_id` 为 `record-67890` 的 DNS 记录
- **THEN** 系统返回记录不存在的错误（因为记录 ID 和站点 ID 不匹配）

### Requirement: 用户可以更新 TEO DNS 记录配置
系统 SHALL 支持用户更新已创建的 TEO DNS 记录的配置信息。用户可以更新记录名、记录类型、记录内容、解析线路、缓存时间、权重和优先级等字段。更新操作通过 `ModifyDnsRecords` 接口完成，更新后系统立即读取记录状态以确保配置生效。

#### Scenario: 成功更新 DNS 记录的内容和 TTL
- **WHEN** 用户更新 `zone_id` 为 `zone-12345`，`record_id` 为 `record-67890` 的记录，将 `content` 更新为 `5.6.7.8`，`ttl` 更新为 `600`
- **THEN** 系统调用 `ModifyDnsRecords` 接口更新记录，读取操作确认记录的 `content` 已更新为 `5.6.7.8`，`ttl` 已更新为 `600`

#### Scenario: 成功更新 DNS 记录的权重和解析线路
- **WHEN** 用户更新 `zone_id` 为 `zone-12345`，`record_id` 为 `record-67890` 的记录，将 `weight` 更新为 `80`，`location` 更新为 `mainland`
- **THEN** 系统调用 `ModifyDnsRecords` 接口更新记录，读取操作确认记录的 `weight` 已更新为 `80`，`location` 已更新为 `mainland`

#### Scenario: 更新不存在的 DNS 记录
- **WHEN** 用户更新 `zone_id` 为 `zone-12345`，`record_id` 为 `record-nonexist` 的记录
- **THEN** 系统返回记录不存在的错误

#### Scenario: 更新 DNS 记录时设置不合法的参数值
- **WHEN** 用户更新 `zone_id` 为 `zone-12345`，`record_id` 为 `record-67890` 的记录，将 `ttl` 设置为 `50`（小于最小值 60）或 `weight` 设置为 `101`（大于最大值 100）
- **THEN** 系统返回参数校验错误

### Requirement: 用户可以删除 TEO DNS 记录
系统 SHALL 支持用户删除已创建的 TEO DNS 记录。用户通过站点 ID 和记录 ID 定位到具体的 DNS 记录并执行删除操作。删除操作通过 `DeleteDnsRecords` 接口完成，删除成功后资源状态从 Terraform state 中移除。

#### Scenario: 成功删除已存在的 DNS 记录
- **WHEN** 用户删除 `zone_id` 为 `zone-12345`，`record_id` 为 `record-67890` 的 DNS 记录
- **THEN** 系统调用 `DeleteDnsRecords` 接口删除记录，Terraform state 中该资源被标记为已删除

#### Scenario: 删除不存在的 DNS 记录
- **WHEN** 用户删除 `zone_id` 为 `zone-12345`，`record_id` 为 `record-nonexist` 的 DNS 记录
- **THEN** 系统返回记录不存在的错误，或者静默成功（取决于 API 行为）

#### Scenario: 删除已被其他资源依赖的 DNS 记录
- **WHEN** 用户删除一个正在被其他配置（如负载均衡器）引用的 DNS 记录
- **THEN** 系统返回记录被引用的错误，或者成功删除（取决于 API 行为）

### Requirement: 系统 MUST 支持多种 DNS 记录类型
系统 SHALL 支持创建和管理多种类型的 DNS 记录，包括 A、AAAA、MX、CNAME、TXT、NS、CAA、SRV 等类型。不同记录类型对参数格式和约束有不同的要求。

#### Scenario: 创建 AAAA 类型 DNS 记录
- **WHEN** 用户配置 `type` 为 `AAAA`，`content` 为 `2001:0db8:85a3:0000:0000:8a2e:0370:7334`
- **THEN** 系统成功创建 IPv6 地址类型的 DNS 记录

#### Scenario: 创建 TXT 类型 DNS 记录
- **WHEN** 用户配置 `type` 为 `TXT`，`content` 为 `"v=spf1 include:_spf.example.com ~all"`
- **THEN** 系统成功创建 TXT 类型的 DNS 记录，用于 SPF 记录配置

#### Scenario: 创建 SRV 类型 DNS 记录
- **WHEN** 用户配置 `type` 为 `SRV`，`name` 为 `_sip._tcp.example.com`，`content` 为 `10 60 5060 sipserver.example.com`
- **THEN** 系统成功创建 SRV 类型的 DNS 记录

### Requirement: 系统 MUST 支持可选参数的默认值处理
系统 SHALL 为未明确配置的可选参数使用云 API 提供的默认值。对于 `ttl` 参数，默认值为 300 秒；对于 `weight` 参数，默认值为 -1（不设置权重）；对于 `priority` 参数，默认值为 0；对于 `location` 参数，默认值为 `Default`（默认解析线路）。

#### Scenario: 创建 DNS 记录时未设置 TTL
- **WHEN** 用户配置时未指定 `ttl` 参数
- **THEN** 系统使用云 API 的默认 TTL 值 300 秒创建记录

#### Scenario: 创建 DNS 记录时未设置权重
- **WHEN** 用户配置时未指定 `weight` 参数
- **THEN** 系统使用云 API 的默认权重值 -1（不设置权重）创建记录

#### Scenario: 创建 DNS 记录时未设置解析线路
- **WHEN** 用户配置时未指定 `location` 参数
- **THEN** 系统使用云 API 的默认解析线路 `Default` 创建记录

### Requirement: 系统 MUST 处理最终一致性和异步操作
系统 SHALL 支持对异步操作和最终一致性的处理。对于创建、更新和删除操作，系统在 API 调用成功后，通过 Read 接口轮询直到记录状态与期望状态一致。系统使用 `helper.Retry()` 提供合理的重试机制，并使用 `defer tccommon.InconsistentCheck()` 确保状态一致性。

#### Scenario: 创建 DNS 记录后读取到一致的记录状态
- **WHEN** 用户创建 DNS 记录成功后，系统立即调用 Read 接口
- **THEN** 系统通过重试机制读取到完整的记录信息，包括 `record_id`、`status` 和其他字段

#### Scenario: 更新 DNS 记录后读取到一致的记录状态
- **WHEN** 用户更新 DNS 记录成功后，系统立即调用 Read 接口
- **THEN** 系统通过重试机制读取到更新后的记录信息，确认所有修改已生效

#### Scenario: 删除 DNS 记录后确认记录不存在
- **WHEN** 用户删除 DNS 记录成功后，系统立即调用 Read 接口
- **THEN** 系统通过重试机制确认记录已被删除，返回记录不存在的状态

### Requirement: 系统 MUST 提供完整的错误处理和日志记录
系统 SHALL 提供完善的错误处理和日志记录机制。系统使用 `defer tccommon.LogElapsed()` 记录每个 API 调用的耗时，使用 `defer tccommon.InconsistentCheck()` 检测状态不一致的情况，并将云 API 返回的错误信息准确地传递给用户。

#### Scenario: API 调用超时时返回超时错误
- **WHEN** 云 API 调用超过预设的超时时间
- **THEN** 系统返回超时错误，包含详细的错误信息和操作耗时

#### Scenario: API 调用失败时返回详细的错误信息
- **WHEN** 云 API 调用因参数错误或权限问题失败
- **THEN** 系统返回云 API 的错误信息，包括错误码和错误描述

#### Scenario: 记录操作耗时被正确记录
- **WHEN** 用户执行创建、读取、更新或删除操作
- **THEN** 系统记录每个操作的耗时，用于监控和调试
