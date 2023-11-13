/*
Provides a resource to create a tdmqRocketmq group

Example Usage

```hcl
resource "tencentcloud_tdmq_rocketmq_group" "group" {
  group_id = &lt;nil&gt;
  namespaces = &lt;nil&gt;
  read_enable = &lt;nil&gt;
  broadcast_enable = &lt;nil&gt;
  cluster_id = &lt;nil&gt;
  remark = &lt;nil&gt;
                    }
```

Import

tdmqRocketmq group can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_rocketmq_group.group group_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmqRocketmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudTdmqRocketmqGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqRocketmqGroupCreate,
		Read:   resourceTencentCloudTdmqRocketmqGroupRead,
		Update: resourceTencentCloudTdmqRocketmqGroupUpdate,
		Delete: resourceTencentCloudTdmqRocketmqGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Group name (8-64 characters).",
			},

			"namespaces": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Namespace. Currently, only one namespace is supported.",
			},

			"read_enable": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable consumption.",
			},

			"broadcast_enable": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable broadcast consumption.",
			},

			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"remark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Remarks (up to 128 characters).",
			},

			"name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Consumer group name.",
			},

			"consumer_num": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The number of online consumers.",
			},

			"t_p_s": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Consumption TPS.",
			},

			"total_accumulative": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The total number of heaped messages.",
			},

			"consumption_mode": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "`0`: Cluster consumption mode; `1`: Broadcast consumption mode; `-1`: Unknown.",
			},

			"retry_partition_num": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The number of partitions in a retry topic.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Creation time in milliseconds.",
			},

			"update_time": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Modification time in milliseconds.",
			},

			"client_protocol": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Client protocol.",
			},

			"consumer_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Consumer type. Enumerated values: ACTIVELY or PASSIVELY.",
			},
		},
	}
}

func resourceTencentCloudTdmqRocketmqGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request     = tdmqRocketmq.NewCreateRocketMQGroupRequest()
		response    = tdmqRocketmq.NewCreateRocketMQGroupResponse()
		clusterId   string
		namespaceId string
		groupId     string
	)
	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
		request.GroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespaces"); ok {
		namespacesSet := v.(*schema.Set).List()
		for i := range namespacesSet {
			namespaces := namespacesSet[i].(string)
			request.Namespaces = append(request.Namespaces, &namespaces)
		}
	}

	if v, ok := d.GetOkExists("read_enable"); ok {
		request.ReadEnable = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("broadcast_enable"); ok {
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqRocketmqClient().CreateRocketMQGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tdmqRocketmq group failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId
	d.SetId(strings.Join([]string{clusterId, namespaceId, groupId}, FILED_SP))

	return resourceTencentCloudTdmqRocketmqGroupRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_group.read")()
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

	group, err := service.DescribeTdmqRocketmqGroupById(ctx, clusterId, namespaceId, groupId)
	if err != nil {
		return err
	}

	if group == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TdmqRocketmqGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if group.GroupId != nil {
		_ = d.Set("group_id", group.GroupId)
	}

	if group.Namespaces != nil {
		_ = d.Set("namespaces", group.Namespaces)
	}

	if group.ReadEnable != nil {
		_ = d.Set("read_enable", group.ReadEnable)
	}

	if group.BroadcastEnable != nil {
		_ = d.Set("broadcast_enable", group.BroadcastEnable)
	}

	if group.ClusterId != nil {
		_ = d.Set("cluster_id", group.ClusterId)
	}

	if group.Remark != nil {
		_ = d.Set("remark", group.Remark)
	}

	if group.Name != nil {
		_ = d.Set("name", group.Name)
	}

	if group.ConsumerNum != nil {
		_ = d.Set("consumer_num", group.ConsumerNum)
	}

	if group.TPS != nil {
		_ = d.Set("t_p_s", group.TPS)
	}

	if group.TotalAccumulative != nil {
		_ = d.Set("total_accumulative", group.TotalAccumulative)
	}

	if group.ConsumptionMode != nil {
		_ = d.Set("consumption_mode", group.ConsumptionMode)
	}

	if group.RetryPartitionNum != nil {
		_ = d.Set("retry_partition_num", group.RetryPartitionNum)
	}

	if group.CreateTime != nil {
		_ = d.Set("create_time", group.CreateTime)
	}

	if group.UpdateTime != nil {
		_ = d.Set("update_time", group.UpdateTime)
	}

	if group.ClientProtocol != nil {
		_ = d.Set("client_protocol", group.ClientProtocol)
	}

	if group.ConsumerType != nil {
		_ = d.Set("consumer_type", group.ConsumerType)
	}

	return nil
}

func resourceTencentCloudTdmqRocketmqGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_group.update")()
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

	immutableArgs := []string{"group_id", "namespaces", "read_enable", "broadcast_enable", "cluster_id", "remark", "name", "consumer_num", "t_p_s", "total_accumulative", "consumption_mode", "retry_partition_num", "create_time", "update_time", "client_protocol", "consumer_type"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("read_enable") {
		if v, ok := d.GetOkExists("read_enable"); ok {
			request.ReadEnable = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("broadcast_enable") {
		if v, ok := d.GetOkExists("broadcast_enable"); ok {
			request.BroadcastEnable = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqRocketmqClient().ModifyRocketMQGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tdmqRocketmq group failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTdmqRocketmqGroupRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_group.delete")()
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
