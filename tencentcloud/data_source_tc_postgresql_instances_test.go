package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

var testDataPostgresqlInstancesName = "data.tencentcloud_postgresql_instances.id_test"

func TestAccTencentCloudPostgresqlInstancesDataSource(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccStepSetRegion(t, "ap-guangzhou")
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataPostgresqlInstanceBasic,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-guangzhou")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testDataPostgresqlInstancesName, "instance_list.#", "1"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.id"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.create_time"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.charge_type"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.engine_version"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.project_id"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.memory"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.storage"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.private_access_ip"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.private_access_port"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.public_access_switch"),
					resource.TestCheckResourceAttrSet(testDataPostgresqlInstancesName, "instance_list.0.charset"),
				),
			},
		},
	})
}

func testAccCheckLBDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	clbService := ClbService{
		client: testAccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_lb" {
			continue
		}

		_, err := clbService.DescribeLoadBalancerById(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("clb instance still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

const testAccTencentCloudDataPostgresqlInstanceBasic = CommonPresetPGSQL + `

data "tencentcloud_postgresql_instances" "id_test"{
  id = local.pgsql_id
}
`
