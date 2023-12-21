package mariadb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmariadb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/mariadb"

	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbInstanceResource_basic -v
func TestAccTencentCloudMariadbInstanceResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMariadbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMariadbInstanceExists("tencentcloud_mariadb_instance.instance"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_instance.instance", "instance_name", "terraform-test"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_instance.instance", "db_version_id", "8.0"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_instance.instance", "node_count", "2"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_instance.instance", "memory", "2"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_instance.instance", "storage", "10"),
					// resource.TestCheckResourceAttr("tencentcloud_mariadb_instance.instance", "period", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_instance.instance", "auto_renew_flag", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_instance.instance", "ipv6_flag", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_instance.instance", "tags.createby", "terrafrom"),
				),
			},
			{
				ResourceName:            "tencentcloud_mariadb_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "auto_voucher", "voucher_ids", "init_params", "dcn_region", "dcn_instance_id"},
			},
		},
	})
}

func testAccCheckMariadbInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcmariadb.NewMariadbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_mariadb_instance" {
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
			return fmt.Errorf("Instance %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckMariadbInstanceExists(r string) resource.TestCheckFunc {
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
			return fmt.Errorf("Instance %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccMariadbInstanceVar = `
variable "subnet_id" {
  default = "` + tcacctest.DefaultMariadbInstanceSubnetId + `"
}

variable "vpc_id" {
  default = "` + tcacctest.DefaultMariadbInstanceVpcId + `"
}
`

const testAccMariadbInstance = testAccMariadbInstanceVar + `
resource "tencentcloud_mariadb_instance" "instance" {
	zones = ["ap-guangzhou-6", "ap-guangzhou-7"]
	node_count = 2
	memory = 2
	storage = 10
	period = 1
	# auto_voucher =
	# voucher_ids =
	vpc_id = var.vpc_id
	subnet_id = var.subnet_id
	# project_id = ""
	db_version_id = "8.0"
	instance_name = "terraform-test"
	# security_group_ids = ""
	auto_renew_flag = 1
	ipv6_flag = 0
	tags = {
	  "createby" = "terrafrom"
	}
	init_params {
	  param = "character_set_server"
	  value = "utf8mb4"
	}
	init_params {
	  param = "lower_case_table_names"
	  value = "0"
	}
	init_params {
	  param = "innodb_page_size"
	  value = "16384"
	}
	init_params {
	  param = "sync_mode"
	  value = "1"
	}
	dcn_region = ""
	dcn_instance_id = ""
}
`
