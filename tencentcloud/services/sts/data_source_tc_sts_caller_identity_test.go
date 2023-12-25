package sts_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudStsCallerIdentityDataSource -v
func TestAccTencentCloudStsCallerIdentityDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceStsCallerIdentity,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_sts_caller_identity.caller_identity"),
				),
			},
		},
	})
}

const testAccDataSourceStsCallerIdentity = `

data "tencentcloud_sts_caller_identity" "caller_identity" {
}

`
