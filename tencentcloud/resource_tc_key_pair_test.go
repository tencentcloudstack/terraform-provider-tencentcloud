package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudKeyPair(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "tencentcloud_key_pair.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKeyPairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyPairPublicKey,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyPairExists("tencentcloud_key_pair.foo"),
					resource.TestCheckResourceAttr("tencentcloud_key_pair.foo", "key_name", "from_terraform"),
				),
			},
			{
				ResourceName:      "tencentcloud_key_pair.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckKeyPairExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("key pair %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("key pair id is not set")
		}
		cvmService := CvmService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		keyPair, err := cvmService.DescribeKeyPairById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				keyPair, err = cvmService.DescribeKeyPairById(ctx, rs.Primary.ID)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if keyPair == nil {
			return fmt.Errorf("key pair is not found")
		}
		return nil
	}
}

func testAccCheckKeyPairDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cvmService := CvmService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_key_pair" {
			continue
		}

		keyPair, err := cvmService.DescribeKeyPairById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				keyPair, err = cvmService.DescribeKeyPairById(ctx, rs.Primary.ID)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if keyPair != nil {
			return fmt.Errorf("key pair still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

const testAccKeyPairPublicKey = `
resource "tencentcloud_key_pair" "foo" {
  key_name   = "from_terraform"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDjd8fTnp7Dcuj4mLaQxf9Zs/ORgUL9fQxRCNKkPgP1paTy1I513maMX126i36Lxxl3+FUB52oVbo/FgwlIfX8hyCnv8MCxqnuSDozf1CD0/wRYHcTWAtgHQHBPCC2nJtod6cVC3kB18KeV4U7zsxmwFeBIxojMOOmcOBuh7+trRw=="
}
`
