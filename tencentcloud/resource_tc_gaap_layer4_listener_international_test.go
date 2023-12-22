package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudInternationalGaapResource_listener4(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInternationalGaapLayer4ListenerBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer4ListenerExists("tencentcloud_gaap_layer4_listener.foo", id, "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "name", "ci-test-gaap-4-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "port", "9090"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "interval", "5"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "connect_timeout", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_bind_set.#", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "proxy_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_gaap_layer4_listener.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var testAccInternationalGaapLayer4ListenerBasic = `
resource tencentcloud_gaap_layer4_listener "foo" {
  protocol        = "TCP"
  name            = "ci-test-gaap-4-listener"
  port            = 9090
  realserver_type = "IP"
  proxy_id        = "link-092unjav"
  health_check    = true

  realserver_bind_set {
    id   = "rs-5qfe73i1"
    ip   = "1.1.1.5"
    port = 80
  }

  realserver_bind_set {
    id     = "rs-7smh4tmf"
    ip     = "119.29.29.35"
    port   = 80
  }
}
`
