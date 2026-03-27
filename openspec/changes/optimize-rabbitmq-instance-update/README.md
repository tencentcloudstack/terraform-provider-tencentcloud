# 优化 RabbitMQ 实例的 update 逻辑

本提案旨在优化 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源的 update 逻辑,支持更多字段的动态更新,减少不必要的实例重建。

## 📋 概述

当前 RabbitMQ 实例资源的 update 逻辑过于严格,有 12 个字段被标记为不可变。但实际上,部分字段(如 `auto_renew_flag`、`enable_public_access`、`band_width`) 应该是可以通过 API 动态更新的。

## 🎯 目标

- ✅ 支持动态更新自动续费标志(`auto_renew_flag`)
- ✅ 支持动态切换公网访问开关(`enable_public_access`)
- ✅ 支持动态调整公网带宽(`band_width`)
- ✅ 优化错误提示,为不可变字段提供更友好的说明
- ✅ 保持向后兼容性,不影响现有配置

## 📁 文件结构

```
optimize-rabbitmq-instance-update/
├── .openspec.yaml                          # OpenSpec 元数据
├── proposal.md                             # 提案文档
├── tasks.md                               # 实现任务清单
└── specs/
    └── rabbitmq-instance-update/
        └── spec.md                        # 技术规范文档
```

## 🚀 快速开始

### 查看提案

```bash
cat proposal.md
```

### 查看任务清单

```bash
cat tasks.md
```

### 查看技术规范

```bash
cat specs/rabbitmq-instance-update/spec.md
```

## 📊 字段变更对比

### 不可变字段(6个)

| 字段名 | 原因 |
|--------|------|
| `zone_ids` | 集群可用区在创建后无法更改 |
| `vpc_id` | 私有网络配置在创建后无法更改 |
| `subnet_id` | 子网配置在创建后无法更改 |
| `time_span` | 购买时长无法直接修改 |
| `pay_mode` | 付费模式转换有特殊限制 |
| `cluster_version` | 集群版本升级需要专用流程 |

### 新增可变字段(3个)

| 字段名 | 类型 | API 字段 |
|--------|------|----------|
| `auto_renew_flag` | `bool` | `AutoRenewFlag` |
| `enable_public_access` | `bool` | `EnablePublicAccess` |
| `band_width` | `int` | `Bandwidth` |

### 保持可变字段(2个)

| 字段名 | 类型 | API 字段 |
|--------|------|----------|
| `cluster_name` | `string` | `ClusterName` |
| `resource_tags` | `list` | `Tags` |

## 🧪 测试示例

### 更新自动续费标志

```hcl
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  cluster_name       = "test-rabbitmq"
  zone_ids           = [1]
  vpc_id             = "vpc-xxxxxx"
  subnet_id          = "subnet-xxxxxx"
  auto_renew_flag    = false  # 从 true 改为 false
}
```

### 切换公网访问

```hcl
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  enable_public_access = true  # 启用公网访问
  band_width          = 20    # 设置带宽为 20 Mbps
}
```

## 🔍 实现细节

详见 `specs/rabbitmq-instance-update/spec.md`

## ✅ 验收标准

- [ ] 自动续费标志可以正常更新
- [ ] 公网访问开关可以正常切换
- [ ] 公网带宽可以正常调整
- [ ] 标签更新功能正常
- [ ] 不可变字段更新时返回友好错误提示
- [ ] 所有测试用例通过
- [ ] 向后兼容性验证通过
- [ ] 文档更新完成

## 📝 相关资源

- [TDMQ RabbitMQ API 文档](https://cloud.tencent.com/document/product/1491)
- [Terraform Provider 文档](https://registry.terraform.io/providers/tencentcloudstack/tencentcloud/latest/docs)

## 🤝 贡献

欢迎提交 Issue 和 Pull Request!

## 📄 许可证

遵循项目主仓库的许可证。
