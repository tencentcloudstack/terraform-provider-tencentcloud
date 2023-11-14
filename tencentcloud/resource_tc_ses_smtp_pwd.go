/*
Provides a resource to create a ses smtp_pwd

Example Usage

```hcl
resource "tencentcloud_ses_smtp_pwd" "smtp_pwd" {
  password = "xX1@#xXXXX"
  email_address = "abc@ef.com"
}
```

Import

ses smtp_pwd can be imported using the id, e.g.

```
terraform import tencentcloud_ses_smtp_pwd.smtp_pwd smtp_pwd_id
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

func resourceTencentCloudSesSmtp_pwd() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSesSmtp_pwdCreate,
		Read:   resourceTencentCloudSesSmtp_pwdRead,
		Update: resourceTencentCloudSesSmtp_pwdUpdate,
		Delete: resourceTencentCloudSesSmtp_pwdDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"password": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "SMTP password. Length limit: 64.",
			},

			"email_address": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Email address. Length limit: 128.",
			},
		},
	}
}

func resourceTencentCloudSesSmtp_pwdCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_smtp_pwd.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = ses.NewUpdateEmailSmtpPassWordRequest()
		response     = ses.NewUpdateEmailSmtpPassWordResponse()
		emailAddress string
	)
	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	if v, ok := d.GetOk("email_address"); ok {
		emailAddress = v.(string)
		request.EmailAddress = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSesClient().UpdateEmailSmtpPassWord(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ses smtp_pwd failed, reason:%+v", logId, err)
		return err
	}

	emailAddress = *response.Response.EmailAddress
	d.SetId(emailAddress)

	return resourceTencentCloudSesSmtp_pwdRead(d, meta)
}

func resourceTencentCloudSesSmtp_pwdRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_smtp_pwd.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SesService{client: meta.(*TencentCloudClient).apiV3Conn}

	smtp_pwdId := d.Id()

	smtp_pwd, err := service.DescribeSesSmtp_pwdById(ctx, emailAddress)
	if err != nil {
		return err
	}

	if smtp_pwd == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SesSmtp_pwd` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if smtp_pwd.Password != nil {
		_ = d.Set("password", smtp_pwd.Password)
	}

	if smtp_pwd.EmailAddress != nil {
		_ = d.Set("email_address", smtp_pwd.EmailAddress)
	}

	return nil
}

func resourceTencentCloudSesSmtp_pwdUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_smtp_pwd.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := ses.NewUpdateEmailSmtpPassWordRequest()

	smtp_pwdId := d.Id()

	request.EmailAddress = &emailAddress

	immutableArgs := []string{"password", "email_address"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("password") {
		if v, ok := d.GetOk("password"); ok {
			request.Password = helper.String(v.(string))
		}
	}

	if d.HasChange("email_address") {
		if v, ok := d.GetOk("email_address"); ok {
			request.EmailAddress = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSesClient().UpdateEmailSmtpPassWord(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update ses smtp_pwd failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSesSmtp_pwdRead(d, meta)
}

func resourceTencentCloudSesSmtp_pwdDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_smtp_pwd.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SesService{client: meta.(*TencentCloudClient).apiV3Conn}
	smtp_pwdId := d.Id()

	if err := service.DeleteSesSmtp_pwdById(ctx, emailAddress); err != nil {
		return err
	}

	return nil
}
