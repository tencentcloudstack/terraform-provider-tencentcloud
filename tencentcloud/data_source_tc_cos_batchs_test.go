package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCosBatchsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBatchsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_cos_batchs.cos_batchs"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cos_batchs.cos_batchs", "jobs.#"),
				),
			},
		},
	})
}

const testAccCosBatchsDataSource = `
data "tencentcloud_cos_batchs" "cos_batchs" {
    uin = "100022975249"
    appid = "1308919341"
}
`
