package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixCvmModifyInstanceDiskTypeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmModifyInstanceDiskType,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_modify_instance_disk_type.modify_instance_disk_type", "id")),
			},
			{
				ResourceName:      "tencentcloud_cvm_modify_instance_disk_type.modify_instance_disk_type",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmModifyInstanceDiskType = `

resource "tencentcloud_cvm_modify_instance_disk_type" "modify_instance_disk_type" {
  instance_id = "ins-r8hr2upy"
  data_disks {
		disk_size = 50
		disk_type = "CLOUD_BASIC"
		disk_id = "disk-hrsd0u81"
		delete_with_instance = true
		snapshot_id = "snap-r9unnd89"
		encrypt = false
		kms_key_id = "kms-abcd1234"
		throughput_performance = 2
		cdc_id = "cdc-b9pbd3px"

  }
  system_disk {
		disk_type = "CLOUD_PREMIUM"
		disk_id = "disk-1drr53sd"
		disk_size = 50
		cdc_id = "cdc-b9pbd3px"

  }
}

`
