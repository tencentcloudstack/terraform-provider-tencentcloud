package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixBiEmbedIntervalResource_basic(t *testing.T) {
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
		},
	})
}

const testAccBiEmbedInterval = testAccBiEmbedToken + `

resource "tencentcloud_bi_embed_interval" "embed_interval" {
  project_id = 11015030
  page_id    = 10520483
  bi_token   = tencentcloud_bi_embed_token.embed_token.bi_token
  scope      = "page"
}

`
