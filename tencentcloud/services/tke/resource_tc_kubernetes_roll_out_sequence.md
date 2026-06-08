Provides a resource to create a TKE roll-out sequence

Example Usage

```hcl
resource "tencentcloud_kubernetes_roll_out_sequence" "example" {
  name    = "tf-example"
  enabled = true

  sequence_flows {
    tags {
      key   = "env"
      value = ["production"]
    }

    soak_time = 3600
  }

  sequence_flows {
    tags {
      key   = "env"
      value = ["staging", "testing"]
    }

    soak_time = 1800
  }
}
```

Import

TKE roll-out sequence can be imported using the id, e.g.

```
terraform import tencentcloud_kubernetes_roll_out_sequence.example 123
```
