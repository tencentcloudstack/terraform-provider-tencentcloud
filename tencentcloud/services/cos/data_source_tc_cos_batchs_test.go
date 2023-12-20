package cos_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCosBatchsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBatchsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cos_batchs.cos_batchs"),
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
