/*
Provides a resource to create a dbbrain verify_user_account

Example Usage

```hcl
resource "tencentcloud_dbbrain_verify_user_account" "verify_user_account" {
  instance_id = ""
  user = ""
  password = ""
  product = ""
}
```

Import

dbbrain verify_user_account can be imported using the id, e.g.

```
terraform import tencentcloud_dbbrain_verify_user_account.verify_user_account verify_user_account_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudDbbrainVerifyUserAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDbbrainVerifyUserAccountCreate,
		Read:   resourceTencentCloudDbbrainVerifyUserAccountRead,
		Delete: resourceTencentCloudDbbrainVerifyUserAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"user": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Database account name.",
			},

			"password": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Database account password.",
			},

			"product": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values： &amp;quot;mysql&amp;quot; - cloud database MySQL; &amp;quot;cynosdb&amp;quot; - cloud database TDSQL-C for MySQL, the default is &amp;quot;mysql&amp;quot;.",
			},
		},
	}
}

func resourceTencentCloudDbbrainVerifyUserAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dbbrain_verify_user_account.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = dbbrain.NewVerifyUserAccountRequest()
		response   = dbbrain.NewVerifyUserAccountResponse()
		instanceId uint64
	)
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user"); ok {
		request.User = helper.String(v.(string))
	}

	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		request.Product = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDbbrainClient().VerifyUserAccount(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Println("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Println("[CRITAL]%s operate dbbrain verifyUserAccount failed, reason:%+v", logId, err)
		return nil
	}

	instanceId = *response.Response.InstanceId
	d.SetId(helper.UInt64ToStr(instanceId))

	return resourceTencentCloudDbbrainVerifyUserAccountRead(d, meta)
}

func resourceTencentCloudDbbrainVerifyUserAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dbbrain_verify_user_account.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDbbrainVerifyUserAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dbbrain_verify_user_account.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
