package cos_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localcos "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cos"

	"context"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudCosBucketObjectResource_source(t *testing.T) {
	t.Parallel()

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
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCosBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucketObject_source(tcacctest.Appid, path),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketObjectExists("tencentcloud_cos_bucket_object.object_source"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket_object.object_source", "content_type", "binary/octet-stream"),
				),
			},
		},
	})
}

func TestAccTencentCloudCosBucketObjectResource_content(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCosBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucketObject_content(tcacctest.Appid),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketObjectExists("tencentcloud_cos_bucket_object.object_content"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket_object.object_content", "content", "aaaaaaaaaaaaaaaa"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket_object.object_content", "content_type", "binary/octet-stream"),
				),
			},
		},
	})
}

func TestAccTencentCloudCosBucketObjectResource_tags(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCosBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucketObject_tags(tcacctest.Appid),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketObjectExists("tencentcloud_cos_bucket_object.object_with_tags"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket_object.object_with_tags", "tags.test", "test"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket_object.object_with_tags", "tags.hello", "world"),
				),
			},
		},
	})
}

func TestAccTencentCloudCosBucketObjectResource_storageClass(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCosBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucketObject_storageClass(tcacctest.Appid),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketObjectExists("tencentcloud_cos_bucket_object.object_storage"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket_object.object_storage", "storage_class", "STANDARD_IA"),
				),
			},
		},
	})
}

func TestAccTencentCloudCosBucketObjectResource_acl(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCosBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucketObject_acl(tcacctest.Appid, "private"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketObjectExists("tencentcloud_cos_bucket_object.object_acl"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket_object.object_acl", "acl", "private"),
				),
			},
			// test update acl
			{
				Config: testAccCosBucketObject_acl(tcacctest.Appid, "public-read"),
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cos object %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cos object id is not set")
		}
		cosService := localcos.NewCosService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

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
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cosService := localcos.NewCosService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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
  bucket       = tencentcloud_cos_bucket.object_bucket.bucket
  key          = "tf-object-source"
  source       = "%s"
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
  bucket       = tencentcloud_cos_bucket.object_bucket.bucket
  key          = "tf-object-content"
  content      = "aaaaaaaaaaaaaaaa"
  content_type = "binary/octet-stream"
}
`, acctest.RandInt(), appid)
}

func testAccCosBucketObject_tags(appid string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cos_bucket" "object_bucket" {
  bucket = "tf-bucket-%d-%s"
}

resource "tencentcloud_cos_bucket_object" "object_with_tags" {
  bucket       = tencentcloud_cos_bucket.object_bucket.bucket
  key          = "tf-object-tags"
  content       = "aaaaaaaaaaaaaaaa"
  content_type = "binary/octet-stream"
  tags = {
    test = "test"
    hello = "world"
  }
}
`, acctest.RandInt(), appid)
}

func testAccCosBucketObject_storageClass(appid string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cos_bucket" "object_bucket" {
  bucket = "tf-bucket-%d-%s"
}

resource "tencentcloud_cos_bucket_object" "object_storage" {
  bucket        = tencentcloud_cos_bucket.object_bucket.bucket
  key           = "tf-object-full"
  content       = "aaaaaaaaaaaaaaaa"
  content_type  = "binary/octet-stream"
  storage_class = "STANDARD_IA"
}
`, acctest.RandInt(), appid)
}

func testAccCosBucketObject_acl(appid, acl string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cos_bucket" "object_bucket" {
  bucket = "tf-bucket-obj-acl-%s"
}

resource "tencentcloud_cos_bucket_object" "object_acl" {
  bucket = tencentcloud_cos_bucket.object_bucket.bucket
  key = "acl.txt"
  content = "aaaaaaaaaaaaaaaa"
  acl     = "%s"
  
}
`, appid, acl)
}
