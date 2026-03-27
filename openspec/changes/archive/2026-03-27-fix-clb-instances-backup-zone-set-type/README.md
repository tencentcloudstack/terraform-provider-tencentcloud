# Fix CLB Instances Backup Zone Set Data Type

## 概述

修复 `tencentcloud_clb_instances` 数据源中 `backup_zone_set` 字段的 schema 类型定义错误。

## 问题描述

当前 `backup_zone_set` 字段使用了不正确的类型定义 `&schema.Schema{Type: schema.TypeMap}`，应该使用 `&schema.Resource` 来定义复杂对象结构，以严格匹配 SDK 返回的 `[]*clb.ZoneInfo` 数据类型。

## 修复内容

- 将 `backup_zone_set` 的 schema 定义从简单的 `TypeMap` 改为使用 `&schema.Resource` 定义详细的字段结构
- 明确定义每个子字段的类型: `zone_id` (int), `zone` (string), `zone_name` (string), `zone_region` (string), `local_zone` (bool)
- 确保 schema 定义与 SDK 的 `ZoneInfo` 结构严格对应

## 文件清单

- `proposal.md` - 变更提案，详细说明修改原因和影响
- `tasks.md` - 实现任务清单，逐步实施指南
- `README.md` - 本文档，变更概述

## 实施步骤

1. 阅读 `proposal.md` 了解变更背景
2. 按照 `tasks.md` 中的步骤逐步实施
3. 修改完成后执行 `go fmt` 格式化代码
4. 编译验证确保无错误

## 相关文档

- API 文档: https://cloud.tencent.com/document/api/214/30685
- SDK 源码: `tencentcloud-sdk-go/tencentcloud/clb/v20180317/models.go`
- 数据源代码: `tencentcloud/services/clb/data_source_tc_clb_instances.go`

## 状态

- 创建日期: 2026-03-27
- 当前状态: 待实施
- 预计工时: 15-30分钟

## 实施命令

准备好后，运行以下命令开始实施:

```bash
/opsx:apply fix-clb-instances-backup-zone-set-type
```
