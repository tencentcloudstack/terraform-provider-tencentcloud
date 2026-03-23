## ADDED Requirements

### Requirement: CDN 域名 import 后所有可读字段写入 state
执行 `terraform import tencentcloud_cdn_domain.example domain.com` 后，所有 API 可返回的字段 SHALL 被正确写入 state，不得因 state 为空而跳过任何字段。

#### Scenario: import 后 checkCdnInfoWritable 守卫字段正确写入
- **WHEN** 执行 `terraform import` 导入已有 CDN 域名
- **THEN** `ip_filter`、`ip_freq_limit`、`compression`、`band_width_alert`、`error_page`、`downstream_capping`、`origin_pull_optimization`、`referer`、`max_age`、`specific_config_mainland`、`specific_config_overseas`、`origin_pull_timeout`、`post_max_size`、`cache_key`、`aws_private_access`、`oss_private_access`、`hw_private_access`、`qn_private_access`、`others_private_access` 均被写入 state（若 API 返回非 nil）

#### Scenario: import 后 switch 字段正确写入
- **WHEN** 执行 `terraform import` 导入已有 CDN 域名
- **THEN** `response_header_cache_switch`、`seo_switch`、`video_seek_switch`、`offline_cache_switch`、`quic_switch` 均被写入 state（若 API 返回非 nil）

#### Scenario: import 后 plan diff 为空
- **WHEN** 执行 `terraform import` 后立即执行 `terraform plan`
- **THEN** 对于 API 可返回的字段，plan 输出中不应出现 diff（`No changes` 或仅有 API 不返回的字段差异）

### Requirement: 存量用户 apply 流程不受影响
修改 `checkCdnInfoWritable` 后，存量用户（state 中已有字段值）的正常 apply 流程 SHALL 与修改前行为完全一致。

#### Scenario: 存量用户 apply 时字段正常写入
- **WHEN** 存量用户执行 `terraform apply`（state 中已有 `ip_filter` 等字段值）
- **THEN** Read 函数正常将 API 返回值写入 state，行为与修改前一致

#### Scenario: API 返回 nil 时不写入 state
- **WHEN** API 返回某字段为 nil（该功能未配置）
- **THEN** 该字段不被写入 state，避免覆盖用户配置为空的情况
