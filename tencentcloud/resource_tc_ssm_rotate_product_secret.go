/*
Provides a resource to create a ssm rotate_product_secret

Example Usage

```hcl
resource "tencentcloud_ssm_rotate_product_secret" "rotate_product_secret" {
  secret_name = ""
}
```

Import

ssm rotate_product_secret can be imported using the id, e.g.

```
terraform import tencentcloud_ssm_rotate_product_secret.rotate_product_secret rotate_product_secret_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudSsmRotateProductSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSsmRotateProductSecretCreate,
		Read:   resourceTencentCloudSsmRotateProductSecretRead,
		Delete: resourceTencentCloudSsmRotateProductSecretDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"secret_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Secret name.",
			},
		},
	}
}

func resourceTencentCloudSsmRotateProductSecretCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_rotate_product_secret.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = ssm.NewRotateProductSecretRequest()
		response   = ssm.NewRotateProductSecretResponse()
		secretName string
	)
	if v, ok := d.GetOk("secret_name"); ok {
		secretName = v.(string)
		request.SecretName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSsmClient().RotateProductSecret(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ssm rotateProductSecret failed, reason:%+v", logId, err)
		return err
	}

	secretName = *response.Response.secretName
	d.SetId(secretName)

	return resourceTencentCloudSsmRotateProductSecretRead(d, meta)
}

func resourceTencentCloudSsmRotateProductSecretRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_rotate_product_secret.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSsmRotateProductSecretDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_rotate_product_secret.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
