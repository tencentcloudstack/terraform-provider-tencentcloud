## Why

`tencentcloud_teo_deploy_config_group_version` 资源用于在 EdgeOne (TEO) 中发布配置组版本。当前资源通过 `DescribeConfigGroupVersionHistory` 接口读取部署记录信息（如 `record_id`、`deploy_time`、`status`、`message`），但缺少发布后环境维度的状态信息。

用户需要感知版本发布后环境的实际生效状态、环境类型（生产/测试）、生效范围以及当前实际生效的配置组版本信息（含 `VersionId`、`VersionNumber`、`SourceVersion`、`GroupType`、`GroupId`、`Description`、`Status`、`CreateTime`），以便通过 Terraform state 跟踪部署结果与环境现状。这些信息由云 API `DescribeEnvironments` 的 `EnvInfo` 结构返回，当前资源已调用该接口（在 `Create` 中用于轮询环境状态）但未将其返回值写入 state。

## What Changes

- 为 `tencentcloud_teo_deploy_config_group_version` 资源新增出参（Computed），来源为 `DescribeEnvironments` 接口的 `EnvInfo` 结构：
  - `total_count`（uint64）：环境总数。
  - `env_type`（string）：环境类型，取值 `production` / `staging`。
  - `scope`（string 列表）：当前环境配置生效范围（生产环境为 `["ALL"]`，测试环境返回测试节点 IP）。
  - `current_config_group_version_infos`（集合）：当前环境各配置组实际生效的版本信息列表，元素包含：
    - `version_id`（string）
    - `version_number`（string）
    - `source_version`（string）
    - `group_type`（string）
    - `group_id`（string）
    - `description`（string）
    - `status`（string）
    - `create_time`（string）
  - `env_create_time`（string）：环境创建时间。
  - `env_update_time`（string）：环境更新时间。

> 注意：现有 `config_group_version_infos`（来自部署记录 `DeployRecord`）字段保持不变，不进行修改。新增的 `current_config_group_version_infos` 是来自环境维度 `DescribeEnvironments` 的当前生效版本信息，二者来源不同、语义不同，故新增独立字段，不复用现有字段。

## Capabilities

### New Capabilities
- `teo-deploy-config-group-version-env-info`: 为 `tencentcloud_teo_deploy_config_group_version` 资源补充来自 `DescribeEnvironments` 的环境维度的只读（Computed）出参，使用户可在 state 中跟踪发布后环境的实际生效状态与当前生效配置组版本信息。

### Modified Capabilities
<!-- 无既有 spec 的需求变更 -->

## Impact

- **代码**:
  - `tencentcloud/services/teo/resource_tc_teo_deploy_config_group_version.go`：在 `Schema` 中新增上述 Computed 字段；在 `Read` 方法中调用 `DescribeTeoEnvironmentsByFilter`（已存在的 service 层方法）并填充新字段。
  - `tencentcloud/services/teo/service_tencentcloud_teo.go`：现有 `DescribeTeoEnvironmentsByFilter` 已返回 `[]*EnvInfo`，无需新增 service 方法，仅在 Read 中补充取值逻辑。
- **测试**: `tencentcloud/services/eo/resource_tc_teo_deploy_config_group_version_test.go` 补充新增字段的断言。
- **文档**: `tencentcloud/services/teo/resource_tc_teo_deploy_config_group_version.md` 同步更新示例（由 `make doc` 生成 `website/docs/`）。
- **API**: 复用 `DescribeEnvironments`（包名 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`），该接口已存在于 vendor，无新增依赖。
- **兼容性**: 仅新增 Computed 字段，不修改既有 schema 字段，完全向后兼容。
