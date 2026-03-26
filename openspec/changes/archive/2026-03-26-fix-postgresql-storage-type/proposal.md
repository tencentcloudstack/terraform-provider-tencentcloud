# 提案: 为 PostgreSQL 资源添加 StorageType 存储类型支持

## 概述

**提案编号**: POSTGRES-STORAGE-001  
**创建时间**: 2026-03-26  
**状态**: 待评审  
**优先级**: 中等

## 目标

在腾讯云 Terraform Provider 的 PostgreSQL 相关资源和数据源中添加 `StorageType` 字段支持,使用户能够:

1. 在创建实例时指定存储类型
2. 在读取实例属性时获取存储类型信息
3. 在查询数据源时根据存储类型过滤结果

## 背景

腾讯云 PostgreSQL 数据库实例支持多种存储类型,包括:

- `PHYSICAL_LOCAL_SSD`: 物理机本地 SSD 硬盘 (默认值)
- `CLOUD_PREMIUM`: 高性能云硬盘
- `CLOUD_SSD`: SSD 云硬盘
- `CLOUD_HSSD`: 增强型 SSD 云硬盘

目前 Terraform Provider 未暴露此字段,用户无法通过 Terraform 配置实例的存储类型。

## 影响范围

### 1. Resource: `tencentcloud_postgresql_instance`

**文件**: `tencentcloud/services/postgresql/resource_tc_postgresql_instance.go`

#### 修改点 1.1: Schema 添加
在 Schema 中添加 `storage_type` 字段 (约第 133 行之后):

```go
"storage_type": {
    Type:        schema.TypeString,
    Optional:    true,
    ForceNew:    true,
    Default:     "PHYSICAL_LOCAL_SSD",
    ValidateFunc: tccommon.ValidateAllowedStringValue([]string{
        "PHYSICAL_LOCAL_SSD",
        "CLOUD_PREMIUM",
        "CLOUD_SSD",
        "CLOUD_HSSD",
    }),
    Description: "Instance storage type. Valid values: `PHYSICAL_LOCAL_SSD` (physical machine local SSD disk, default), `CLOUD_PREMIUM` (premium cloud disk), `CLOUD_SSD` (SSD cloud disk), `CLOUD_HSSD` (enhanced SSD cloud disk).",
},
```

**API 文档**: https://cloud.tencent.com/document/api/409/56107

#### 修改点 1.2: Create 阶段 - API 请求
在 `resourceTencentCloudPostgresqlInstanceCreate` 函数中 (约第 560-585 行):

- 从 schema 获取 `storage_type` 值
- 将其传递给 `CreatePostgresqlInstance` 方法

**服务层修改**: `service_tencentcloud_postgresql.go` 中的 `CreatePostgresqlInstance` 方法需要添加 `storageType` 参数并设置到 SDK 请求中。

#### 修改点 1.3: Read 阶段 - 属性读取
在 `resourceTencentCloudPostgresqlInstanceRead` 函数中 (约第 928-933 行):

调用 `DescribeDBInstanceAttribute` 接口后,从返回值中读取 `DBInstanceStorageType` 字段并设置到 state:

```go
if ins.DBInstanceStorageType != nil {
    _ = d.Set("storage_type", ins.DBInstanceStorageType)
}
```

**API 返回字段**: `DBInstanceStorageType` (与入参 `StorageType` 对应)

---

### 2. DataSource: `tencentcloud_postgresql_db_instance_versions`

**文件**: `tencentcloud/services/postgresql/data_source_tc_postgresql_db_instance_versions.go`

#### 修改点 2.1: Schema 添加
在 Schema 中添加 `storage_type` 字段 (约第 75 行之后):

```go
"storage_type": {
    Type:        schema.TypeString,
    Optional:    true,
    Description: "Instance storage type. Valid values: `PHYSICAL_LOCAL_SSD` (default), `CLOUD_PREMIUM`, `CLOUD_SSD`, `CLOUD_HSSD`.",
},
```

**API 文档**: https://cloud.tencent.com/document/api/409/89018

#### 修改点 2.2: API 调用传参
在 `dataSourceTencentCloudPostgresqlDbInstanceVersionsRead` 函数中 (约第 92 行):

- 获取 `storage_type` 参数
- 传递给 `DescribePostgresqlDbInstanceVersionsByFilter` 方法

**服务层修改**: `service_tencentcloud_postgresql.go` 中的相应方法需要支持 `StorageType` 参数。

---

### 3. DataSource: `tencentcloud_postgresql_db_instance_classes`

**文件**: `tencentcloud/services/postgresql/data_source_tc_postgresql_db_instance_classes.go`

#### 修改点 3.1: Schema 添加
在 Schema 中添加 `storage_type` 字段 (约第 36 行之后):

```go
"storage_type": {
    Type:        schema.TypeString,
    Optional:    true,
    Default:     "PHYSICAL_LOCAL_SSD",
    Description: "Instance storage type. Valid values: `PHYSICAL_LOCAL_SSD` (default), `CLOUD_PREMIUM`, `CLOUD_SSD`, `CLOUD_HSSD`.",
},
```

**API 文档**: https://cloud.tencent.com/document/api/409/89019

#### 修改点 3.2: API 调用传参
在 `dataSourceTencentCloudPostgresqlDbInstanceClassesRead` 函数中 (约第 94-106 行):

```go
if v, ok := d.GetOk("storage_type"); ok {
    paramMap["StorageType"] = helper.String(v.(string))
}
```

---

### 4. DataSource: `tencentcloud_postgresql_specinfos`

