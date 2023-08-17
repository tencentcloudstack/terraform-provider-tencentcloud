/*
Provides a resource to create a eb put_events

Example Usage

```hcl
resource "tencentcloud_eb_event_bus" "foo" {
  event_bus_name = "tf-event_bus"
  description    = "event bus desc"
  enable_store   = false
  save_days      = 1
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_eb_put_events" "put_events" {
  event_list {
    source = "ckafka.cloud.tencent"
    data = jsonencode(
      {
        "topic" : "test-topic",
        "Partition" : 1,
        "offset" : 37,
        "msgKey" : "test",
        "msgBody" : "Hello from Ckafka again!"
      }
    )
    type    = "connector:ckafka"
    subject = "qcs::ckafka:ap-guangzhou:uin/1250000000:ckafkaId/uin/1250000000/ckafka-123456"
    time    = 1691572461939

  }
  event_bus_id = tencentcloud_eb_event_bus.foo.id
}
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	eb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/eb/v20210416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudEbPutEvents() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEbPutEventsCreate,
		Read:   resourceTencentCloudEbPutEventsRead,
		Delete: resourceTencentCloudEbPutEventsDelete,

		Schema: map[string]*schema.Schema{
			"event_list": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "event list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Event source information, new product reporting must comply with EB specifications.",
						},
						"data": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Event data, the content is controlled by the system that created the event, the current datacontenttype only supports application/json;charset=utf-8, so this field is a json string.",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Event type, customizable, optional. The cloud service writes COS:Created:PostObject by default, use: to separate the type field.",
						},
						"subject": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Detailed description of the event source, customizable, optional. The cloud service defaults to the standard qcs resource representation syntax: qcs::dts:ap-guangzhou:appid/uin:xxx.",
						},
						"time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The timestamp in milliseconds when the event occurred,time.Now().UnixNano()/1e6.",
						},
					},
				},
			},

			"event_bus_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "event bus Id.",
			},
		},
	}
}

func resourceTencentCloudEbPutEventsCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eb_put_events.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = eb.NewPutEventsRequest()
		eventBusId string
	)
	if v, ok := d.GetOk("event_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			event := eb.Event{}
			if v, ok := dMap["source"]; ok {
				event.Source = helper.String(v.(string))
			}
			if v, ok := dMap["data"]; ok {
				event.Data = helper.String(v.(string))
			}
			if v, ok := dMap["type"]; ok {
				event.Type = helper.String(v.(string))
			}
			if v, ok := dMap["subject"]; ok {
				event.Subject = helper.String(v.(string))
			}
			if v, ok := dMap["time"]; ok {
				event.Time = helper.IntInt64(v.(int))
			}
			request.EventList = append(request.EventList, &event)
		}
	}

	if v, ok := d.GetOk("event_bus_id"); ok {
		eventBusId = v.(string)
		request.EventBusId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseEbClient().PutEvents(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate eb putEvents failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(eventBusId)

	return resourceTencentCloudEbPutEventsRead(d, meta)
}

func resourceTencentCloudEbPutEventsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eb_put_events.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudEbPutEventsDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eb_put_events.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
