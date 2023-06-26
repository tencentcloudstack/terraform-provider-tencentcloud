package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbDescribeInstanceSlowQueriesDataSource_basic -v
func TestAccTencentCloudCynosdbDescribeInstanceSlowQueriesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbDescribeInstanceSlowQueriesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_describe_instance_slow_queries.describe_instance_slow_queries")),
			},
		},
	})
}

const testAccCynosdbDescribeInstanceSlowQueriesDataSource = `
data "tencentcloud_cynosdb_describe_instance_slow_queries" "describe_instance_slow_queries" {
  cluster_id = "cynosdbmysql-bws8h88b"
  start_time = "2023-06-01 12:00:00"
  end_time   = "2023-06-19 14:00:00"
}
`
