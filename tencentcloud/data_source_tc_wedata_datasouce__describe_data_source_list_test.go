package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudWedataDatasouce_DescribeDataSourceListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataDatasouce_DescribeDataSourceListDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_datasouce__describe_data_source_list.datasouce__describe_data_source_list")),
			},
		},
	})
}

const testAccWedataDatasouce_DescribeDataSourceListDataSource = `

data "tencentcloud_wedata_datasouce__describe_data_source_list" "datasouce__describe_data_source_list" {
  order_fields {
		name = ""
		direction = ""

  }
  filters {
		name = ""
		values = 

  }
  }

`
