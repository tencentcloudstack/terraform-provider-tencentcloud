package dnspod_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudSubdomainValidateTxtValueOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSubdomainValidateTxtValueOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_subdomain_validate_txt_value_operation.subdomain_validate_txt_value_operation", "id"),
					resource.TestCheckResourceAttr("tencentcloud_subdomain_validate_txt_value_operation.subdomain_validate_txt_value_operation", "domain", "iac-tf.cloud"),
					resource.TestCheckResourceAttr("tencentcloud_subdomain_validate_txt_value_operation.subdomain_validate_txt_value_operation", "domain_zone", "www.iac-tf.cloud"),
					resource.TestCheckResourceAttr("tencentcloud_subdomain_validate_txt_value_operation.subdomain_validate_txt_value_operation", "record_type", "TXT"),
					resource.TestCheckResourceAttrSet("tencentcloud_subdomain_validate_txt_value_operation.subdomain_validate_txt_value_operation", "subdomain"),
					resource.TestCheckResourceAttrSet("tencentcloud_subdomain_validate_txt_value_operation.subdomain_validate_txt_value_operation", "value"),
				),
			},
		},
	})
}

const testAccSubdomainValidateTxtValueOperation = `
resource "tencentcloud_subdomain_validate_txt_value_operation" "subdomain_validate_txt_value_operation" {
  domain_zone = "www.iac-tf.cloud"
}
`
