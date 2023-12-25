package tsf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfReleaseApiGroupResource_basic -v
func TestAccTencentCloudTsfReleaseApiGroupResource_basic(t *testing.T) {
	// t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTsfUnitNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfReleaseApiGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_release_api_group.release_api_group", "id"),
				),
			},
		},
	})
}

const testAccTsfReleaseApiGroup = `

resource "tencentcloud_tsf_release_api_group" "release_api_group" {
   group_id = "grp-qp0rj3zi"
}
`
