package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixSsmRotationDetailDataSource_basic -v
func TestAccTencentCloudNeedFixSsmRotationDetailDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSsmRotationDetailDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssm_rotation_detail.example"),
				),
			},
		},
	})
}

const testAccSsmRotationDetailDataSource = `
data "tencentcloud_ssm_rotation_detail" "example" {
  secret_name = "tf_example"
}
`
