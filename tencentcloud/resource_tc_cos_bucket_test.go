package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_cos_bucket
	resource.AddTestSweepers("tencentcloud_cos_bucket", &resource.Sweeper{
		Name: "tencentcloud_cos_bucket",
		F:    testSweepCosBuckets,
	})
}

func testSweepCosBuckets(region string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	sharedClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(*TencentCloudClient)

	cosService := CosService{
		client: client.apiV3Conn,
	}
	buckets, err := cosService.ListBuckets(ctx)
	if err != nil {
		return fmt.Errorf("list buckets error: %s", err.Error())
	}

	prefix := regexp.MustCompile("^(tf|test)-")

	for _, v := range buckets {
		bucket := *v.Name
		if !prefix.MatchString(bucket) {
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

		if err = cosService.DeleteBucket(ctx, bucket, true, true); err != nil {
			log.Printf("[ERROR] delete bucket %s error: %s", bucket, err.Error())
		}
	}
	return nil
}

func TestAccTencentCloudCosBucketResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucket_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.bucket_basic"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_basic", "enable_intelligent_tiering", "true"),
				),
			},
			// test update bucket acl
			{
				Config: testAccCosBucket_basicUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.bucket_basic"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_basic", "encryption_algorithm", "AES256"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_basic", "versioning_enable", "true"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_basic", "acceleration_enable", "true"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_basic", "force_clean", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_cos_bucket.bucket_basic", "cos_bucket_url"),
				),
			},
			{
				ResourceName:            "tencentcloud_cos_bucket.bucket_basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_clean"},
			},
		},
	})
}

func TestAccTencentCloudCosBucketResource_ACL(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucket_ACL(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.bucket_acl"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_acl", "acl", "public-read"),
					resource.TestCheckResourceAttrSet("tencentcloud_cos_bucket.bucket_acl", "acl_body"),
				),
			},
			// test update bucket acl
			{
				Config: testAccCosBucket_ACLUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.bucket_acl"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_acl", "acl", "private"),
					resource.TestCheckResourceAttrSet("tencentcloud_cos_bucket.bucket_acl", "acl_body"),
				),
			},
			{
				ResourceName:            "tencentcloud_cos_bucket.bucket_acl",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl_body"},
			},
		},
	})
}

func TestAccTencentCloudCosBucketResource_tags(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucket_tags(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.bucket_tags"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_tags", "tags.test", "test"),
				),
			},
			{
				Config: testAccCosBucket_tagsReplace(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.bucket_tags"),
					resource.TestCheckNoResourceAttr("tencentcloud_cos_bucket.bucket_tags", "tags.test"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_tags", "tags.abc", "abc"),
				),
			},
			{
				Config: testAccCosBucket_tagsDelete(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.bucket_tags"),
					resource.TestCheckNoResourceAttr("tencentcloud_cos_bucket.bucket_tags", "tags.abc"),
				),
			},
		},
	})
}

func TestAccTencentCloudCosBucketResource_cors(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucket_cors(),
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
				Config: testAccCosBucket_corsUpdate(),
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

func TestAccTencentCloudCosBucketResource_lifecycle(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucket_lifecycle(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.bucket_lifecycle"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.id", "rule1"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.filter_prefix", "test/"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.expiration.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.expiration.*",
						map[string]string{
							"days": "365",
						}),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.transition.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.transition.*",
						map[string]string{
							"days":          "30",
							"storage_class": "STANDARD_IA",
						}),
					resource.TestCheckTypeSetElemNestedAttrs("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.transition.*",
						map[string]string{
							"days":          "60",
							"storage_class": "ARCHIVE",
						}),
				),
			},
			// test update bucket lifecycle
			{
				Config: testAccBucket_lifecycleUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.bucket_lifecycle"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.filter_prefix", "test/"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.expiration.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.expiration.*",
						map[string]string{
							"days": "300",
						}),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.transition.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.transition.*",
						map[string]string{
							"days":          "30",
							"storage_class": "STANDARD_IA",
						}),
					resource.TestCheckTypeSetElemNestedAttrs("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.transition.*",
						map[string]string{
							"days":          "90",
							"storage_class": "ARCHIVE",
						}),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.non_current_expiration.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_lifecycle", "lifecycle_rules.0.non_current_transition.#", "2"),
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

