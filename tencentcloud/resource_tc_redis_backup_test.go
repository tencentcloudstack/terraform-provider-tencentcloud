package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisBackupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisBackup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_backup.backup", "id")),
			},
			{
				ResourceName:      "tencentcloud_redis_backup.backup",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisBackup = `

resource "tencentcloud_redis_backup" "backup" {
  limit_type = "NoLimit"
  vpc_comparison_symbol = "In"
  ip_comparison_symbol = "In"
  limit_vpc {
		region = "ap-guangzhou"
		vpc_list = 

  }
  limit_ip = 
}

`
