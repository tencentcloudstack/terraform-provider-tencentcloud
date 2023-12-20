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
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
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
					resource.TestCheckResourceAttr(ImageInstance, "instance_id", tcacctest.DefaultCvmId),
					resource.TestCheckResourceAttr(ImageInstance, "data_disk_ids.#", "1"),
					resource.TestCheckResourceAttr(ImageInstance, "image_description", "create image with instance"),
				),
			},
			{
				Config: testAccImageWithInstanceUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ImageInstance, "image_name", "image-instance-update-keep"),
					resource.TestCheckResourceAttr(ImageInstance, "instance_id", tcacctest.DefaultCvmId),
					resource.TestCheckResourceAttr(ImageInstance, "data_disk_ids.#", "1"),
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
		resource "tencentcloud_image" "image_instance" {
  			image_name   		= "image-instance-keep"
  			instance_id  		= var.cvm_id
  			data_disk_ids 		= [var.disk_id]
  			image_description 	= "create image with instance"
		}`

	testAccImageWithInstanceUpdate = tcacctest.DefaultCvmImageVariable + `
		resource "tencentcloud_image" "image_instance" {
  			image_name   		= "image-instance-update-keep"
  			instance_id  		= var.cvm_id
  			data_disk_ids 		= [var.disk_id]
  			image_description 	= "update image with instance"
		}`
)
