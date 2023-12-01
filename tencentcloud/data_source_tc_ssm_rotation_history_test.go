package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixSsmRotationHistoryDataSource_basic -v
func TestAccTencentCloudNeedFixSsmRotationHistoryDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSsmRotationHistoryDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssm_rotation_history.example"),
				),
			},
		},
	})
}

const testAccSsmRotationHistoryDataSource = `
data "tencentcloud_ssm_rotation_history" "example" {
  secret_name = "keep_terraform"
}
`
