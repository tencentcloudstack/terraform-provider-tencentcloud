package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMongodbInstanceBackupsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstanceBackupsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mongodb_instance_backups.instance_backups")),
			},
		},
	})
}

const testAccMongodbInstanceBackupsDataSource = `

data "tencentcloud_mongodb_instance_backups" "instance_backups" {
  instance_id = "cmgo-gwqk8669"
  backup_method = 0
}

`
