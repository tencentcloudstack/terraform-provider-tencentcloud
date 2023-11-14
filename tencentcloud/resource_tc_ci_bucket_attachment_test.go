package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCiBucketAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCiBucketAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ci_bucket_attachment.bucket_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_ci_bucket_attachment.bucket_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCiBucketAttachment = `

resource "tencentcloud_ci_bucket_attachment" "bucket_attachment" {
  bucket = "terraform-ci-xxxxxx"
  ci_status = &lt;nil&gt;
}

`
