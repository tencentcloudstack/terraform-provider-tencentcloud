package tencentcloud

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsInputResource_basic(t *testing.T) {
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
				Config: fmt.Sprintf(testAccMpsInput, randomNum, randomNum),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_input.input", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_input.input", "flow_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_input.input", "input_group.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_input.input", "input_group.0.input_name", fmt.Sprintf("tf_mps_input_group_%d", randomNum)),
					resource.TestCheckResourceAttr("tencentcloud_mps_input.input", "input_group.0.protocol", "SRT"),
					resource.TestCheckResourceAttr("tencentcloud_mps_input.input", "input_group.0.description", "input name Description"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_input.input", "input_group.0.allow_ip_list.*", "0.0.0.0/0"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_input.input", "input_group.0.srt_settings.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_input.input", "input_group.0.srt_settings.0.mode", "LISTENER"),
					resource.TestCheckResourceAttr("tencentcloud_mps_input.input", "input_group.0.srt_settings.0.latency", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_mps_input.input", "input_group.0.srt_settings.0.stream_id", "#!::u=johnny,r=resource,h=xxx.com,t=stream,m=play"),
				),
			},
			{
				Config: fmt.Sprintf(testAccMpsInput_update, randomNum, randomNum),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_input.input", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_input.input", "flow_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_input.input", "input_group.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_input.input", "input_group.0.input_name", fmt.Sprintf("tf_mps_input_group_%d_changed", randomNum)),
					resource.TestCheckResourceAttr("tencentcloud_mps_input.input", "input_group.0.protocol", "SRT"),
					resource.TestCheckResourceAttr("tencentcloud_mps_input.input", "input_group.0.description", "input name Description changed"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_input.input", "input_group.0.allow_ip_list.*", "1.1.1.1/8"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_input.input", "input_group.0.srt_settings.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_input.input", "input_group.0.srt_settings.0.latency", "999"),
					resource.TestCheckResourceAttr("tencentcloud_mps_input.input", "input_group.0.srt_settings.0.recv_latency", "999"),
					resource.TestCheckResourceAttr("tencentcloud_mps_input.input", "input_group.0.srt_settings.0.stream_id", "#!::u=johnny,r=resource,h=xxx.com,t=stream,m=play"),
				),
			},
			{
				ResourceName: "tencentcloud_mps_input.input",
				ImportState:  true,
			},
		},
	})
}

const testAccMpsInput_flow = `
resource "tencentcloud_mps_flow" "flow" {
  flow_name     = "tf_test_mps_flow_input_%d"
  max_bandwidth = 10000000

  event_id = "018b195dc9321eb401a604179263"
}

`

const testAccMpsInput = testAccMpsInput_flow + `

resource "tencentcloud_mps_input" "input" {
  flow_id = tencentcloud_mps_flow.flow.id
  input_group {
    input_name    = "tf_mps_input_group_%d"
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


`

const testAccMpsInput_update = testAccMpsInput_flow + `

resource "tencentcloud_mps_input" "input" {
	flow_id = tencentcloud_mps_flow.flow.id
	input_group {
		input_name    = "tf_mps_input_group_%d_changed"
		protocol      = "SRT"
		description   = "input name Description changed"
		allow_ip_list = ["1.1.1.1/8"]
		srt_settings {
		  mode              = "LISTENER"
		  stream_id         = "#!::u=johnny,r=resource,h=xxx.com,t=stream,m=play"
		  latency           = 999
		  recv_latency      = 999
		  peer_latency      = 1000
		  peer_idle_timeout = 1000
		}
	  }
  }
  

`
