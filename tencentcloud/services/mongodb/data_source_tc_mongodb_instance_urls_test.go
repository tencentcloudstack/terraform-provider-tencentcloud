package mongodb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudMongodbInstanceUrlsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstanceUrlsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mongodb_instance_urls.mongodb_instance_urls"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_instance_urls.mongodb_instance_urls", "urls.#"),
				),
			},
		},
	})
}

const testAccMongodbInstanceUrlsDataSource = testAccMongodbInstance + `
data "tencentcloud_mongodb_instance_urls" "mongodb_instance_urls" {
  instance_id = tencentcloud_mongodb_instance.mongodb.id
}
`
