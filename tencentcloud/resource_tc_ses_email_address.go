/*
Provides a resource to create a ses email_address

Example Usage

```hcl
resource "tencentcloud_ses_email_address" "email_address" {
  email_address = &lt;nil&gt;
  email_sender_name = &lt;nil&gt;
}
```

Import

ses email_address can be imported using the id, e.g.

```
terraform import tencentcloud_ses_email_address.email_address email_address_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudSesEmailAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSesEmailAddressCreate,
		Read:   resourceTencentCloudSesEmailAddressRead,
		Update: resourceTencentCloudSesEmailAddressUpdate,
		Delete: resourceTencentCloudSesEmailAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"email_address": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Your sender address. (You can create up to 10 sender addresses for each domain.).",
			},

			"email_sender_name": {
				Optional: true,
				Type:     schema.TypeString,
				Description: "	Sender name.",
			},
		},
	}
}

func resourceTencentCloudSesEmailAddressCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_email_address.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = ses.NewCreateEmailAddressRequest()
		response     = ses.NewCreateEmailAddressResponse()
		emailAddress string
	)
	if v, ok := d.GetOk("email_address"); ok {
		emailAddress = v.(string)
		request.EmailAddress = helper.String(v.(string))
	}

	if v, ok := d.GetOk("email_sender_name"); ok {
		request.EmailSenderName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSesClient().CreateEmailAddress(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ses emailAddress failed, reason:%+v", logId, err)
		return err
	}

	emailAddress = *response.Response.EmailAddress
	d.SetId(emailAddress)

	return resourceTencentCloudSesEmailAddressRead(d, meta)
}

func resourceTencentCloudSesEmailAddressRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_email_address.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SesService{client: meta.(*TencentCloudClient).apiV3Conn}

	emailAddressId := d.Id()

	emailAddress, err := service.DescribeSesEmailAddressById(ctx, emailAddress)
	if err != nil {
		return err
	}

	if emailAddress == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SesEmailAddress` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if emailAddress.EmailAddress != nil {
		_ = d.Set("email_address", emailAddress.EmailAddress)
	}

	if emailAddress.EmailSenderName != nil {
		_ = d.Set("email_sender_name", emailAddress.EmailSenderName)
	}

	return nil
}

func resourceTencentCloudSesEmailAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_email_address.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"email_address", "email_sender_name"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudSesEmailAddressRead(d, meta)
}

func resourceTencentCloudSesEmailAddressDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_email_address.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SesService{client: meta.(*TencentCloudClient).apiV3Conn}
	emailAddressId := d.Id()

	if err := service.DeleteSesEmailAddressById(ctx, emailAddress); err != nil {
		return err
	}

	return nil
}
