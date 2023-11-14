package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMongodbInstanceBackupDownloadTaskResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstanceBackupDownloadTask,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance_backup_download_task.instance_backup_download_task", "id")),
			},
			{
				ResourceName:      "tencentcloud_mongodb_instance_backup_download_task.instance_backup_download_task",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMongodbInstanceBackupDownloadTask = `

resource "tencentcloud_mongodb_instance_backup_download_task" "instance_backup_download_task" {
  instance_id = ""
  backup_name = "myBackupName"
  backup_sets {
		replica_set_id = ""

  }
}

`
