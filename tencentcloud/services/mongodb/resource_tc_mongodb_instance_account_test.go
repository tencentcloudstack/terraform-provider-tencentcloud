package mongodb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMongodbInstanceAccountResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstanceAccount,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance_account.instance_account", "id")),
			},
			{
				ResourceName:            "tencentcloud_mongodb_instance_account.instance_account",
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"password", "mongo_user_password"},
			},
		},
	})
}

const testAccMongodbInstanceAccount = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "tf-mongodb-test"
  memory         = local.memory
  volume         = local.volume
  engine_version = local.engine_version
  machine_type   = local.machine_type
  security_groups = [local.security_group_id]
  available_zone = "ap-guangzhou-3"
  project_id     = 0
  password       = "test1234"
  vpc_id         = var.vpc_id
  subnet_id      = var.subnet_id
}

resource "tencentcloud_mongodb_instance_account" "instance_account" {
  instance_id = tencentcloud_mongodb_instance.mongodb.id
  user_name = "test_account"
  password = "test1234"
  mongo_user_password = "test1234"
  user_desc = "test account"
  auth_role {
    mask = 0
    namespace = "*"
  }
}
`
