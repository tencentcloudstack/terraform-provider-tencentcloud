package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

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
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_binlog_download_url.binlog_download_url")),
			},
		},
	})
}

const testAccCynosdbBinlogDownloadUrlDataSource = `

data "tencentcloud_cynosdb_binlog_download_url" "binlog_download_url" {
  cluster_id = "cynosdbmysql-123"
  binlog_id = 100
  }

`
