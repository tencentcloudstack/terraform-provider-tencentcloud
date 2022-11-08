/*
Provides a resource to create a ses email_address

Example Usage

```hcl
resource "tencentcloud_ses_email_address" "email_address" {
  email_address     = "aaa@iac-tf.cloud"
  email_sender_name = "aaa"
}

```
Import

ses email_address can be imported using the id, e.g.
```
$ terraform import tencentcloud_ses_email_address.email_address aaa@iac-tf.cloud
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSesEmailAddress() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudSesEmailAddressRead,
		Create: resourceTencentCloudSesEmailAddressCreate,
		Delete: resourceTencentCloudSesEmailAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"email_address": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Your sender address. (You can create up to 10 sender addresses for each domain.).",
			},

			"email_sender_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Sender name.",
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
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create ses email_address failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(emailAddress)
	return resourceTencentCloudSesEmailAddressRead(d, meta)
}

func resourceTencentCloudSesEmailAddressRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_email_address.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SesService{client: meta.(*TencentCloudClient).apiV3Conn}

	emailAddress := d.Id()

	emailSender, err := service.DescribeSesEmailAddress(ctx, emailAddress)

	if err != nil {
		return err
	}

	if emailSender == nil {
		d.SetId("")
		return fmt.Errorf("resource `email_address` %s does not exist", emailAddress)
	}

	if emailSender.EmailAddress != nil {
		_ = d.Set("email_address", emailSender.EmailAddress)
	}

	if emailSender.EmailSenderName != nil {
		_ = d.Set("email_sender_name", emailSender.EmailSenderName)
	}

	return nil
}

func resourceTencentCloudSesEmailAddressDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_email_address.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SesService{client: meta.(*TencentCloudClient).apiV3Conn}

	emailAddressId := d.Id()

	if err := service.DeleteSesEmail_addressById(ctx, emailAddressId); err != nil {
		return err
	}

	return nil
}
