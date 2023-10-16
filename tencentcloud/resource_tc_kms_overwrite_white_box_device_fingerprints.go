/*
Provides a resource to create a kms overwrite_white_box_device_fingerprints

Example Usage

```hcl
resource "tencentcloud_kms_overwrite_white_box_device_fingerprints" "example" {
  key_id = "23e80852-1e38-11e9-b129-5cb9019b4b01"
}
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudKmsOverwriteWhiteBoxDeviceFingerprints() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKmsOverwriteWhiteBoxDeviceFingerprintsCreate,
		Read:   resourceTencentCloudKmsOverwriteWhiteBoxDeviceFingerprintsRead,
		Delete: resourceTencentCloudKmsOverwriteWhiteBoxDeviceFingerprintsDelete,

		Schema: map[string]*schema.Schema{
			"key_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "CMK unique identifier.",
			},
			"device_fingerprints": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Device fingerprint list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identity": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "identity.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudKmsOverwriteWhiteBoxDeviceFingerprintsCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_overwrite_white_box_device_fingerprints.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		request = kms.NewOverwriteWhiteBoxDeviceFingerprintsRequest()
		keyId   string
	)

	if v, ok := d.GetOk("key_id"); ok {
		request.KeyId = helper.String(v.(string))
		keyId = v.(string)
	}

	if v, ok := d.GetOk("deviceFingerprints"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			deviceFingerprint := kms.DeviceFingerprint{}
			if v, ok := dMap["identity"]; ok {
				deviceFingerprint.Identity = helper.String(v.(string))
			}

			if v, ok := dMap["description"]; ok {
				deviceFingerprint.Description = helper.String(v.(string))
			}

			request.DeviceFingerprints = append(request.DeviceFingerprints, &deviceFingerprint)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseKmsClient().OverwriteWhiteBoxDeviceFingerprints(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate kms overwriteWhiteBoxDeviceFingerprints failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(keyId)

	return resourceTencentCloudKmsOverwriteWhiteBoxDeviceFingerprintsRead(d, meta)
}

func resourceTencentCloudKmsOverwriteWhiteBoxDeviceFingerprintsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_overwrite_white_box_device_fingerprints.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudKmsOverwriteWhiteBoxDeviceFingerprintsDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kms_overwrite_white_box_device_fingerprints.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
