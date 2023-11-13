/*
Provides a resource to create a ses receiver_detail

Example Usage

```hcl
resource "tencentcloud_ses_receiver_detail" "receiver_detail" {
  receiver_id = 123
  datas {
		email = "abc@ef.com"
		template_data = "{&quot;name&quot;:&quot;xxx&quot;,&quot;age&quot;:&quot;xx&quot;}"

  }
  emails =
}
```

Import

ses receiver_detail can be imported using the id, e.g.

```
terraform import tencentcloud_ses_receiver_detail.receiver_detail receiver_detail_id
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

func resourceTencentCloudSesReceiverDetail() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSesReceiverDetailCreate,
		Read:   resourceTencentCloudSesReceiverDetailRead,
		Update: resourceTencentCloudSesReceiverDetailUpdate,
		Delete: resourceTencentCloudSesReceiverDetailDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"receiver_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Recipient group ID.",
			},

			"datas": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Recipient email and template parameters in array format. The number of recipients is limited to within 20,000.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Recipient email addresses.",
						},
						"template_data": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Variable parameters in the template, please use json.dump to format the JSON object as a string type. The object is a set of key-value pairs, where each key represents a variable in the template, and the variables in the template are represented by {{key}}, and the corresponding values will be replaced with {{value}} when sent.Note: Parameter values cannot be complex data such as HTML. The total length of TemplateData (the entire JSON structure) should be less than 800 bytes.",
						},
					},
				},
			},

			"emails": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Email address.",
			},
		},
	}
}

func resourceTencentCloudSesReceiverDetailCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_receiver_detail.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		createReceiverDetailWithDataRequest  = ses.NewCreateReceiverDetailWithDataRequest()
		createReceiverDetailWithDataResponse = ses.NewCreateReceiverDetailWithDataResponse()
	)
	if v, ok := d.GetOkExists("receiver_id"); ok {
		request.ReceiverId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("datas"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			receiverInputData := ses.ReceiverInputData{}
			if v, ok := dMap["email"]; ok {
				receiverInputData.Email = helper.String(v.(string))
			}
			if v, ok := dMap["template_data"]; ok {
				receiverInputData.TemplateData = helper.String(v.(string))
			}
			request.Datas = append(request.Datas, &receiverInputData)
		}
	}

	if v, ok := d.GetOk("emails"); ok {
		emailsSet := v.(*schema.Set).List()
		for i := range emailsSet {
			emails := emailsSet[i].(string)
			request.Emails = append(request.Emails, &emails)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSesClient().CreateReceiverDetailWithData(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ses receiverDetail failed, reason:%+v", logId, err)
		return err
	}

	receiverDetail = *response.Response.ReceiverDetail
	d.SetId(receiverDetail)

	return resourceTencentCloudSesReceiverDetailRead(d, meta)
}

func resourceTencentCloudSesReceiverDetailRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_receiver_detail.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SesService{client: meta.(*TencentCloudClient).apiV3Conn}

	receiverDetailId := d.Id()

	receiverDetail, err := service.DescribeSesReceiverDetailById(ctx, receiverDetail)
	if err != nil {
		return err
	}

	if receiverDetail == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SesReceiverDetail` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if receiverDetail.ReceiverId != nil {
		_ = d.Set("receiver_id", receiverDetail.ReceiverId)
	}

	if receiverDetail.Datas != nil {
		datasList := []interface{}{}
		for _, datas := range receiverDetail.Datas {
			datasMap := map[string]interface{}{}

			if receiverDetail.Datas.Email != nil {
				datasMap["email"] = receiverDetail.Datas.Email
			}

			if receiverDetail.Datas.TemplateData != nil {
				datasMap["template_data"] = receiverDetail.Datas.TemplateData
			}

			datasList = append(datasList, datasMap)
		}

		_ = d.Set("datas", datasList)

	}

	if receiverDetail.Emails != nil {
		_ = d.Set("emails", receiverDetail.Emails)
	}

	return nil
}

func resourceTencentCloudSesReceiverDetailUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_receiver_detail.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"receiver_id", "datas", "emails"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudSesReceiverDetailRead(d, meta)
}

func resourceTencentCloudSesReceiverDetailDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_receiver_detail.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SesService{client: meta.(*TencentCloudClient).apiV3Conn}
	receiverDetailId := d.Id()

	if err := service.DeleteSesReceiverDetailById(ctx, receiverDetail); err != nil {
		return err
	}

	return nil
}
