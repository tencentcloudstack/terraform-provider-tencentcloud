/*
Provides a resource to create a tat invoker

Example Usage

```hcl
resource "tencentcloud_tat_invoker" "invoker" {
  invoker_id = ""
}
```

Import

tat invoker can be imported using the id, e.g.

```
terraform import tencentcloud_tat_invoker.invoker invoker_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTatInvoker() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTatInvokerCreate,
		Read:   resourceTencentCloudTatInvokerRead,
		Delete: resourceTencentCloudTatInvokerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"invoker_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of the invoker to be enabled.",
			},
		},
	}
}

func resourceTencentCloudTatInvokerCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tat_invoker.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tat.NewEnableInvokerRequest()
		response  = tat.NewEnableInvokerResponse()
		invokerId string
	)
	if v, ok := d.GetOk("invoker_id"); ok {
		invokerId = v.(string)
		request.InvokerId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTatClient().EnableInvoker(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tat invoker failed, reason:%+v", logId, err)
		return err
	}

	invokerId = *response.Response.InvokerId
	d.SetId(invokerId)

	return resourceTencentCloudTatInvokerRead(d, meta)
}

func resourceTencentCloudTatInvokerRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tat_invoker.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TatService{client: meta.(*TencentCloudClient).apiV3Conn}

	invokerId := d.Id()

	invoker, err := service.DescribeTatInvokerById(ctx, invokerId)
	if err != nil {
		return err
	}

	if invoker == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TatInvoker` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if invoker.InvokerId != nil {
		_ = d.Set("invoker_id", invoker.InvokerId)
	}

	return nil
}

func resourceTencentCloudTatInvokerDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tat_invoker.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TatService{client: meta.(*TencentCloudClient).apiV3Conn}
	invokerId := d.Id()

	if err := service.DeleteTatInvokerById(ctx, invokerId); err != nil {
		return err
	}

	return nil
}
