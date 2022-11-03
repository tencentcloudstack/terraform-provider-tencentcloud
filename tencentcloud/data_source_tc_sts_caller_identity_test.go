package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudStsCallerIdentityDataSource -v
func TestAccTencentCloudStsCallerIdentityDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceStsCallerIdentity,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_sts_caller_identity.caller_identity"),
				),
			},
		},
	})
}

const testAccDataSourceStsCallerIdentity = `

data "tencentcloud_sts_caller_identity" "caller_identity" {
}

`
