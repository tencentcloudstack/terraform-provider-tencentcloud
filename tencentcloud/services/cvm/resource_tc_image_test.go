package cvm_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"

	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	ImageSnap     = "tencentcloud_image.image_snap"
	ImageInstance = "tencentcloud_image.image_instance"
)

func TestAccTencentCloudImageResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckImageDestroy,
		Steps: []resource.TestStep{
			// use snapshot id
			{
				Config: testAccImageWithSnapShot,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImageExists(ImageSnap),
					resource.TestCheckResourceAttr(ImageSnap, "image_name", "image-snapshot-keep"),
					resource.TestCheckResourceAttr(ImageSnap, "snapshot_ids.#", "1"),
					resource.TestCheckResourceAttr(ImageSnap, "force_poweroff", "true"),
					resource.TestCheckResourceAttr(ImageSnap, "image_description", "create image with snapshot"),
				),
			},
			{
				Config: testAccImageWithSnapShotUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ImageSnap, "image_name", "image-snapshot-update-keep"),
					resource.TestCheckResourceAttr(ImageSnap, "snapshot_ids.#", "1"),
					resource.TestCheckResourceAttr(ImageSnap, "force_poweroff", "false"),
					resource.TestCheckResourceAttr(ImageSnap, "image_description", "update image with snapshot"),
				),
			},
			{
				ResourceName:            ImageSnap,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_poweroff"},
			},
			// use instance id
			{
				Config: testAccImageWithInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImageExists(ImageInstance),
					resource.TestCheckResourceAttr(ImageInstance, "image_name", "image-instance-keep"),
					resource.TestCheckResourceAttrSet(ImageInstance, "instance_id"),
					resource.TestCheckResourceAttr(ImageInstance, "data_disk_ids.#", "1"),
					resource.TestCheckResourceAttr(ImageInstance, "image_description", "create image with instance"),
				),
			},
			{
				Config: testAccImageWithInstanceUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ImageInstance, "image_name", "image-instance-update-keep"),
					resource.TestCheckResourceAttr(ImageInstance, "image_description", "update image with instance"),
				),
			},
		},
	})
}

func testAccCheckImageDestroy(s *terraform.State) error {
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

func testAccCheckImageExists(n string) resource.TestCheckFunc {
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

const (
	testAccImageWithSnapShot = tcacctest.DefaultCvmImageVariable + `
		resource "tencentcloud_image" "image_snap" {
  			image_name   		= "image-snapshot-keep"
  			snapshot_ids 		= [var.snap_id]
			force_poweroff 		= true
			image_description 	= "create image with snapshot"
		}`

	testAccImageWithSnapShotUpdate = tcacctest.DefaultCvmImageVariable + `
		resource "tencentcloud_image" "image_snap" {
  			image_name   		= "image-snapshot-update-keep"
  			snapshot_ids 		= [var.snap_id]
  			force_poweroff   	= false
  			image_description 	= "update image with snapshot"
		}`

	testAccImageWithInstance = tcacctest.DefaultCvmImageVariable + `

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
				name = "instance-family"
				values = ["S4","SA2"]
			}
			filter {
				values = ["ap-guangzhou-7"]
				name = "zone"
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
			subnet_id = tencentcloud_subnet.subnet.id
			system_disk_type = "CLOUD_PREMIUM"
			project_id = 0
			instance_name = "cvm-image"
			availability_zone = "ap-guangzhou-7"
			image_id = data.tencentcloud_images.default.images.0.image_id
			instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
			vpc_id = tencentcloud_vpc.vpc.id
			data_disks {
				delete_with_instance = true
				data_disk_type = "CLOUD_PREMIUM"
				data_disk_size = 100
			}
		}

		resource "tencentcloud_image" "image_instance" {
  			image_name   		= "image-instance-keep"
  			instance_id  		= tencentcloud_instance.cvm_image.id
  			data_disk_ids 		= [tencentcloud_instance.cvm_image.data_disks.0.data_disk_id]
  			image_description 	= "create image with instance"
		}`

	testAccImageWithInstanceUpdate = tcacctest.DefaultCvmImageVariable + `
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
				name = "instance-family"
				values = ["S4","SA2"]
			}
			filter {
				values = ["ap-guangzhou-7"]
				name = "zone"
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
			subnet_id = tencentcloud_subnet.subnet.id
			system_disk_type = "CLOUD_PREMIUM"
			project_id = 0
			instance_name = "cvm-image"
			availability_zone = "ap-guangzhou-7"
			image_id = data.tencentcloud_images.default.images.0.image_id
			instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
			vpc_id = tencentcloud_vpc.vpc.id
			data_disks {
				delete_with_instance = true
				data_disk_type = "CLOUD_PREMIUM"
				data_disk_size = 100
			}
		}

		resource "tencentcloud_image" "image_instance" {
  			image_name   		= "image-instance-update-keep"
  			instance_id  		= tencentcloud_instance.cvm_image.id
  			data_disk_ids 		= [tencentcloud_instance.cvm_image.data_disks.0.data_disk_id]
  			image_description 	= "update image with instance"
		}`
)
