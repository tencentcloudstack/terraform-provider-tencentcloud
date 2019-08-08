package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudCosBucketDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucketDataSource_basic(appid),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.bucket_basic"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.0.cors_rules.#", "0"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.0.lifecycle_rules.#", "0"),
					resource.TestCheckResourceAttr("data.tencentcloud_cos_buckets.bucket_list", "bucket_list.0.website.#", "0"),
				),
			},
		},
	})
}

func TestAccTencentCloudCosBucketDataSource_full(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucketDataSource_full(appid),
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
	bucket = "tf-bucket-%d-%s"
}

data "tencentcloud_cos_buckets" "bucket_list" {
	bucket_prefix = "${tencentcloud_cos_bucket.bucket_basic.bucket}"
}
`, acctest.RandInt(), appid)
}

func testAccCosBucketDataSource_full(appid string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cos_bucket" "bucket_full" {
  bucket = "tf-bucket-%d-%s"
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
  }
  website {
    index_document = "index.html"
    error_document = "error.html"
  }
}

data "tencentcloud_cos_buckets" "bucket_list" {
  bucket_prefix = "${tencentcloud_cos_bucket.bucket_full.bucket}"
}
`, acctest.RandInt(), appid)
}
