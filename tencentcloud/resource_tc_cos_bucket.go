/*
Provides a COS resource to create a COS bucket and set its attributes.

Example Usage

Private Bucket

```hcl
resource "tencentcloud_cos_bucket" "mycos" {
  bucket = "mycos-1258798060"
  acl    = "private"
}
```

Static Website

```hcl
resource "tencentcloud_cos_bucket" "mycos" {
  bucket = "mycos-1258798060"

  website = {
    index_document = "index.html"
    error_document = "error.html"
  }
}
```

Using CORS

```hcl
resource "tencentcloud_cos_bucket" "mycos" {
  bucket = "mycos-1258798060"
  acl    = "public-read-write"

  cors_rules {
    allowed_origins = ["http://*.abc.com"]
    allowed_methods = ["PUT", "POST"]
    allowed_headers = ["*"]
    max_age_seconds = 300
    expose_headers  = ["Etag"]
  }
}
```

Using object lifecycle

```hcl
resource "tencentcloud_cos_bucket" "mycos" {
  bucket = "mycos-1258798060"
  acl    = "public-read-write"

  lifecycle_rules {
    filter_prefix = "path1/"

    transition {
      date          = "2019-06-01"
      storage_class = "STANDARD_IA"
    }

    expiration {
      days = 90
    }
  }
}
```

Setting log status

```hcl
resource "tencentcloud_cam_role" "cosLogGrant" {
  name          = "CLS_QcsRole"
document      = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"principal\":{\"service\":[\"cls.cloud.tencent.com\"]}}]}"

  description   = "cos log enable grant"
}


data "tencentcloud_cam_policies" "cosAccess" {
  name      = "QcloudCOSAccessForCLSRole"
}


resource "tencentcloud_cam_role_policy_attachment" "cosLogGrant" {
  role_id   = tencentcloud_cam_role.cosLogGrant.id
  policy_id = data.tencentcloud_cam_policies.cosAccess.policy_list.0.policy_id
}


resource "tencentcloud_cos_bucket" "mylog" {
  bucket = "mylog-1258798060"
  acl    = "private"
}

resource "tencentcloud_cos_bucket" "mycos" {
  bucket = "mycos-1258798060"
  acl    = "private"
  log_enable = true
  log_target_bucket = "mylog-1258798060"
  log_prefix = "MyLogPrefix"
}
```

Import

COS bucket can be imported, e.g.

```
$ terraform import tencentcloud_cos_bucket.bucket bucket-name
```
*/
package tencentcloud

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

const (
	tencentCloudCosStorageClassStandard   = "STANDARD"
	tencentCloudCosStorageClassStandardIA = "STANDARD_IA"
	tencentCloudCosStorageClassArchive    = "ARCHIVE"
)

var (
	availableCosStorageClass = []string{
		tencentCloudCosStorageClassStandard,
		tencentCloudCosStorageClassStandardIA,
		tencentCloudCosStorageClassArchive,
	}
)

func resourceTencentCloudCosBucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCosBucketCreate,
		Read:   resourceTencentCloudCosBucketRead,
		Update: resourceTencentCloudCosBucketUpdate,
		Delete: resourceTencentCloudCosBucketDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCosBucketName,
				Description:  "The name of a bucket to be created. Bucket format should be [custom name]-[appid], for example `mycos-1258798060`.",
			},
			"acl": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  s3.ObjectCannedACLPrivate,
				ValidateFunc: validateAllowedStringValue([]string{
					s3.ObjectCannedACLPrivate,
					s3.ObjectCannedACLPublicRead,
					s3.ObjectCannedACLPublicReadWrite,
				}),
				Description: "The canned ACL to apply. Valid values: private, public-read, and public-read-write. Defaults to private.",
			},
			"encryption_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The server-side encryption algorithm to use. Valid value is `AES256`.",
			},
			"versioning_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable bucket versioning.",
			},
			"cors_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A rule of Cross-Origin Resource Sharing (documented below).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_origins": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Specifies which origins are allowed.",
						},
						"allowed_methods": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Specifies which methods are allowed. Can be `GET`, `PUT`, `POST`, `DELETE` or `HEAD`.",
						},
						"allowed_headers": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Specifies which headers are allowed.",
						},
						"max_age_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies time in seconds that browser can cache the response for a preflight request.",
						},
						"expose_headers": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Specifies expose header in the response.",
						},
					},
				},
			},
			"lifecycle_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A configuration of object lifecycle management (documented below).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filter_prefix": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Object key prefix identifying one or more objects to which the rule applies.",
						},
						"transition": {
							Type:        schema.TypeSet,
							Optional:    true,
							Set:         transitionHash,
							Description: "Specifies a period in the object's transitions (documented below).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validateCosBucketLifecycleTimestamp,
										Description:  "Specifies the date after which you want the corresponding action to take effect.",
									},
									"days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validateIntegerMin(0),
										Description:  "Specifies the number of days after object creation when the specific rule action takes effect.",
									},
									"storage_class": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateAllowedStringValue(availableCosStorageClass),
										Description:  "Specifies the storage class to which you want the object to transition. Available values include `STANDARD`, `STANDARD_IA` and `ARCHIVE`.",
									},
								},
							},
						},
						"expiration": {
							Type:        schema.TypeSet,
							Optional:    true,
							Set:         expirationHash,
							MaxItems:    1,
							Description: "Specifies a period in the object's expire (documented below).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validateCosBucketLifecycleTimestamp,
										Description:  "Specifies the date after which you want the corresponding action to take effect.",
									},
									"days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validateIntegerMin(0),
										Description:  "Specifies the number of days after object creation when the specific rule action takes effect.",
									},
								},
							},
						},
					},
				},
			},
			"website": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "A website object(documented below).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index_document": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "COS returns this index document when requests are made to the root domain or any of the subfolders.",
						},
						"error_document": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "An absolute path to the document to return in case of a 4XX error.",
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The tags of a bucket.",
			},
			"log_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicate the access log of this bucket to be saved or not. Default is `false`. If set `true`, the access log will be saved with `log_target_bucket`. To enable log, the full access of log service must be granted. [Full Access Role Policy](https://intl.cloud.tencent.com/document/product/436/16920).",
			},
			"log_target_bucket": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The target bucket name which saves the access log of this bucket per 5 minutes. The log access file format is `log_target_bucket`/`log_prefix`{YYYY}/{MM}/{DD}/{time}_{random}_{index}.gz. Only valid when `log_enable` is `true`. User must have full access on this bucket.",
			},
			"log_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The prefix log name which saves the access log of this bucket per 5 minutes. Eg. `MyLogPrefix/`. The log access file format is `log_target_bucket`/`log_prefix`{YYYY}/{MM}/{DD}/{time}_{random}_{index}.gz. Only valid when `log_enable` is `true`.",
			},
			//computed
			"cos_bucket_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of this cos bucket.",
			},
		},
	}
}

func resourceTencentCloudCosBucketCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_bucket.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	bucket := d.Get("bucket").(string)
	acl := d.Get("acl").(string)

	cosService := CosService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := cosService.PutBucket(ctx, bucket, acl)
	if err != nil {
		return err
	}

	d.SetId(bucket)

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		if err := cosService.SetBucketTags(ctx, bucket, tags); err != nil {
			return err
		}
	}

	return resourceTencentCloudCosBucketUpdate(d, meta)
}

func resourceTencentCloudCosBucketRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_bucket.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	bucket := d.Id()
	cosService := CosService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := cosService.HeadBucket(ctx, bucket)
	if err != nil {
		if awsError, ok := err.(awserr.RequestFailure); ok && awsError.StatusCode() == 404 {
			log.Printf("[WARN]%s bucket (%s) not found, error code (404)", logId, bucket)
			d.SetId("")
			return nil
		} else {
			return err
		}
	}

	cosBuckeUrl := fmt.Sprintf("%s.cos.%s.myqcloud.com", d.Id(), meta.(*TencentCloudClient).apiV3Conn.Region)
	_ = d.Set("cos_bucket_url", cosBuckeUrl)
	// set bucket in the import case
	if _, ok := d.GetOk("bucket"); !ok {
		_ = d.Set("bucket", d.Id())
	}

	// read the cors
	corsRules, err := cosService.GetBucketCors(ctx, bucket)
	if err != nil {
		return err
	}
	if err = d.Set("cors_rules", corsRules); err != nil {
		return fmt.Errorf("setting cors_rules error: %v", err)
	}

	// read the lifecycle
	lifecycleRules, err := cosService.GetBucketLifecycle(ctx, bucket)
	if err != nil {
		return err
	}
	if err = d.Set("lifecycle_rules", lifecycleRules); err != nil {
		return fmt.Errorf("setting lifecycle_rules error: %v", err)
	}

	// read the website
	website, err := cosService.GetBucketWebsite(ctx, bucket)
	if err != nil {
		return err
	}
	if err = d.Set("website", website); err != nil {
		return fmt.Errorf("setting website error: %v", err)
	}

	// read the encryption algorithm
	encryption, err := cosService.GetBucketEncryption(ctx, bucket)
	if err != nil {
		return err
	}
	if err = d.Set("encryption_algorithm", encryption); err != nil {
		return fmt.Errorf("setting encryption error: %v", err)
	}

	// read the versioning
	versioning, err := cosService.GetBucketVersioning(ctx, bucket)
	if err != nil {
		return err
	}
	if err = d.Set("versioning_enable", versioning); err != nil {
		return fmt.Errorf("setting versioning_enable error: %v", err)
	}

	//read the log
	logEnable, logTargetBucket, logPrefix, err := cosService.GetBucketLogStatus(ctx, bucket)
	if err != nil {
		return err
	}
	_ = d.Set("log_enable", logEnable)
	_ = d.Set("log_target_bucket", logTargetBucket)
	_ = d.Set("log_prefix", logPrefix)

	// read the tags
	tags, err := cosService.GetBucketTags(ctx, bucket)
	if err != nil {
		return fmt.Errorf("get tags failed: %v", err)
	}
	if len(tags) > 0 {
		_ = d.Set("tags", tags)
	}

	return nil
}

func resourceTencentCloudCosBucketUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_bucket.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	client := meta.(*TencentCloudClient).apiV3Conn.UseCosClient()

	d.Partial(true)

	if d.HasChange("acl") {
		err := resourceTencentCloudCosBucketAclUpdate(ctx, client, d)
		if err != nil {
			return err
		}
		d.SetPartial("acl")
	}

	if d.HasChange("cors_rules") {
		err := resourceTencentCloudCosBucketCorsUpdate(ctx, client, d)
		if err != nil {
			return err
		}
		d.SetPartial("cors_rules")
	}

	if d.HasChange("lifecycle_rules") {
		err := resourceTencentCloudCosBucketLifecycleUpdate(ctx, client, d)
		if err != nil {
			return err
		}
		d.SetPartial("lifecycle_rules")
	}

	if d.HasChange("website") {
		err := resourceTencentCloudCosBucketWebsiteUpdate(ctx, client, d)
		if err != nil {
			return err
		}
		d.SetPartial("website")
	}

	if d.HasChange("encryption_algorithm") {
		err := resourceTencentCloudCosBucketEncryptionUpdate(ctx, client, d)
		if err != nil {
			return err
		}
		d.SetPartial("encryption_algorithm")
	}

	if d.HasChange("versioning_enable") {
		err := resourceTencentCloudCosBucketVersioningUpdate(ctx, client, d)
		if err != nil {
			return err
		}
		d.SetPartial("versioning_enable")
	}

	if d.HasChange("tags") {
		bucket := d.Id()

		cosService := CosService{client: meta.(*TencentCloudClient).apiV3Conn}
		if err := cosService.SetBucketTags(ctx, bucket, helper.GetTags(d, "tags")); err != nil {
			return err
		}

		d.SetPartial("tags")
	}

	if d.HasChange("log_enable") || d.HasChange("log_target_bucket") || d.HasChange("log_prefix") {
		err := resourceTencentCloudCosBucketLogStatusUpdate(ctx, client, d)
		if err != nil {
			return err
		}
		d.SetPartial("log_enable")
		d.SetPartial("log_target_bucket")
		d.SetPartial("log_prefix")
	}

	d.Partial(false)

	// wait for update cache
	// if not, the data may be outdated.
	time.Sleep(3 * time.Second)

	return resourceTencentCloudCosBucketRead(d, meta)
}

func resourceTencentCloudCosBucketDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_bucket.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	bucket := d.Id()
	cosService := CosService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := cosService.DeleteBucket(ctx, bucket)
	if err != nil {
		return err
	}

	// wait for update cache
	// if not, head bucket may be successful
	time.Sleep(3 * time.Second)

	return nil
}

func resourceTencentCloudCosBucketEncryptionUpdate(ctx context.Context, client *s3.S3, d *schema.ResourceData) error {
	logId := getLogId(ctx)

	bucket := d.Get("bucket").(string)
	encryption := d.Get("encryption_algorithm").(string)
	if encryption == "" {
		request := s3.DeleteBucketEncryptionInput{
			Bucket: aws.String(bucket),
		}
		response, err := client.DeleteBucketEncryption(&request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete bucket encryption", request.String(), err.Error())
			return fmt.Errorf("cos delete bucket error: %s, bucket: %s", err.Error(), bucket)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, "delete bucket encryption", request.String(), response.String())

		return nil
	}

	request := s3.PutBucketEncryptionInput{
		Bucket: aws.String(bucket),
	}
	request.ServerSideEncryptionConfiguration = &s3.ServerSideEncryptionConfiguration{}
	rules := make([]*s3.ServerSideEncryptionRule, 0)
	defaultRule := &s3.ServerSideEncryptionByDefault{
		SSEAlgorithm: aws.String(encryption),
	}
	rule := &s3.ServerSideEncryptionRule{
		ApplyServerSideEncryptionByDefault: defaultRule,
	}
	rules = append(rules, rule)
	request.ServerSideEncryptionConfiguration.Rules = rules

	response, err := client.PutBucketEncryption(&request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "put bucket encryption", request.String(), err.Error())
		return fmt.Errorf("cos put bucket encryption error: %s, bucket: %s", err.Error(), bucket)
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "put bucket encryption", request.String(), response.String())

	return nil
}

func resourceTencentCloudCosBucketVersioningUpdate(ctx context.Context, client *s3.S3, d *schema.ResourceData) error {
	logId := getLogId(ctx)

	bucket := d.Get("bucket").(string)
	versioning := d.Get("versioning_enable").(bool)
	status := "Suspended"
	if versioning {
		status = "Enabled"
	}
	request := s3.PutBucketVersioningInput{
		Bucket: aws.String(bucket),
		VersioningConfiguration: &s3.VersioningConfiguration{
			Status: aws.String(status),
		},
	}
	response, err := client.PutBucketVersioning(&request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "put bucket encryption", request.String(), err.Error())
		return fmt.Errorf("cos put bucket encryption error: %s, bucket: %s", err.Error(), bucket)
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "put bucket encryption", request.String(), response.String())

	return nil
}

func resourceTencentCloudCosBucketAclUpdate(ctx context.Context, client *s3.S3, d *schema.ResourceData) error {
	logId := getLogId(ctx)

	bucket := d.Get("bucket").(string)
	acl := d.Get("acl").(string)
	request := s3.PutBucketAclInput{
		Bucket: aws.String(bucket),
		ACL:    aws.String(acl),
	}
	response, err := client.PutBucketAcl(&request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "put bucket acl", request.String(), err.Error())
		return fmt.Errorf("cos put bucket error: %s, bucket: %s", err.Error(), bucket)
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "put bucket acl", request.String(), response.String())

	return nil
}

func resourceTencentCloudCosBucketCorsUpdate(ctx context.Context, client *s3.S3, d *schema.ResourceData) error {
	logId := getLogId(ctx)

	bucket := d.Get("bucket").(string)
	cors := d.Get("cors_rules").([]interface{})

	if len(cors) == 0 {
		request := s3.DeleteBucketCorsInput{
			Bucket: aws.String(bucket),
		}
		response, err := client.DeleteBucketCors(&request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete bucket cors", request.String(), err.Error())
			return fmt.Errorf("cos delete bucket cors error: %s, bucket: %s", err.Error(), bucket)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, "delete bucket cors", request.String(), response.String())
	} else {
		rules := make([]*s3.CORSRule, 0, len(cors))
		for _, item := range cors {
			corsMap := item.(map[string]interface{})
			rule := &s3.CORSRule{}
			for k, v := range corsMap {
				if k == "max_age_seconds" {
					rule.MaxAgeSeconds = aws.Int64(int64(v.(int)))
				} else {
					vMap := make([]*string, len(v.([]interface{})))
					for i, value := range v.([]interface{}) {
						if str, ok := value.(string); ok {
							vMap[i] = aws.String(str)
						}
					}
					switch k {
					case "allowed_origins":
						rule.AllowedOrigins = vMap
					case "allowed_methods":
						rule.AllowedMethods = vMap
					case "allowed_headers":
						rule.AllowedHeaders = vMap
					case "expose_headers":
						rule.ExposeHeaders = vMap
					}
				}
			}
			rules = append(rules, rule)
		}
		request := s3.PutBucketCorsInput{
			Bucket: aws.String(bucket),
			CORSConfiguration: &s3.CORSConfiguration{
				CORSRules: rules,
			},
		}
		response, err := client.PutBucketCors(&request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "put bucket cors", request.String(), err.Error())
			return fmt.Errorf("cos put bucket cors error: %s, bucket: %s", err.Error(), bucket)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, "put bucket cors", request.String(), response.String())
	}
	return nil
}

