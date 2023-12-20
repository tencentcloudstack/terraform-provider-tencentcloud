package cos_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localcos "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cos"

	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudCosBucketPolicy(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCosBucketPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucketPolicyBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosBucketPolicyExists("tencentcloud_cos_bucket_policy.foo"),
					resource.TestCheckResourceAttrSet("tencentcloud_cos_bucket_policy.foo", "bucket"),
					resource.TestCheckResourceAttrSet("tencentcloud_cos_bucket_policy.foo", "policy"),
				),
			}, {
				Config: testAccCosBucketPolicyUpdate(),
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
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cosService := localcos.NewCosService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cos_bucket_policy" {
			continue
		}

		policy, err := cosService.DescribePolicyByBucket(ctx, rs.Primary.ID)
		if err == nil && policy != "" {
			return fmt.Errorf("[CHECK][cos bucket policy][Desctroy] check: cos bucket policy still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCosBucketPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][cos bucket policy][Exists] check: cos bucket policy %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][cos bucket policy][Exists] check: cos bucket policy id is not set")
		}
		cosService := localcos.NewCosService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		policy, err := cosService.DescribePolicyByBucket(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if policy == "" {
			return fmt.Errorf("[CHECK][cos bucket policy][Exists] check: cos bucket policy %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCosBucketPolicyBasic() string {
	return tcacctest.UserInfoData + `
resource "tencentcloud_cos_bucket" "bucket" {
  bucket = "test-tf-policy-${local.app_id}"
  acl    = "private"
}

resource "tencentcloud_cos_bucket_policy" "foo" {
  bucket        = tencentcloud_cos_bucket.bucket.bucket
  policy        = <<EOF
{
  "Statement": [
    {
      "Principal": {
        "service": [
          "cvm.cloud.tencent.com"
        ]
      },
      "Effect": "Allow",
      "Action": [
        "name/cos:HeadBucket",
        "name/cos:ListMultipartUploads",
        "name/cos:ListParts",
        "name/cos:GetObject",
        "name/cos:HeadObject",
        "name/cos:OptionsObject"
      ],
      "Resource": [
        "qcs::cos:ap-guangzhou:uid/${local.app_id}:test-tf-policy-${local.app_id}/*"
      ]
    }
  ],
  "version": "2.0"
}
EOF
}
`
}
func testAccCosBucketPolicyUpdate() string {
	return tcacctest.UserInfoData + `
resource "tencentcloud_cos_bucket" "bucket" {
  bucket = "test-tf-policy-${local.app_id}"
  acl    = "private"
}

resource "tencentcloud_cos_bucket_policy" "foo" {
  bucket        = tencentcloud_cos_bucket.bucket.bucket
  policy        = <<EOF
{
  "Statement": [
    {
      "Principal": {
        "service": [
          "cvm.cloud.tencent.com"
        ]
      },
      "Effect": "Deny",
      "Action": [
        "name/cos:HeadBucket",
        "name/cos:ListMultipartUploads",
        "name/cos:ListParts",
        "name/cos:GetObject",
        "name/cos:HeadObject",
        "name/cos:OptionsObject"
      ],
      "Resource": [
        "qcs::cos:ap-guangzhou:uid/${local.app_id}:test-tf-policy-${local.app_id}/*"
      ]
    }
  ],
  "version": "2.0"
}
EOF
}
`
}
