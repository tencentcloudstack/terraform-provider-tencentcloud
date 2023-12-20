package cos_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCosBucketVersionResource_basic -v
func TestAccTencentCloudCosBucketVersionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucketVersion,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cos_bucket_version.bucket_version", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket_version.bucket_version", "status", "Enabled"),
				),
			},
			{
				ResourceName:      "tencentcloud_cos_bucket_version.bucket_version",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCosBucketVersionUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cos_bucket_version.bucket_version", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket_version.bucket_version", "status", "Suspended"),
				),
			},
		},
	})
}

const testAccCosBucketVersionVar = `
variable "bucket" {
	default = "` + tcacctest.DefaultCiBucket + `"
}

`

const testAccCosBucketVersion = testAccCosBucketVersionVar + `

resource "tencentcloud_cos_bucket_version" "bucket_version" {
	bucket = var.bucket
	status = "Enabled"
}

`

const testAccCosBucketVersionUp = testAccCosBucketVersionVar + `

resource "tencentcloud_cos_bucket_version" "bucket_version" {
	bucket = var.bucket
	status = "Suspended"
}

`
