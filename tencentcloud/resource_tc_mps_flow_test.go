package tencentcloud

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsFlowResource_srt(t *testing.T) {
	t.Parallel()
	randIns := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNum := randIns.Intn(1000)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMpsFlow_srt, randomNum, randomNum),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_flow.flow", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_flow.flow", "flow_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_flow.flow", "event_id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_flow.flow", "max_bandwidth", "10000000"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_flow.flow", "input_group.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_flow.flow", "input_group.0.input_name", "test_inputname"),
					resource.TestCheckResourceAttr("tencentcloud_mps_flow.flow", "input_group.0.protocol", "SRT"),
					resource.TestCheckResourceAttr("tencentcloud_mps_flow.flow", "input_group.0.description", "input name Description"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_flow.flow", "input_group.0.allow_ip_list.*", "0.0.0.0/0"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_flow.flow", "input_group.0.srt_settings.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_flow.flow", "input_group.0.srt_settings.0.mode", "LISTENER"),
					resource.TestCheckResourceAttr("tencentcloud_mps_flow.flow", "input_group.0.srt_settings.0.stream_id", "#!::u=johnny,r=resource,h=xxx.com,t=stream,m=play"),
					resource.TestCheckResourceAttr("tencentcloud_mps_flow.flow", "input_group.0.srt_settings.0.latency", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_mps_flow.flow", "input_group.0.srt_settings.0.recv_latency", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_mps_flow.flow", "input_group.0.srt_settings.0.peer_idle_timeout", "1000"),
				),
			},
			{
				Config: fmt.Sprintf(testAccMpsFlow_srt_update, randomNum, randomNum),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_flow.flow", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_flow.flow", "flow_name", fmt.Sprintf("tf_test_mps_flow_srt_%d_changed", randomNum)),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_flow.flow", "event_id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_flow.flow", "max_bandwidth", "10000000"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_flow.flow", "input_group.#"),
				),
			},
			{
				ResourceName: "tencentcloud_mps_flow.flow",
				ImportState:  true,
			},
		},
	})
}

func TestAccTencentCloudMpsFlowResource_rtp(t *testing.T) {
	t.Parallel()
	randIns := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNum := randIns.Intn(1000)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMpsFlow_rtp, randomNum, randomNum),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_flow.flow_rtp", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_flow.flow_rtp", "flow_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_flow.flow_rtp", "event_id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_flow.flow_rtp", "max_bandwidth", "10000000"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_flow.flow_rtp", "input_group.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_flow.flow_rtp", "input_group.0.input_name", "test_inputname"),
					resource.TestCheckResourceAttr("tencentcloud_mps_flow.flow_rtp", "input_group.0.protocol", "RTP"),
					resource.TestCheckResourceAttr("tencentcloud_mps_flow.flow_rtp", "input_group.0.description", "input name Description"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_flow.flow_rtp", "input_group.0.allow_ip_list.*", "0.0.0.0/0"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_flow.flow_rtp", "input_group.0.rtp_settings.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_flow.flow_rtp", "input_group.0.rtp_settings.0.fec", "none"),
					resource.TestCheckResourceAttr("tencentcloud_mps_flow.flow_rtp", "input_group.0.rtp_settings.0.idle_timeout", "1000"),
				),
			},
			{
				ResourceName: "tencentcloud_mps_flow.flow_rtp",
				ImportState:  true,
			},
		},
	})
}

const testAccMpsFlow_srt = `
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

`

const testAccMpsFlow_srt_update = `
resource "tencentcloud_mps_event" "event" {
	event_name = "tf_test_event_srt_%d"
	description = "tf test mps event description"
  }

resource "tencentcloud_mps_flow" "flow" {
  flow_name = "tf_test_mps_flow_srt_%d_changed"
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


`

const testAccMpsFlow_rtp = `
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

`
