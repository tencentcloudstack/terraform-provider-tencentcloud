package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudKmsGetParametersForImportDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsGetParametersForImportDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_kms_get_parameters_for_import.get_parameters_for_import")),
			},
		},
	})
}

const testAccKmsGetParametersForImportDataSource = `

data "tencentcloud_kms_get_parameters_for_import" "get_parameters_for_import" {
      wrapping_algorithm = "RSAES_OAEP_SHA_1"
  wrapping_key_spec = "RSA_2048"
    }

`
