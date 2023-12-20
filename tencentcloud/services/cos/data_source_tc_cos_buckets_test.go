package cos_test

import (
	"fmt"
	"regexp"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCosBucketDataSource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucketDataSource_basic(tcacctest.Appid),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.bucket_basic"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.0.cors_rules.#", "0"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.0.lifecycle_rules.#", "0"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.0.website.#", "0"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.0.cos_bucket_url"),
				),
			},
		},
	})
}

func TestAccTencentCloudCosBucketDataSource_tags(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucketDataSource_tags(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.0.tags.fixed_resource", "do_not_remove"),
				),
			},
		},
	})
}

func TestAccTencentCloudCosBucketDataSource_full(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucketDataSource_full(tcacctest.Appid),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.bucket_full"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.0.cors_rules.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.0.cors_rules.0.allowed_headers.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.0.cors_rules.0.allowed_headers.0", "*"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.0.cors_rules.0.allowed_methods.#", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.0.cors_rules.0.allowed_methods.0", "GET"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.0.cors_rules.0.allowed_methods.1", "POST"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.0.cors_rules.0.allowed_origins.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list",
						"bucket_list.0.cors_rules.0.allowed_origins.0", "https://www.test.com"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.0.cors_rules.0.expose_headers.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.0.cors_rules.0.expose_headers.0", "x-cos-test"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.0.cors_rules.0.max_age_seconds", "300"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.0.lifecycle_rules.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list",
						"bucket_list.0.lifecycle_rules.0.filter_prefix", "test/"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list",
						"bucket_list.0.lifecycle_rules.0.expiration.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list",
						"bucket_list.0.lifecycle_rules.0.transition.#", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list",
						"bucket_list.0.lifecycle_rules.0.non_current_expiration.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list",
						"bucket_list.0.lifecycle_rules.0.non_current_transition.#", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.0.website.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.0.website.0.index_document", "index.html"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.0.website.0.error_document", "error.html"),
				),
			},
		},
	})
}

func testAccCosBucketDataSource_basic(appid string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cos_bucket" "bucket_basic" {
  bucket = "tf-baisc-%d-%s"
}

data "tencentcloud_cos_buckets" "bucket_list" {
  bucket_prefix = tencentcloud_cos_bucket.bucket_basic.bucket
}
`, acctest.RandInt(), appid)
}

func testAccCosBucketDataSource_tags() string {
	return fmt.Sprintf(`
%s
data "tencentcloud_cos_buckets" "bucket_list" {
  tags = var.fixed_tags
}
`, tcacctest.FixedTagVariable)
}

func testAccCosBucketDataSource_full(appid string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cos_bucket" "bucket_full" {
  bucket = "tf-full-%d-%s"

  cors_rules {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "POST"]
    allowed_origins = ["https://www.test.com"]
    expose_headers  = ["x-cos-test"]
    max_age_seconds = 300
  }

  lifecycle_rules {
    filter_prefix = "test/"

    expiration {
      days = 365
    }

    transition {
      days          = 30
      storage_class = "STANDARD_IA"
    }

    transition {
      days          = 60
      storage_class = "ARCHIVE"
    }

	non_current_expiration {
      non_current_days = 600
    }

	non_current_transition {
      non_current_days = 90
      storage_class = "STANDARD_IA"
    }

    non_current_transition {
      non_current_days = 180
      storage_class = "ARCHIVE"
    }
  }

  website {
    index_document = "index.html"
    error_document = "error.html"
  }
}

data "tencentcloud_cos_buckets" "bucket_list" {
  bucket_prefix = tencentcloud_cos_bucket.bucket_full.bucket
}
`, acctest.RandInt(), appid)
}