func TestAccTencentCloudCosBucketResource_website(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucket_website(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.bucket_website"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_website", "website.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_website", "website.0.index_document", "index.html"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_website", "website.0.error_document", "error.html"),
					resource.TestCheckResourceAttrSet("tencentcloud_cos_bucket.bucket_website", "website.0.endpoint"),
				),
			},
			// test update bucket website
			{
				Config: testAccBucket_websiteUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.bucket_website"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_website", "website.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_website", "website.0.index_document", "testindex.html"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_website", "website.0.error_document", "testerror.html"),
					resource.TestCheckResourceAttrSet("tencentcloud_cos_bucket.bucket_website", "website.0.endpoint"),
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

func TestAccTencentCloudCosBucketResource_MAZ(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucket_MAZ(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.bucket_maz"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.bucket_maz", "multi_az", "true"),
				),
			},
			{
				ResourceName:            "tencentcloud_cos_bucket.bucket_maz",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl", "multi_az"},
			},
		},
	})
}

func TestAccTencentCloudCosBucketResource_originPull(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucket_originPull(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.with_origin"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_origin", "origin_pull_rules.0.priority", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_origin", "origin_pull_rules.0.sync_back_to_source", "false"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_origin", "origin_pull_rules.0.host", "abc.example.com"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_origin", "origin_pull_rules.0.prefix", "/"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_origin", "origin_pull_rules.0.protocol", "FOLLOW"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_origin", "origin_pull_rules.0.follow_query_string", "true"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_origin", "origin_pull_rules.0.follow_redirection", "true"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_cos_bucket.with_origin", "origin_pull_rules.0.follow_http_headers.*", "origin"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_cos_bucket.with_origin", "origin_pull_rules.0.follow_http_headers.*", "host"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_cos_bucket.with_origin", "origin_pull_rules.0.follow_http_headers.*", "expires"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_origin", "origin_pull_rules.0.custom_http_headers.x-custom-header", "custom_value"),
				),
			},
			{
				Config: testAccBucket_originPullUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.with_origin"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_origin", "origin_pull_rules.0.priority", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_origin", "origin_pull_rules.0.sync_back_to_source", "false"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_origin", "origin_pull_rules.0.host", "test.abc.example.com"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_origin", "origin_pull_rules.0.prefix", "/test"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_origin", "origin_pull_rules.0.protocol", "FOLLOW"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_origin", "origin_pull_rules.0.follow_query_string", "true"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_origin", "origin_pull_rules.0.follow_redirection", "true"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_cos_bucket.with_origin", "origin_pull_rules.0.follow_http_headers.*", "origin"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_cos_bucket.with_origin", "origin_pull_rules.0.follow_http_headers.*", "host"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_cos_bucket.with_origin", "origin_pull_rules.0.follow_http_headers.*", "expires"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_origin", "origin_pull_rules.0.custom_http_headers.x-custom-header", "test"),
				),
			},
			{
				ResourceName:            "tencentcloud_cos_bucket.with_origin",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl"},
			},
		},
	})
}

// TODO this case is now disabled until domain configured
/*
func TestAccTencentCloudCosBucket_originDomain(t *testing.T) {

	t.Parallel()

	randomName := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucket_originDomain(appid, randomName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.with_domain"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_domain", "origin_domain_rules.0.status", "ENABLED"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_domain", "origin_domain_rules.0.domain", "www.example.com"),
				),
			},
			{
				Config: testAccBucket_originDomainUpdate(appid, randomName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.with_domain"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_domain", "origin_domain_rules.0.status", "DISABLED"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_domain", "origin_domain_rules.0.domain", "www.example.com"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_domain", "origin_domain_rules.1.status", "ENABLED"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_domain", "origin_domain_rules.1.domain", "test.example1.com"),
				),
			},
			{
				ResourceName:            "tencentcloud_cos_bucket.with_domain",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl"},
			},
		},
	})
}
*/

func TestAccTencentCloudCosBucketResource_replication(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketReplication(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.with_replication"),
					resource.TestMatchResourceAttr("tencentcloud_cos_bucket.with_replication", "replica_role", regexp.MustCompile(`^qcs::cam::uin/\d+:uin/\d+$`)),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_replication", "replica_rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_replication", "replica_rules.0.id", "test-rep1"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_replication", "replica_rules.0.status", "Enabled"),
				),
			},
			{
				Config: testAccBucketReplicationUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.with_replication"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_replication", "replica_rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_replication", "replica_rules.0.status", "Disabled"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_replication", "replica_rules.0.prefix", "dist"),
				),
			},
			{
				Config: testAccBucketReplicationRemove(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosBucketExists("tencentcloud_cos_bucket.with_replication"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_replication", "replica_role", ""),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket.with_replication", "replica_rules.#", "0"),
				),
			},
			{
				ResourceName:            "tencentcloud_cos_bucket.with_replication",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl", "replica_role"},
			},
		},
	})
}

func testAccCheckCosBucketExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

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
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

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

func testAccCosBucket_basic() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_cos_bucket" "bucket_basic" {
  bucket = "tf-bucket-basic-${local.app_id}"
  acl    = "public-read"
  enable_intelligent_tiering = true
}
`, userInfoData)
}

func testAccCosBucket_basicUpdate() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_cos_bucket" "bucket_basic" {
  bucket               = "tf-bucket-basic-${local.app_id}"
  acl                  = "private"
  encryption_algorithm = "AES256"
  versioning_enable    = true
  acceleration_enable  = true
  force_clean          = true
}
`, userInfoData)
}

func testAccCosBucket_ACL() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_cos_bucket" "bucket_acl" {
  bucket	= "tf-bucket-acl-${local.app_id}"
  acl       = "public-read"
  acl_body 	= <<EOF
<AccessControlPolicy>
    <Owner>
        <ID>qcs::cam::uin/${local.uin}:uin/${local.uin}</ID>
		<DisplayName>qcs::cam::uin/${local.uin}:uin/${local.uin}</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>
            <Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="Group">
                <URI>http://cam.qcloud.com/groups/global/AllUsers</URI>
            </Grantee>
            <Permission>READ</Permission>
        </Grant>
        <Grant>
            <Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="CanonicalUser">
                <ID>qcs::cam::uin/${local.uin}:uin/${local.uin}</ID>
				<DisplayName>qcs::cam::uin/${local.uin}:uin/${local.uin}</DisplayName>
            </Grantee>
            <Permission>FULL_CONTROL</Permission>
        </Grant>
    </AccessControlList>
</AccessControlPolicy>
EOF
}
`, userInfoData)
}

func testAccCosBucket_ACLUpdate() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_cos_bucket" "bucket_acl" {
  bucket	= "tf-bucket-acl-${local.app_id}"
  acl 		= "private"
  acl_body	= <<EOF
<AccessControlPolicy>
    <Owner>
        <ID>qcs::cam::uin/${local.uin}:uin/${local.uin}</ID>
		<DisplayName>qcs::cam::uin/${local.uin}:uin/${local.uin}</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>
            <Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="CanonicalUser">
                <ID>qcs::cam::uin/${local.uin}:uin/${local.uin}</ID>
				<DisplayName>qcs::cam::uin/${local.uin}:uin/${local.uin}</DisplayName>
            </Grantee>
            <Permission>FULL_CONTROL</Permission>
        </Grant>
    </AccessControlList>
</AccessControlPolicy>
EOF
}
`, userInfoData)
}

func testAccCosBucket_tags() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_cos_bucket" "bucket_tags" {
  bucket = "tf-bucket-tags-${local.app_id}"
  acl    = "public-read"

  tags = {
    "test" = "test"
  }
}
`, userInfoData)
}

func testAccCosBucket_tagsReplace() string {
	return fmt.Sprintf(`
%s
resource "tencentcloud_cos_bucket" "bucket_tags" {
  bucket = "tf-bucket-tags-${local.app_id}"
  acl    = "public-read"

  tags = {
    "abc" = "abc"
  }
}
`, userInfoData)
}

func testAccCosBucket_tagsDelete() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_cos_bucket" "bucket_tags" {
  bucket = "tf-bucket-tags-${local.app_id}"
  acl    = "public-read"
}
`, userInfoData)
}

func testAccCosBucket_cors() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_cos_bucket" "bucket_cors" {
  bucket = "tf-bucket-cors-${local.app_id}"
  acl    = "public-read"

  cors_rules {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "POST"]
    allowed_origins = ["https://www.test.com"]
    expose_headers  = ["x-cos-test"]
    max_age_seconds = 300
  }
}
`, userInfoData)
}

func testAccCosBucket_corsUpdate() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_cos_bucket" "bucket_cors" {
  bucket = "tf-bucket-cors-${local.app_id}"
  acl    = "public-read"
  cors_rules {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "POST", "PUT"]
    allowed_origins = ["https://www.example.com"]
    expose_headers  = ["x-cos-test"]
    max_age_seconds = 100
  }
}
`, userInfoData)
}

