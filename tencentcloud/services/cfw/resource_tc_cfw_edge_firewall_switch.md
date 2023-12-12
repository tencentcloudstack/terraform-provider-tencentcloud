Provides a resource to create a cfw edge_firewall_switch

Example Usage

If not set subnet_id

```hcl
data "tencentcloud_cfw_edge_fw_switches" "example" {}

resource "tencentcloud_cfw_edge_firewall_switch" "example" {
  public_ip   = data.tencentcloud_cfw_edge_fw_switches.example.data[0].public_ip
  switch_mode = 1
  enable      = 0
}
```

If set subnet id

```hcl
data "tencentcloud_cfw_edge_fw_switches" "example" {}

resource "tencentcloud_cfw_edge_firewall_switch" "example" {
  public_ip   = data.tencentcloud_cfw_edge_fw_switches.example.data[0].public_ip
  subnet_id   = "subnet-id"
  switch_mode = 1
  enable      = 1
}
```