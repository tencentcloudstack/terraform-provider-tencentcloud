# 任务清单：fix-bh-resource-add-ip-fields

## 1. schema 新增 public_ip_set 和 private_ip_set 字段

**文件**: `tencentcloud/services/bh/resource_tc_bh_resource.go`

- [x] 在 `resource_id` 字段之后新增：
  ```go
  "public_ip_set": {
      Type:     schema.TypeList,
      Computed: true,
      Elem:     &schema.Schema{Type: schema.TypeString},
      Description: "Public IP address list of the bastion host instance.",
  },
  "private_ip_set": {
      Type:     schema.TypeList,
      Computed: true,
      Elem:     &schema.Schema{Type: schema.TypeString},
      Description: "Private IP address list of the bastion host instance.",
  },
  ```

---

## 2. Read 模块新增赋值

**文件**: `tencentcloud/services/bh/resource_tc_bh_resource.go`

- [x] 在 `external_access` 赋值之后新增：
  ```go
  if respData.PublicIpSet != nil {
      _ = d.Set("public_ip_set", helper.StringsInterfaces(respData.PublicIpSet))
  }
  if respData.PrivateIpSet != nil {
      _ = d.Set("private_ip_set", helper.StringsInterfaces(respData.PrivateIpSet))
  }
  ```
- [x] 执行 `go fmt ./tencentcloud/services/bh/`

---

## 3. 编译验证

- [x] `go build ./tencentcloud/services/bh/` 确认编译通过

---

## 总结

- **预计工作量**：极小（约 5 分钟）
- **风险等级**：极低（纯增量 Computed 字段）
- **破坏性变更**：无
- **状态**: 已完成

---

## 4. 简化 IP 列表赋值方式（post-apply 补充）

**背景**: `[]*string` 类型可直接传给 `d.Set`，无需 `helper.StringsInterfaces` 转换，与 bh 包内其他资源（`dasb_device`、`dasb_acl`）保持一致风格。

**文件**: `tencentcloud/services/bh/resource_tc_bh_resource.go`

- [x] 将 `helper.StringsInterfaces(respData.PublicIpSet)` 改为 `respData.PublicIpSet`
- [x] 将 `helper.StringsInterfaces(respData.PrivateIpSet)` 改为 `respData.PrivateIpSet`
- [x] 执行 `go fmt ./tencentcloud/services/bh/`
- [x] `go build ./tencentcloud/services/bh/` 确认编译通过
