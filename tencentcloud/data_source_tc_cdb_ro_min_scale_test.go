package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbRoMinScaleDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbRoMinScaleDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cdb_ro_min_scale.ro_min_scale")),
			},
		},
	})
}

const testAccCdbRoMinScaleDataSource = `

data "tencentcloud_cdb_ro_min_scale" "ro_min_scale" {
  ro_instance_id = ""
  master_instance_id = ""
    }

`
