/*
Provides a resource to create a tdmq rabbitmq_exchange

Example Usage

```hcl
resource "tencentcloud_tdmq_rabbitmq_exchange" "rabbitmq_exchange" {
  exchange = ""
  vhost_id = ""
  type = ""
  cluster_id = ""
  remark = ""
  alternate_exchange = ""
  delayed_type = ""
}

```
Import

tdmq rabbitmq_exchange can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rabbitmq_exchange.rabbitmq_exchange rabbitmqExchange_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTdmqRabbitmqExchange() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTdmqRabbitmqExchangeRead,
		Create: resourceTencentCloudTdmqRabbitmqExchangeCreate,
		Update: resourceTencentCloudTdmqRabbitmqExchangeUpdate,
		Delete: resourceTencentCloudTdmqRabbitmqExchangeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"exchange": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "exchange name.",
			},

			"vhost_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "vhost.",
			},

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "exchange type.",
			},

			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "cluster id.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "exchange comment.",
			},

			"alternate_exchange": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "alternate exchange name.",
			},

			"delayed_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "delayed exchange type, the value must be one of Direct, Fanout, Topic.",
			},
		},
	}
}

func resourceTencentCloudTdmqRabbitmqExchangeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_exchange.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tdmq.NewCreateAMQPExchangeRequest()
		clusterId string
		vHostId   string
		exchange  string
	)

	if v, ok := d.GetOk("exchange"); ok {
		exchange = v.(string)
		request.Exchange = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vhost_id"); ok {
		vHostId = v.(string)
		request.VHosts = []*string{helper.String(v.(string))}
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("alternate_exchange"); ok {
		request.AlternateExchange = helper.String(v.(string))
	}

	if v, ok := d.GetOk("delayed_type"); ok {
		request.DelayedType = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().CreateAMQPExchange(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmq rabbitmqExchange failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(clusterId + FILED_SP + vHostId + FILED_SP + exchange)
	return resourceTencentCloudTdmqRabbitmqExchangeRead(d, meta)
}

func resourceTencentCloudTdmqRabbitmqExchangeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_exchange.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	vHostId := idSplit[1]
	exchange := idSplit[2]

	rabbitmqExchange, err := service.DescribeTdmqRabbitmqExchange(ctx, clusterId, vHostId, exchange)

	if err != nil {
		return err
	}

	if rabbitmqExchange == nil {
		d.SetId("")
		return fmt.Errorf("resource `rabbitmqExchange` %s does not exist", exchange)
	}

	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("vhost_id", vHostId)
	_ = d.Set("exchange", exchange)

	if rabbitmqExchange.Type != nil {
		_ = d.Set("type", rabbitmqExchange.Type)
	}

	if rabbitmqExchange.Remark != nil {
		_ = d.Set("remark", rabbitmqExchange.Remark)
	}

	if rabbitmqExchange.AlternateExchange != nil {
		_ = d.Set("alternate_exchange", rabbitmqExchange.AlternateExchange)
	}

	if rabbitmqExchange.DelayType != nil {
		_ = d.Set("delayed_type", rabbitmqExchange.DelayType)
	}

	return nil
}

func resourceTencentCloudTdmqRabbitmqExchangeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_exchange.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tdmq.NewModifyAMQPExchangeRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	vHostId := idSplit[1]
	exchange := idSplit[2]

	request.ClusterId = &clusterId
	request.VHostId = &vHostId
	request.Exchange = &exchange

	if d.HasChange("exchange") {
		return fmt.Errorf("`exchange` do not support change now.")
	}

	if d.HasChange("vhost_id") {
		return fmt.Errorf("`vhost_id` do not support change now.")
	}

	if d.HasChange("type") {
		return fmt.Errorf("`type` do not support change now.")
	}

	if d.HasChange("cluster_id") {
		return fmt.Errorf("`cluster_id` do not support change now.")
	}

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}
	}

	if d.HasChange("alternate_exchange") {
		return fmt.Errorf("`alternate_exchange` do not support change now.")
	}

	if d.HasChange("delayed_type") {
		return fmt.Errorf("`delayed_type` do not support change now.")
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().ModifyAMQPExchange(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmq rabbitmqExchange failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTdmqRabbitmqExchangeRead(d, meta)
}

func resourceTencentCloudTdmqRabbitmqExchangeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_exchange.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	vHostId := idSplit[1]
	exchange := idSplit[2]

	if err := service.DeleteTdmqRabbitmqExchangeById(ctx, clusterId, vHostId, exchange); err != nil {
		return err
	}

	return nil
}
