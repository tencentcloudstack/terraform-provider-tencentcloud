package tencentcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCssAuthenticateDomainOwnerOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssAuthenticateDomainOwnerOperationCreate,
		Read:   resourceTencentCloudCssAuthenticateDomainOwnerOperationRead,
		Delete: resourceTencentCloudCssAuthenticateDomainOwnerOperationDelete,
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The domain name to verify.",
			},

			"verify_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Default:     CSS_VERIFY_TYPE_DB_CHECK,
				Description: "Authentication type. Possible values:`dnsCheck`: Immediately verify whether the resolution record of the configured dns is consistent with the content to be verified, and save the record if successful.`fileCheck`: Immediately verify whether the web file is consistent with the content to be verified, and save the record if successful.`dbCheck`: Check if authentication has been successful.",
			},
		},
	}
}

func resourceTencentCloudCssAuthenticateDomainOwnerOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_css_authenticate_domain_owner_operation.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = css.NewAuthenticateDomainOwnerRequest()
		response = css.NewAuthenticateDomainOwnerResponse()
		name     string
		status   int64
	)
	if v, ok := d.GetOk("domain_name"); ok {
		request.DomainName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("verify_type"); ok {
		request.VerifyType = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().AuthenticateDomainOwner(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate css authenticateDomainOwnerOperation failed, reason:%+v", logId, err)
		return nil
	}

	if response != nil {
		name = *response.Response.MainDomain
		status = *response.Response.Status
		if status < CSS_OWERSHIP_VIRIFIED {
			return fmt.Errorf("Authenticate domain failed. Please check your domain and retry. MainDomain:[%s] Content:[%s], ", name, *response.Response.Content)
		}
	}

	d.SetId(name)

	return nil
}

func resourceTencentCloudCssAuthenticateDomainOwnerOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_authenticate_domain_owner_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCssAuthenticateDomainOwnerOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_authenticate_domain_owner_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
