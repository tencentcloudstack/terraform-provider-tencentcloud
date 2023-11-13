package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCiHotLinkResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCiHotLink,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ci_hot_link.hot_link", "id")),
			},
			{
				ResourceName:      "tencentcloud_ci_hot_link.hot_link",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCiHotLink = `

resource "tencentcloud_ci_hot_link" "hot_link" {
  bucket = "terraform-ci-xxxxxx"
  hot_link {
		url = &lt;nil&gt;
		type = &lt;nil&gt;

  }
}

`
