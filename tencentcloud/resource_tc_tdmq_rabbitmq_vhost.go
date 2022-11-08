/*
Provides a resource to create a tdmq rabbitmq_vhost

Example Usage

```hcl
resource "tencentcloud_tdmq_rabbitmq_vhost" "rabbitmq_vhost" {
  cluster_id = ""
  vhost_id = ""
  msg_ttl = ""
  remark = ""
}

```
Import

tdmq rabbitmq_vhost can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rabbitmq_vhost.rabbitmq_vhost rabbitmqVhost_id
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

func resourceTencentCloudTdmqRabbitmqVhost() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTdmqRabbitmqVhostRead,
		Create: resourceTencentCloudTdmqRabbitmqVhostCreate,
		Update: resourceTencentCloudTdmqRabbitmqVhostUpdate,
		Delete: resourceTencentCloudTdmqRabbitmqVhostDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "cluster id.",
			},

			"vhost_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "vhost name, can only contain letters, numbers, '-' and '_'.",
			},

			"msg_ttl": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "retention time for unconsumed messages, the unit is ms, range is 60s~15days.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "cluster description.",
			},
		},
	}
}

func resourceTencentCloudTdmqRabbitmqVhostCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_vhost.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tdmq.NewCreateAMQPVHostRequest()
		clusterId string
		vHostId   string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vhost_id"); ok {
		vHostId = v.(string)
		request.VHostId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("msg_ttl"); ok {
		request.MsgTtl = helper.Uint64(uint64(v.(int)))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().CreateAMQPVHost(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmq rabbitmqVhost failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(clusterId + FILED_SP + vHostId)
	return resourceTencentCloudTdmqRabbitmqVhostRead(d, meta)
}

func resourceTencentCloudTdmqRabbitmqVhostRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_vhost.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	vHostId := idSplit[1]

	rabbitmqVhost, err := service.DescribeTdmqRabbitmqVhost(ctx, clusterId, vHostId)

	if err != nil {
		return err
	}

	if rabbitmqVhost == nil {
		d.SetId("")
		return fmt.Errorf("resource `rabbitmqVhost` %s does not exist", vHostId)
	}

	_ = d.Set("cluster_id", clusterId)

	_ = d.Set("vhost_id", vHostId)

	if rabbitmqVhost.MsgTtl != nil {
		_ = d.Set("msg_ttl", rabbitmqVhost.MsgTtl)
	}

	if rabbitmqVhost.Remark != nil {
		_ = d.Set("remark", rabbitmqVhost.Remark)
	}

	return nil
}

func resourceTencentCloudTdmqRabbitmqVhostUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_vhost.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tdmq.NewModifyAMQPVHostRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	vHostId := idSplit[1]

	request.ClusterId = &clusterId
	request.VHostId = &vHostId

	if d.HasChange("cluster_id") {
		return fmt.Errorf("`cluster_id` do not support change now.")
	}

	if d.HasChange("vhost_id") {
		return fmt.Errorf("`vhost_id` do not support change now.")
	}

	if v, ok := d.GetOk("msg_ttl"); ok {
		request.MsgTtl = helper.Uint64(uint64(v.(int)))
	}

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().ModifyAMQPVHost(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmq rabbitmqVhost failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTdmqRabbitmqVhostRead(d, meta)
}

func resourceTencentCloudTdmqRabbitmqVhostDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_vhost.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	vHostId := idSplit[1]

	if err := service.DeleteTdmqRabbitmqVhostById(ctx, clusterId, vHostId); err != nil {
		return err
	}

	return nil
}
