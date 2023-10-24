package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixBiEmbedTokenApplyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBiEmbedToken,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_bi_embed_token_apply.embed_token", "id")),
			},
		},
	})
}

const testAccBiEmbedToken = `

resource "tencentcloud_bi_embed_token_apply" "embed_token" {
  project_id   = 11015030
  page_id      = 10520483
  scope        = "page"
  expire_time  = "240"
  user_corp_id = ""
  user_id      = ""
}

`
