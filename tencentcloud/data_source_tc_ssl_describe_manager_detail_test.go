package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSslDescribeManagerDetailDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeManagerDetailDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_manager_detail.describe_manager_detail")),
			},
		},
	})
}

const testAccSslDescribeManagerDetailDataSource = `

data "tencentcloud_ssl_describe_manager_detail" "describe_manager_detail" {
  manager_id = "12895"
}

`
