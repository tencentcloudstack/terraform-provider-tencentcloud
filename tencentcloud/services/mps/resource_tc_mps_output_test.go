package mps_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsOutputResource_basic(t *testing.T) {
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
				Config: fmt.Sprintf(testAccMpsOutput, randomNum),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_output.output", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_output.output", "flow_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_output.output", "output.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_output.output", "output.0.output_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_output.output", "output.0.description"),
					resource.TestCheckResourceAttr("tencentcloud_mps_output.output", "output.0.protocol", "RTP"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_output.output", "output.0.rtp_settings.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_output.output", "output.0.rtp_settings.0.fec", "none"),
					resource.TestCheckResourceAttr("tencentcloud_mps_output.output", "output.0.rtp_settings.0.idle_timeout", "1000"),
				),
			},
			{
				Config: fmt.Sprintf(testAccMpsOutput_update, randomNum),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_output.output", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_output.output", "flow_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_output.output", "output.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_output.output", "output.0.output_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_output.output", "output.0.description"),
					resource.TestCheckResourceAttr("tencentcloud_mps_output.output", "output.0.protocol", "RTP"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_output.output", "output.0.rtp_settings.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_output.output", "output.0.rtp_settings.0.destinations.0.ip", "203.205.141.88"),
					resource.TestCheckResourceAttr("tencentcloud_mps_output.output", "output.0.rtp_settings.0.destinations.0.port", "65533"),
				),
			},
			{
				ResourceName:      "tencentcloud_mps_output.output",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsOutput = `

resource "tencentcloud_mps_output" "output" {
  flow_id = "018b22baaf3909831f1704344703" //keep_mps_flow
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

const testAccMpsOutput_update = `

resource "tencentcloud_mps_output" "output" {
  flow_id = "018b22baaf3909831f1704344703" //keep_mps_flow
  output {
    output_name   = "tf_mps_output_group_%d_changed"
    description   = "tf mps output group changed"
    protocol      = "RTP"
    output_region = "ap-guangzhou"
    rtp_settings {
      destinations {
        ip   = "203.205.141.88"
        port = 65533
      }
      fec          = "none"
      idle_timeout = 1000
    }
  }
}


`
