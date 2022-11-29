package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTcmPrometheusAttachment_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcmPrometheusAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_prometheus_attachment.prometheus_attachment", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_tcm_prometheus_attachment.prometheusAttachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTcmPrometheusAttachment = `

resource "tencentcloud_tcm_prometheus_attachment" "prometheus_attachment" {
	mesh_id = "mesh-rofjmux7"
	prometheus {
	  vpc_id = "vpc-xxx"
	  subnet_id = "subnet-xxx"
	  region = "ap-guangzhou"
	  instance_id = ""
	  custom_prom {
		  is_public_addr = false
		  vpc_id = "vpc-xxx"
		  url = "http://x.x.x.x:9090"
		  auth_type = "none, basic"
		  username = "test"
		  password = "test"
	  }
	}
  }

`
