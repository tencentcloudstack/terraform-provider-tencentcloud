package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClickhouseSpecDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClickhouseSpecDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_clickhouse_spec.spec")),
			},
		},
	})
}

const testAccClickhouseSpecDataSource = `
data "tencentcloud_clickhouse_spec" "spec" {
  zone       = "ap-guangzhou-7"
  pay_mode   = "PREPAID"
  is_elastic = false
}
`
