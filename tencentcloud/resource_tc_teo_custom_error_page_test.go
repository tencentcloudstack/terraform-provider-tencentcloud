package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTeoCustomErrorPageResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoCustomErrorPage,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_custom_error_page.custom_error_page", "id")),
			},
			{
				ResourceName:      "tencentcloud_teo_custom_error_page.custom_error_page",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoCustomErrorPage = `

resource "tencentcloud_teo_custom_error_page" "custom_error_page" {
  zone_id = &lt;nil&gt;
  entity = &lt;nil&gt;
    name = &lt;nil&gt;
  content = &lt;nil&gt;
}

`
