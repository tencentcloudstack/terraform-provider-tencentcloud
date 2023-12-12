Provides a resource to create a mps output

Example Usage

Create a output group with RTP

```hcl
resource "tencentcloud_mps_output" "output" {
  flow_id = "your_flow_id"
  output {
    output_name   = "your_output_name"
    description   = "tf mps output group"
    protocol      = "RTP"
    output_region = "ap-guangzhou"
    rtp_settings {
      destinations {
        ip   = "203.205.141.84"
        port = 65535
      }
      fec          = "none"
      idle_timeout = 1000
    }
  }
}
```

Import

mps output can be imported using the id, e.g.

```
terraform import tencentcloud_mps_output.output flow_id#output_id
```