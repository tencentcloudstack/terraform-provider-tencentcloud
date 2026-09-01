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

			"dkim_option": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "DKIM key length. 0: 1024-bit, 1: 2048-bit.",
			},

			"tag_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Tag key.",
			},

			"tag_value": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Tag value.",
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

	if v, ok := d.GetOkExists("dkim_option"); ok {
		request.DKIMOption = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("tag_key"); ok {
		tagKey := v.(string)
		if v, ok := d.GetOk("tag_value"); ok {
			tagList := make([]*ses.TagList, 0, 1)
			tagList = append(tagList, &ses.TagList{
				TagKey:   helper.String(tagKey),
				TagValue: helper.String(v.(string)),
			})
			request.TagList = tagList
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSesClient().CreateEmailIdentity(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("create ses domain failed, response is nil"))
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

	response, err := service.DescribeSesDomain(ctx, emailIdentity)

	if err != nil {
		return err
	}

	if response == nil {
		log.Printf("[CRUD] ses_domain id=%s", d.Id())
		d.SetId("")
		return fmt.Errorf("resource `ses_domain` %s does not exist", emailIdentity)
	}

	_ = d.Set("email_identity", emailIdentity)

	if response.Attributes != nil {
		attributesList := make([]interface{}, 0, len(response.Attributes))
		for _, v := range response.Attributes {
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

	if response.DKIMOption != nil {
		_ = d.Set("dkim_option", *response.DKIMOption)
	}

	if len(response.TagList) > 0 {
		if response.TagList[0].TagKey != nil {
			_ = d.Set("tag_key", *response.TagList[0].TagKey)
		}
		if response.TagList[0].TagValue != nil {
			_ = d.Set("tag_value", *response.TagList[0].TagValue)
		}
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
