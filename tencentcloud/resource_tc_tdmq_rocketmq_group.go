/*
Provides a resource to create a tdmqRocketmq group

Example Usage

```hcl
resource "tencentcloud_tdmq_rocketmq_cluster" "cluster" {
	cluster_name = "test_rocketmq"
	remark = "test recket mq"
}

resource "tencentcloud_tdmq_rocketmq_namespace" "namespace" {
  cluster_id = tencentcloud_tdmq_rocketmq_cluster.cluster.cluster_id
  namespace_name = "test_namespace"
  ttl = 65000
  retention_time = 65000
  remark = "test namespace"
}

resource "tencentcloud_tdmq_rocketmq_group" "group" {
  group_name = "test_rocketmq_group"
  namespace = tencentcloud_tdmq_rocketmq_namespace.namespace.namespace_name
  read_enable = true
  broadcast_enable = true
  cluster_id = tencentcloud_tdmq_rocketmq_cluster.cluster.cluster_id
  remark = "test rocketmq group"
}
```
Import

tdmqRocketmq group can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rocketmq_group.group group_id
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
	tdmqRocketmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTdmqRocketmqGroup() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTdmqRocketmqGroupRead,
		Create: resourceTencentCloudTdmqRocketmqGroupCreate,
		Update: resourceTencentCloudTdmqRocketmqGroupUpdate,
		Delete: resourceTencentCloudTdmqRocketmqGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Group name (8-64 characters).",
			},

			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Namespace. Currently, only one namespace is supported.",
			},

			"read_enable": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable consumption.",
			},

			"broadcast_enable": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable broadcast consumption.",
			},

			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster ID.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remarks (up to 128 characters).",
			},

			"consumer_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of online consumers.",
			},

			"tps": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Consumption TPS.",
			},

			"total_accumulative": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of heaped messages.",
			},

			"consumption_mode": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "`0`: Cluster consumption mode; `1`: Broadcast consumption mode; `-1`: Unknown.",
			},

			"retry_partition_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of partitions in a retry topic.",
			},

			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Creation time in milliseconds.",
			},

			"update_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Modification time in milliseconds.",
			},

			"client_protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Client protocol.",
			},

			"consumer_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Consumer type. Enumerated values: ACTIVELY or PASSIVELY.",
			},
		},
	}
}

func resourceTencentCloudTdmqRocketmqGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmqRocketmq_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = tdmqRocketmq.NewCreateRocketMQGroupRequest()
		clusterId     string
		namespaceName string
		groupName     string
	)

	if v, ok := d.GetOk("group_name"); ok {
		groupName = v.(string)
		request.GroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		namespaceName = v.(string)
		request.Namespaces = []*string{&namespaceName}
	}

	if v, _ := d.GetOk("read_enable"); v != nil {
		request.ReadEnable = helper.Bool(v.(bool))
	}

	if v, _ := d.GetOk("broadcast_enable"); v != nil {
		request.BroadcastEnable = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {

		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().CreateRocketMQGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmqRocketmq group failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(clusterId + FILED_SP + namespaceName + FILED_SP + groupName)
	return resourceTencentCloudTdmqRocketmqGroupRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmqRocketmq_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqRocketmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	namespaceName := idSplit[1]
	groupName := idSplit[2]

	groupList, err := service.DescribeTdmqRocketmqGroup(ctx, clusterId, namespaceName, groupName)

	if err != nil {
		return err
	}

	if len(groupList) == 0 {
		d.SetId("")
		return fmt.Errorf("resource `group` %s does not exist", groupName)
	}
	group := groupList[0]
	_ = d.Set("group_name", group.Name)
	_ = d.Set("namespace", namespaceName)
	_ = d.Set("read_enable", group.ReadEnabled)
	_ = d.Set("broadcast_enable", group.BroadcastEnabled)
	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("remark", group.Remark)
	_ = d.Set("consumer_num", group.ConsumerNum)
	_ = d.Set("tps", group.TPS)
	_ = d.Set("total_accumulative", group.TotalAccumulative)
	_ = d.Set("consumption_mode", group.ConsumptionMode)
	_ = d.Set("retry_partition_num", group.RetryPartitionNum)
	_ = d.Set("create_time", group.CreateTime)
	_ = d.Set("update_time", group.UpdateTime)
	_ = d.Set("client_protocol", group.ClientProtocol)
	_ = d.Set("consumer_type", group.ConsumerType)

	return nil
}

func resourceTencentCloudTdmqRocketmqGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmqRocketmq_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tdmqRocketmq.NewModifyRocketMQGroupRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	namespaceId := idSplit[1]
	groupId := idSplit[2]

	request.ClusterId = &clusterId
	request.NamespaceId = &namespaceId
	request.GroupId = &groupId

	if d.HasChange("group_id") {

		return fmt.Errorf("`group_id` do not support change now.")

	}

	if d.HasChange("namespaces") {

		return fmt.Errorf("`namespaces` do not support change now.")

	}

	if d.HasChange("read_enable") {
		if v, ok := d.GetOk("read_enable"); ok {
			request.ReadEnable = helper.Bool(v.(bool))
		}

	}

	if d.HasChange("broadcast_enable") {
		if v, ok := d.GetOk("broadcast_enable"); ok {
			request.BroadcastEnable = helper.Bool(v.(bool))
		}

	}

	if d.HasChange("cluster_id") {

		return fmt.Errorf("`cluster_id` do not support change now.")

	}

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}

	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().ModifyRocketMQGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmqRocketmq group failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTdmqRocketmqGroupRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmqRocketmq_group.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqRocketmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	namespaceId := idSplit[1]
	groupId := idSplit[2]

	if err := service.DeleteTdmqRocketmqGroupById(ctx, clusterId, namespaceId, groupId); err != nil {
		return err
	}

	return nil
}
