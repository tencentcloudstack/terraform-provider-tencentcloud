package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcTenantCcnDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcTenantCcnDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_tenant_ccn.tenant_ccn")),
			},
		},
	})
}

const testAccVpcTenantCcnDataSource = `

data "tencentcloud_vpc_tenant_ccn" "tenant_ccn" {
  name = "ccn-ids"
  values = 
  offset = 0
  limit = 20
      }

`
