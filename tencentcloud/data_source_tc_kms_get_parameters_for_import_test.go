package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixKmsGetParametersForImportDataSource_basic -v
func TestAccTencentCloudNeedFixKmsGetParametersForImportDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsGetParametersForImportDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_kms_get_parameters_for_import.example"),
				),
			},
		},
	})
}

const testAccKmsGetParametersForImportDataSource = `
data "tencentcloud_kms_get_parameters_for_import" "example" {
  key_id             = "786aea8c-4aec-11ee-b601-525400281a45"
  wrapping_algorithm = "RSAES_OAEP_SHA_1"
  wrapping_key_spec  = "RSA_2048"
}
`
