Provides a resource to create a mps flow

Example Usage

Create a mps RTP flow

```hcl
resource "tencentcloud_mps_event" "event" {
	event_name = "tf_test_event_srt_%d"
	description = "tf test mps event description"
  }

resource "tencentcloud_mps_flow" "flow" {
  flow_name = "tf_test_mps_flow_srt_%d"
  max_bandwidth = 10000000
  input_group {
		input_name = "test_inputname"
		protocol = "SRT"
		description = "input name Description"
		allow_ip_list = ["0.0.0.0/0"]
		srt_settings {
			mode = "LISTENER"
			stream_id = "#!::u=johnny,r=resource,h=xxx.com,t=stream,m=play"
			latency = 1000
			recv_latency = 1000
			peer_latency =  1000
			peer_idle_timeout =  1000
		}
  }
  event_id = tencentcloud_mps_event.event.id
}
```

Create a mps RTP flow

```hcl
resource "tencentcloud_mps_event" "event_rtp" {
	event_name = "tf_test_event_rtp_%d"
	description = "tf test mps event description"
  }

resource "tencentcloud_mps_flow" "flow_rtp" {
  flow_name = "tf_test_mps_flow_rtp_%d"
  max_bandwidth = 10000000
  input_group {
		input_name = "test_inputname"
		protocol = "RTP"
		description = "input name Description"
		allow_ip_list = ["0.0.0.0/0"]
		rtp_settings {
			fec = "none"
			idle_timeout = 1000
		}
  }
  event_id = tencentcloud_mps_event.event_rtp.id
}
```

Create a mps RTP flow and start it

Before you start a mps flow, you need to create a output first.

```hcl
resource "tencentcloud_mps_event" "event_rtp" {
  event_name  = "your_event_name"
  description = "tf test mps event description"
}

resource "tencentcloud_mps_flow" "flow_rtp" {
  flow_name     = "your_flow_name"
  max_bandwidth = 10000000
  input_group {
    input_name    = "test_inputname"
    protocol      = "RTP"
    description   = "input name Description"
    allow_ip_list = ["0.0.0.0/0"]
    rtp_settings {
      fec          = "none"
      idle_timeout = 1000
    }
  }
  event_id = tencentcloud_mps_event.event_rtp.id
}

resource "tencentcloud_mps_output" "output" {
  flow_id = tencentcloud_mps_flow.flow_rtp.id
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

resource "tencentcloud_mps_start_flow_operation" "operation" {
  flow_id    = tencentcloud_mps_flow.flow_rtp.id
  start      = true
  depends_on = [tencentcloud_mps_output.output]
}
```

Import

mps flow can be imported using the id, e.g.

```
terraform import tencentcloud_mps_flow.flow flow_id
```