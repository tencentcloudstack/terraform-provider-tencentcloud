package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudSsmSecretsDataSource(t *testing.T) {
	dataSourceName := "data.tencentcloud_ssm_secrets.secret"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccTencentCloudSsmSecretsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "secret_list.0.secret_name", "unit-test"),
					resource.TestCheckResourceAttr(dataSourceName, "secret_list.0.description", "test secret"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secret_list.0.kms_key_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secret_list.0.create_uin"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secret_list.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secret_list.0.create_time"),
				),
			},
		},
	})
}

const TestAccTencentCloudSsmSecretsDataSourceConfig = `
resource "tencentcloud_ssm_secret" "secret" {
  secret_name = "unit-test"
  description = "test secret"
  init_secret {
    version_id = "v1"
    secret_string = "123456789"
  }

  tags = {
    test-tag = "test"
  }
}

data "tencentcloud_ssm_secrets" "secret" {
  secret_name = tencentcloud_ssm_secret.secret.secret_name
  state = 1
  
  tags = {
    test-tag = "test"
  }
}
`
