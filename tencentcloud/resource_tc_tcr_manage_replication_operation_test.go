package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTcrManageReplicationOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTcrManageReplicationOperation, "sync", time.Now().Nanosecond()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_manage_replication_operation.my_replica", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_manage_replication_operation.my_replica", "source_registry_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_manage_replication_operation.my_replica", "destination_registry_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_manage_replication_operation.my_replica", "rule.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_manage_replication_operation.my_replica", "rule.0.name"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_manage_replication_operation.my_replica", "rule.0.override", "true"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_manage_replication_operation.my_replica", "rule.0.filters.#", "3"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_manage_replication_operation.my_replica", "rule.0.filters.0.type", "name"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_manage_replication_operation.my_replica", "rule.0.filters.1.type", "tag"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_manage_replication_operation.my_replica", "rule.0.filters.2.type", "resource"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_manage_replication_operation.my_replica", "description", "this is the tcr sync operation"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_manage_replication_operation.my_replica", "destination_region_id", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_manage_replication_operation.my_replica", "peer_replication_option.#"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_manage_replication_operation.my_replica", "peer_replication_option.0.enable_peer_replication", "false"),
				),
			},
		},
	})
}

const testAccTcrManageReplicationOperation = defaultTCRInstanceData + `

resource "tencentcloud_tcr_instance" "mytcr_dest" {
	name        = "tf-test-tcr-%s"
	instance_type = "premium"
	delete_bucket = true
  }

resource "tencentcloud_tcr_namespace" "myns_dest" {
	instance_id 	 = tencentcloud_tcr_instance.mytcr_dest.id
	name			 = "tf_test_ns_dest"
	is_public		 = true
	is_auto_scan	 = true
	is_prevent_vul   = true
	severity		 = "medium"
	cve_whitelist_items	{
		cve_id = "cve-xxxxx"
	}
}

resource "tencentcloud_tcr_manage_replication_operation" "my_replica" {
  source_registry_id = local.tcr_id
  destination_registry_id = tencentcloud_tcr_instance.mytcr_dest.id
  rule {
		name = "test_sync_%d"
		dest_namespace = tencentcloud_tcr_namespace.myns_dest.name
		override = true
		filters {
			type = "name"
			value = join("/", [var.tcr_namespace, "**"])
		}
		filters {
			type = "tag"
			value = ""
		}
		filters {
			type = "resource"
			value = ""
		}
  }
  description = "this is the tcr sync operation"
  destination_region_id = 1 // "ap-guangzhou"
  peer_replication_option {
		peer_registry_uin = ""
		peer_registry_token = ""
		enable_peer_replication = false
  }
}

`
