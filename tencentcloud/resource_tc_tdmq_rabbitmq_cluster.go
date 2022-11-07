/*
Provides a resource to create a tdmq rabbitmq_cluster

Example Usage

```hcl
resource "tencentcloud_tdmq_rabbitmq_cluster" "rabbitmq_cluster" {
  name = ""
  remark = ""
}

```
Import

tdmq rabbitmq_cluster can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rabbitmq_cluster.rabbitmq_cluster rabbitmqCluster_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTdmqRabbitmqCluster() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTdmqRabbitmqClusterRead,
		Create: resourceTencentCloudTdmqRabbitmqClusterCreate,
		Update: resourceTencentCloudTdmqRabbitmqClusterUpdate,
		Delete: resourceTencentCloudTdmqRabbitmqClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "cluster name.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "cluster description, 128 characters or less.",
			},
		},
	}
}

func resourceTencentCloudTdmqRabbitmqClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_cluster.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tdmq.NewCreateAMQPClusterRequest()
		response  *tdmq.CreateAMQPClusterResponse
		clusterId string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().CreateAMQPCluster(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmq rabbitmqCluster failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId

	d.SetId(clusterId)
	return resourceTencentCloudTdmqRabbitmqClusterRead(d, meta)
}

func resourceTencentCloudTdmqRabbitmqClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_cluster.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	clusterId := d.Id()

	rabbitmq, err := service.DescribeTdmqRabbitmqCluster(ctx, clusterId)

	if err != nil {
		return err
	}

	rabbitmqCluster := rabbitmq.ClusterInfo
	if rabbitmqCluster == nil {
		d.SetId("")
		return fmt.Errorf("resource `rabbitmqCluster` %s does not exist", clusterId)
	}

	if rabbitmqCluster.ClusterName != nil {
		_ = d.Set("name", rabbitmqCluster.ClusterName)
	}

	if rabbitmqCluster.Remark != nil {
		_ = d.Set("remark", rabbitmqCluster.Remark)
	}

	return nil
}

func resourceTencentCloudTdmqRabbitmqClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_cluster.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tdmq.NewModifyAMQPClusterRequest()

	clusterId := d.Id()

	request.ClusterId = &clusterId

	if d.HasChange("name") {
		return fmt.Errorf("`name` do not support change now.")
	}

	if v, ok := d.GetOk("name"); ok {
		request.ClusterName = helper.String(v.(string))
	} else {
		return fmt.Errorf("`name` cannot be empty.")
	}

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().ModifyAMQPCluster(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmq rabbitmqCluster failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTdmqRabbitmqClusterRead(d, meta)
}

func resourceTencentCloudTdmqRabbitmqClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_cluster.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	clusterId := d.Id()

	if err := service.DeleteTdmqRabbitmqClusterById(ctx, clusterId); err != nil {
		return err
	}

	return nil
}
