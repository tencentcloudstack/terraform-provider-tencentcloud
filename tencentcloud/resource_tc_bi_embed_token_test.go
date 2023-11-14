package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudBiEmbedTokenResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBiEmbedToken,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_bi_embed_token.embed_token", "id")),
			},
			{
				ResourceName:      "tencentcloud_bi_embed_token.embed_token",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccBiEmbedToken = `

resource "tencentcloud_bi_embed_token" "embed_token" {
  project_id = 123
  page_id = 123
  scope = "page"
  expire_time = "240"
  extra_param = ""
  user_corp_id = "abc"
  user_id = "abc"
}

`
