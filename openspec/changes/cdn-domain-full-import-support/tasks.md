## 1. 修改 Read 函数核心逻辑

- [x] 1.1 修改 `checkCdnInfoWritable` 函数（`resource_tc_cdn_domain.go` 第 4210 行）：将 `return val != nil && ok` 改为 `return val != nil`，移除 `helper.InterfacesHeadMap` state 存在性检查
- [x] 1.2 修复 `response_header_cache_switch` 字段（第 2953 行）：将 `if _, ok := d.GetOk("response_header_cache_switch"); ok && dc.ResponseHeaderCache != nil` 改为 `if dc.ResponseHeaderCache != nil`
- [x] 1.3 修复 `seo_switch` 字段（第 2963 行）：将 `if _, ok := d.GetOk("seo_switch"); ok && dc.Seo != nil` 改为 `if dc.Seo != nil`
- [x] 1.4 修复 `video_seek_switch` 字段（第 2987 行）：将 `if _, ok := d.GetOk("video_seek_switch"); ok && dc.VideoSeek != nil` 改为 `if dc.VideoSeek != nil`
- [x] 1.5 修复 `offline_cache_switch` 字段（第 3082 行）：将 `if _, ok := d.GetOk("offline_cache_switch"); ok && dc.OfflineCache != nil` 改为 `if dc.OfflineCache != nil`
- [x] 1.6 修复 `quic_switch` 字段（第 3085 行）：将 `if _, ok := d.GetOk("quic_switch"); ok && dc.Quic != nil` 改为 `if dc.Quic != nil`

## 2. 更新测试文件

- [x] 2.1 更新 `resource_tc_cdn_domain_test.go` 第一处 `ImportStateVerifyIgnore`：移除 `ip_filter`、`ip_freq_limit`、`compression`、`band_width_alert`、`error_page`、`downstream_capping`、`origin_pull_optimization`、`post_max_size`、`referer`、`max_age`、`cache_key`、`aws_private_access`、`oss_private_access`、`hw_private_access`、`qn_private_access`，保留 `full_url_cache`、`https_config`、`response_header`、`status_code_cache`
- [x] 2.2 更新 `resource_tc_cdn_domain_test.go` 第二处 `ImportStateVerifyIgnore`：移除 `cache_key`，保留 `https_config`、`authentication`、`full_url_cache`

## 3. 验证

- [x] 3.1 运行 `go build ./tencentcloud/services/cdn/...` 验证编译通过
- [x] 3.2 运行 `go vet ./tencentcloud/services/cdn/...` 验证无 lint 错误
