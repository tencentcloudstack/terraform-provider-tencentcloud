package mongodb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNNeedFixMongodbInstanceBackupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
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
