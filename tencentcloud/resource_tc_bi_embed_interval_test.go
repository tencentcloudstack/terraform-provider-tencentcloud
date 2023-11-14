package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudBiEmbedIntervalResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBiEmbedInterval,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_bi_embed_interval.embed_interval", "id")),
			},
			{
				ResourceName:      "tencentcloud_bi_embed_interval.embed_interval",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccBiEmbedInterval = `

resource "tencentcloud_bi_embed_interval" "embed_interval" {
  project_id = 123
  page_id = 123
  b_i_token = "abc"
  extra_param = ""
  scope = "page"
}

`
