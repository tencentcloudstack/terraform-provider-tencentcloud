/*
Provides a resource to create a kms cloud_resource

Example Usage

```hcl
resource "tencentcloud_kms_cloud_resource" "cloud_resource" {
  key_id = "23e80852-1e38-11e9-b129-5cb9019b4b01"
  product_id = "ssm"
  resource_id = "ins-123456"
}
```

Import

kms cloud_resource can be imported using the id, e.g.

```
terraform import tencentcloud_kms_cloud_resource.cloud_resource cloud_resource_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudKmsCloudResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKmsCloudResourceCreate,
		Read:   resourceTencentCloudKmsCloudResourceRead,
		Delete: resourceTencentCloudKmsCloudResourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"key_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "CMK unique identifier.",
			},

			"product_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "A unique identifier for the cloud product.",
			},

			"resource_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The resource/instance ID of the cloud product.",
			},
		},
	}
}

func resourceTencentCloudKmsCloudResourceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_cloud_resource.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = kms.NewBindCloudResourceRequest()
		response = kms.NewBindCloudResourceResponse()
		keyId    string
	)
	if v, ok := d.GetOk("key_id"); ok {
		keyId = v.(string)
		request.KeyId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product_id"); ok {
		request.ProductId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_id"); ok {
		request.ResourceId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseKmsClient().BindCloudResource(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create kms cloudResource failed, reason:%+v", logId, err)
		return err
	}

	keyId = *response.Response.KeyId
	d.SetId(keyId)

	return resourceTencentCloudKmsCloudResourceRead(d, meta)
}

func resourceTencentCloudKmsCloudResourceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_cloud_resource.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := KmsService{client: meta.(*TencentCloudClient).apiV3Conn}

	cloudResourceId := d.Id()

	cloudResource, err := service.DescribeKmsCloudResourceById(ctx, keyId)
	if err != nil {
		return err
	}

	if cloudResource == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `KmsCloudResource` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if cloudResource.KeyId != nil {
		_ = d.Set("key_id", cloudResource.KeyId)
	}

	if cloudResource.ProductId != nil {
		_ = d.Set("product_id", cloudResource.ProductId)
	}

	if cloudResource.ResourceId != nil {
		_ = d.Set("resource_id", cloudResource.ResourceId)
	}

	return nil
}

func resourceTencentCloudKmsCloudResourceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_cloud_resource.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := KmsService{client: meta.(*TencentCloudClient).apiV3Conn}
	cloudResourceId := d.Id()

	if err := service.DeleteKmsCloudResourceById(ctx, keyId); err != nil {
		return err
	}

	return nil
}
