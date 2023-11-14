package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCwpDescribeMachinesSimpleDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCwpDescribeMachinesSimpleDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cwp_describe_machines_simple.describe_machines_simple")),
			},
		},
	})
}

const testAccCwpDescribeMachinesSimpleDataSource = `

data "tencentcloud_cwp_describe_machines_simple" "describe_machines_simple" {
  machine_type = ""
  machine_region = ""
  filters {
		name = ""
		values = 

  }
  project_ids = 
  }

`
