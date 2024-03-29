package dayuv2_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAntiddosSchedulingDomainUserNameResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiddosSchedulingDomainUserName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_antiddos_scheduling_domain_user_name.scheduling_domain_user_name", "id"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_scheduling_domain_user_name.scheduling_domain_user_name", "domain_name", "68tlc0iy.dayugslb.com"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_scheduling_domain_user_name.scheduling_domain_user_name", "domain_user_name", "tf-test"),
				),
			},
			{
				Config: testAccAntiddosSchedulingDomainUserNameUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_antiddos_scheduling_domain_user_name.scheduling_domain_user_name", "id"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_scheduling_domain_user_name.scheduling_domain_user_name", "domain_name", "68tlc0iy.dayugslb.com"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_scheduling_domain_user_name.scheduling_domain_user_name", "domain_user_name", "tf-test-update"),
				),
			},
			{
				ResourceName:      "tencentcloud_antiddos_scheduling_domain_user_name.scheduling_domain_user_name",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccAntiddosSchedulingDomainUserName = `
resource "tencentcloud_antiddos_scheduling_domain_user_name" "scheduling_domain_user_name" {
	domain_name = "68tlc0iy.dayugslb.com"
	domain_user_name = "tf-test"
}
`

const testAccAntiddosSchedulingDomainUserNameUpdate = `
resource "tencentcloud_antiddos_scheduling_domain_user_name" "scheduling_domain_user_name" {
	domain_name = "68tlc0iy.dayugslb.com"
	domain_user_name = "tf-test-update"
}
`
