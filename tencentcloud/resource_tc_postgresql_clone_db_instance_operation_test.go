package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccPostgresqlCloneDbInstanceOperationObject = "tencentcloud_postgresql_clone_db_instance_operation.clone_db_instance_operation"

func TestAccTencentCloudPostgresqlCloneDbInstanceOperationResource_clone(t *testing.T) {
	t.Parallel()
	recoverTime := time.Now().AddDate(0, 0, -2).Format("2006-01-02 15:04:05")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccPostgresqlCloneDbInstanceOperation, recoverTime),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAccPostgresqlCloneDbInstanceOperationObject, "id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlCloneDbInstanceOperationObject, "db_instance_id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlCloneDbInstanceOperationObject, "spec_code"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlCloneDbInstanceOperationObject, "storage"),
					resource.TestCheckResourceAttr(testAccPostgresqlCloneDbInstanceOperationObject, "period", "1"),
					resource.TestCheckResourceAttr(testAccPostgresqlCloneDbInstanceOperationObject, "auto_renew_flag", "0"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlCloneDbInstanceOperationObject, "vpc_id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlCloneDbInstanceOperationObject, "subnet_id"),
					resource.TestCheckResourceAttr(testAccPostgresqlCloneDbInstanceOperationObject, "name", "tf_test_pg_ins_clone"),
					resource.TestCheckResourceAttr(testAccPostgresqlCloneDbInstanceOperationObject, "instance_charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlCloneDbInstanceOperationObject, "security_group_ids.#"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlCloneDbInstanceOperationObject, "db_node_set.#"),
					resource.TestCheckResourceAttr(testAccPostgresqlCloneDbInstanceOperationObject, "db_node_set.0.role", "Primary"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlCloneDbInstanceOperationObject, "db_node_set.0.zone"),
					resource.TestCheckResourceAttr(testAccPostgresqlCloneDbInstanceOperationObject, "db_node_set.1.role", "Standby"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlCloneDbInstanceOperationObject, "db_node_set.1.zone"),
					resource.TestCheckResourceAttr(testAccPostgresqlCloneDbInstanceOperationObject, "project_id", "0"),
					resource.TestCheckResourceAttr(testAccPostgresqlCloneDbInstanceOperationObject, "tag_list.0.tag_key", "issued_by"),
					resource.TestCheckResourceAttr(testAccPostgresqlCloneDbInstanceOperationObject, "tag_list.0.tag_value", "terraform_test"),
					resource.TestCheckResourceAttr(testAccPostgresqlCloneDbInstanceOperationObject, "auto_voucher", "0"),
					resource.TestCheckResourceAttr(testAccPostgresqlCloneDbInstanceOperationObject, "recovery_target_time", recoverTime),
				),
			},
		},
	})
}

const testAccPostgresqlCloneDbInstanceOperation = defaultSecurityGroupData + CommonPresetPGSQL + defaultVpcSubnets + `
locals {
	pg_az = "ap-guangzhou-6"
}
data "tencentcloud_postgresql_specinfos" "foo" {
	availability_zone = local.pg_az
  }

resource "tencentcloud_postgresql_clone_db_instance_operation" "clone_db_instance_operation" {
  db_instance_id = local.pgsql_id
  spec_code = data.tencentcloud_postgresql_specinfos.foo.list.0.id
  storage = data.tencentcloud_postgresql_specinfos.foo.list.0.storage_min
  period = 1
  auto_renew_flag = 0
  vpc_id = local.vpc_id
  subnet_id = local.subnet_id
  name = "tf_test_pg_ins_clone"
  instance_charge_type = "POSTPAID_BY_HOUR"
  security_group_ids = [local.sg_id]
  project_id = 0
  db_node_set {
	role = "Primary"
	zone = local.pg_az
  }
  db_node_set {
	role = "Standby"
	zone = local.pg_az
  }
  tag_list {
		tag_key = "issued_by"
		tag_value = "terraform_test"
  }

  auto_voucher = 0
  recovery_target_time = "%s"
}

`
