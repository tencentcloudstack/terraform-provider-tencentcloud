package cam

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	camv20190116 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCamMessageReceiver() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamMessageReceiverCreate,
		Read:   resourceTencentCloudCamMessageReceiverRead,
		Delete: resourceTencentCloudCamMessageReceiverDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Username of the message recipient.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Recipient's notes.",
			},

			"country_code": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The international area code for mobile phone numbers is 86 for domestic areas.",
			},

			"phone_number": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Mobile phone number, for example: 132****2492.",
			},

			"email": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Email address, for example: 57*****@qq.com.",
			},

			// computed
			"uid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "UID.",
			},

			"is_receiver_owner": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether it is the primary contact person.",
			},

			"phone_flag": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether the mobile phone number is verified.",
			},

			"email_flag": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether the email is verified.",
			},

			"wechat_flag": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether WeChat is allowed to receive notifications.",
			},

			"uin": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Account uin.",
			},
		},
	}
}

func resourceTencentCloudCamMessageReceiverCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_message_receiver.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = camv20190116.NewCreateMessageReceiverRequest()
		name    string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
		name = v.(string)
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("country_code"); ok {
		request.CountryCode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("phone_number"); ok {
		request.PhoneNumber = helper.String(v.(string))
	}

	if v, ok := d.GetOk("email"); ok {
		request.Email = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCamV20190116Client().CreateMessageReceiverWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create cam message receiver failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create cam message receiver failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(name)
	return resourceTencentCloudCamMessageReceiverRead(d, meta)
}

func resourceTencentCloudCamMessageReceiverRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_message_receiver.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = CamService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		name    = d.Id()
	)

	respData, err := service.DescribeCamMessageReceiverById(ctx, name)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_cam_message_receiver` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.Remark != nil {
		_ = d.Set("remark", respData.Remark)
	}

	if respData.CountryCode != nil {
		_ = d.Set("country_code", respData.CountryCode)
	}

	if respData.PhoneNumber != nil {
		_ = d.Set("phone_number", respData.PhoneNumber)
	}

	if respData.Email != nil {
		_ = d.Set("email", respData.Email)
	}

	if respData.Uid != nil {
		_ = d.Set("uid", respData.Uid)
	}

	if respData.IsReceiverOwner != nil {
		_ = d.Set("is_receiver_owner", respData.IsReceiverOwner)
	}

	if respData.PhoneFlag != nil {
		_ = d.Set("phone_flag", respData.PhoneFlag)
	}

	if respData.EmailFlag != nil {
		_ = d.Set("email_flag", respData.EmailFlag)
	}

	if respData.WechatFlag != nil {
		_ = d.Set("wechat_flag", respData.WechatFlag)
	}

	if respData.Uin != nil {
		_ = d.Set("uin", respData.Uin)
	}

	return nil
}

func resourceTencentCloudCamMessageReceiverDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_message_receiver.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = camv20190116.NewDeleteMessageReceiverRequest()
		name    = d.Id()
	)

	request.Name = &name
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCamV20190116Client().DeleteMessageReceiverWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete cam message receiver failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
