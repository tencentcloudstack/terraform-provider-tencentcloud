package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainDiagDbInstancesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainDiagDbInstancesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_diag_db_instances.diag_db_instances")),
			},
		},
	})
}

const testAccDbbrainDiagDbInstancesDataSource = `

data "tencentcloud_dbbrain_diag_db_instances" "diag_db_instances" {
  is_supported = 
  product = ""
  instance_names = 
  instance_ids = 
  regions = 
    }

`
