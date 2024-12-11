package dnspod_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudSubdomainValidateStatusDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccSubdomainValidateStatusDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_subdomain_validate_status.subdomain_validate_status"),
				resource.TestCheckResourceAttr("data.tencentcloud_subdomain_validate_status.subdomain_validate_status", "status", "1"),
			),
		}},
	})
}

func TestAccTencentCloudSubdomainValidateStatusDataSource_notReady(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccSubdomainValidateStatusDataSourceNotReady,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_subdomain_validate_status.subdomain_validate_status"),
				resource.TestCheckResourceAttr("data.tencentcloud_subdomain_validate_status.subdomain_validate_status", "status", "0"),
			),
		}},
	})
}

const testAccSubdomainValidateStatusDataSource = `
data "tencentcloud_subdomain_validate_status" "subdomain_validate_status" {
  domain_zone = "www.iac-tf.cloud"
}
`

const testAccSubdomainValidateStatusDataSourceNotReady = `
data "tencentcloud_subdomain_validate_status" "subdomain_validate_status" {
  domain_zone = "www.iac-tf.com"
}
`