func testAccBucket_lifecycle() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_cos_bucket" "bucket_lifecycle" {
  bucket = "tf-bucket-lifecycle-${local.app_id}"
  acl    = "public-read"
  versioning_enable = true
  lifecycle_rules {
    id = "rule1"
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
}
`, userInfoData)
}

func testAccBucket_lifecycleUpdate() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_cos_bucket" "bucket_lifecycle" {
  bucket = "tf-bucket-lifecycle-${local.app_id}"
  acl    = "public-read"
  versioning_enable = true
  lifecycle_rules {
    id = "rule1"
    filter_prefix = "test/"
    expiration {
      days = 300
    }
    transition {
      days          = 30
      storage_class = "STANDARD_IA"
    }
    transition {
      days          = 90
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
}
`, userInfoData)
}

func testAccBucket_website() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_cos_bucket" "bucket_website" {
  bucket = "tf-bucket-website-${local.app_id}"
  acl    = "public-read"
  website {
    index_document = "index.html"
    error_document = "error.html"
  }
}
`, userInfoData)
}

func testAccBucket_websiteUpdate() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_cos_bucket" "bucket_website" {
  bucket = "tf-bucket-website-${local.app_id}"
  acl    = "public-read"
  website {
    index_document = "testindex.html"
    error_document = "testerror.html"
  }
}
`, userInfoData)
}

func testAccBucket_MAZ() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_cos_bucket" "bucket_maz" {
  bucket   = "tf-bucket-maz-${local.app_id}"
  acl      = "public-read"
  multi_az = true
  versioning_enable = true
}
`, userInfoData)
}

func testAccBucket_originPull() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_cos_bucket" "with_origin" {
  bucket = "tf-bucket-origin-${local.app_id}"
  acl    = "private"
  origin_pull_rules {
    priority = 1
    sync_back_to_source = false
    host = "abc.example.com"
    prefix = "/"
    protocol = "FOLLOW" // "HTTP" "HTTPS"
    follow_query_string = true
    follow_redirection = true
    follow_http_headers = ["origin", "host", "expires"]
    custom_http_headers = {
	  "x-custom-header" = "custom_value"
    }
  }
}
`, userInfoData)
}

func testAccBucket_originPullUpdate() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_cos_bucket" "with_origin" {
  bucket = "tf-bucket-origin-${local.app_id}"
  acl    = "private"
  origin_pull_rules {
    priority = 1
    sync_back_to_source = false
    host = "test.abc.example.com"
    prefix = "/test"
    protocol = "FOLLOW" // "HTTP" "HTTPS"
    follow_query_string = true
    follow_redirection = true
    follow_http_headers = ["origin", "host", "expires"]
    custom_http_headers = {
	  "x-custom-header" = "test"
    }
  }
}
`, userInfoData)
}

func testAccBucketReplication() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_cos_bucket" "replica1" {
  bucket = "tf-replica-foo-${local.app_id}"
  acl    = "private"
  versioning_enable = true
}

resource "tencentcloud_cos_bucket" "with_replication" {
  bucket = "tf-bucket-replica-${local.app_id}"
  acl    = "private"
  versioning_enable = true
  replica_role = "qcs::cam::uin/${local.owner_uin}:uin/${local.uin}"
  replica_rules {
	id = "test-rep1"
    status = "Enabled"
    destination_bucket = "qcs::cos:%s::${tencentcloud_cos_bucket.replica1.bucket}"
  }
}
`, userInfoData, defaultRegion)
}

func testAccBucketReplicationUpdate() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_cos_bucket" "replica1" {
  bucket = "tf-replica-foo-${local.app_id}"
  acl    = "private"
  versioning_enable = true
}

resource "tencentcloud_cos_bucket" "with_replication" {
  bucket = "tf-bucket-replica-${local.app_id}"
  acl    = "private"
  versioning_enable = true
  replica_role = "qcs::cam::uin/${local.owner_uin}:uin/${local.uin}"
  replica_rules {
	id = "test-rep1"
    status = "Disabled"
    prefix = "dist"
    destination_bucket = "qcs::cos:%s::${tencentcloud_cos_bucket.replica1.bucket}"
  }
}
`, userInfoData, defaultRegion)
}

func testAccBucketReplicationRemove() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_cos_bucket" "replica1" {
  bucket = "tf-replica-foo-${local.app_id}"
  acl    = "private"
  versioning_enable = true
}

resource "tencentcloud_cos_bucket" "with_replication" {
  bucket = "tf-bucket-replica-${local.app_id}"
  acl    = "private"
  versioning_enable = true
}
`, userInfoData)
}
