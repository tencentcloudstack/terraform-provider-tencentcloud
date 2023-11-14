/*
Provides a resource to create a ses black_list

Example Usage

```hcl
resource "tencentcloud_ses_black_list" "black_list" {
  email_address_list =
}
```

Import

ses black_list can be imported using the id, e.g.

```
terraform import tencentcloud_ses_black_list.black_list black_list_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
	"log"
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
			"email_address_list": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of email addresses to be unblocklisted. Enter at least one address.",
			},
		},
	}
}

func resourceTencentCloudSesBlackListCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_black_list.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = ses.NewDeleteBlackListRequest()
		response  = ses.NewDeleteBlackListResponse()
		blackList string
	)
	if v, ok := d.GetOk("email_address_list"); ok {
		emailAddressListSet := v.(*schema.Set).List()
		for i := range emailAddressListSet {
			emailAddressList := emailAddressListSet[i].(string)
			request.EmailAddressList = append(request.EmailAddressList, &emailAddressList)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSesClient().DeleteBlackList(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ses BlackList failed, reason:%+v", logId, err)
		return err
	}

	blackList = *response.Response.BlackList
	d.SetId(blackList)

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
