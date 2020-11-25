/*
Provides a COS resource to create a COS bucket policy and set its attributes.

Example Usage

```hcl
resource "tencentcloud_cos_bucket_policy" "cos_policy" {
  bucket = "mycos-1258798060"

  policy = <<EOF
{
  "version": "2.0",
  "Statement": [
    {
      "Principal": {
        "qcs": [
          "qcs::cam::uin/<your-account-id>:uin/<your-account-id>"
        ]
      },
      "Action": [
        "name/cos:DeleteBucket",
        "name/cos:PutBucketACL"
      ],
      "Effect": "allow",
      "Resource": [
        "qcs::cos:<bucket region>:uid/<your-account-id>:<bucket name>/*"
      ]
    }
  ]
}
EOF
}
```

Import

COS bucket policy can be imported, e.g.

```
$ terraform import tencentcloud_cos_bucket_policy.bucket bucket-name
```
*/
package tencentcloud

import (
	"context"
	"encoding/json"
	"log"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTencentCloudCosBucketPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCosBucketPolicyCreate,
		Read:   resourceTencentCloudCosBucketPolicyRead,
		Update: resourceTencentCloudCosBucketPolicyUpdate,
		Delete: resourceTencentCloudCosBucketPolicyDelete,
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
			"policy": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, olds, news string, d *schema.ResourceData) bool {
					var oldJson interface{}
					err := json.Unmarshal([]byte(olds), &oldJson)
					if err != nil {
						return olds == news
					}
					var newJson interface{}
					err = json.Unmarshal([]byte(news), &newJson)
					if err != nil {
						return olds == news
					}
					flag := reflect.DeepEqual(oldJson, newJson)
					return flag
				},
				Description: "The text of the policy. this field is required. the syntax refers to https://cloud.tencent.com/document/product/436/18023.",
			},
		},
	}
}

func resourceTencentCloudCosBucketPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_bucket_policy.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	bucket := d.Get("bucket").(string)
	policy := d.Get("policy").(string)

	cosService := CosService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := cosService.PutBucketPolicy(ctx, bucket, policy)
	if err != nil {
		return err
	}
	d.SetId(bucket)

	return resourceTencentCloudCosBucketPolicyRead(d, meta)
}

func resourceTencentCloudCosBucketPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_bucket_policy.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	bucket := d.Id()
	cosService := CosService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result string
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		policy, e := cosService.DescribePolicyByBucket(ctx, bucket)
		if e != nil {
			return retryError(e)
		}
		result = policy
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read cos bucket policy failed, reason:%s\n", logId, err.Error())
		return err
	}
	result, err = removeSid(result)
	if err != nil {
		log.Printf("[CRITAL]%s read cos bucket policy failed, reason:%s\n", logId, err.Error())
		return err
	}
	if result == "" {
		d.SetId("")
		return nil
	}
	_ = d.Set("policy", result)
	_ = d.Set("bucket", bucket)

	return nil
}

func resourceTencentCloudCosBucketPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_bucket_policy.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cosService := CosService{client: meta.(*TencentCloudClient).apiV3Conn}
	bucket := d.Id()

	if d.HasChange("policy") {
		policy := d.Get("policy").(string)
		err := cosService.PutBucketPolicy(ctx, bucket, policy)
		if err != nil {
			return err
		}
	}

	// wait for update cache
	// if not, the data may be outdated.
	time.Sleep(3 * time.Second)

	return resourceTencentCloudCosBucketPolicyRead(d, meta)
}

func resourceTencentCloudCosBucketPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_bucket_policy.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	bucket := d.Id()
	cosService := CosService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := cosService.DeleteBucketPolicy(ctx, bucket)
	if err != nil {
		return err
	}

	// wait for update cache
	// if not, head bucket may be successful
	time.Sleep(3 * time.Second)

	return nil
}

//In the returned JSON, the SDK automatically adds the Sid, which needs to be removed
func removeSid(v string) (result string, err error) {
	m := make(map[string]interface{})
	err = json.Unmarshal([]byte(v), &m)
	if err != nil {
		return
	}
	var stateMend []interface{}
	if v, ok := m["Statement"]; ok {
		stateMend = v.([]interface{})
	} else if v, ok := m["statement"]; ok {
		stateMend = v.([]interface{})
	}
	for index, v := range stateMend {
		mp := v.(map[string]interface{})
		delete(mp, "Sid")
		stateMend[index] = mp
	}
	if _, ok := m["Statement"]; ok {
		m["Statement"] = stateMend
	} else if _, ok := m["statement"]; ok {
		m["statement"] = stateMend
	}
	s, err := json.Marshal(m)
	return string(s), err
}
