package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresModifyDBInstanceReadOnlyGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresModifyDBInstanceReadOnlyGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgres_modify_d_b_instance_read_only_group.modify_d_b_instance_read_only_group", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgres_modify_d_b_instance_read_only_group.modify_d_b_instance_read_only_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresModifyDBInstanceReadOnlyGroup = `

resource "tencentcloud_postgres_modify_d_b_instance_read_only_group" "modify_d_b_instance_read_only_group" {
  d_b_instance_id = "postgres-6r233v55"
  read_only_group_id = "pgrogrp-test1"
  new_read_only_group_id = "pgrogrp-test2"
  tags = {
    "createdBy" = "terraform"
  }
}

`
