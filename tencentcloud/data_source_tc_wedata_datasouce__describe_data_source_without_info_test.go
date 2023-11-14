package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudWedataDatasouce_DescribeDataSourceWithoutInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataDatasouce_DescribeDataSourceWithoutInfoDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_datasouce__describe_data_source_without_info.datasouce__describe_data_source_without_info")),
			},
		},
	})
}

const testAccWedataDatasouce_DescribeDataSourceWithoutInfoDataSource = `

data "tencentcloud_wedata_datasouce__describe_data_source_without_info" "datasouce__describe_data_source_without_info" {
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
