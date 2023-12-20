package ccn_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixCcnTenantInstancesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCcnTenantInstancesDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("tencentcloud_ccn_tenant_instances.tenant_ccn")),
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
