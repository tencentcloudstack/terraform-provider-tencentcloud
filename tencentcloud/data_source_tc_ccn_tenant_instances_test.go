package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixVpcTenantCcnDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcTenantCcnDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("tencentcloud_ccn_tenant_instances.tenant_ccn")),
			},
		},
	})
}

const testAccVpcTenantCcnDataSource = `

data "tencentcloud_ccn_tenant_instances" "tenant_ccn" {
  ccn_ids = ["ccn-39lqkygf"]
  is_security_lock = ["true"]
}

`
