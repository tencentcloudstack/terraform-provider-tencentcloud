package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	ImageSnap     = "tencentcloud_image.image_snap"
	ImageInstance = "tencentcloud_image.image_instance"
)

func TestAccTencentCloudImageResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckImageDestroy,
		Steps: []resource.TestStep{
			// use snapshot id
			{
				Config: testAccImageWithSnapShot,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImageExists(ImageSnap),
					resource.TestCheckResourceAttr(ImageSnap, "image_name", "image-snapshot-keep"),
					resource.TestCheckResourceAttr(ImageSnap, "snapshot_ids.#", "2"),
					resource.TestCheckResourceAttr(ImageSnap, "force_poweroff", "true"),
					resource.TestCheckResourceAttr(ImageSnap, "image_description", "create image with snapshot"),
				),
			},
			{
				Config: testAccImageWithSnapShotUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ImageSnap, "image_name", "image-snapshot-update-keep"),
					resource.TestCheckResourceAttr(ImageSnap, "snapshot_ids.#", "2"),
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
					resource.TestCheckResourceAttr(ImageInstance, "instance_id", "ins-2ju245xg"),
					resource.TestCheckResourceAttr(ImageInstance, "data_disk_ids.#", "1"),
					resource.TestCheckResourceAttr(ImageInstance, "image_description", "create image with instance"),
				),
			},
			{
				Config: testAccImageWithInstanceUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ImageInstance, "image_name", "image-instance-update-keep"),
					resource.TestCheckResourceAttr(ImageInstance, "instance_id", "ins-2ju245xg"),
					resource.TestCheckResourceAttr(ImageInstance, "data_disk_ids.#", "1"),
					resource.TestCheckResourceAttr(ImageInstance, "image_description", "update image with instance"),
				),
			},
		},
	})
}

func testAccCheckImageDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	cvmService := CvmService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_image" {
			continue
		}

		_, has, err := cvmService.DescribeImageById(ctx, rs.Primary.ID, true)
		if err != nil || has {
			err = resource.Retry(3*readRetryTimeout, func() *resource.RetryError {
				_, has, err = cvmService.DescribeImageById(ctx, rs.Primary.ID, true)
				if nil != err {
					return retryError(err)
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
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("image %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("image id is not set")
		}
		cvmService := CvmService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
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
	testAccImageWithSnapShot = `
		resource "tencentcloud_image" "image_snap" {
  			image_name   		= "image-snapshot-keep"
  			snapshot_ids 		= ["snap-nbp3xy1d", "snap-nvzu3dmh"]
			force_poweroff 		= true
			image_description 	= "create image with snapshot"
		}`

	testAccImageWithSnapShotUpdate = `
		resource "tencentcloud_image" "image_snap" {
  			image_name   		= "image-snapshot-update-keep"
  			snapshot_ids 		= ["snap-nbp3xy1d", "snap-nvzu3dmh"]
  			force_poweroff   	= false
  			image_description 	= "update image with snapshot"
		}`

	testAccImageWithInstance = `
		resource "tencentcloud_image" "image_instance" {
  			image_name   		= "image-instance-keep"
  			instance_id  		= "ins-2ju245xg"
  			data_disk_ids 		= ["disk-gii0vtwi"]
  			image_description 	= "create image with instance"
		}`

	testAccImageWithInstanceUpdate = `
		resource "tencentcloud_image" "image_instance" {
  			image_name   		= "image-instance-update-keep"
  			instance_id  		= "ins-2ju245xg"
  			data_disk_ids 		= ["disk-gii0vtwi"]
  			image_description 	= "update image with instance"
		}`
)
