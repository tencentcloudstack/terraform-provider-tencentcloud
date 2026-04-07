# 变更提案：新增 tencentcloud_bh_acl 资源

## 变更类型

**新功能** — 新增 BH（堡垒机）访问权限（ACL）资源，支持完整的 CRUD 操作。

## Why

### 背景

腾讯云运维安全中心（堡垒机，BH）提供访问权限（ACL）管理能力，允许管理员精细控制哪些用户/用户组可以访问哪些资产/资产组，以及可以进行哪些操作（文件传输、剪贴板、磁盘映射等）。

目前 `tencentcloud` Provider 已有 `tencentcloud_dasb_acl`（旧版 DASB 产品）资源，但缺少对应的新版 BH 产品的 `tencentcloud_bh_acl` 资源，用户无法通过 Terraform 管理 BH 访问权限。

### 接口信息

| 操作 | 接口名 | 文档 |
|------|--------|------|
| 创建 | `CreateAcl` | https://cloud.tencent.com/document/api/1025/74411 |
| 查询 | `DescribeAcls` | https://cloud.tencent.com/document/api/1025/74409 |
| 修改 | `ModifyAcl` | https://cloud.tencent.com/document/api/1025/74408 |
| 删除 | `DeleteAcls` | https://cloud.tencent.com/document/api/1025/74410 |

SDK 版本：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bh/v20230418`

### 影响

缺少此资源导致：
- 用户无法通过 Terraform IaC 方式管理 BH 访问权限，需手动在控制台操作
- 无法与其他 BH 资源（用户、资产等）联动实现完整的 BH 权限生命周期管理

## What Changes

### 资源 Schema 主要字段

`CreateAcl` 必填字段：
- `name`（String，Required）：权限名称，最大 32 字符
- `allow_disk_redirect`（Bool，Required）：是否开启磁盘映射
- `allow_any_account`（Bool，Required）：是否允许任意账号登录

`CreateAcl` / `ModifyAcl` 可选字段（全部 Optional）：
- `allow_clip_file_up/down`：剪贴板文件上/下行
- `allow_clip_text_up/down`：剪贴板文本上/下行
- `allow_file_up/down`：SFTP 文件上/下传
- `allow_disk_file_up/down`：RDP 磁盘映射文件上/下传
- `allow_shell_file_up/down`：rz sz 文件上/下传
- `allow_file_del`：SFTP 文件删除
- `allow_access_credential`：是否允许使用访问串
- `allow_keyboard_logger`：是否允许键盘记录
- `max_file_up_size` / `max_file_down_size`：文件传输大小限制（预留）
- `max_access_credential_duration`：访问串有效期（秒，86400 整数倍）
- `user_id_set`（TypeSet of Int）：关联用户 ID
- `user_group_id_set`（TypeSet of Int）：关联用户组 ID
- `device_id_set`（TypeSet of Int）：关联资产 ID
- `device_group_id_set`（TypeSet of Int）：关联资产组 ID
- `app_asset_id_set`（TypeSet of Int）：关联应用资产 ID
- `account_set`（TypeSet of String）：关联账号
- `cmd_template_id_set`（TypeSet of Int）：关联高危命令模板 ID
- `ac_template_id_set`（TypeSet of String）：关联高危 DB 模板 ID
- `validate_from` / `validate_to`（String）：权限生效/失效时间，ISO8601
- `department_id`（String）：所属部门 ID

Computed 字段（只读，从 API 读取）：
- `acl_id`（Int）：权限 ID（即资源 ID，由 `CreateAcl` 返回的 `Id`）

### 新增文件

| 文件 | 说明 |
|------|------|
| `tencentcloud/services/bh/resource_tc_bh_acl.go` | 资源主文件（Schema + CRUD） |
| `tencentcloud/services/bh/resource_tc_bh_acl.md` | 资源文档示例 |
| `tencentcloud/services/bh/resource_tc_bh_acl_test.go` | 验收测试 |

### 修改文件

| 文件 | 修改内容 |
|------|---------|
| `tencentcloud/services/bh/service_tencentcloud_bh.go` | 新增 `DescribeBhAclById` 服务方法 |
| `tencentcloud/provider.go` | 注册新资源 `tencentcloud_bh_acl` |

### 向后兼容性

✅ 纯新增，不影响任何现有资源。

### 代码风格参考

严格参照 `tencentcloud_igtm_strategy` 资源的代码风格：
- `resource.Retry` + `tccommon.WriteRetryTimeout/ReadRetryTimeout`
- `defer tccommon.LogElapsed(...)()` + `defer tccommon.InconsistentCheck(d, meta)()`
- `tccommon.NewResourceLifeCycleHandleFuncContext`
- service 层封装查询逻辑（`DescribeBhAclById`）
- `helper.Uint64ToStr` / `helper.StrToUint64Point` 处理 ID 转换
