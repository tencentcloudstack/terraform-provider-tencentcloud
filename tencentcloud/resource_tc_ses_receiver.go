/*
Provides a resource to create a ses receiver

Example Usage

```hcl
resource "tencentcloud_ses_receiver" "receiver" {
  receivers_name = "Recipient group name"
  desc = "Recipient group description"
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
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudSesReceiver() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSesReceiverCreate,
		Read:   resourceTencentCloudSesReceiverRead,
		Update: resourceTencentCloudSesReceiverUpdate,
		Delete: resourceTencentCloudSesReceiverDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"receivers_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Recipient group name.",
			},

			"desc": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Recipient group description.",
			},
		},
	}
}

func resourceTencentCloudSesReceiverCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_receiver.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = ses.NewCreateReceiverRequest()
		response = ses.NewCreateReceiverResponse()
		receiver string
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

	receiver = *response.Response.Receiver
	d.SetId(receiver)

	return resourceTencentCloudSesReceiverRead(d, meta)
}

func resourceTencentCloudSesReceiverRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_receiver.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SesService{client: meta.(*TencentCloudClient).apiV3Conn}

	receiverId := d.Id()

	Receiver, err := service.DescribeSesReceiverById(ctx, receiver)
	if err != nil {
		return err
	}

	if Receiver == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SesReceiver` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if Receiver.ReceiversName != nil {
		_ = d.Set("receivers_name", Receiver.ReceiversName)
	}

	if Receiver.Desc != nil {
		_ = d.Set("desc", Receiver.Desc)
	}

	return nil
}

func resourceTencentCloudSesReceiverUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_receiver.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"receivers_name", "desc"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudSesReceiverRead(d, meta)
}

func resourceTencentCloudSesReceiverDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_receiver.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SesService{client: meta.(*TencentCloudClient).apiV3Conn}
	receiverId := d.Id()

	if err := service.DeleteSesReceiverById(ctx, receiver); err != nil {
		return err
	}

	return nil
}
