package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfDescribeUsableUnitNamespacesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfDescribeUsableUnitNamespacesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_describe_usable_unit_namespaces.describe_usable_unit_namespaces")),
			},
		},
	})
}

const testAccTsfDescribeUsableUnitNamespacesDataSource = `

data "tencentcloud_tsf_describe_usable_unit_namespaces" "describe_usable_unit_namespaces" {
  search_word = ""
  }

`
