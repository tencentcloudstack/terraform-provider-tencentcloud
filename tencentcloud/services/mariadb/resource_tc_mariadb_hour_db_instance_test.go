package mariadb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmariadb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/mariadb"

	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbHourDbInstance_basic -v
func TestAccTencentCloudMariadbHourDbInstance_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMariadbHourDbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbHourDbInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMariadbHourDbInstanceExists("tencentcloud_mariadb_hour_db_instance.basic"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_hour_db_instance.basic", "instance_name", "db-test-2"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_hour_db_instance.basic", "db_version_id", "8.0"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_hour_db_instance.basic", "memory", "2"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_hour_db_instance.basic", "node_count", "2"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_hour_db_instance.basic", "storage", "10"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_hour_db_instance.basic", "subnet_id", tcacctest.DefaultMariadbSubnetId),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_hour_db_instance.basic", "vpc_id", tcacctest.DefaultMariadbVpcId),
					resource.TestCheckResourceAttrSet("tencentcloud_mariadb_hour_db_instance.basic", "zones.#"),
				),
			},
			{
				ResourceName:      "tencentcloud_mariadb_hour_db_instance.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMariadbHourDbInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcmariadb.NewMariadbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_mariadb_hour_db_instance" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}

		instance, err := service.DescribeMariadbDbInstance(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if instance != nil {
			return fmt.Errorf("db hour Instance %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckMariadbHourDbInstanceExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}

		service := svcmariadb.NewMariadbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		instance, err := service.DescribeMariadbDbInstance(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if instance == nil {
			return fmt.Errorf("db hour Instance %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccMariadbHourDbInstanceVar = `
variable "subnet_id" {
  default = "` + tcacctest.DefaultMariadbSubnetId + `"
}

variable "vpc_id" {
  default = "` + tcacctest.DefaultMariadbVpcId + `"
}
`

const testAccMariadbHourDbInstance = testAccMariadbHourDbInstanceVar + `
resource "tencentcloud_mariadb_hour_db_instance" "basic" {
  db_version_id = "8.0"
  instance_name = "db-test-2"
  memory        = 2
  node_count    = 2
  storage       = 10
  subnet_id     = var.subnet_id
  vpc_id        = var.vpc_id
  zones         = ["ap-guangzhou-6","ap-guangzhou-7"]
  tags          = {
	createdBy   = "terraform"
  }
}
`
