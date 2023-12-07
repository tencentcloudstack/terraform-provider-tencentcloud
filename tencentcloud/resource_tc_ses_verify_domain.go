package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSesVerifyDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSesVerifyDomainCreate,
		Read:   resourceTencentCloudSesVerifyDomainRead,
		Delete: resourceTencentCloudSesVerifyDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"email_identity": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Domain name requested for verification.",
			},
		},
	}
}

func resourceTencentCloudSesVerifyDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_verify_domain.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request       = ses.NewUpdateEmailIdentityRequest()
		emailIdentity string
	)
	if v, ok := d.GetOk("email_identity"); ok {
		emailIdentity = v.(string)
		request.EmailIdentity = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSesClient().UpdateEmailIdentity(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ses verifyDomain failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(emailIdentity)

	service := SesService{client: meta.(*TencentCloudClient).apiV3Conn}
	if err := service.CheckEmailIdentityById(ctx, emailIdentity); err != nil {
		return err
	}

	return resourceTencentCloudSesVerifyDomainRead(d, meta)
}

func resourceTencentCloudSesVerifyDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_verify_domain.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSesVerifyDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_verify_domain.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
