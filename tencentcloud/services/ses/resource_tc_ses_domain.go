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

func ResourceTencentCloudSesDomain() *schema.Resource {
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

			"attributes": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "DNS configuration details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record Type CNAME | A | TXT | MX.",
						},
						"send_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain name.",
						},
						"expected_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Values that need to be configured.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudSesDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ses_domain.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request       = ses.NewCreateEmailIdentityRequest()
		emailIdentity string
	)

	if v, ok := d.GetOk("email_identity"); ok {
		emailIdentity = v.(string)
		request.EmailIdentity = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSesClient().CreateEmailIdentity(request)
		if e != nil {
			return tccommon.RetryError(e)
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
	defer tccommon.LogElapsed("resource.tencentcloud_ses_domain.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := SesService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	emailIdentity := d.Id()

	attributes, err := service.DescribeSesDomain(ctx, emailIdentity)

	if err != nil {
		return err
	}

	if attributes == nil {
		d.SetId("")
		return fmt.Errorf("resource `domain` %s does not exist", emailIdentity)
	}

	_ = d.Set("email_identity", emailIdentity)

	if attributes != nil {
		attributesList := make([]interface{}, 0, len(attributes))
		for _, v := range attributes {
			attributesMap := map[string]interface{}{}

			if v.Type != nil {
				attributesMap["type"] = v.Type
			}

			if v.SendDomain != nil {
				attributesMap["send_domain"] = v.SendDomain
			}

			if v.ExpectedValue != nil {
				attributesMap["expected_value"] = v.ExpectedValue
			}

			attributesList = append(attributesList, attributesMap)
		}

		_ = d.Set("attributes", attributesList)
	}

	return nil
}

func resourceTencentCloudSesDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ses_domain.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := SesService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	emailIdentity := d.Id()

	if err := service.DeleteSesDomainById(ctx, emailIdentity); err != nil {
		return err
	}

	return nil
}
