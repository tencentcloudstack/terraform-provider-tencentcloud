package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsContentReviewTemplatesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsContentReviewTemplatesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mps_content_review_templates.content_review_templates")),
			},
		},
	})
}

const testAccMpsContentReviewTemplatesDataSource = `

data "tencentcloud_mps_content_review_templates" "content_review_templates" {
  definitions = 
  type = ""
  }

`
