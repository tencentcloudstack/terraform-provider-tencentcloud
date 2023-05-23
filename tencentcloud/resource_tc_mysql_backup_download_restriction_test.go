package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlBackupDownloadRestrictionResource_basic -v
func TestAccTencentCloudMysqlBackupDownloadRestrictionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlBackupDownloadRestriction,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_mysql_backup_download_restriction.backup_download_restriction", "limit_type", "Customize"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_backup_download_restriction.backup_download_restriction", "vpc_comparison_symbol", "In"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_backup_download_restriction.backup_download_restriction", "ip_comparison_symbol", "In"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_backup_download_restriction.backup_download_restriction", "limit_vpc.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_backup_download_restriction.backup_download_restriction", "limit_vpc.0.region", "ap-guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_backup_download_restriction.backup_download_restriction", "limit_vpc.0.vpc_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_backup_download_restriction.backup_download_restriction", "limit_vpc.0.vpc_list.0", "vpc-4owdpnwr"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_backup_download_restriction.backup_download_restriction", "limit_ip.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_backup_download_restriction.backup_download_restriction", "limit_ip.0", "127.0.0.1"),
				),
			},
			{
				ResourceName:      "tencentcloud_mysql_backup_download_restriction.backup_download_restriction",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMysqlBackupDownloadRestriction = `

resource "tencentcloud_mysql_backup_download_restriction" "backup_download_restriction" {
	limit_type = "Customize"
	vpc_comparison_symbol = "In"
	ip_comparison_symbol = "In"
	limit_vpc {
		  region = "ap-guangzhou"
		  vpc_list = ["vpc-4owdpnwr"]
	}
	limit_ip = ["127.0.0.1"]
}  

`
