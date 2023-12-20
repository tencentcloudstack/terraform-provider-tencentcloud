package cynosdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCynosdbReadOnlyInstanceExclusiveAccessResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbReadOnlyInstanceExclusiveAccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_read_only_instance_exclusive_access.read_only_instance_exclusive_access", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_read_only_instance_exclusive_access.read_only_instance_exclusive_access", "cluster_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_read_only_instance_exclusive_access.read_only_instance_exclusive_access", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_read_only_instance_exclusive_access.read_only_instance_exclusive_access", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_read_only_instance_exclusive_access.read_only_instance_exclusive_access", "subnet_id"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_read_only_instance_exclusive_access.read_only_instance_exclusive_access", "port", "1234"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_read_only_instance_exclusive_access.read_only_instance_exclusive_access", "security_group_ids.#"),
				),
			},
		},
	})
}

const testAccCynosdbReadOnlyInstanceExclusiveAccess = tcacctest.CommonCynosdb + tcacctest.DefaultVpcSubnets + `

resource "tencentcloud_cynosdb_read_only_instance_exclusive_access" "read_only_instance_exclusive_access" {
  cluster_id = var.cynosdb_cluster_id
  instance_id = var.cynosdb_cluster_instance_id
  vpc_id = local.vpc_id
  subnet_id = local.subnet_id
  port = 1234
  security_group_ids = [var.cynosdb_cluster_security_group_id]
}

`