**文件**: `tencentcloud/services/postgresql/data_source_tc_postgresql_specinfos.go`

#### 修改点 4.1: Schema 添加
在 Schema 中添加 `storage_type` 字段 (约第 20 行之后):

```go
"storage_type": {
    Type:        schema.TypeString,
    Optional:    true,
    Default:     "PHYSICAL_LOCAL_SSD",
    Description: "Instance storage type. Valid values: `PHYSICAL_LOCAL_SSD` (default), `CLOUD_PREMIUM`, `CLOUD_SSD`, `CLOUD_HSSD`.",
},
```

**API 文档**: https://cloud.tencent.com/document/api/409/16776

#### 修改点 4.2: API 调用传参
在 `dataSourceTencentCloudPostgresqlSpecinfosRead` 函数中 (约第 90 行):

- 获取 `storage_type` 参数
- 传递给 `DescribeSpecinfos` 方法

---

### 5. DataSource: `tencentcloud_postgresql_instances`

**文件**: `tencentcloud/services/postgresql/data_source_tc_postgresql_instances.go`

**注意**: 根据 API 文档分析,`DescribeDBInstances` 接口**不支持** `StorageType` 作为查询参数。

#### 修改点 5.1: 输出结果中添加存储类型
在 `instance_list` 和 `db_instance_set` 的 schema 中添加 `storage_type` 作为 Computed 字段:

```go
// 在 instance_list 和 db_instance_set 的 Elem.Schema 中添加
"storage_type": {
    Type:        schema.TypeString,
    Computed:    true,
    Description: "Instance storage type.",
},
```

#### 修改点 5.2: 读取响应数据
在数据处理逻辑中 (约第 542-605 行和 609-896 行),从 API 响应中读取存储类型字段并填充到结果中。

**需要确认**: 检查 `DescribeDBInstances` API 返回值中是否包含存储类型字段 (可能是 `DBInstanceStorageType` 或其他名称)。

---

## 服务层修改

**文件**: `tencentcloud/services/postgresql/service_tencentcloud_postgresql.go`

### 5.1 CreatePostgresqlInstance 方法
需要添加 `storageType` 参数,并设置到 SDK 请求中:

```go
func (me *PostgresqlService) CreatePostgresqlInstance(
    ctx context.Context,
    // ... 现有参数 ...
    storageType string, // 新增参数
) (instanceId string, err error) {
    // ...
    request.StorageType = &storageType
    // ...
}
```

### 5.2 数据源相关方法
更新以下方法以支持 `StorageType` 参数:

- `DescribePostgresqlDbInstanceVersionsByFilter`
- `DescribePostgresqlDbInstanceClassesByFilter`
- `DescribeSpecinfos`

---

## 测试计划

### 1. 单元测试
- 验证 schema 定义正确性
- 验证参数传递和验证逻辑

### 2. 集成测试
需要创建以下测试用例:

```hcl
# 测试用例 1: 创建使用 CLOUD_SSD 存储类型的实例
resource "tencentcloud_postgresql_instance" "test_cloud_ssd" {
  name              = "tf-test-cloud-ssd"
  storage_type      = "CLOUD_SSD"
  availability_zone = "ap-guangzhou-3"
  # ... 其他必需参数
}

# 测试用例 2: 数据源查询特定存储类型的版本
data "tencentcloud_postgresql_db_instance_versions" "cloud_hssd" {
  storage_type = "CLOUD_HSSD"
}

# 测试用例 3: 查询特定存储类型的规格
data "tencentcloud_postgresql_db_instance_classes" "premium" {
  zone              = "ap-guangzhou-7"
  db_engine         = "postgresql"
  db_major_version  = "15"
  storage_type      = "CLOUD_PREMIUM"
}
```

### 3. 兼容性测试
- 验证不指定 `storage_type` 时使用默认值 `PHYSICAL_LOCAL_SSD`
- 验证现有配置不受影响 (向后兼容)

---

## 风险评估

### 低风险
1. **破坏性变更**: 无。所有新增字段均为 Optional 或 Computed,不会影响现有配置。
2. **API 兼容性**: 所有涉及的腾讯云 API 均已支持 `StorageType` 参数。

### 注意事项
1. **ForceNew**: `resource_tc_postgresql_instance` 中的 `storage_type` 设置为 `ForceNew: true`,意味着修改存储类型会重建实例。
2. **数据源限制**: `tencentcloud_postgresql_instances` 可能无法根据存储类型过滤实例 (取决于 API 支持情况)。

---

## 实施步骤

1. ✅ 创建 OpenSpec 提案文档
2. ⏳ 更新 `service_tencentcloud_postgresql.go` 服务层方法
3. ⏳ 修改 `resource_tc_postgresql_instance.go`
4. ⏳ 修改 4 个数据源文件
5. ⏳ 为每个修改文件执行 `go fmt`
6. ⏳ 编写测试用例
7. ⏳ 执行测试验证
8. ⏳ 更新文档

---

## 参考文档

- [CreateInstances API](https://cloud.tencent.com/document/api/409/56107)
- [DescribeDBVersions API](https://cloud.tencent.com/document/api/409/89018)
- [DescribeClasses API](https://cloud.tencent.com/document/api/409/89019)
- [DescribeProductConfig API](https://cloud.tencent.com/document/api/409/16776)
- [DescribeDBInstances API](https://cloud.tencent.com/document/api/409/16773)

---

## 审批流程

- [ ] 技术审核
- [ ] 测试验证
- [ ] 文档更新
- [ ] 合并到主分支
