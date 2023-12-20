package cos_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCosBucketMultipartUploadsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucketMultipartUploadsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cos_bucket_multipart_uploads.cos_bucket_multipart_uploads"),
				),
			},
		},
	})
}

const testAccCosBucketMultipartUploadsDataSource = `
data "tencentcloud_cos_bucket_multipart_uploads" "cos_bucket_multipart_uploads" {
    bucket = "keep-test-1308919341"
}
`
