package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTcmPrometheusAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcmPrometheusAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tcm_prometheus_attachment.prometheus_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_tcm_prometheus_attachment.prometheus_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTcmPrometheusAttachment = `

resource "tencentcloud_tcm_prometheus_attachment" "prometheus_attachment" {
  mesh_i_d = "mesh-xxxxxxxx"
  prometheus {
		vpc_id = "vpc-xxx"
		subnet_id = "subnet-xxx"
		region = "sh"
		instance_id = "prom-xxx"
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
