package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudCosBucketPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosBucketPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucketPolicyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosBucketPolicyExists("tencentcloud_cos_bucket_policy.foo"),
					resource.TestCheckResourceAttrSet("tencentcloud_cos_bucket_policy.foo", "bucket"),
					resource.TestCheckResourceAttrSet("tencentcloud_cos_bucket_policy.foo", "policy"),
				),
			}, {
				Config: testAccCosBucketPolicyUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosBucketPolicyExists("tencentcloud_cos_bucket_policy.foo"),
					resource.TestCheckResourceAttrSet("tencentcloud_cos_bucket_policy.foo", "bucket"),
					resource.TestCheckResourceAttrSet("tencentcloud_cos_bucket_policy.foo", "policy"),
				),
			},
			{
				ResourceName:      "tencentcloud_cos_bucket_policy.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCosBucketPolicyDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	cosService := CosService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cos_bucket_policy" {
			continue
		}

		policy, err := cosService.DescribePolicyByBucket(ctx, rs.Primary.ID)
		if err == nil && policy != "" {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][cos bucket policy][Desctroy] check: cos bucket policy still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCosBucketPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][cos bucket policy][Exists] check: cos bucket policy %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][cos bucket policy][Exists] check: cos bucket policy id is not set")
		}
		cosService := CosService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		policy, err := cosService.DescribePolicyByBucket(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if policy == "" {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][cos bucket policy][Exists] check: cos bucket policy %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccCosBucketPolicyBasic = `
resource "tencentcloud_cos_bucket_policy" "foo" {
  bucket        = "bucket-for-terraform-state-1259649581"
  policy        = "{\"version\":\"2.0\",\"Statement\":[{\"Action\":[\"name/cos:DeleteBucket\"],\"Effect\":\"allow\",\"Resource\":[\"qcs::cos:ap-guangzhou:uid/1259649581:bucket-for-terraform-state-1259649581/*\"],\"Principal\":{\"qcs\":[\"qcs::cam::uin/100010835595:uin/100014918835\"]}}]}"
}
`

const testAccCosBucketPolicyUpdate = `
resource "tencentcloud_cos_bucket_policy" "foo" {
  bucket        = "bucket-for-terraform-state-1259649581"
  policy        = "{\"version\":\"2.0\",\"Statement\":[{\"Action\":[\"name/cos:PutBucketACL\"],\"Effect\":\"allow\",\"Resource\":[\"qcs::cos:ap-guangzhou:uid/1259649581:bucket-for-terraform-state-1259649581/*\"],\"Principal\":{\"qcs\":[\"qcs::cam::uin/100010835595:uin/100014918835\"]}}]}"
}
`
