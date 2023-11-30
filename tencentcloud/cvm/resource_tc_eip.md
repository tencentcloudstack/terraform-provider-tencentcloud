Provides an EIP resource.

Example Usage

Paid by the bandwidth package
```hcl
resource "tencentcloud_eip" "foo" {
  name                 = "awesome_gateway_ip"
  bandwidth_package_id = "bwp-jtvzuky6"
  internet_charge_type = "BANDWIDTH_PACKAGE"
  type                 = "EIP"
}
```

AntiDDos Eip
```
resource "tencentcloud_eip" "foo" {
  name                 = "awesome_gateway_ip"
  bandwidth_package_id = "bwp-4ocyia9s"
  internet_charge_type = "BANDWIDTH_PACKAGE"
  type                 = "AntiDDoSEIP"
  anti_ddos_package_id = "xxxxxxxx"

  tags = {
    "test" = "test"
  }
}
```

Eip With Network Egress
```
resource "tencentcloud_eip" "foo" {
  name                       = "egress_eip"
  egress                     = "center_egress2"
  internet_charge_type       = "BANDWIDTH_PACKAGE"
  internet_service_provider  = "CMCC"
  internet_max_bandwidth_out = 1
  type                       = "EIP"
}
```

Import

EIP can be imported using the id, e.g.

```
$ terraform import tencentcloud_eip.foo eip-nyvf60va
```