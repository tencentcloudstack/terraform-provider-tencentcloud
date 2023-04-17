package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNNeedFixMongodbInstanceBackupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstanceBackup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance_backup.instance_backup", "id")),
			},
		},
	})
}

const testAccMongodbInstanceBackup = `

resource "tencentcloud_mongodb_instance_backup" "instance_backup" {
  instance_id = "cmgo-jbrmgzfl"
  backup_method = 0
  backup_remark = "my backup"
}

`
