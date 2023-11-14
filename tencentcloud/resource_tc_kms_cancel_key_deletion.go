/*
Provides a resource to create a kms cancel_key_deletion

Example Usage

```hcl
resource "tencentcloud_kms_cancel_key_deletion" "cancel_key_deletion" {
  key_id = "23e80852-1e38-11e9-b129-5cb9019b4b01"
}
```

Import

kms cancel_key_deletion can be imported using the id, e.g.

```
terraform import tencentcloud_kms_cancel_key_deletion.cancel_key_deletion cancel_key_deletion_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudKmsCancelKeyDeletion() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKmsCancelKeyDeletionCreate,
		Read:   resourceTencentCloudKmsCancelKeyDeletionRead,
		Delete: resourceTencentCloudKmsCancelKeyDeletionDelete,
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
		},
	}
}

func resourceTencentCloudKmsCancelKeyDeletionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_cancel_key_deletion.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = kms.NewCancelKeyDeletionRequest()
		response = kms.NewCancelKeyDeletionResponse()
		keyId    string
	)
	if v, ok := d.GetOk("key_id"); ok {
		keyId = v.(string)
		request.KeyId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseKmsClient().CancelKeyDeletion(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate kms cancelKeyDeletion failed, reason:%+v", logId, err)
		return err
	}

	keyId = *response.Response.KeyId
	d.SetId(keyId)

	return resourceTencentCloudKmsCancelKeyDeletionRead(d, meta)
}

func resourceTencentCloudKmsCancelKeyDeletionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_cancel_key_deletion.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudKmsCancelKeyDeletionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_cancel_key_deletion.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
