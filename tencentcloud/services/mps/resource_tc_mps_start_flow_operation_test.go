package mps_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsStartFlowOperationResource_basic(t *testing.T) {
	t.Parallel()
	randIns := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNum := randIns.Intn(1000)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMpsStartFlowOperation_start, randomNum, randomNum, randomNum),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_start_flow_operation.operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_start_flow_operation.operation", "flow_id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_start_flow_operation.operation", "start", "true"),
				),
			},
			{
				Config: fmt.Sprintf(testAccMpsStartFlowOperation_stop, randomNum, randomNum, randomNum),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_start_flow_operation.operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_start_flow_operation.operation", "flow_id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_start_flow_operation.operation", "start", "false"),
				),
			},
		},
	})
}

const testAccMpsStartFlowOperation_basic = `
resource "tencentcloud_mps_event" "event_rtp" {
  event_name  = "tf_test_event_rtp_oper_%d"
  description = "tf test mps event description"
}

resource "tencentcloud_mps_flow" "flow_rtp" {
  flow_name     = "tf_test_flow_rtp_oper_%d"
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
    output_name   = "tf_mps_output_group_%d"
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

`

const testAccMpsStartFlowOperation_start = testAccMpsStartFlowOperation_basic + `
resource "tencentcloud_mps_start_flow_operation" "operation" {
  flow_id    = tencentcloud_mps_flow.flow_rtp.id
  start      = true
  depends_on = [tencentcloud_mps_output.output]
}

`

const testAccMpsStartFlowOperation_stop = testAccMpsStartFlowOperation_basic + `
resource "tencentcloud_mps_start_flow_operation" "operation" {
  flow_id    = tencentcloud_mps_flow.flow_rtp.id
  start      = false
  depends_on = [tencentcloud_mps_output.output]
}

`
