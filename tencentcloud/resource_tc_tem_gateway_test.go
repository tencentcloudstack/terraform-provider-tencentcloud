package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTemGatewayResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTemGateway,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tem_gateway.gateway", "id")),
			},
			{
				ResourceName:      "tencentcloud_tem_gateway.gateway",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTemGateway = `

resource "tencentcloud_tem_gateway" "gateway" {
  ingress {
		ingress_name = "en-xxx"
		environment_id = "en-xxx"
		cluster_namespace = "default"
		address_ip_version = "IPV4"
		rewrite_type = "AUTO"
		mixed = false
		tls {
			hosts = 
			secret_name = &lt;nil&gt;
			certificate_id = &lt;nil&gt;
		}
		rules {
			host = &lt;nil&gt;
			protocol = "http"
			http {
				paths {
					path = &lt;nil&gt;
					backend {
						service_name = &lt;nil&gt;
						service_port = &lt;nil&gt;
					}
				}
			}
		}
		clb_id = "xxx"

  }
}

`
