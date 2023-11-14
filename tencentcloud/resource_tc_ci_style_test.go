package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCiStyleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCiStyle,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ci_style.style", "id")),
			},
			{
				ResourceName:      "tencentcloud_ci_style.style",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCiStyle = `

resource "tencentcloud_ci_style" "style" {
  style_name = &lt;nil&gt;
  style_body = &lt;nil&gt;
}

`
