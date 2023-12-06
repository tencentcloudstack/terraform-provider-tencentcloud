package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClickhouseInstanceShardsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClickhouseInstanceShardsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_clickhouse_instance_shards.instance_shards")),
			},
		},
	})
}

const testAccClickhouseInstanceShardsDataSource = `
data "tencentcloud_clickhouse_instance_shards" "instance_shards" {
  instance_id = "cdwch-datuhk3z"
}
`
