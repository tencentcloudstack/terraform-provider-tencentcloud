# OpenSpec: PostgreSQL StorageType 字段支持

## 项目信息

**项目名称**: 为 PostgreSQL 资源添加 StorageType 存储类型支持  
**提案编号**: POSTGRES-STORAGE-001  
**创建时间**: 2026-03-26  
**负责人**: AI Terraform Provider Expert  
**当前状态**: 提案阶段

---

## 文档结构

```
.openspec/
├── README.md           # 本文档 - 项目概览和使用指南
├── proposal.md         # 详细技术提案
├── tasks.md           # 任务分解和进度跟踪
└── IMPLEMENTATION.md  # 实施记录 (开发时创建)
```

---

## 快速开始

### 1. 阅读提案
查看 `proposal.md` 了解:
- 项目背景和目标
- 涉及的资源和修改点
- API 文档参考
- 风险评估

### 2. 查看任务清单
查看 `tasks.md` 了解:
- 详细的任务分解 (共 17 个任务)
- 每个任务的预计时间
- 优先级划分

### 3. 开始实施
按照以下顺序进行开发:

1. **服务层修改** (4 个任务) - 约 90 分钟
   - 修改 `service_tencentcloud_postgresql.go` 中的 4 个方法

2. **Resource 修改** (4 个任务) - 约 42 分钟
   - 修改 `resource_tc_postgresql_instance.go`

3. **DataSource 修改** (12 个任务) - 约 108 分钟
   - 修改 4 个 data source 文件

4. **测试验证** (3 个任务) - 约 85 分钟
   - 编译检查
   - Linter 检查
   - 本地测试 (可选)

5. **文档更新** (2 个任务) - 约 50 分钟
   - 更新实施文档
   - 创建示例配置

**总预计时间**: 4-6 小时

---

## 项目目标

### 功能目标
1. ✅ 用户可在创建 PostgreSQL 实例时指定存储类型
2. ✅ 用户可在读取实例属性时获取存储类型
3. ✅ 用户可在数据源查询时根据存储类型过滤

### 技术目标
1. ✅ 支持所有 4 种存储类型枚举值
2. ✅ 保持向后兼容性 (默认值为 `PHYSICAL_LOCAL_SSD`)
3. ✅ 符合 Terraform Provider 最佳实践

---

## 涉及的文件

| 文件 | 类型 | 修改内容 |
|------|------|----------|
| `service_tencentcloud_postgresql.go` | 服务层 | 添加 StorageType 参数支持 |
| `resource_tc_postgresql_instance.go` | Resource | Schema + Create + Read |
| `data_source_tc_postgresql_db_instance_versions.go` | DataSource | Schema + 参数传递 |
| `data_source_tc_postgresql_db_instance_classes.go` | DataSource | Schema + 参数传递 |
| `data_source_tc_postgresql_specinfos.go` | DataSource | Schema + 参数传递 |
| `data_source_tc_postgresql_instances.go` | DataSource | Schema (Computed 输出) |

**文件总数**: 6  
**代码行数**: 预计新增约 150 行,修改约 50 行

---

## StorageType 枚举值

| 值 | 描述 | 默认 |
|----|------|------|
| `PHYSICAL_LOCAL_SSD` | 物理机本地 SSD 硬盘 | ✅ |
| `CLOUD_PREMIUM` | 高性能云硬盘 | |
| `CLOUD_SSD` | SSD 云硬盘 | |
| `CLOUD_HSSD` | 增强型 SSD 云硬盘 | |

---

## API 映射关系

### CreateInstances (创建实例)
- **入参**: `StorageType` (String, Optional, Default: `PHYSICAL_LOCAL_SSD`)
- **文档**: https://cloud.tencent.com/document/api/409/56107

### DescribeDBInstanceAttribute (查询实例详情)
- **出参**: `DBInstanceStorageType` (String)
- **映射**: 对应入参 `StorageType`

### DescribeDBVersions (查询数据库版本)
- **入参**: `StorageType` (String, Optional)
- **文档**: https://cloud.tencent.com/document/api/409/89018

### DescribeClasses (查询售卖规格)
- **入参**: `StorageType` (String, Optional, Default: `PHYSICAL_LOCAL_SSD`)
- **文档**: https://cloud.tencent.com/document/api/409/89019

### DescribeProductConfig (查询售卖配置)
- **入参**: `StorageType` (String, Optional, Default: `PHYSICAL_LOCAL_SSD`)
- **文档**: https://cloud.tencent.com/document/api/409/16776

