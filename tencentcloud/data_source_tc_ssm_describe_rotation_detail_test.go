package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSsmDescribeRotationDetailDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSsmDescribeRotationDetailDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssm_describe_rotation_detail.describe_rotation_detail")),
			},
		},
	})
}

const testAccSsmDescribeRotationDetailDataSource = `

data "tencentcloud_ssm_describe_rotation_detail" "describe_rotation_detail" {
  secret_name = ""
        }

`
