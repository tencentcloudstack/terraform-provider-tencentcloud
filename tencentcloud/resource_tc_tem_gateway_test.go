package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixTemGateway_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTemGateway,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tem_gateway.gateway", "id"),
				),
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
    ingress_name = "demo"
    environment_id = "en-853mggjm"
    cluster_namespace = "default"
    address_i_p_version = "IPV4"
    rewrite_type = "NONE"
    mixed = false
    rules {
      host = "test.com"
      protocol = "http"
      http {
        paths {
          path = "/"
          backend {
            service_name = "demo"
            service_port = 80
          }
        }
      }
    }
    rules {
      host = "hello.com"
      protocol = "http"
      http {
        paths {
          path = "/"
          backend {
            service_name = "hello"
            service_port = 36000
          }
        }
      }
    }
  }
}
`
