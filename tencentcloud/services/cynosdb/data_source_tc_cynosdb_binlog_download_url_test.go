package cynosdb_test

import (
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbBinlogDownloadUrlDataSource_basic -v
func TestAccTencentCloudCynosdbBinlogDownloadUrlDataSource_basic(t *testing.T) {
	t.Parallel()
	loc, _ := time.LoadLocation("Asia/Chongqing")
	startTime := time.Now().AddDate(0, 0, -7).In(loc).Format("2006-01-02 15:04:05")
	endTime := time.Now().In(loc).Format("2006-01-02 15:04:05")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCynosdbBinlogDownloadUrlDataSource, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_binlog_download_url.binlog_download_url"),
				),
			},
		},
	})
}

const testAccCynosdbBinlogDownloadUrlDataSource = tcacctest.CommonCynosdb + `
data "tencentcloud_cynosdb_describe_instance_slow_queries" "describe_instance_slow_queries" {
  cluster_id = var.cynosdb_cluster_id
  start_time = "%s"
  end_time   = "%s"
}

data "tencentcloud_cynosdb_binlog_download_url" "binlog_download_url" {
  cluster_id = var.cynosdb_cluster_id
  binlog_id  = data.tencentcloud_cynosdb_describe_instance_slow_queries.describe_instance_slow_queries.binlogs.0.binlog_id
}
`
