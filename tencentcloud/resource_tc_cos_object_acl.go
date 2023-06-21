/*
Provides a resource to create a cos object_acl

Example Usage

```hcl
resource "tencentcloud_cos_object_acl" "object_acl" {
  bucket = "mycos-1258798060"
  key = "exampleobject"
  x_cos_acl = "private"
}
```

Import

cos object_acl can be imported using the id, e.g.

```
terraform import tencentcloud_cos_object_acl.object_acl object_acl_id
```
*/
package tencentcloud

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cos "github.com/tencentyun/cos-go-sdk-v5"
)

func resourceTencentCloudCosObjectAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCosObjectAclCreate,
		Read:   resourceTencentCloudCosObjectAclRead,
		Update: resourceTencentCloudCosObjectAclUpdate,
		Delete: resourceTencentCloudCosObjectAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Bucket format should be [custom name]-[appid], for example `mycos-1258798060`.",
			},

			"key": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "object name in bucket.",
			},

			"x_cos_acl": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ACL attribute of the object.",
			},
		},
	}
}

func resourceTencentCloudCosObjectAclCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_object_acl.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var (
		bucket string
		object string
	)
	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	}

	if v, ok := d.GetOk("key"); ok {
		object = v.(string)
	}
	opt := &cos.ObjectPutACLOptions{
		Header: &cos.ACLHeaderOptions{
			XCosACL: d.Get("x_cos_acl").(string),
		},
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTencentCosClient(bucket).Object.PutACL(ctx, object, opt)
		if e != nil {
			return retryError(e)
		} else {
			request, _ := xml.Marshal(opt)
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, "PutDomainCertificate", request, result.Response.Body)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cos objectAcl failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(bucket + FILED_SP + object)

	return resourceTencentCloudCosObjectAclRead(d, meta)
}

func resourceTencentCloudCosObjectAclRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_object_acl.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	idItems := strings.Split(d.Id(), FILED_SP)

	_, response, err := meta.(*TencentCloudClient).apiV3Conn.UseTencentCosClient(idItems[0]).Object.GetACL(ctx, idItems[1])
	if err != nil {
		return fmt.Errorf("cos [GetACL] error: %s, bucket: %s, object: %s", err.Error(), idItems[0], idItems[1])
	}
	if response.StatusCode == 404 {
		log.Printf("[WARN] [GetACL] returns %d, %s", 404, err)
		return nil
	}

	_ = d.Set("x_cos_acl", response.Header.Get("x-cos-acl"))

	return nil
}

func resourceTencentCloudCosObjectAclUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_object_acl.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if !d.HasChange("x_cos_acl") {
		return resourceTencentCloudCosObjectAclRead(d, meta)
	}
	idItems := strings.Split(d.Id(), FILED_SP)
	opt := &cos.ObjectPutACLOptions{
		Header: &cos.ACLHeaderOptions{
			XCosACL: d.Get("x_cos_acl").(string),
		},
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTencentCosClient(idItems[0]).Object.PutACL(ctx, idItems[1], opt)
		if e != nil {
			return retryError(e)
		} else {
			request, _ := xml.Marshal(opt)
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, "PutDomainCertificate", request, result.Response.Body)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cos objectAcl failed, reason:%+v", logId, err)
		return err
	}
	return resourceTencentCloudCosObjectAclRead(d, meta)
}

func resourceTencentCloudCosObjectAclDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_object_acl.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
