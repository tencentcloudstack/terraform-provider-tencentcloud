Provides a resource to manage CFW (Cloud Firewall) cluster firewall bypass configuration.

Example Usage

VPC_FW type

```hcl
resource "tencentcloud_cfw_cluster_fw_bypass_config" "vpc_fw_example" {
  fw_type = "VPC_FW"
  ccn_id  = "ccn-xxxxxxxx"
  enable  = false
}
```

NAT_FW type

```hcl
resource "tencentcloud_cfw_cluster_fw_bypass_config" "nat_fw_example" {
  fw_type    = "NAT_FW"
  ccn_id     = "ccn-xxxxxxxx"
  nat_ins_id = "nat-xxxxxxxx"
  enable     = false
}
```

Import

CFW cluster firewall bypass config can be imported using the composite ID.

For VPC_FW type, the format is `{fw_type}#{ccn_id}`:

```
terraform import tencentcloud_cfw_cluster_fw_bypass_config.vpc_fw_example VPC_FW#ccn-xxxxxxxx
```

For NAT_FW type, the format is `{fw_type}#{ccn_id}#{nat_ins_id}`:

```
terraform import tencentcloud_cfw_cluster_fw_bypass_config.nat_fw_example NAT_FW#ccn-xxxxxxxx#nat-xxxxxxxx
```
