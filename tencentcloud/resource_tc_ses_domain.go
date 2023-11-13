/*
Provides a resource to create a ses domain

Example Usage

```hcl
resource "tencentcloud_ses_domain" "domain" {
  email_identity = "mail.qcloud.com"
}
```

Import

ses domain can be imported using the id, e.g.

```
terraform import tencentcloud_ses_domain.domain domain_id
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

func resourceTencentCloudSesDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSesDomainCreate,
		Read:   resourceTencentCloudSesDomainRead,
		Update: resourceTencentCloudSesDomainUpdate,
		Delete: resourceTencentCloudSesDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"email_identity": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Your sender domain. You are advised to use a third-level domain, for example, mail.qcloud.com.",
			},
		},
	}
}

func resourceTencentCloudSesDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_domain.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = ses.NewCreateEmailIdentityRequest()
		response      = ses.NewCreateEmailIdentityResponse()
		emailIdentity string
	)
	if v, ok := d.GetOk("email_identity"); ok {
		emailIdentity = v.(string)
		request.EmailIdentity = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSesClient().CreateEmailIdentity(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ses domain failed, reason:%+v", logId, err)
		return err
	}

	emailIdentity = *response.Response.EmailIdentity
	d.SetId(emailIdentity)

	return resourceTencentCloudSesDomainRead(d, meta)
}

func resourceTencentCloudSesDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_domain.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SesService{client: meta.(*TencentCloudClient).apiV3Conn}

	domainId := d.Id()

	domain, err := service.DescribeSesDomainById(ctx, emailIdentity)
	if err != nil {
		return err
	}

	if domain == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SesDomain` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if domain.EmailIdentity != nil {
		_ = d.Set("email_identity", domain.EmailIdentity)
	}

	return nil
}

func resourceTencentCloudSesDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_domain.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := ses.NewUpdateEmailIdentityRequest()

	domainId := d.Id()

	request.EmailIdentity = &emailIdentity

	immutableArgs := []string{"email_identity"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("email_identity") {
		if v, ok := d.GetOk("email_identity"); ok {
			request.EmailIdentity = helper.String(v.(string))
		}
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
		log.Printf("[CRITAL]%s update ses domain failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSesDomainRead(d, meta)
}

func resourceTencentCloudSesDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_domain.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SesService{client: meta.(*TencentCloudClient).apiV3Conn}
	domainId := d.Id()

	if err := service.DeleteSesDomainById(ctx, emailIdentity); err != nil {
		return err
	}

	return nil
}
