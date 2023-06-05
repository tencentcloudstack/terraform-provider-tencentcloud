package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcSecurityGroupReferencesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSecurityGroupReferencesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_security_group_references.security_group_references")),
			},
		},
	})
}

const testAccVpcSecurityGroupReferencesDataSource = `

data "tencentcloud_vpc_security_group_references" "security_group_references" {
  security_group_ids = ["sg-edmur627"]
}
`
