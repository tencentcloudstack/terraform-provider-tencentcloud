package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMysqlRoGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlRoGroupCreate,
		Read:   resourceTencentCloudMysqlRoGroupRead,
		Update: resourceTencentCloudMysqlRoGroupUpdate,
		Delete: resourceTencentCloudMysqlRoGroupDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, in the format: cdbro-3i70uj0k.",
			},

			"ro_group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of the RO group.",
			},

			"ro_group_info": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Details of the RO group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ro_group_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "RO group name.",
						},
						"ro_max_delay_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "RO instance maximum latency threshold. The unit is seconds, the minimum value is 1. Note that the RO group must have enabled instance delay culling policy for this value to be valid.",
						},
						"ro_offline_delay": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Whether to enable delayed culling of instances. Supported values are: 1 - on; 0 - not on. Note that if you enable instance delay culling, you must set the delay threshold (RoMaxDelayTime) parameter.",
						},
						"min_ro_in_group": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The minimum number of reserved instances. It can be set to any value less than or equal to the number of RO instances under this RO group. Note that if the setting value is greater than the number of RO instances, it will not be removed; if it is set to 0, all instances whose latency exceeds the limit will be removed.",
						},
						"weight_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "weight mode. Supported values include: `system` - automatically assigned by the system; `custom` - user-defined settings. Note that if the `custom` mode is set, the RO instance weight configuration (RoWeightValues) parameter must be set.",
						},
						"replication_delay_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Delayed replication time.",
						},
					},
				},
			},

			"ro_weight_values": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "The weight of the instance within the RO group. If the weight mode of the RO group is changed to user-defined mode (custom), this parameter must be set, and the weight value of each RO instance needs to be set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "RO instance ID.",
						},
						"weight": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Weights. The value range is [0, 100].",
						},
					},
				},
			},

			"is_balance_ro_load": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to rebalance the load of RO instances in the RO group. Supported values include: 1 - rebalance load; 0 - do not rebalance load. The default value is 0. Note that when it is set to rebalance the load, the RO instance in the RO group will have a momentary disconnection of the database connection, please ensure that the application can reconnect to the database.",
			},
		},
	}
}

func resourceTencentCloudMysqlRoGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_ro_group.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	var roGroupId string
	if v, ok := d.GetOk("ro_group_id"); ok {
		roGroupId = v.(string)
	}

	d.SetId(instanceId + FILED_SP + roGroupId)

	return resourceTencentCloudMysqlRoGroupUpdate(d, meta)
}

func resourceTencentCloudMysqlRoGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_ro_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	roGroupId := idSplit[1]

	roGroup, err := service.DescribeMysqlRoGroupById(ctx, instanceId, roGroupId)
	if err != nil {
		return err
	}

	if roGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MysqlRoGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if roGroup.RoGroupId != nil {
		_ = d.Set("ro_group_id", roGroup.RoGroupId)
	}

	if roGroup != nil {
		roGroupInfoMap := map[string]interface{}{}

		if roGroup.RoGroupName != nil {
			roGroupInfoMap["ro_group_name"] = roGroup.RoGroupName
		}

		if roGroup.RoMaxDelayTime != nil {
			roGroupInfoMap["ro_max_delay_time"] = roGroup.RoMaxDelayTime
		}

		if roGroup.RoOfflineDelay != nil {
			roGroupInfoMap["ro_offline_delay"] = roGroup.RoOfflineDelay
		}

		if roGroup.MinRoInGroup != nil {
			roGroupInfoMap["min_ro_in_group"] = roGroup.MinRoInGroup
		}

		if roGroup.WeightMode != nil {
			roGroupInfoMap["weight_mode"] = roGroup.WeightMode
		}

		_ = d.Set("ro_group_info", []interface{}{roGroupInfoMap})
	}

	if roGroup.RoInstances != nil {
		roWeightValuesList := []interface{}{}
		for _, roWeightValues := range roGroup.RoInstances {
			roWeightValuesMap := map[string]interface{}{}

			if roWeightValues.InstanceId != nil {
				roWeightValuesMap["instance_id"] = roWeightValues.InstanceId
			}

			if roWeightValues.Weight != nil {
				roWeightValuesMap["weight"] = roWeightValues.Weight
			}

			roWeightValuesList = append(roWeightValuesList, roWeightValuesMap)
		}

		_ = d.Set("ro_weight_values", roWeightValuesList)

	}

	return nil
}

func resourceTencentCloudMysqlRoGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_ro_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := mysql.NewModifyRoGroupInfoRequest()
	response := mysql.NewModifyRoGroupInfoResponse()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	roGroupId := idSplit[1]

	request.RoGroupId = &roGroupId

	if d.HasChange("ro_group_info") {
		if dMap, ok := helper.InterfacesHeadMap(d, "ro_group_info"); ok {
			roGroupAttr := mysql.RoGroupAttr{}
			if v, ok := dMap["ro_group_name"]; ok {
				roGroupAttr.RoGroupName = helper.String(v.(string))
			}
			if v, ok := dMap["ro_max_delay_time"]; ok {
				roGroupAttr.RoMaxDelayTime = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["ro_offline_delay"]; ok {
				roGroupAttr.RoOfflineDelay = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["min_ro_in_group"]; ok {
				roGroupAttr.MinRoInGroup = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["weight_mode"]; ok {
				roGroupAttr.WeightMode = helper.String(v.(string))
			}
			if v, ok := dMap["replication_delay_time"]; ok {
				roGroupAttr.ReplicationDelayTime = helper.IntInt64(v.(int))
			}
			request.RoGroupInfo = &roGroupAttr
		}
	}

	if d.HasChange("ro_weight_values") {
		if v, ok := d.GetOk("ro_weight_values"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				roWeightValue := mysql.RoWeightValue{}
				if v, ok := dMap["instance_id"]; ok {
					roWeightValue.InstanceId = helper.String(v.(string))
				}
				if v, ok := dMap["weight"]; ok {
					roWeightValue.Weight = helper.IntInt64(v.(int))
				}
				request.RoWeightValues = append(request.RoWeightValues, &roWeightValue)
			}
		}
	}

	if d.HasChange("is_balance_ro_load") {
		if v, ok := d.GetOkExists("is_balance_ro_load"); ok {
			request.IsBalanceRoLoad = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().ModifyRoGroupInfo(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mysql roGroup failed, reason:%+v", logId, err)
		return err
	}

	asyncRequestId := *response.Response.AsyncRequestId
	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		taskStatus, message, err := service.DescribeAsyncRequestInfo(ctx, asyncRequestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
			return nil
		}
		if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("%s create mysql rollback status is %s", instanceId, taskStatus))
		}
		err = fmt.Errorf("%s create mysql rollback status is %s,we won't wait for it finish ,it show message:%s", instanceId, taskStatus, message)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mysql rollback fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudMysqlRoGroupRead(d, meta)
}

func resourceTencentCloudMysqlRoGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_ro_group.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
