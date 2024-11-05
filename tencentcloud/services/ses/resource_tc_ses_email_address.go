package ses

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSesEmailAddress() *schema.Resource {
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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Your sender address(You can create up to 10 sender addresses for each domain).",
			},

			"email_sender_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Sender name.",
			},

			"smtp_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Password for SMTP, Length limit 64.",
			},
		},
	}
}

func resourceTencentCloudSesEmailAddressCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ses_email_address.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSesClient().CreateEmailAddress(request)
		if e != nil {
			return tccommon.RetryError(e)
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

	// set smtp pwd
	if v, ok := d.GetOk("smtp_password"); ok {
		pwdRequest := ses.NewUpdateEmailSmtpPassWordRequest()
		pwdRequest.EmailAddress = &emailAddress
		pwdRequest.Password = helper.String(v.(string))
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSesClient().UpdateEmailSmtpPassWord(pwdRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, pwdRequest.GetAction(), pwdRequest.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update ses email address smtp password failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudSesEmailAddressRead(d, meta)
}

func resourceTencentCloudSesEmailAddressRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ses_email_address.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := SesService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

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

func resourceTencentCloudSesEmailAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ses_email_address.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		emailAddress = d.Id()
	)

	if d.HasChange("smtp_password") {
		if v, ok := d.GetOk("smtp_password"); ok {
			pwdRequest := ses.NewUpdateEmailSmtpPassWordRequest()
			pwdRequest.EmailAddress = &emailAddress
			pwdRequest.Password = helper.String(v.(string))
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSesClient().UpdateEmailSmtpPassWord(pwdRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, pwdRequest.GetAction(), pwdRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s update ses email address smtp password failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	return resourceTencentCloudSesEmailAddressRead(d, meta)
}

func resourceTencentCloudSesEmailAddressDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ses_email_address.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := SesService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	emailAddressId := d.Id()

	if err := service.DeleteSesEmail_addressById(ctx, emailAddressId); err != nil {
		return err
	}

	return nil
}
