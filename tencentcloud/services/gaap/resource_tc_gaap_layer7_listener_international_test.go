package gaap_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudInternationalGaapResource_listener7(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInternationalGaapLayer7ListenerBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer7ListenerExists("tencentcloud_gaap_layer7_listener.foo", id, "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "protocol", "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "port", "8080"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "client_certificate_ids.#", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "proxy_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_gaap_layer7_listener.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var testAccInternationalGaapLayer7ListenerBasic = `
resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 8080
  proxy_id = "link-092unjav"
}
`
