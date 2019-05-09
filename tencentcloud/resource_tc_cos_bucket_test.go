package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const appid string = "1258798060"

func init() {
	resource.AddTestSweepers("tencentcloud_cos_bucket", &resource.Sweeper{
		Name: "tencentcloud_cos_bucket",
		F:    testSweepCosBuckets,
	})
}

func testSweepCosBuckets(region string) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	sharedClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("geting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(TencentCloudClient)

	cosService := CosService{
		client: client.apiV3Conn,
	}
	buckets, err := cosService.ListBuckets(ctx)
	if err != nil {
		return fmt.Errorf("list buckets error: %s", err.Error())
	}

	for _, v := range buckets {
		bucket := *v.Name
		if !strings.HasPrefix(bucket, "tf-bucket-") {
			continue
		}

		// delete all object in the bucket before deleting bucket
		if objects, err := cosService.ListObjects(ctx, bucket); err != nil {
			log.Printf("[ERROR] list objects error: %s", err.Error())
		} else if len(objects) > 0 {
			for _, o := range objects {
				if err := cosService.DeleteObject(ctx, bucket, *o.Key); err != nil {
					log.Printf("[ERROR] delete object %s error: %s", *o.Key, err.Error())
				}
			}
		}
		log.Printf("[INFO] deleting cos bucket: %s", bucket)

		if err = cosService.DeleteBucket(ctx, bucket); err != nil {
			log.Printf("[ERROR] delete bucket %s error: %s", bucket, err.Error())
		}
	}
	return nil
}

func TestAccTencentCloudCosBucket_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucket_basic(appid),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.bucket_basic"),
				),
			},
			// test update bucket acl
			{
				Config: testAccCosBucket_basicUpdate(appid),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.bucket_basic"),
				),
			},
			{
				ResourceName:            "tencentcloud_cos_bucket.bucket_basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl"},
			},
		},
	})
}

func TestAccTencentCloudCosBucket_cors(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucket_cors(appid),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.bucket_cors"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_cors", "cors_rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_cors", "cors_rules.0.allowed_headers.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_cors", "cors_rules.0.allowed_headers.0", "*"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_cors", "cors_rules.0.allowed_methods.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_cors", "cors_rules.0.allowed_methods.0", "GET"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_cors", "cors_rules.0.allowed_methods.1", "POST"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_cors", "cors_rules.0.allowed_origins.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_cors", "cors_rules.0.allowed_origins.0", "https://www.test.com"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_cors", "cors_rules.0.expose_headers.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_cors", "cors_rules.0.expose_headers.0", "x-cos-test"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_cors", "cors_rules.0.max_age_seconds", "300"),
				),
			},
			// test updata bucket cors
			{
				Config: testAccCosBucket_corsUpdate(appid),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_cors", "cors_rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_cors", "cors_rules.0.allowed_methods.#", "3"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_cors", "cors_rules.0.allowed_methods.0", "GET"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_cors", "cors_rules.0.allowed_methods.1", "POST"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_cors", "cors_rules.0.allowed_methods.2", "PUT"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_cors", "cors_rules.0.allowed_origins.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_cors", "cors_rules.0.allowed_origins.0", "https://www.example.com"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_cors", "cors_rules.0.max_age_seconds", "100"),
				),
			},
			{
				ResourceName:            "tencentcloud_cos_bucket.bucket_cors",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl"},
			},
		},
	})
}

func TestAccTencentCloudCosBucket_lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucket_lifecycle(appid),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.bucket_lifecycle"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.filter_prefix", "test/"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.expiration.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.expiration.3672460294.days", "365"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.transition.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.transition.2000431762.days", "30"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.transition.2000431762.storage_class", "STANDARD_IA"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.transition.3491203533.days", "60"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.transition.3491203533.storage_class", "ARCHIVE"),
				),
			},
			// test update bucket lifecycle
			{
				Config: testAccBucket_lifecycleUpdate(appid),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.bucket_lifecycle"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.filter_prefix", "test/"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.expiration.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.expiration.2736768241.days", "300"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.transition.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.transition.2000431762.days", "30"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.transition.2000431762.storage_class", "STANDARD_IA"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.transition.1139768587.days", "90"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.transition.1139768587.storage_class", "ARCHIVE"),
				),
			},
			{
				ResourceName:            "tencentcloud_cos_bucket.bucket_lifecycle",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl"},
			},
		},
	})
}

