package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseInstanceDiskNumDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseInstanceDiskNumDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_lighthouse_instance_disk_num.instance_disk_num")),
			},
		},
	})
}

const testAccLighthouseInstanceDiskNumDataSource = DefaultLighthoustVariables + `

data "tencentcloud_lighthouse_instance_disk_num" "instance_disk_num" {
  instance_ids = [var.lighthouse_id]
}
`
