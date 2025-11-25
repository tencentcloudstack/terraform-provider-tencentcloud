Provides a resource to create a CFW nat firewall switch

Example Usage

Turn off switch

```hcl
data "tencentcloud_cfw_nat_fw_switches" "example" {
  nat_ins_id = "cfwnat-18d2ba18"
}

resource "tencentcloud_cfw_nat_firewall_switch" "example" {
  nat_ins_id = data.tencentcloud_cfw_nat_fw_switches.example.id
  subnet_id  = data.tencentcloud_cfw_nat_fw_switches.example.data[0].subnet_id
  enable     = 0
}
```

Or turn on switch

```hcl
data "tencentcloud_cfw_nat_fw_switches" "example" {
  nat_ins_id = "cfwnat-18d2ba18"
}

resource "tencentcloud_cfw_nat_firewall_switch" "example" {
  nat_ins_id = data.tencentcloud_cfw_nat_fw_switches.example.id
  subnet_id  = data.tencentcloud_cfw_nat_fw_switches.example.data[0].subnet_id
  enable     = 1
}
```

Import

CFW nat firewall switch can be imported using the natInsId#subnetId, e.g.

```
terraform import tencentcloud_cfw_nat_firewall_switch.example cfwnat-18d2ba18#subnet-ef7wyymr
```
