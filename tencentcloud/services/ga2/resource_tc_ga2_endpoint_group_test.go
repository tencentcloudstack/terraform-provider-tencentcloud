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
					resource.TestCheckResourceAttrSet("tencentcloud_ga2_endpoint_group.example", "endpoint_group_id"),
					resource.TestCheckResourceAttr("tencentcloud_ga2_endpoint_group.example", "endpoint_group_type", "DEFAULT"),
					resource.TestCheckResourceAttr("tencentcloud_ga2_endpoint_group.example", "endpoint_group_configuration.0.name", "tf-example"),
				),
			},
			{
				Config: testAccGa2EndpointGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_ga2_endpoint_group.example", "endpoint_group_configuration.0.name", "tf-example-update"),
					resource.TestCheckResourceAttr("tencentcloud_ga2_endpoint_group.example", "endpoint_group_configuration.0.description", "updated description"),
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
  global_accelerator_id = "ga2-xxxxxxxx"
  listener_id           = "lis-xxxxxxxx"
  endpoint_group_type   = "DEFAULT"

  endpoint_group_configuration {
    name                  = "tf-example"
    endpoint_group_region = "ap-guangzhou"
    description           = "tf example"
    enable_health_check   = true
    check_type            = "HTTP"
    check_port            = "80"
    check_path            = "/"
    check_method          = "GET"
    connect_timeout       = 5000
    health_check_interval = 30
    healthy_threshold     = 3
    unhealthy_threshold   = 3
    forward_protocol      = "HTTP"

    endpoint_configurations {
      endpoint_type    = "PublicIp"
      endpoint_service = "1.1.1.1"
      weight           = 10
    }
  }
}
`

const testAccGa2EndpointGroupUpdate = `
resource "tencentcloud_ga2_endpoint_group" "example" {
  global_accelerator_id = "ga2-xxxxxxxx"
  listener_id           = "lis-xxxxxxxx"
  endpoint_group_type   = "DEFAULT"

  endpoint_group_configuration {
    name                  = "tf-example-update"
    endpoint_group_region = "ap-guangzhou"
    description           = "updated description"
    enable_health_check   = true
    check_type            = "HTTP"
    check_port            = "80"
    check_path            = "/"
    check_method          = "GET"
    connect_timeout       = 5000
    health_check_interval = 30
    healthy_threshold     = 3
    unhealthy_threshold   = 3
    forward_protocol      = "HTTP"

    endpoint_configurations {
      endpoint_type    = "PublicIp"
      endpoint_service = "1.1.1.1"
      weight           = 10
    }

    endpoint_configurations {
      endpoint_type    = "Domain"
      endpoint_service = "example.com"
      weight           = 20
    }
  }
}
`
