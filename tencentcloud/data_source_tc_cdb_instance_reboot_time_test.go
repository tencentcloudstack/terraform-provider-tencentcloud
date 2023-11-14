package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbInstanceRebootTimeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbInstanceRebootTimeDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cdb_instance_reboot_time.instance_reboot_time")),
			},
		},
	})
}

const testAccCdbInstanceRebootTimeDataSource = `

data "tencentcloud_cdb_instance_reboot_time" "instance_reboot_time" {
  instance_ids = 
  }

`
