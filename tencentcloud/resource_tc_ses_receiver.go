/*
Provides a resource to create a ses receiver

Example Usage

```hcl
resource "tencentcloud_ses_receiver" "receiver" {
  receivers_name = "terraform_test"
  desc = "description"

  data {
    email = "abc@abc.com"
  }

  data {
    email = "abcd@abcd.com"
  }
}
```

Create a template with `template_data`
```hcl
resource "tencentcloud_ses_receiver" "receiver" {
  receivers_name = "terraform_test"
  desc = "description"

  data {
    email = "abc@abc.com"
    template_data = "{\"name\":\"xxx\",\"age\":\"xx\"}"
  }

  data {
    email = "abcd@abcd.com"
    template_data = "{\"name\":\"xxx\",\"age\":\"xx\"}"
  }
}
```

Import

ses receiver can be imported using the id, e.g.

```
terraform import tencentcloud_ses_receiver.receiver receiver_id
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSesReceiver() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSesReceiverCreate,
		Read:   resourceTencentCloudSesReceiverRead,
		Delete: resourceTencentCloudSesReceiverDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"receivers_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Recipient group name.",
			},

			"desc": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Recipient group description.",
			},

			"data": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeSet,
				Description: "Recipient email and template parameters in array format. The number of recipients is limited to within 20,000. If there is an object in the `data` list that inputs `template_data`, then other objects are also required.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Recipient email addresses.",
						},
						"template_data": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Variable parameters in the template, please use json.dump to format the JSON object as a string type. The object is a set of key-value pairs, where each key represents a variable in the template, and the variables in the template are represented by {{key}}, and the corresponding values will be replaced with {{value}} when sent.Note: Parameter values cannot be complex data such as HTML. The total length of TemplateData (the entire JSON structure) should be less than 800 bytes.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudSesReceiverCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_receiver.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = ses.NewCreateReceiverRequest()
		response   = ses.NewCreateReceiverResponse()
		receiverId uint64
	)
	if v, ok := d.GetOk("receivers_name"); ok {
		request.ReceiversName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("desc"); ok {
		request.Desc = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSesClient().CreateReceiver(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ses Receiver failed, reason:%+v", logId, err)
		return err
	}

	receiverId = *response.Response.ReceiverId

	if v, ok := d.GetOk("data"); ok {
		datas := v.(*schema.Set).List()
		dataList := make([]*ses.ReceiverInputData, 0, len(datas))
		emilList := make([]*string, 0, len(datas))
		for _, item := range datas {
			var email string
			var templateData string
			dMap := item.(map[string]interface{})
			if v, ok := dMap["email"]; ok {
				email = v.(string)
			}
			if v, ok := dMap["template_data"]; ok {
				templateData = v.(string)
			}

			if templateData != "" {
				receiver := ses.ReceiverInputData{}
				receiver.Email = &email
				receiver.TemplateData = &templateData
				dataList = append(dataList, &receiver)
			} else {
				emilList = append(emilList, &email)
			}
		}

		if len(dataList) > 0 {
			request := ses.NewCreateReceiverDetailWithDataRequest()
			request.ReceiverId = &receiverId
			request.Datas = dataList
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseSesClient().CreateReceiverDetailWithData(request)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s create ses receiverDetail failed, reason:%+v", logId, err)
				return err
			}
		}
		if len(emilList) > 0 {
			request := ses.NewCreateReceiverDetailRequest()
			request.ReceiverId = &receiverId
			request.Emails = emilList
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseSesClient().CreateReceiverDetail(request)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s create ses receiverDetail failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	d.SetId(strconv.Itoa(int(receiverId)))

	return resourceTencentCloudSesReceiverRead(d, meta)
}

func resourceTencentCloudSesReceiverRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_receiver.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SesService{client: meta.(*TencentCloudClient).apiV3Conn}

	receiverId := d.Id()
	receiver, err := service.DescribeSesReceiverById(ctx, receiverId)
	if err != nil {
		return err
	}

	if receiver == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SesReceiver` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if receiver.ReceiversName != nil {
		_ = d.Set("receivers_name", receiver.ReceiversName)
	}

	if receiver.Desc != nil {
		_ = d.Set("desc", receiver.Desc)
	}

	// if receiver.Data != nil {
	// 	dataList := []interface{}{}
	// 	for _, data := range receiver.Data {
	// 		dataMap := map[string]interface{}{}

	// 		if data.Email != nil {
	// 			dataMap["email"] = data.Email
	// 		}

	// 		if data.TemplateData != nil {
	// 			dataMap["template_data"] = data.TemplateData
	// 		}

	// 		dataList = append(dataList, dataMap)
	// 	}

	// 	_ = d.Set("data", dataList)
	// }

	return nil
}

func resourceTencentCloudSesReceiverDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_receiver.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SesService{client: meta.(*TencentCloudClient).apiV3Conn}
	receiverId := d.Id()

	if err := service.DeleteSesReceiverById(ctx, receiverId); err != nil {
		return err
	}

	return nil
}
