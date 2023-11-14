/*
Provides a resource to create a kms cancel_key_archive

Example Usage

```hcl
resource "tencentcloud_kms_cancel_key_archive" "cancel_key_archive" {
  key_id = "23e80852-1e38-11e9-b129-5cb9019b4b01"
}
```

Import

kms cancel_key_archive can be imported using the id, e.g.

```
terraform import tencentcloud_kms_cancel_key_archive.cancel_key_archive cancel_key_archive_id
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

func resourceTencentCloudKmsCancelKeyArchive() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKmsCancelKeyArchiveCreate,
		Read:   resourceTencentCloudKmsCancelKeyArchiveRead,
		Delete: resourceTencentCloudKmsCancelKeyArchiveDelete,
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

func resourceTencentCloudKmsCancelKeyArchiveCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_cancel_key_archive.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = kms.NewCancelKeyArchiveRequest()
		response = kms.NewCancelKeyArchiveResponse()
		keyId    string
	)
	if v, ok := d.GetOk("key_id"); ok {
		keyId = v.(string)
		request.KeyId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseKmsClient().CancelKeyArchive(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate kms cancelKeyArchive failed, reason:%+v", logId, err)
		return err
	}

	keyId = *response.Response.KeyId
	d.SetId(keyId)

	return resourceTencentCloudKmsCancelKeyArchiveRead(d, meta)
}

func resourceTencentCloudKmsCancelKeyArchiveRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_cancel_key_archive.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudKmsCancelKeyArchiveDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_cancel_key_archive.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