func TestAccTencentCloudCosBucket_website(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucket_website(appid),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.bucket_website"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_website", "website.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_website", "website.0.index_document", "index.html"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_website", "website.0.error_document", "error.html"),
				),
			},
			// test update bucket website
			{
				Config: testAccBucket_websiteUpdate(appid),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.bucket_website"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_website", "website.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_website", "website.0.index_document", "testindex.html"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_website", "website.0.error_document", "testerror.html"),
				),
			},
			{
				ResourceName:            "tencentcloud_cos_bucket.bucket_website",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl"},
			},
		},
	})
}

func testAccCheckCosBucketExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := GetLogId(nil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cos bucket %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cos bucket id is not set")
		}
		cosService := CosService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		err := cosService.HeadBucket(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckCosBucketDestroy(s *terraform.State) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	cosService := CosService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cos_bucket" {
			continue
		}

		err := cosService.HeadBucket(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("cos bucket still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCosBucket_basic(appid string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cos_bucket" "bucket_basic" {
	bucket = "tf-bucket-basic-%s"
	acl = "public-read"
}
`, appid)
}

func testAccCosBucket_basicUpdate(appid string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cos_bucket" "bucket_basic" {
	bucket = "tf-bucket-basic-%s"
	acl = "private"
}
`, appid)
}

func testAccCosBucket_cors(appid string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cos_bucket" "bucket_cors" {
	bucket = "tf-bucket-cors-%s"
	acl = "public-read"
	cors_rules {
		allowed_headers = ["*"]
		allowed_methods = ["GET","POST"]
		allowed_origins = ["https://www.test.com"]
		expose_headers = ["x-cos-test"]
		max_age_seconds = 300
	}
}
`, appid)
}

func testAccCosBucket_corsUpdate(appid string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cos_bucket" "bucket_cors" {
	bucket = "tf-bucket-cors-%s"
	acl = "public-read"
	cors_rules {
		allowed_headers = ["*"]
		allowed_methods = ["GET","POST","PUT"]
		allowed_origins = ["https://www.example.com"]
		expose_headers = ["x-cos-test"]
		max_age_seconds = 100
	}
}
`, appid)
}

func testAccBucket_lifecycle(appid string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cos_bucket" "bucket_lifecycle" {
	bucket = "tf-bucket-lifecycle-%s"
	acl = "public-read"
	lifecycle_rules {
		filter_prefix = "test/"
		expiration {
			days = 365
		}
		transition {
			days = 30
			storage_class = "STANDARD_IA"
		}
		transition {
			days = 60
			storage_class = "ARCHIVE"
		}
	}
}
`, appid)
}

func testAccBucket_lifecycleUpdate(appid string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cos_bucket" "bucket_lifecycle" {
	bucket = "tf-bucket-lifecycle-%s"
	acl = "public-read"
	lifecycle_rules {
		filter_prefix = "test/"
		expiration {
			days = 300
		}
		transition {
			days = 30
			storage_class = "STANDARD_IA"
		}
		transition {
			days = 90
			storage_class = "ARCHIVE"
		}
	}
}
`, appid)
}

func testAccBucket_website(appid string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cos_bucket" "bucket_website" {
	bucket = "tf-bucket-website-%s"
	acl = "public-read"
	website {
		index_document = "index.html"
		error_document = "error.html"
	}
}
`, appid)
}

func testAccBucket_websiteUpdate(appid string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cos_bucket" "bucket_website" {
	bucket = "tf-bucket-website-%s"
	acl = "public-read"
	website {
		index_document = "testindex.html"
		error_document = "testerror.html"
	}
}
`, appid)
}
