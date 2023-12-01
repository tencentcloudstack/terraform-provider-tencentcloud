package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSsmSecretsDataSource -v
func TestAccTencentCloudSsmSecretsDataSource(t *testing.T) {
	t.Parallel()
	dataSourceName := "data.tencentcloud_ssm_secrets.example"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccTencentCloudSsmSecretsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "secret_list.0.secret_name", "tf_example_ssm_secret"),
					resource.TestCheckResourceAttr(dataSourceName, "secret_list.0.description", "desc."),
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
data "tencentcloud_ssm_secrets" "example" {
  secret_name = tencentcloud_ssm_secret.example.secret_name
  state       = 1
}

resource "tencentcloud_ssm_secret" "example" {
  secret_name = "tf_example_ssm_secret"
  description = "desc."

  tags = {
    createdBy = "terraform"
  }
}
`
