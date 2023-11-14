package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSsmDescribeRotationHistoryDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSsmDescribeRotationHistoryDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssm_describe_rotation_history.describe_rotation_history")),
			},
		},
	})
}

const testAccSsmDescribeRotationHistoryDataSource = `

data "tencentcloud_ssm_describe_rotation_history" "describe_rotation_history" {
  secret_name = ""
  }

`
