/*
Provides a resource to create a ses black_list

~> **NOTE:** Used to remove email addresses from blacklists.

Example Usage

```hcl
resource "tencentcloud_ses_black_list" "black_list" {
  email_address = "terraform-tf@gmail.com"
}
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSesBlackList() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSesBlackListCreate,
		Read:   resourceTencentCloudSesBlackListRead,
		Delete: resourceTencentCloudSesBlackListDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"email_address": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Email addresses to be unblocklisted.",
			},
		},
	}
}

func resourceTencentCloudSesBlackListCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_black_list.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = ses.NewDeleteBlackListRequest()
		emailAddress string
	)
	if v, ok := d.GetOk("email_address"); ok {
		emailAddress = v.(string)
		request.EmailAddressList = append(request.EmailAddressList, helper.String(v.(string)))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSesClient().DeleteBlackList(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ses BlackList failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(emailAddress)

	return resourceTencentCloudSesBlackListRead(d, meta)
}

func resourceTencentCloudSesBlackListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_black_list.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSesBlackListDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_black_list.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
