# Proposal: Add BH Auth Mode Config Resource

## Background

运维安全中心（堡垒机，BH）产品目前缺少认证方式配置资源，用户无法通过 Terraform 管理堡垒机的认证模式设置（包括双因子认证和国密双因子认证）。

认证方式配置是堡垒机安全设置的重要组成部分，支持以下认证类型：
- **双因子认证**（AuthMode）：OTP、短信、USB Key
- **国密双因子认证**（AuthModeGM）：OTP、短信、USB Key

该配置属于全局配置类型（Config Resource），不依赖特定资源 ID，而是作为整个堡垒机实例的安全策略设置。

**注意**：API 文档中的 ResourceType 字段已标注为"暂停使用"，因此不纳入 Terraform Schema。

## Goals

为 Terraform Provider 添加 `tencentcloud_bh_auth_mode_config` 资源，实现：

1. **Create/Update 功能**：通过 `ModifyAuthModeSetting` 接口配置认证方式
2. **Read 功能**：通过 `DescribeSecuritySetting` 接口查询当前认证配置
3. **Delete 功能**：重置为默认配置（AuthMode=0，关闭认证）
4. **Import 功能**：支持导入现有配置

## Non-Goals

- 不涉及堡垒机的其他安全配置（Login、Password、LDAP 等）
- 不处理多实例场景（该配置为全局级别）
- 不实现历史配置版本管理

## Scope

### 新增文件

1. **资源文件**：`tencentcloud/services/bh/resource_tc_bh_auth_mode_config.go`
   - 实现完整的 CRUD 操作
   - 参考 `resource_tc_igtm_strategy.go` 代码风格
   - 使用 TF 生成的 UUID 作为资源 ID

2. **Provider 注册**：修改 `tencentcloud/provider.go`
   - 在 BH 资源组中注册 `tencentcloud_bh_auth_mode_config`

### API 映射

| Terraform 操作 | API 接口 | 说明 |
|---------------|---------|------|
| Create | ModifyAuthModeSetting | 初始化认证配置 |
| Read | DescribeSecuritySetting | 查询 SecuritySetting.AuthMode |
| Update | ModifyAuthModeSetting | 更新认证配置 |
| Delete | ModifyAuthModeSetting | 重置为默认值（AuthMode=0） |
| Import | DescribeSecuritySetting | 导入现有配置 |

### API 文档

- **查询接口**：[DescribeSecuritySetting](https://cloud.tencent.com/document/api/1025/125050)
- **修改接口**：[ModifyAuthModeSetting](https://cloud.tencent.com/document/api/1025/125048)

## Success Criteria

- [x] 资源文件创建完成，包含完整的 CRUD 函数
- [x] Schema 定义正确，包含所有必要字段
- [x] Provider 注册成功
- [x] 代码风格与 `resource_tc_igtm_strategy.go` 一致
- [x] 使用 UUID 作为资源 ID（全局配置特性）
- [x] 编译通过，无语法错误
- [x] 格式化完成（go fmt）

## Timeline

- **Phase 1**: 创建资源文件和 Schema 定义（1 Task）
- **Phase 2**: 注册到 Provider（1 Task）
- **Phase 3**: 代码验证和格式化（1 Task）

预计完成时间：即时完成
