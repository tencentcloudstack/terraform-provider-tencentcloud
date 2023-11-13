package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudEbSearchDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEbSearchDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_eb_search.search")),
			},
		},
	})
}

const testAccEbSearchDataSource = `

data "tencentcloud_eb_search" "search" {
  start_time = 
  end_time = 
  event_bus_id = ""
  group_field = ""
  filter {
		key = ""
		operator = ""
		value = ""
		type = ""
		filters {
			key = ""
			operator = ""
			value = ""
		}

  }
  order_fields = 
  order_by = ""
  }

`
