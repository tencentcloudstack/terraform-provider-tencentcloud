package ses

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSesSendEmail() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSesSendEmailCreate,
		Read:   resourceTencentCloudSesSendEmailRead,
		Delete: resourceTencentCloudSesSendEmailDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"from_email_address": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Sender address. Enter a sender address, for example, noreply@mail.qcloud.com.To display the sender name, enter the address in the following format:Sender.",
			},

			"destination": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Recipient email addresses. You can send an email to up to 50 recipients at a time. Note: the email content will display all recipient addresses. To send one-to-one emails to several recipients, please call the API multiple times to send the emails.",
			},

			"subject": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Email subject.",
			},

			"reply_to_addresses": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Reply-to address. You can enter a valid personal email address that can receive emails. If this parameter is left empty, reply emails will fail to be sent.",
			},

			"cc": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Cc recipient email address, up to 20 people can be copied.",
			},

			"bcc": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The email address of the cc recipient can support up to 20 cc recipients.",
			},

			"template": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Template parameters for template-based sending. As Simple has been disused, Template is required.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Template ID. If you do not have any template, please create one.",
						},
						"template_data": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Variable parameters in the template. Please use json.dump to format the JSON object into a string type.The object is a set of key-value pairs. Each key denotes a variable, which is represented by {{key}}. The key will be replaced with the correspondingvalue (represented by {{value}}) when sending the email.Note: The parameter value cannot be data of a complex type such as HTML.Example: {name:xxx,age:xx}.",
						},
					},
				},
			},

			"attachments": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Parameters for the attachments to be sent. The TencentCloud API supports a request packet of up to 8 MB in size,and the size of the attachment content will increase by 1.5 times after Base64 encoding. Therefore,you need to keep the total size of all attachments below 4 MB. If the entire request exceeds 8 MB,the API will return an error.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Attachment name, which cannot exceed 255 characters. Some attachment types are not supported. For details, see [Attachment Types.](https://www.tencentcloud.com/document/product/1084/42373?has_map=1).",
						},
						"content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Base64-encoded attachment content. You can send attachments of up to 4 MB in the total size.Note: The TencentCloud API supports a request packet of up to 8 MB in size, and the size of the attachmentcontent will increase by 1.5 times after Base64 encoding. Therefore, you need to keep the total size of allattachments below 4 MB. If the entire request exceeds 8 MB, the API will return an error.",
						},
					},
				},
			},

			"unsubscribe": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Unsubscribe link option.  0: Do not add unsubscribe link; 1: English 2: Simplified Chinese;  3: Traditional Chinese; 4: Spanish; 5: French;  6: German; 7: Japanese; 8: Korean;  9: Arabic; 10: Thai.",
			},

			"trigger_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Email triggering type. 0 (default): non-trigger-based, suitable for marketing emails and non-immediate emails;1: trigger-based, suitable for immediate emails such as emails containing verification codes.If the size of an email exceeds a specified value,the system will automatically choose the non-trigger-based type.",
			},
		},
	}
}

func resourceTencentCloudSesSendEmailCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ses_send_email.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request   = ses.NewSendEmailRequest()
		response  = ses.NewSendEmailResponse()
		messageId string
	)
	if v, ok := d.GetOk("from_email_address"); ok {
		request.FromEmailAddress = helper.String(v.(string))
	}

	if v, ok := d.GetOk("destination"); ok {
		destinationSet := v.(*schema.Set).List()
		for i := range destinationSet {
			destination := destinationSet[i].(string)
			request.Destination = append(request.Destination, &destination)
		}
	}

	if v, ok := d.GetOk("subject"); ok {
		request.Subject = helper.String(v.(string))
	}

	if v, ok := d.GetOk("reply_to_addresses"); ok {
		request.ReplyToAddresses = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cc"); ok {
		ccSet := v.(*schema.Set).List()
		for i := range ccSet {
			cc := ccSet[i].(string)
			request.Cc = append(request.Cc, &cc)
		}
	}

	if v, ok := d.GetOk("bcc"); ok {
		bccSet := v.(*schema.Set).List()
		for i := range bccSet {
			bcc := bccSet[i].(string)
			request.Bcc = append(request.Bcc, &bcc)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "template"); ok {
		template := ses.Template{}
		if v, ok := dMap["template_id"]; ok {
			template.TemplateID = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["template_data"]; ok {
			template.TemplateData = helper.String(v.(string))
		}
		request.Template = &template
	}

	if v, ok := d.GetOk("attachments"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			attachment := ses.Attachment{}
			if v, ok := dMap["file_name"]; ok {
				attachment.FileName = helper.String(v.(string))
			}
			if v, ok := dMap["content"]; ok {
				attachment.Content = helper.String(v.(string))
			}
			request.Attachments = append(request.Attachments, &attachment)
		}
	}

	if v, ok := d.GetOk("unsubscribe"); ok {
		request.Unsubscribe = helper.String(v.(string))
	}

	if v, _ := d.GetOk("trigger_type"); v != nil {
		request.TriggerType = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSesClient().SendEmail(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ses sendEmail failed, reason:%+v", logId, err)
		return err
	}

	messageId = *response.Response.MessageId
	d.SetId(messageId)

	return resourceTencentCloudSesSendEmailRead(d, meta)
}

func resourceTencentCloudSesSendEmailRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ses_send_email.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSesSendEmailDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ses_send_email.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
