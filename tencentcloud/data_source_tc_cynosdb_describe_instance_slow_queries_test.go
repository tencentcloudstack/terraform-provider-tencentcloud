package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbDescribeInstanceSlowQueriesDataSource_basic -v
func TestAccTencentCloudCynosdbDescribeInstanceSlowQueriesDataSource_basic(t *testing.T) {
	t.Parallel()
	loc, _ := time.LoadLocation("Asia/Chongqing")
	startTime := time.Now().AddDate(0, 0, -7).In(loc).Format("2006-01-02 15:04:05")
	endTime := time.Now().In(loc).Format("2006-01-02 15:04:05")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCynosdbDescribeInstanceSlowQueriesDataSource, startTime, endTime),
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_describe_instance_slow_queries.describe_instance_slow_queries")),
			},
		},
	})
}

const testAccCynosdbDescribeInstanceSlowQueriesDataSource = CommonCynosdb + `
data "tencentcloud_cynosdb_describe_instance_slow_queries" "describe_instance_slow_queries" {
  cluster_id = var.cynosdb_cluster_id
  start_time = "%s"
  end_time   = "%s"
}
`
