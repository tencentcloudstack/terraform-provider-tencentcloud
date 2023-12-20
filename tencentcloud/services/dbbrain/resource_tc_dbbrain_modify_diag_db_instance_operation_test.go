package dbbrain_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDbbrainModifyDiagDbInstanceOperationResource_basic_off(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDbbrainModifyDiagDbInstanceConf_off, tcacctest.DefaultDbBrainInstanceId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_modify_diag_db_instance_operation.off", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_modify_diag_db_instance_operation.off", "instance_ids.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_modify_diag_db_instance_operation.off", "instance_confs.#"),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_modify_diag_db_instance_operation.off", "instance_confs.0.daily_inspection", "No"),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_modify_diag_db_instance_operation.off", "instance_confs.0.overview_display", "No"),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_modify_diag_db_instance_operation.off", "product", "mysql"),
				),
			},
		},
	})
}

func TestAccTencentCloudDbbrainModifyDiagDbInstanceOperationResource_basic_on(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDbbrainModifyDiagDbInstanceConf_on, tcacctest.DefaultDbBrainInstanceId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_modify_diag_db_instance_operation.on", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_modify_diag_db_instance_operation.on", "instance_ids.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_modify_diag_db_instance_operation.on", "instance_confs.#"),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_modify_diag_db_instance_operation.on", "instance_confs.0.daily_inspection", "Yes"),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_modify_diag_db_instance_operation.on", "instance_confs.0.overview_display", "Yes"),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_modify_diag_db_instance_operation.on", "product", "mysql"),
				),
			},
		},
	})
}

const testAccDbbrainModifyDiagDbInstanceConf_off = `

resource "tencentcloud_dbbrain_modify_diag_db_instance_operation" "off" {
  instance_confs {
	daily_inspection = "No"
	overview_display = "No"
  }
  product = "mysql"
  instance_ids = ["%s"]
}

`

const testAccDbbrainModifyDiagDbInstanceConf_on = `

resource "tencentcloud_dbbrain_modify_diag_db_instance_operation" "on" {
  instance_confs {
	daily_inspection = "Yes"
	overview_display = "Yes"
  }
  product = "mysql"
  instance_ids = ["%s"]
}

`
