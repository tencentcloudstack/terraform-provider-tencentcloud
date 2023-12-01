package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDescribeManagerDetailDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_SSL)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeManagerDetailDataSource,
				Check: resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_manager_detail.describe_manager_detail"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_describe_manager_detail.describe_manager_detail", "manager_id", "12895"),
				),
			},
		},
	})
}

const testAccSslDescribeManagerDetailDataSource = `

data "tencentcloud_ssl_describe_manager_detail" "describe_manager_detail" {
  manager_id = "12895"
}

`
