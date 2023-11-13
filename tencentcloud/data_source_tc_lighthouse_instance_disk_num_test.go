package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLighthouseInstanceDiskNumDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseInstanceDiskNumDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_lighthouse_instance_disk_num.instance_disk_num")),
			},
		},
	})
}

const testAccLighthouseInstanceDiskNumDataSource = `

data "tencentcloud_lighthouse_instance_disk_num" "instance_disk_num" {
  instance_ids = 
}

`
