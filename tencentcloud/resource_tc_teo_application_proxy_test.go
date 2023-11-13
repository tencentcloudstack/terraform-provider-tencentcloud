package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTeoApplicationProxyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoApplicationProxy,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_application_proxy.application_proxy", "id")),
			},
			{
				ResourceName:      "tencentcloud_teo_application_proxy.application_proxy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoApplicationProxy = `

resource "tencentcloud_teo_application_proxy" "application_proxy" {
  zone_id = &lt;nil&gt;
    proxy_name = &lt;nil&gt;
  proxy_type = &lt;nil&gt;
  plat_type = &lt;nil&gt;
    security_type = &lt;nil&gt;
  accelerate_type = &lt;nil&gt;
  session_persist_time = &lt;nil&gt;
  status = &lt;nil&gt;
        i_pv6 {
		switch = &lt;nil&gt;

  }
  }

`