### DescribeDBInstances (查询实例列表)
- **限制**: ⚠️ 不支持 `StorageType` 作为查询过滤参数
- **出参**: 需确认是否包含存储类型字段
- **文档**: https://cloud.tencent.com/document/api/409/16773

---

## 使用示例

### 创建实例 (指定存储类型)

```hcl
resource "tencentcloud_postgresql_instance" "example" {
  name              = "tf-postgresql-test"
  storage_type      = "CLOUD_SSD"  # 使用 SSD 云硬盘
  availability_zone = "ap-guangzhou-3"
  vpc_id            = "vpc-12345678"
  subnet_id         = "subnet-12345678"
  memory            = 4
  storage           = 100
  engine_version    = "13.3"
  root_password     = "YourPassword123!"
  charge_type       = "POSTPAID_BY_HOUR"
}
```

### 查询支持的版本 (按存储类型)

```hcl
data "tencentcloud_postgresql_db_instance_versions" "cloud_hssd" {
  storage_type = "CLOUD_HSSD"
}

output "supported_versions" {
  value = data.tencentcloud_postgresql_db_instance_versions.cloud_hssd.version_set
}
```

### 查询售卖规格 (按存储类型)

```hcl
data "tencentcloud_postgresql_db_instance_classes" "premium" {
  zone              = "ap-guangzhou-7"
  db_engine         = "postgresql"
  db_major_version  = "15"
  storage_type      = "CLOUD_PREMIUM"
}

output "available_specs" {
  value = data.tencentcloud_postgresql_db_instance_classes.premium.class_info_set
}
```

---

## 注意事项

### ⚠️ 破坏性变更
- **无** - 所有新字段为 Optional/Computed,向后兼容

### ⚠️ ForceNew 行为
- `storage_type` 字段设置为 `ForceNew: true`
- 修改存储类型会导致实例重建

### ⚠️ API 限制
- `DescribeDBInstances` 可能不支持按存储类型过滤
- 需要确认该 API 是否返回存储类型字段

---

## 开发规范

### 代码格式
- **每修改完一个文件,立即执行 `go fmt`**

### 参数命名
- Schema 字段: `storage_type` (snake_case)
- Go 变量: `storageType` (camelCase)
- SDK 字段: `StorageType` (PascalCase)

### 验证函数
```go
ValidateFunc: tccommon.ValidateAllowedStringValue([]string{
    "PHYSICAL_LOCAL_SSD",
    "CLOUD_PREMIUM",
    "CLOUD_SSD",
    "CLOUD_HSSD",
})
```

### 默认值
```go
Default: "PHYSICAL_LOCAL_SSD"
```

---

## 测试策略

### 编译检查
```bash
go build ./tencentcloud/services/postgresql/...
```

### Linter 检查
```bash
golangci-lint run ./tencentcloud/services/postgresql/...
```

### 手动测试 (可选)
1. 创建不同存储类型的实例
2. 验证数据源查询结果
3. 验证实例属性读取

---

## 进度跟踪

使用 `tasks.md` 中的复选框跟踪进度:

- [ ] 阶段 1: 服务层修改 (0/4)
- [ ] 阶段 2: Resource 修改 (0/4)
- [ ] 阶段 3: DataSource - db_instance_versions (0/3)
- [ ] 阶段 4: DataSource - db_instance_classes (0/3)
- [ ] 阶段 5: DataSource - specinfos (0/3)
- [ ] 阶段 6: DataSource - instances (0/4)
- [ ] 阶段 7: 测试与验证 (0/3)
- [ ] 阶段 8: 文档更新 (0/2)

**总进度**: 0/26 (0%)

---

## 参考资料

### 腾讯云 API 文档
- [CreateInstances](https://cloud.tencent.com/document/api/409/56107)
- [DescribeDBVersions](https://cloud.tencent.com/document/api/409/89018)
- [DescribeClasses](https://cloud.tencent.com/document/api/409/89019)
- [DescribeProductConfig](https://cloud.tencent.com/document/api/409/16776)
- [DescribeDBInstances](https://cloud.tencent.com/document/api/409/16773)

### Terraform Provider 开发
- [Terraform Plugin SDK v2](https://developer.hashicorp.com/terraform/plugin/sdkv2)
- [Schema Best Practices](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas)

---

## 联系方式

如有问题,请参考:
- `proposal.md` - 详细技术方案
- `tasks.md` - 任务分解和时间估算
- API 文档 - 字段定义和使用说明

---

**最后更新**: 2026-03-26  
**文档版本**: 1.0
