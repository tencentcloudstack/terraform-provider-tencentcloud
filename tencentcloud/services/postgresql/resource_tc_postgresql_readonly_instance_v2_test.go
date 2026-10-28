package postgresql_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudPostgresqlReadonlyInstanceV2Resource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlReadonlyInstanceV2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_readonly_instance_v2.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_readonly_instance_v2.example", "zone"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_readonly_instance_v2.example", "master_db_instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_readonly_instance_v2.example", "spec_code"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_readonly_instance_v2.example", "storage"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_readonly_instance_v2.example", "instance_count"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_readonly_instance_v2.example", "period"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_readonly_instance_v2.example", "db_instance_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_postgresql_readonly_instance_v2.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresqlReadonlyInstanceV2 = `
resource "tencentcloud_postgresql_readonly_instance_v2" "example" {
  zone                  = "ap-guangzhou-3"
  master_db_instance_id = "postgres-xxxxxx"
  spec_code             = "pgxz.sn1.2g.2c.ha"
  storage               = 30
  instance_count        = 1
  period                = 1
  vpc_id                = "vpc-xxxxxxxx"
  subnet_id             = "subnet-xxxxxxxx"
  instance_charge_type  = "POSTPAID_BY_HOUR"
  name                  = "tf-example-readonly"
  project_id            = 0

  tags = {
    createdBy = "terraform"
  }
}
`
