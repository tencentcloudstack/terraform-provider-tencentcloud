package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixCcnTenantInstancesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCcnTenantInstancesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("tencentcloud_ccn_tenant_instances.tenant_ccn")),
			},
		},
	})
}

const testAccCcnTenantInstancesDataSource = `

data "tencentcloud_ccn_tenant_instances" "tenant_ccn" {
  ccn_ids = ["ccn-39lqkygf"]
  is_security_lock = ["true"]
}

`
