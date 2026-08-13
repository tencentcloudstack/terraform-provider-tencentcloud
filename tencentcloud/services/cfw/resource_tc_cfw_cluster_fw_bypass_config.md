Provides a resource to manage CFW (Cloud Firewall) cluster firewall bypass configuration.

Example Usage

VPC_FW type

```hcl
resource "tencentcloud_cfw_cluster_fw_bypass_config" "example" {
  fw_type = "VPC_FW"
  ccn_id  = "ccn-p3mlp0tj"
  enable  = false
}
```

NAT_FW type

```hcl
resource "tencentcloud_cfw_cluster_fw_bypass_config" "example" {
  fw_type    = "NAT_FW"
  ccn_id     = "ccn-p3mlp0tj"
  nat_ins_id = "nat-h1i1mf4n"
  enable     = true
}
```

Import

CFW cluster firewall bypass config can be imported using the composite ID.

For VPC_FW type, the format is `{fw_type}#{ccn_id}`:

```
terraform import tencentcloud_cfw_cluster_fw_bypass_config.example VPC_FW#ccn-p3mlp0tj
```

For NAT_FW type, the format is `{fw_type}#{nat_ins_id},{ccn_id}`:

```
terraform import tencentcloud_cfw_cluster_fw_bypass_config.example NAT_FW#nat-h1i1mf4n,ccn-p3mlp0tj
```
