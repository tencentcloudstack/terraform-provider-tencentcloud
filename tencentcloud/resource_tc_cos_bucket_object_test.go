package tencentcloud

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudCosBucketObject_source(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "tf-test-cos-object")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	err = ioutil.WriteFile(tmpFile.Name(), []byte("test terraform tencentcloud cos object"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Compatible with windows path format
	path := tmpFile.Name()
	if runtime.GOOS == "windows" {
		path = strings.Replace(path, "\\", "\\\\", -1)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucketObject_source(appid, path),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketObjectExists("tencentcloud_cos_bucket_object.object_source"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket_object.object_source", "content_type", "binary/octet-stream"),
				),
			},
		},
	})
}

func TestAccTencentCloudCosBucketObject_content(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucketObject_content(appid),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketObjectExists("tencentcloud_cos_bucket_object.object_content"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket_object.object_content", "content", "aaaaaaaaaaaaaaaa"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket_object.object_content", "content_type", "binary/octet-stream"),
				),
			},
		},
	})
}

func TestAccTencentCloudCosBucketObject_storageClass(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucketObject_storageClass(appid),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketObjectExists("tencentcloud_cos_bucket_object.object_storage"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket_object.object_storage", "storage_class", "STANDARD_IA"),
				),
			},
		},
	})
}

func TestAccTencentCloudCosBucketObject_acl(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucketObject_acl(appid, "private"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketObjectExists("tencentcloud_cos_bucket_object.object_acl"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket_object.object_acl", "acl", "private"),
				),
			},
			// test update acl
			{
				Config: testAccCosBucketObject_acl(appid, "public-read"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketObjectExists("tencentcloud_cos_bucket_object.object_acl"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket_object.object_acl", "acl", "public-read"),
				),
			},
		},
	})
}

func testAccCheckCosBucketObjectExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(nil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cos object %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cos object id is not set")
		}
		cosService := CosService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		bucket := rs.Primary.Attributes["bucket"]
		key := rs.Primary.Attributes["key"]
		_, err := cosService.HeadObject(ctx, bucket, key)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckCosBucketObjectDestroy(s *terraform.State) error {
	logId := getLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	cosService := CosService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cos_bucket_object" {
			continue
		}

		bucket := rs.Primary.Attributes["bucket"]
		key := rs.Primary.Attributes["key"]
		_, err := cosService.HeadObject(ctx, bucket, key)
		if err == nil {
			return fmt.Errorf("cos object still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCosBucketObject_source(appid string, source string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cos_bucket" "object_bucket" {
	bucket = "tf-bucket-%d-%s"
}

resource "tencentcloud_cos_bucket_object" "object_source" {
	bucket = "${tencentcloud_cos_bucket.object_bucket.bucket}"
	key = "tf-object-source"
	source = "%s"
	content_type = "binary/octet-stream"
}
`, acctest.RandInt(), appid, source)
}

func testAccCosBucketObject_content(appid string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cos_bucket" "object_bucket" {
	bucket = "tf-bucket-%d-%s"
}

resource "tencentcloud_cos_bucket_object" "object_content" {
	bucket = "${tencentcloud_cos_bucket.object_bucket.bucket}"
	key = "tf-object-content"
	content = "aaaaaaaaaaaaaaaa"
	content_type = "binary/octet-stream"
}
`, acctest.RandInt(), appid)
}

func testAccCosBucketObject_storageClass(appid string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cos_bucket" "object_bucket" {
	bucket = "tf-bucket-%d-%s"
}

resource "tencentcloud_cos_bucket_object" "object_storage" {
	bucket = "${tencentcloud_cos_bucket.object_bucket.bucket}"
	key = "tf-object-full"
	content = "aaaaaaaaaaaaaaaa"
	content_type = "binary/octet-stream"
	storage_class = "STANDARD_IA"
}
`, acctest.RandInt(), appid)
}

func testAccCosBucketObject_acl(appid, acl string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cos_bucket" "object_bucket" {
	bucket = "tf-bucket-acl-%s"
}

resource "tencentcloud_cos_bucket_object" "object_acl" {
	bucket = "${tencentcloud_cos_bucket.object_bucket.bucket}"
	key = "tf-object-acl"
	content = "aaaaaaaaaaaaaaaa"
	acl = "%s"
}
`, appid, acl)
}