func resourceTencentCloudCosBucketLifecycleUpdate(ctx context.Context, client *s3.S3, d *schema.ResourceData) error {
	logId := getLogId(ctx)

	bucket := d.Get("bucket").(string)
	lifecycleRules := d.Get("lifecycle_rules").([]interface{})
	if len(lifecycleRules) == 0 {
		request := s3.DeleteBucketLifecycleInput{
			Bucket: aws.String(bucket),
		}
		response, err := client.DeleteBucketLifecycle(&request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete bucket lifecycle", request.String(), err.Error())
			return fmt.Errorf("cos delete bucket lifecycle error: %s, bucket: %s", err.Error(), bucket)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, "delete bucket lifecycle", request.String(), response.String())
	} else {
		rules := make([]*s3.LifecycleRule, 0, len(lifecycleRules))
		for i, lifecycleRule := range lifecycleRules {
			r := lifecycleRule.(map[string]interface{})
			rule := &s3.LifecycleRule{}
			rule.Status = aws.String(s3.ExpirationStatusEnabled)
			prefix := r["filter_prefix"].(string)
			rule.Filter = &s3.LifecycleRuleFilter{
				Prefix: &prefix,
			}

			// Transitions
			transitions := d.Get(fmt.Sprintf("lifecycle_rules.%d.transition", i)).(*schema.Set).List()
			if len(transitions) > 0 {
				rule.Transitions = make([]*s3.Transition, 0, len(transitions))
				for _, transition := range transitions {
					transitionValue := transition.(map[string]interface{})
					t := &s3.Transition{}
					if val, ok := transitionValue["date"].(string); ok && val != "" {
						date, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", val))
						if err != nil {
							return fmt.Errorf("parsing cos bucket lifecycle transition date(%s) error: %s", val, err.Error())
						}
						t.Date = aws.Time(date)
					} else if val, ok := transitionValue["days"].(int); ok && val >= 0 {
						t.Days = aws.Int64(int64(val))
					}
					if val, ok := transitionValue["storage_class"].(string); ok && val != "" {
						t.StorageClass = aws.String(val)
					}

					rule.Transitions = append(rule.Transitions, t)
				}
			}

			// Expiration
			expirations := d.Get(fmt.Sprintf("lifecycle_rules.%d.expiration", i)).(*schema.Set).List()
			if len(expirations) > 0 {
				expiration := expirations[0].(map[string]interface{})
				e := &s3.LifecycleExpiration{}

				if val, ok := expiration["data"].(string); ok && val != "" {
					date, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", val))
					if err != nil {
						return fmt.Errorf("parsing cos bucket lifecycle expiration data(%s) error: %s", val, err.Error())
					}
					e.Date = aws.Time(date)
				} else if val, ok := expiration["days"].(int); ok && val > 0 {
					e.Days = aws.Int64(int64(val))
				}

				rule.Expiration = e
			}
			rules = append(rules, rule)
		}

		request := s3.PutBucketLifecycleConfigurationInput{
			Bucket: aws.String(bucket),
			LifecycleConfiguration: &s3.BucketLifecycleConfiguration{
				Rules: rules,
			},
		}
		response, err := client.PutBucketLifecycleConfiguration(&request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "put bucket lifecycle", request.String(), err.Error())
			return fmt.Errorf("cos put bucket lifecycle error: %s, bucket: %s", err.Error(), bucket)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, "put bucket lifecycle", request.String(), response.String())
	}

	return nil
}

