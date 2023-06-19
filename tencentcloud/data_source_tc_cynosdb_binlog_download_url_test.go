package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbBinlogDownloadUrlDataSource_basic -v
func TestAccTencentCloudCynosdbBinlogDownloadUrlDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbBinlogDownloadUrlDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_binlog_download_url.binlog_download_url"),
				),
			},
		},
	})
}

const testAccCynosdbBinlogDownloadUrlDataSource = `
data "tencentcloud_cynosdb_binlog_download_url" "binlog_download_url" {
  cluster_id = "cynosdbmysql-bws8h88b"
  binlog_id  = 6202249
}
`
