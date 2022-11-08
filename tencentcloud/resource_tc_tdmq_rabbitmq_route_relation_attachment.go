/*
Provides a resource to create a tdmq rabbitmq_route_relation_attachment

Example Usage

```hcl
resource "tencentcloud_tdmq_rabbitmq_route_relation_attachment" "rabbitmq_route_relation_attachment" {
  cluster_id = ""
  vhost_id = ""
  source_exchange = ""
  dest_type = ""
  dest_value = ""
  remark = ""
  routing_key = ""
}

```
Import

tdmq rabbitmq_route_relation_attachment can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rabbitmq_route_relation_attachment.rabbitmq_route_relation_attachment rabbitmqRouteRelationAttachment_id
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

func resourceTencentCloudTdmqRabbitmqRouteRelationAttachment() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTdmqRabbitmqRouteRelationAttachmentRead,
		Create: resourceTencentCloudTdmqRabbitmqRouteRelationAttachmentCreate,
		Delete: resourceTencentCloudTdmqRabbitmqRouteRelationAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "cluster id.",
			},

			"vhost_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "vhost id.",
			},

			"source_exchange": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "source exchange name.",
			},

			"dest_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "destination type, the optional value is Queue or Exchange.",
			},

			"dest_value": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "destination value.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "route relation comment.",
			},

			"routing_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "route key, default value is `default`.",
			},
		},
	}
}

func resourceTencentCloudTdmqRabbitmqRouteRelationAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_route_relation_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = tdmq.NewCreateAMQPRouteRelationRequest()
		clusterId       string
		vHostId         string
		routeRelationId string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vhost_id"); ok {
		vHostId = v.(string)
		request.VHostId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_exchange"); ok {
		request.SourceExchange = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dest_type"); ok {
		request.DestType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dest_value"); ok {
		request.DestValue = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("routing_key"); ok {
		request.RoutingKey = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().CreateAMQPRouteRelation(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmq rabbitmqRouteRelationAttachment failed, reason:%+v", logId, err)
		return err
	}

	// routeRelationId = *response.Response.RouteRelationId

	d.SetId(clusterId + FILED_SP + vHostId + FILED_SP + routeRelationId)
	return resourceTencentCloudTdmqRabbitmqRouteRelationAttachmentRead(d, meta)
}

func resourceTencentCloudTdmqRabbitmqRouteRelationAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_route_relation_attachment.read")()
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
	routeRelationId := idSplit[2]

	routeRelation, err := service.DescribeTdmqRabbitmqRouteRelationAttachment(ctx, clusterId, vHostId, routeRelationId)

	if err != nil {
		return err
	}

	if routeRelation == nil {
		d.SetId("")
		return fmt.Errorf("resource `rabbitmqRouteRelationAttachment` %s does not exist", routeRelationId)
	}

	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("vhost_id", vHostId)

	if routeRelation.SourceExchange != nil {
		_ = d.Set("source_exchange", routeRelation.SourceExchange)
	}

	if routeRelation.DestType != nil {
		_ = d.Set("dest_type", routeRelation.DestType)
	}

	if routeRelation.DestValue != nil {
		_ = d.Set("dest_value", routeRelation.DestValue)
	}

	if routeRelation.Remark != nil {
		_ = d.Set("remark", routeRelation.Remark)
	}

	if routeRelation.RoutingKey != nil {
		_ = d.Set("routing_key", routeRelation.RoutingKey)
	}

	return nil
}

func resourceTencentCloudTdmqRabbitmqRouteRelationAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_route_relation_attachment.delete")()
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
	routeRelationId := idSplit[2]

	if err := service.DeleteTdmqRabbitmqRouteRelationAttachmentById(ctx, clusterId, vHostId, routeRelationId); err != nil {
		return err
	}

	return nil
}
