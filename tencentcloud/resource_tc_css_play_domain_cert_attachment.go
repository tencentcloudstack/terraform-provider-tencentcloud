/*
Provides a resource to create a css play_domain_cert_attachment.
This resource is used for binding the play domain and specified certification together.

Example Usage

```hcl
resource "tencentcloud_css_play_domain_cert_attachment" "play_domain_cert_attachment" {
  cloud_cert_id = &lt;nil&gt;

domain_name = ""
status =

  }
```

Import

css play_domain_cert_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_css_play_domain_cert_attachment.play_domain_cert_attachment domainName#cloudCertId
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCssPlayDomainCertAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssPlayDomainCertAttachmentCreate,
		Read:   resourceTencentCloudCssPlayDomainCertAttachmentRead,
		Delete: resourceTencentCloudCssPlayDomainCertAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cloud_cert_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Tencent cloud ssl certificate Id. Refer to `tencentcloud_ssl_certificate` to create or obtain the resource ID.",
			},
			"domain_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "domain name.",
			},
			"status": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Whether to enable the https rule for the domain name. 1: enable, 0: disabled, -1: remain unchanged.",
			},
			"certificate_alias": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "certificate remarks. Synonymous with CertName.",
			},
			"cert_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "certificate type. 0: Self-owned certificate, 1: Tencent Cloud ssl managed certificate.",
			},
			"cert_expire_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "certificate expiration time.",
			},
			"cert_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "certificate ID.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the rule was last updated.",
			},
		},
	}
}

func resourceTencentCloudCssPlayDomainCertAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_play_domain_cert_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request     = css.NewModifyLiveDomainCertBindingsRequest()
		response    = css.NewModifyLiveDomainCertBindingsResponse()
		cloudCertId string
		domainName  string
	)

	if v, ok := d.GetOk("cloud_cert_id"); ok {
		cloudCertId = v.(string)
		request.CloudCertId = helper.String(cloudCertId)
	}

	infos := []*css.LiveCertDomainInfo{}
	if v, ok := d.GetOk("domain_name"); ok {
		domainName = v.(string)
		infos[0].DomainName = helper.String(domainName)
	}

	if v, _ := d.GetOk("status"); v != nil {
		infos[0].Status = helper.IntInt64(v.(int))
	}

	request.DomainInfos = infos

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().ModifyLiveDomainCertBindings(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create css playDomainCertAttachment failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.Errors != nil {
		return fmt.Errorf("[CRITAL]%s create css playDomainCertAttachment failed, reason:%+v", logId, response.Response.Errors)
	}

	d.SetId(strings.Join([]string{domainName, cloudCertId}, FILED_SP))

	return resourceTencentCloudCssPlayDomainCertAttachmentRead(d, meta)
}

func resourceTencentCloudCssPlayDomainCertAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_play_domain_cert_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	domainName := idSplit[0]
	cloudCertId := idSplit[1]

	playDomainCertAttachment, err := service.DescribeCssPlayDomainCertAttachmentById(ctx, domainName, cloudCertId)
	if err != nil {
		return err
	}

	if playDomainCertAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CssPlayDomainCertAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if playDomainCertAttachment.CloudCertId != nil {
		_ = d.Set("cloud_cert_id", playDomainCertAttachment.CloudCertId)
	}

	if playDomainCertAttachment.DomainName != nil {
		_ = d.Set("domain_name", playDomainCertAttachment.DomainName)
	}

	if playDomainCertAttachment.Status != nil {
		_ = d.Set("status", playDomainCertAttachment.Status)
	}

	if playDomainCertAttachment.CertificateAlias != nil {
		_ = d.Set("certificate_alias", playDomainCertAttachment.CertificateAlias)
	}

	if playDomainCertAttachment.CertType != nil {
		_ = d.Set("cert_type", playDomainCertAttachment.CertType)
	}

	if playDomainCertAttachment.CertExpireTime != nil {
		_ = d.Set("cert_expire_time", playDomainCertAttachment.CertExpireTime)
	}

	if playDomainCertAttachment.CertId != nil {
		_ = d.Set("cert_id", playDomainCertAttachment.CertId)
	}

	if playDomainCertAttachment.UpdateTime != nil {
		_ = d.Set("update_time", playDomainCertAttachment.UpdateTime)
	}

	return nil
}

func resourceTencentCloudCssPlayDomainCertAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_play_domain_cert_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	domainName := idSplit[0]

	if err := service.DeleteCssPlayDomainCertAttachmentById(ctx, domainName); err != nil {
		return err
	}

	return nil
}
