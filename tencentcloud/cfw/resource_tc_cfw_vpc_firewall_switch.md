Provides a resource to create a cfw vpc_firewall_switch

Example Usage

Turn off switch

```hcl
data "tencentcloud_cfw_vpc_fw_switches" "example" {
  vpc_ins_id = "cfwg-c8c2de41"
}

resource "tencentcloud_cfw_vpc_firewall_switch" "example" {
  vpc_ins_id = data.tencentcloud_cfw_vpc_fw_switches.example.id
  switch_id  = data.tencentcloud_cfw_vpc_fw_switches.example.switch_list[0].switch_id
  enable     = 0
}
```

Or turn on switch

```hcl
data "tencentcloud_cfw_vpc_fw_switches" "example" {
  vpc_ins_id = "cfwg-c8c2de41"
}

resource "tencentcloud_cfw_vpc_firewall_switch" "example" {
  vpc_ins_id = data.tencentcloud_cfw_vpc_fw_switches.example.id
  switch_id  = data.tencentcloud_cfw_vpc_fw_switches.example.switch_list[0].switch_id
  enable     = 1
}
```
Import

cfw vpc_firewall_switch can be imported using the id, e.g.

```
terraform import tencentcloud_cfw_vpc_firewall_switch.example cfwg-c8c2de41#cfws-f2c63ded84
```