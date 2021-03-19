package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudSsmSecretVersionsDataSource(t *testing.T) {
	dataSourceName := "data.tencentcloud_ssm_secret_versions.secret_version"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccTencentCloudSsmSecretVersionsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "secret_version_list.0.version_id", "v2"),
					resource.TestCheckResourceAttr(dataSourceName, "secret_version_list.0.secret_binary", "MTIzMTIzMTIzMTIzMTIzQQ=="),
				),
			},
		},
	})
}

const TestAccTencentCloudSsmSecretVersionsDataSourceConfig = `
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

resource "tencentcloud_ssm_secret_version" "v2" {
  secret_name = tencentcloud_ssm_secret.secret.secret_name
  version_id = "v2"
  secret_binary = "MTIzMTIzMTIzMTIzMTIzQQ=="
}

data "tencentcloud_ssm_secret_versions" "secret_version" {
  secret_name = tencentcloud_ssm_secret_version.v2.secret_name
  version_id = tencentcloud_ssm_secret_version.v2.version_id
}
`
