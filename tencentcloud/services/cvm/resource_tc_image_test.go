package cvm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
)

func TestAccTencentCloudCvmImageResource_UseSnapshotId(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmImageResource_UseSnapshotIdCreate,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmImageExists("tencentcloud_image.image_snap"), resource.TestCheckResourceAttr("tencentcloud_image.image_snap", "image_name", "image-snapshot-keep"), resource.TestCheckResourceAttr("tencentcloud_image.image_snap", "snapshot_ids.#", "1"), resource.TestCheckResourceAttr("tencentcloud_image.image_snap", "force_poweroff", "true"), resource.TestCheckResourceAttr("tencentcloud_image.image_snap", "image_description", "create image with instance")),
			},
			{
				Config: testAccCvmImageResource_UseSnapshotIdChange1,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmImageExists("tencentcloud_image.image_snap"), resource.TestCheckResourceAttr("tencentcloud_image.image_snap", "image_description", "update image with instance"), resource.TestCheckResourceAttr("tencentcloud_image.image_snap", "image_name", "image-snapshot-update-keep"), resource.TestCheckResourceAttr("tencentcloud_image.image_snap", "force_poweroff", "false")),
			},
			{
				ResourceName:            "tencentcloud_image.image_snap",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_poweroff"},
			},
		},
	})
}

const testAccCvmImageResource_UseSnapshotIdCreate = `

resource "tencentcloud_image" "image_snap" {
    image_name = "image-snapshot-keep"
    snapshot_ids = ["snap-gem0ivcj"]
    force_poweroff = true
    image_description = "create image with instance"
}

`
const testAccCvmImageResource_UseSnapshotIdChange1 = `

resource "tencentcloud_image" "image_snap" {
    image_name = "image-snapshot-update-keep"
    snapshot_ids = ["snap-gem0ivcj"]
    force_poweroff = false
    image_description = "update image with instance"
}

`

func TestAccTencentCloudCvmImageResource_UseInstanceId(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmImageResource_UseInstanceIdCreate,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmImageExists("tencentcloud_image.image_instance"), resource.TestCheckResourceAttr("tencentcloud_image.image_instance", "image_name", "image-instance-keep"), resource.TestCheckResourceAttrSet("tencentcloud_image.image_instance", "instance_id"), resource.TestCheckResourceAttr("tencentcloud_image.image_instance", "data_disk_ids.#", "1"), resource.TestCheckResourceAttr("tencentcloud_image.image_instance", "image_description", "create image with instance")),
			},
			{
				Config: testAccCvmImageResource_UseInstanceIdChange1,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmImageExists("tencentcloud_image.image_instance"), resource.TestCheckResourceAttr("tencentcloud_image.image_instance", "image_name", "image-instance-update-keep"), resource.TestCheckResourceAttr("tencentcloud_image.image_instance", "image_description", "update image with instance")),
			},
		},
	})
}

const testAccCvmImageResource_UseInstanceIdCreate = `

data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    
    filter {
        values = ["ap-guangzhou-7"]
        name = "zone"
    }
    filter {
        name = "instance-family"
        values = ["S4","SA2"]
    }
    cpu_core_count = 2
    memory_size = 2
    exclude_sold_out = true
}
resource "tencentcloud_vpc" "vpc" {
    name = "image-vpc"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_subnet" "subnet" {
    vpc_id = tencentcloud_vpc.vpc.id
    name = "image-subnet"
    cidr_block = "10.0.0.0/16"
    availability_zone = "ap-guangzhou-7"
}
resource "tencentcloud_instance" "cvm_image" {
    availability_zone = "ap-guangzhou-7"
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    
    data_disks {
        delete_with_instance = true
        data_disk_type = "CLOUD_PREMIUM"
        data_disk_size = 100
    }
    subnet_id = tencentcloud_subnet.subnet.id
    system_disk_type = "CLOUD_PREMIUM"
    image_id = data.tencentcloud_images.default.images.0.image_id
    vpc_id = tencentcloud_vpc.vpc.id
    project_id = 0
    instance_name = "cvm-image"
}
resource "tencentcloud_image" "image_instance" {
    instance_id = tencentcloud_instance.cvm_image.id
    data_disk_ids = [tencentcloud_instance.cvm_image.data_disks.0.data_disk_id]
    image_description = "create image with instance"
    image_name = "image-instance-keep"
}

`
const testAccCvmImageResource_UseInstanceIdChange1 = `

data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    
    filter {
        values = ["ap-guangzhou-7"]
        name = "zone"
    }
    filter {
        values = ["S4","SA2"]
        name = "instance-family"
    }
    cpu_core_count = 2
    memory_size = 2
    exclude_sold_out = true
}
resource "tencentcloud_vpc" "vpc" {
    name = "image-vpc"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_subnet" "subnet" {
    vpc_id = tencentcloud_vpc.vpc.id
    name = "image-subnet"
    cidr_block = "10.0.0.0/16"
    availability_zone = "ap-guangzhou-7"
}
resource "tencentcloud_instance" "cvm_image" {
    availability_zone = "ap-guangzhou-7"
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    
    data_disks {
        delete_with_instance = true
        data_disk_type = "CLOUD_PREMIUM"
        data_disk_size = 100
    }
    subnet_id = tencentcloud_subnet.subnet.id
    system_disk_type = "CLOUD_PREMIUM"
    image_id = data.tencentcloud_images.default.images.0.image_id
    vpc_id = tencentcloud_vpc.vpc.id
    project_id = 0
    instance_name = "cvm-image"
}
resource "tencentcloud_image" "image_instance" {
    instance_id = tencentcloud_instance.cvm_image.id
    data_disk_ids = [tencentcloud_instance.cvm_image.data_disks.0.data_disk_id]
    image_description = "update image with instance"
    image_name = "image-instance-update-keep"
}

`

func testAccCheckCvmImageDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cvmService := svccvm.NewCvmService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_image" {
			continue
		}

		_, has, err := cvmService.DescribeImageById(ctx, rs.Primary.ID, true)
		if err != nil || has {
			err = resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
				_, has, err = cvmService.DescribeImageById(ctx, rs.Primary.ID, true)
				if nil != err {
					return tccommon.RetryError(err)
				}
				if has {
					return resource.RetryableError(fmt.Errorf("image still exists: %s", rs.Primary.ID))
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if has {
			return fmt.Errorf("image still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckCvmImageExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("image %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("image id is not set")
		}
		cvmService := svccvm.NewCvmService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		_, has, err := cvmService.DescribeImageById(ctx, rs.Primary.ID, false)
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("image doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}
