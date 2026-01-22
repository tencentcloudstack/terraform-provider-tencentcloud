Provides a resource to create a VPC bandwidth package

Example Usage

```hcl
resource "tencentcloud_vpc_bandwidth_package" "example" {
  network_type           = "BGP"
  charge_type            = "TOP5_POSTPAID_BY_MONTH"
  bandwidth_package_name = "tf-example"
  tags = {
    createdBy = "Terraform"
  }
}
```

PrePaid Bandwidth Package

```hcl
resource "tencentcloud_vpc_bandwidth_package" "example" {
  network_type           = "BGP"
  charge_type            = "FIXED_PREPAID_BY_MONTH"
  bandwidth_package_name = "tf-example"
  time_span              = 3
  internet_max_bandwidth = 100
  tags = {
    createdBy = "Terraform"
  }
}
````

Bandwidth Package With Egress

```hcl
resource "tencentcloud_vpc_bandwidth_package" "example" {
  network_type           = "SINGLEISP_CMCC"
  charge_type            = "ENHANCED95_POSTPAID_BY_MONTH"
  bandwidth_package_name = "tf-example"
  internet_max_bandwidth = 400
  egress                 = "center_egress2"
  tags                   = {
    createdBy = "Terraform"
  }
}
```

Import

vpc bandwidth_package can be imported using the id, e.g.
```
$ terraform import tencentcloud_vpc_bandwidth_package.example bwp-hq8h7qpy
```