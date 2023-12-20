package cvm_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTestingCvmInstanceResource_Basic(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.cvm_basic"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { tcacctest.AccPreCheck(t) },
		IDRefreshName: id,
		Providers:     tcacctest.AccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { tcacctest.AccStepPreConfigSetTempAKSK(t, tcacctest.ACCOUNT_TYPE_COMMON) },
				Config:    testAccTencentCloudTestingCvmInstanceBasic,
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
				PreConfig: func() { tcacctest.AccStepPreConfigSetTempAKSK(t, tcacctest.ACCOUNT_TYPE_COMMON) },
				Config:    testAccTencentCloudTestingCvmInstanceModifyInstanceType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "instance_type"),
				),
			},
		},
	})
}

func TestAccTencentCloudTestingCvmInstanceResource_WithDataDisk(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { tcacctest.AccPreCheck(t) },
		IDRefreshName: id,
		Providers:     tcacctest.AccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { tcacctest.AccStepPreConfigSetTempAKSK(t, tcacctest.ACCOUNT_TYPE_COMMON) },
				Config:    testAccTencentCloudTestingInstanceWithDataDisk,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "system_disk_size", "100"),
					resource.TestCheckResourceAttr(id, "system_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_size", "100"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_snapshot_id", ""),
				),
			},
			{
				PreConfig: func() { tcacctest.AccStepPreConfigSetTempAKSK(t, tcacctest.ACCOUNT_TYPE_COMMON) },
				Config:    testAccTencentCloudTestingInstanceWithDataDiskUpdate,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "system_disk_size", "100"),
					resource.TestCheckResourceAttr(id, "system_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_size", "150"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_snapshot_id", ""),
					resource.TestCheckResourceAttr(id, "data_disks.1.data_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.1.data_disk_size", "150"),
				),
			},
		},
	})
}

const testAccTencentCloudTestingCvmInstanceBasic = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_instance" "cvm_basic" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_testing_zone
  image_id          = var.cvm_testing_image_id
  instance_type     = "S2.MEDIUM2"
  vpc_id            = var.cvm_testing_vpc_id
  subnet_id         = var.cvm_testing_subnet_id
  system_disk_type  = "CLOUD_PREMIUM"
  project_id        = 0
}
`
const testAccTencentCloudTestingCvmInstanceModifyInstanceType = tcacctest.DefaultInstanceVariable + `

resource "tencentcloud_instance" "cvm_basic" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_testing_zone
  image_id          = var.cvm_testing_image_id
  instance_type     = "S2.MEDIUM2"
  vpc_id            = var.cvm_testing_vpc_id
  subnet_id         = var.cvm_testing_subnet_id
  system_disk_type  = "CLOUD_PREMIUM"
  project_id        = 0
}
`

const testAccTencentCloudTestingInstanceWithDataDisk = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_testing_zone
  image_id          = var.cvm_testing_image_id
  instance_type     = "S2.MEDIUM2"

  system_disk_type = "CLOUD_PREMIUM"
  system_disk_size = 100

  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 100
    delete_with_instance  = true
  } 
   
  disable_security_service = true
  disable_monitor_service  = true
}
`

const testAccTencentCloudTestingInstanceWithDataDiskUpdate = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_testing_zone
  image_id          = var.cvm_testing_image_id
  instance_type     = "S2.MEDIUM2"

  system_disk_type = "CLOUD_PREMIUM"
  system_disk_size = 100

  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  } 
   
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }

  disable_security_service = true
  disable_monitor_service  = true
}
`
