Provides a resource to create a Tencent Cloud Global Accelerator V2 (GA2) instance.

Example Usage

```hcl
resource "tencentcloud_ga2_global_accelerator" "example" {
  name                 = "tf-example"
  instance_charge_type = "POSTPAID"
  description          = "tf example global accelerator"

  tags = {
    createdBy = "Terraform"
  }
}
```

Cross-border global accelerator

```hcl
resource "tencentcloud_ga2_global_accelerator" "example" {
  name                      = "tf-example"
  instance_charge_type      = "POSTPAID"
  description               = "tf example cross-border accelerator"
  cross_border_type         = "HighQuality"
  cross_border_promise_flag = true

  tags = {
    createdBy = "Terraform"
  }
}
```

Import

GA2 global accelerator instance can be imported using the id, e.g.

```
terraform import tencentcloud_ga2_global_accelerator.example ga-ar31grog
```
