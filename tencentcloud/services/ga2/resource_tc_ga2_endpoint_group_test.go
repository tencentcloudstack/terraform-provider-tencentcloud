package ga2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudGa2EndpointGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGa2EndpointGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ga2_endpoint_group.example", "id"),
				),
			},
			{
				Config: testAccGa2EndpointGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ga2_endpoint_group.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_ga2_endpoint_group.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccGa2EndpointGroup = `
resource "tencentcloud_ga2_endpoint_group" "example" {
  global_accelerator_id = "ga-4mredmiu"
  listener_id           = "lsr-1vd1fdwf"
  endpoint_group_type   = "VIRTUAL"

  endpoint_group_configuration {
    name                  = "tf-example"
    endpoint_group_region = "ap-guangzhou"
    description           = "tf example endpoint group"
    enable_health_check   = false
    forward_protocol      = "HTTP"

    endpoint_configurations {
      endpoint_type    = "CustomPublicIp"
      endpoint_service = "1.1.1.1"
      weight           = 10
    }

    endpoint_configurations {
      endpoint_type    = "CustomDomain"
      endpoint_service = "example.com"
      weight           = 20
    }

    port_overrides {
      listener_port = 8080
      endpoint_port = 9090
    }
  }
}
`

const testAccGa2EndpointGroupUpdate = `
resource "tencentcloud_ga2_endpoint_group" "example" {
  global_accelerator_id = "ga-4mredmiu"
  listener_id           = "lsr-1vd1fdwf"
  endpoint_group_type   = "VIRTUAL"

  endpoint_group_configuration {
    name                  = "tf-example"
    endpoint_group_region = "ap-guangzhou"
    description           = "tf example endpoint group"
    enable_health_check   = false
    forward_protocol      = "HTTP"

    endpoint_configurations {
      endpoint_type    = "CustomPublicIp"
      endpoint_service = "1.1.1.1"
      weight           = 10
    }

    endpoint_configurations {
      endpoint_type    = "CustomDomain"
      endpoint_service = "example.com"
      weight           = 20
    }

    port_overrides {
      listener_port = 8080
      endpoint_port = 9000
    }
  }
}
`
