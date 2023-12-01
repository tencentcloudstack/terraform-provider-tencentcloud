package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSsmSecretVersionsDataSource -v
func TestAccTencentCloudSsmSecretVersionsDataSource(t *testing.T) {
	t.Parallel()
	dataSourceName := "data.tencentcloud_ssm_secret_versions.example"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccTencentCloudSsmSecretVersionsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "secret_version_list.0.version_id", "v1"),
					resource.TestCheckResourceAttr(dataSourceName, "secret_version_list.0.secret_binary", "MTIzMTIzMTIzMTIzMTIzQQ=="),
				),
			},
		},
	})
}

const TestAccTencentCloudSsmSecretVersionsDataSourceConfig = `
data "tencentcloud_ssm_secret_versions" "example" {
  secret_name = tencentcloud_ssm_secret_version.v1.secret_name
  version_id  = tencentcloud_ssm_secret_version.v1.version_id
}

resource "tencentcloud_ssm_secret" "example" {
  secret_name = "tf-example-ssm-unit-test"
  description = "desc."

  tags = {
    createdBy = "terraform"
  }
}

resource "tencentcloud_ssm_secret_version" "v1" {
  secret_name   = tencentcloud_ssm_secret.example.secret_name
  version_id    = "v1"
  secret_binary = "MTIzMTIzMTIzMTIzMTIzQQ=="
}
`
