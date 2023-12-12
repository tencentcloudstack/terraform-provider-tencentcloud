package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverInstanceSslResource_basic -v
func TestAccTencentCloudSqlserverInstanceSslResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverInstanceSslEnable,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_instance_ssl.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_instance_ssl.example", "instance_id", "mssql-qelbzgwf"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_instance_ssl.example", "type", "enable"),
				),
			},
			{
				Config: testAccSqlserverInstanceSslRenew,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_instance_ssl.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_instance_ssl.example", "instance_id", "mssql-qelbzgwf"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_instance_ssl.example", "type", "renew"),
				),
			},
			{
				Config: testAccSqlserverInstanceSslDisable,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_instance_ssl.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_instance_ssl.example", "instance_id", "mssql-qelbzgwf"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_instance_ssl.example", "type", "disable"),
				),
			},
		},
	})
}

const testAccSqlserverInstanceSslEnable = `
resource "tencentcloud_sqlserver_instance_ssl" "example" {
  instance_id = "mssql-qelbzgwf"
  type        = "enable"
}
`

const testAccSqlserverInstanceSslRenew = `
resource "tencentcloud_sqlserver_instance_ssl" "example" {
  instance_id = "mssql-qelbzgwf"
  type        = "renew"
}
`

const testAccSqlserverInstanceSslDisable = `
resource "tencentcloud_sqlserver_instance_ssl" "example" {
  instance_id = "mssql-qelbzgwf"
  type        = "disable"
}
`
