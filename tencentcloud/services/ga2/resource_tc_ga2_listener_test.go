package ga2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudGa2ListenerResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGa2Listener,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ga2_listener.example", "id"),
				),
			},
			{
				Config: testAccGa2ListenerUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ga2_listener.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_ga2_listener.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccGa2Listener = `
resource "tencentcloud_ga2_listener" "example" {
  global_accelerator_id = "ga-4mredmiu"
  name                  = "tf-example"
  protocol              = "TCP"

  port_ranges {
    from_port = 80
    to_port   = 80
  }

  description     = "tf example listener"
  client_affinity = "Open"
  idle_timeout    = 900
}
`

const testAccGa2ListenerUpdate = `
resource "tencentcloud_ga2_listener" "example" {
  global_accelerator_id = "ga-4mredmiu"
  name                  = "tf-example"
  protocol              = "TCP"

  port_ranges {
    from_port = 80
    to_port   = 80
  }

  description     = "tf example listener"
  client_affinity = "Open"
  idle_timeout    = 900
}
`
