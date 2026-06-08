Provides a resource to create a TKE roll-out sequence

Example Usage

```hcl
resource "tencentcloud_kubernetes_roll_out_sequence" "example" {
  name    = "tf-example"
  enabled = true

  sequence_flows {
    tags {
      key   = "Env"
      value = ["Test"]
    }

    soak_time = 300
  }

  sequence_flows {
    tags {
      key   = "Env"
      value = ["Pre-Production"]
    }

    tags {
      key   = "Protection-Level"
      value = ["Medium"]
    }

    soak_time = 600
  }

  sequence_flows {
    tags {
      key   = "Env"
      value = ["Production"]
    }

    tags {
      key   = "Protection-Level"
      value = ["High"]
    }

    soak_time = 600
  }
}
```

Import

TKE roll-out sequence can be imported using the sequenceId, e.g.

```
terraform import tencentcloud_kubernetes_roll_out_sequence.example 29
```
