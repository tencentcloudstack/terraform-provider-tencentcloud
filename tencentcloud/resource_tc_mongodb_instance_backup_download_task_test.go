package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
  instance_id = "cmgo-gwqk8669"
  backup_name = "cmgo-gwqk8669_2023-05-10 10:45"
  backup_sets {
    replica_set_id = "cmgo-gwqk8669_0"
  }
}

`
