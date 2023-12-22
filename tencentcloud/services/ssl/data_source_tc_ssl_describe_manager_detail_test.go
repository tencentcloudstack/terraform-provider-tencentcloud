package ssl_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDescribeManagerDetailDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_SSL)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeManagerDetailDataSource,
				Check: resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_manager_detail.describe_manager_detail"),
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
