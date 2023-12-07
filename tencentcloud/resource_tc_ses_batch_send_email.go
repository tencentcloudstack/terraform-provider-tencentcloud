package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSesBatchSendEmail() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSesBatchSendEmailCreate,
		Read:   resourceTencentCloudSesBatchSendEmailRead,
		Delete: resourceTencentCloudSesBatchSendEmailDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"from_email_address": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Sender address. Enter a sender address such as noreply@mail.qcloud.com. To display the sender name, enter the address in the following format:sender &amp;amp;lt;email address&amp;amp;gt;. For example:Tencent Cloud team &amp;amp;lt;noreply@mail.qcloud.com&amp;amp;gt;.",
			},

			"receiver_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Recipient group ID.",
			},

			"subject": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Email subject.",
			},

			"task_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Task type. 1: immediate; 2: scheduled; 3: recurring.",
			},

			"reply_to_addresses": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Reply-to address. You can enter a valid personal email address that can receive emails. If this parameter is left empty, reply emails will fail to be sent.",
			},

			"template": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Template when emails are sent using a template.",
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
				Description: "Attachment parameters to set when you need to send attachments. This parameter is currently unavailable.",
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

			"cycle_param": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Parameter required for a recurring sending task.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"begin_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Start time of the task.",
						},
						"interval_time": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Task recurrence in hours.",
						},
						"term_cycle": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies whether to end the cycle. This parameter is used to update the task. Valid values: 0: No; 1: Yes.",
						},
					},
				},
			},

			"timed_param": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Parameter required for a scheduled sending task.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"begin_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Start time of a scheduled sending task.",
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

			"ad_location": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Whether to add an ad tag. 0: Add no tag; 1: Add before the subject; 2: Add after the subject.",
			},
		},
	}
}

func resourceTencentCloudSesBatchSendEmailCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_batch_send_email.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = ses.NewBatchSendEmailRequest()
		response = ses.NewBatchSendEmailResponse()
		taskId   uint64
	)
	if v, ok := d.GetOk("from_email_address"); ok {
		request.FromEmailAddress = helper.String(v.(string))
	}

	if v, _ := d.GetOk("receiver_id"); v != nil {
		request.ReceiverId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("subject"); ok {
		request.Subject = helper.String(v.(string))
	}

	if v, _ := d.GetOk("task_type"); v != nil {
		request.TaskType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("reply_to_addresses"); ok {
		request.ReplyToAddresses = helper.String(v.(string))
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

	if dMap, ok := helper.InterfacesHeadMap(d, "cycle_param"); ok {
		cycleEmailParam := ses.CycleEmailParam{}
		if v, ok := dMap["begin_time"]; ok {
			cycleEmailParam.BeginTime = helper.String(v.(string))
		}
		if v, ok := dMap["interval_time"]; ok {
			cycleEmailParam.IntervalTime = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["term_cycle"]; ok {
			cycleEmailParam.TermCycle = helper.IntUint64(v.(int))
		}
		request.CycleParam = &cycleEmailParam
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "timed_param"); ok {
		timedEmailParam := ses.TimedEmailParam{}
		if v, ok := dMap["begin_time"]; ok {
			timedEmailParam.BeginTime = helper.String(v.(string))
		}
		request.TimedParam = &timedEmailParam
	}

	if v, ok := d.GetOk("unsubscribe"); ok {
		request.Unsubscribe = helper.String(v.(string))
	}

	if v, _ := d.GetOk("ad_location"); v != nil {
		request.ADLocation = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSesClient().BatchSendEmail(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ses batchSendEmail failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.TaskId
	d.SetId(helper.UInt64ToStr(taskId))

	return resourceTencentCloudSesBatchSendEmailRead(d, meta)
}

func resourceTencentCloudSesBatchSendEmailRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_batch_send_email.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSesBatchSendEmailDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_batch_send_email.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
