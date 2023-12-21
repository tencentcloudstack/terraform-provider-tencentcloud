package postgresql_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	svcclb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/clb"
)

var testDataPostgresqlInstancesName = "data.tencentcloud_postgresql_instances.id_test"

func TestAccTencentCloudPostgresqlInstancesDataSource(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccStepSetRegion(t, "ap-guangzhou")
			tcacctest.AccPreCheck(t)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckLBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataPostgresqlInstanceBasic,
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-guangzhou")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
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

	clbService := svcclb.NewClbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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

const testAccTencentCloudDataPostgresqlInstanceBasic = tcacctest.CommonPresetPGSQL + `

data "tencentcloud_postgresql_instances" "id_test"{
  id = local.pgsql_id
}
`