func resourceTencentCloudCosBucketWebsiteUpdate(ctx context.Context, client *s3.S3, d *schema.ResourceData) error {
	logId := getLogId(ctx)

	bucket := d.Get("bucket").(string)
	website := d.Get("website").([]interface{})

	if len(website) == 0 {
		request := s3.DeleteBucketWebsiteInput{
			Bucket: aws.String(bucket),
		}
		response, err := client.DeleteBucketWebsite(&request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete bucket website", request.String(), err.Error())
			return fmt.Errorf("cos delete bucket website error: %s, bucket: %s", err.Error(), bucket)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, "delete bucket website", request.String(), response.String())
	} else {
		var w map[string]interface{}
		if website[0] != nil {
			w = website[0].(map[string]interface{})
		} else {
			w = make(map[string]interface{})
		}
		var indexDocument, errorDocument string
		if v, ok := w["index_document"]; ok {
			indexDocument = v.(string)
		}
		if v, ok := w["error_document"]; ok {
			errorDocument = v.(string)
		}
		request := s3.PutBucketWebsiteInput{
			Bucket: aws.String(bucket),
			WebsiteConfiguration: &s3.WebsiteConfiguration{
				IndexDocument: &s3.IndexDocument{
					Suffix: aws.String(indexDocument),
				},
				ErrorDocument: &s3.ErrorDocument{
					Key: aws.String(errorDocument),
				},
			},
		}
		response, err := client.PutBucketWebsite(&request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "put bucket website", request.String(), err.Error())
			return fmt.Errorf("cos put bucket website error: %s, bucket: %s", err.Error(), bucket)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, "put bucket website", request.String(), response.String())
	}

	return nil
}

func resourceTencentCloudCosBucketLogStatusUpdate(ctx context.Context, client *s3.S3, d *schema.ResourceData) error {
	logId := getLogId(ctx)

	bucket := d.Id()

	logSwitch := d.Get("log_enable").(bool)
	if logSwitch {
		if d.HasChange("log_target_bucket") || d.HasChange("log_prefix") {
			targetBucket := d.Get("log_target_bucket").(string)
			logPrefix := d.Get("log_prefix").(string)
			//check
			if targetBucket == "" || logPrefix == "" {
				return fmt.Errorf("log_target_bucket and log_prefix should set valid value when log_enable is true")
			}

			//set log target bucket and prefix
			//grant are solved by the tencentcloud_cam_role_attachment resource
			request := &s3.PutBucketLoggingInput{
				Bucket: aws.String(bucket),
				BucketLoggingStatus: &s3.BucketLoggingStatus{
					LoggingEnabled: &s3.LoggingEnabled{
						TargetBucket: aws.String(targetBucket),
						TargetPrefix: aws.String(logPrefix),
					},
				},
			}

			resp, err := client.PutBucketLogging(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, "cos enable log error", request.String(), err.Error())
				return fmt.Errorf("cos enable log error: %s, bucket: %s", err.Error(), bucket)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], resp[%s]\n",
				logId, "cos enable log success", request.String(), resp.String())
		}
	} else {
		targetBucket := d.Get("log_target_bucket").(string)
		logPrefix := d.Get("log_prefix").(string)
		//check
		if targetBucket != "" || logPrefix != "" {
			return fmt.Errorf("log_target_bucket and log_prefix should set null when log_enable is false")
		}
		// set disabled, put empty request
		request := &s3.PutBucketLoggingInput{
			Bucket:              aws.String(bucket),
			BucketLoggingStatus: &s3.BucketLoggingStatus{},
		}

		resp, err := client.PutBucketLogging(request)
		if err != nil {
			return fmt.Errorf("cos disable log error: %s, bucket: %s", err.Error(), bucket)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], resp[%s]\n",
			logId, "cos enable log success", request.String(), resp.String())
	}

	return nil
}

func expirationHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["date"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["days"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}
	return hashcode.String(buf.String())
}

func transitionHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["date"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["days"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}
	if v, ok := m["storage_class"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	return hashcode.String(buf.String())
}
