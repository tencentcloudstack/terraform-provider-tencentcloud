package ssm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixSsmRotationHistoryDataSource_basic -v
func TestAccTencentCloudNeedFixSsmRotationHistoryDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSsmRotationHistoryDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ssm_rotation_history.example"),
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
