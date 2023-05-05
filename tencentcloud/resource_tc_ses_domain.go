/*
Provides a resource to create a ses domain

Example Usage

```hcl
resource "tencentcloud_ses_domain" "domain" {
    email_identity = "iac.cloud"
}

```
Import

ses domain can be imported using the id, e.g.
```
$ terraform import tencentcloud_ses_domain.domain iac.cloud
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSesDomain() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudSesDomainRead,
		Create: resourceTencentCloudSesDomainCreate,
		Delete: resourceTencentCloudSesDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"email_identity": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create ses domain failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(emailIdentity)
	return resourceTencentCloudSesDomainRead(d, meta)
}

func resourceTencentCloudSesDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_domain.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SesService{client: meta.(*TencentCloudClient).apiV3Conn}

	emailIdentity := d.Id()

	domain, err := service.DescribeSesDomain(ctx, emailIdentity)

	if err != nil {
		return err
	}

	if domain == nil {
		d.SetId("")
		return fmt.Errorf("resource `domain` %s does not exist", emailIdentity)
	}

	_ = d.Set("email_identity", emailIdentity)

	return nil
}

func resourceTencentCloudSesDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_domain.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SesService{client: meta.(*TencentCloudClient).apiV3Conn}

	emailIdentity := d.Id()

	if err := service.DeleteSesDomainById(ctx, emailIdentity); err != nil {
		return err
	}

	return nil
}
