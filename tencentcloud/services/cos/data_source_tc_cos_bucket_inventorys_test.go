package cos_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCosBucketInventorysDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucketInventorysDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cos_bucket_inventorys.cos_bucket_inventorys"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cos_bucket_inventorys.cos_bucket_inventorys", "inventorys.#"),
				),
			},
		},
	})
}

const testAccCosBucketInventorysDataSource = `
data "tencentcloud_cos_bucket_inventorys" "cos_bucket_inventorys" {
    bucket = "keep-test-1308919341"
}
`
