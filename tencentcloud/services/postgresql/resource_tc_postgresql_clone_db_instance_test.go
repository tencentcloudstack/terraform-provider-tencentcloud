package postgresql_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudNeedFixPostgresqlCloneDbInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlCloneDbInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_clone_db_instance.postgresql_clone_db_instance", "id")),
			},
		},
	})
}

const testAccPostgresqlCloneDbInstance = `
resource "tencentcloud_postgresql_clone_db_instance" "example" {
  db_instance_id       = "postgres-evsqpyap"
  name                 = "tf-example-clone"
  spec_code            = "pg.it.medium4"
  storage              = 200
  period               = 1
  auto_renew_flag      = 0
  vpc_id               = "vpc-a6zec4mf"
  subnet_id            = "subnet-b8hintyy"
  instance_charge_type = "POSTPAID_BY_HOUR"
  security_group_ids   = ["sg-8stavs03"]
  project_id           = 0
  recovery_target_time = "2024-10-12 18:17:00"
  db_node_set {
    role = "Primary"
    zone = "ap-guangzhou-6"
  }

  db_node_set {
    role = "Standby"
    zone = "ap-guangzhou-6"
  }

  tag_list {
    tag_key   = "createBy"
    tag_value = "terraform"
  }
}
`
