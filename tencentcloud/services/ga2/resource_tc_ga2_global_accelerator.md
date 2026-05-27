Provides a resource to create a GA2 (Global Accelerator 2) global accelerator instance.

Example Usage

```hcl
resource "tencentcloud_ga2_global_accelerator" "example" {
  name                    = "tf-example"
  instance_charge_type    = "POSTPAID"
  description             = "terraform example global accelerator"
  cross_border_type       = "HighQuality"
  cross_border_promise_flag = true

  tags = {
    createdBy = "terraform"
  }
}
```

Import

GA2 global accelerator can be imported using the ID, e.g.

```
terraform import tencentcloud_ga2_global_accelerator.example ga2-xxxxxxxx
```
