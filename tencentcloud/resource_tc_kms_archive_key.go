/*
Provides a resource to create a kms archive_key

Example Usage

```hcl
resource "tencentcloud_kms_archive_key" "archive_key" {
  key_id = "23e80852-1e38-11e9-b129-5cb9019b4b01"
}
```

Import

kms archive_key can be imported using the id, e.g.

```
terraform import tencentcloud_kms_archive_key.archive_key archive_key_id
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

func resourceTencentCloudKmsArchiveKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKmsArchiveKeyCreate,
		Read:   resourceTencentCloudKmsArchiveKeyRead,
		Delete: resourceTencentCloudKmsArchiveKeyDelete,
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

func resourceTencentCloudKmsArchiveKeyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_archive_key.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = kms.NewArchiveKeyRequest()
		response = kms.NewArchiveKeyResponse()
		keyId    string
	)
	if v, ok := d.GetOk("key_id"); ok {
		keyId = v.(string)
		request.KeyId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseKmsClient().ArchiveKey(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate kms archiveKey failed, reason:%+v", logId, err)
		return err
	}

	keyId = *response.Response.KeyId
	d.SetId(keyId)

	return resourceTencentCloudKmsArchiveKeyRead(d, meta)
}

func resourceTencentCloudKmsArchiveKeyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_archive_key.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudKmsArchiveKeyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_archive_key.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
