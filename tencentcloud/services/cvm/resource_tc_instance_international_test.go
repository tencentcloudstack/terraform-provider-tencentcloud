package cvm_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudInstanceResource_basic(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.cvm_basic"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { tcacctest.AccPreCheck(t) },
		IDRefreshName: id,
		Providers:     tcacctest.AccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInternationalTencentCloudInstanceBasic,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "private_ip"),
					resource.TestCheckResourceAttrSet(id, "vpc_id"),
					resource.TestCheckResourceAttrSet(id, "subnet_id"),
					resource.TestCheckResourceAttrSet(id, "project_id"),
				),
			},
			{
				ResourceName:            id,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"disable_monitor_service", "disable_security_service", "hostname", "password", "force_delete"},
			},
		},
	})
}

const testAccInternationalTencentCloudInstanceBasic = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_instance" "cvm_basic" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_international_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  vpc_id            = var.cvm_international_vpc_id
  subnet_id         = var.cvm_international_subnet_id
  system_disk_type  = "CLOUD_PREMIUM"
  project_id        = 0
}
`
