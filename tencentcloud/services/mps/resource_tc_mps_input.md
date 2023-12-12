Provides a resource to create a mps input

Example Usage

Create mps input group with SRT

```hcl
resource "tencentcloud_mps_input" "input" {
  flow_id = tencentcloud_mps_flow.flow.id
  input_group {
    input_name    = "your_input_name"
    protocol      = "SRT"
    description   = "input name Description"
    allow_ip_list = ["0.0.0.0/0"]
    srt_settings {
      mode              = "LISTENER"
      stream_id         = "#!::u=johnny,r=resource,h=xxx.com,t=stream,m=play"
      latency           = 1000
      recv_latency      = 1000
      peer_latency      = 1000
      peer_idle_timeout = 1000
    }
  }
}
```

Import

mps input can be imported using the id, e.g.

```
terraform import tencentcloud_mps_input.input input_id
```