package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudRumSignDataSource_basic -v
func TestAccTencentCloudRumSignDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumSignDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_rum_sign.sign"),
				),
			},
		},
	})
}

const testAccRumSignDataSource = `

data "tencentcloud_rum_sign" "sign" {
  timeout = 1800
  file_type = 1
}

`
