package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCiPosterTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCiPosterTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ci_poster_template.poster_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_ci_poster_template.poster_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCiPosterTemplate = `

resource "tencentcloud_ci_poster_template" "poster_template" {
  input = &lt;nil&gt;
  name = &lt;nil&gt;
  category_ids = &lt;nil&gt;
}

`
