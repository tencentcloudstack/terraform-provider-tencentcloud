package dnspod_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDnspodDomainInstancesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodDomainInstancesDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dnspod_domain_instances.domain_instances")),
			},
		},
	})
}

const testAccDnspodDomainInstancesDataSource = `

data "tencentcloud_dnspod_domain_instances" "domain_instances" {
	domain = "tencentcloud-terraform-provider.cn"
}
`
