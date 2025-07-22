package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoContentIdentifierResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoContentIdentifier,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_content_identifier.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_content_identifier.example", "plan_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_content_identifier.example", "description"),
				),
			},
			{
				Config: testAccTeoContentIdentifierUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_content_identifier.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_content_identifier.example", "plan_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_content_identifier.example", "description"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_content_identifier.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoContentIdentifier = `
resource "tencentcloud_teo_content_identifier" "example" {
  plan_id     = "edgeone-3bzvsgjkfw6g"
  description = "example"
  tags {
    tag_key   = "tagKey"
    tag_value = "tagValue"
  }
}
`

const testAccTeoContentIdentifierUpdate = `
resource "tencentcloud_teo_content_identifier" "example" {
  plan_id     = "edgeone-3bzvsgjkfw6g"
  description = "example update"
  tags {
    tag_key   = "tagKey"
    tag_value = "tagValue"
  }
}
`
