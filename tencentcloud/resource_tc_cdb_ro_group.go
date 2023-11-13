/*
Provides a resource to create a cdb ro_group

Example Usage

```hcl
resource "tencentcloud_cdb_ro_group" "ro_group" {
  ro_group_id = ""
  ro_group_info {
		ro_group_name = ""
		ro_max_delay_time =
		ro_offline_delay =
		min_ro_in_group =
		weight_mode = ""
		replication_delay_time =

  }
  ro_weight_values {
		instance_id = ""
		weight =

  }
  is_balance_ro_load =
  replication_delay_time =
}
```

Import

cdb ro_group can be imported using the id, e.g.

```
terraform import tencentcloud_cdb_ro_group.ro_group ro_group_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"log"
	"strings"
	"time"
)

func resourceTencentCloudCdbRoGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdbRoGroupCreate,
		Read:   resourceTencentCloudCdbRoGroupRead,
		Update: resourceTencentCloudCdbRoGroupUpdate,
		Delete: resourceTencentCloudCdbRoGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
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
							Description: "Weight mode. Supported values include: `system` - automatically assigned by the system; `custom` - user-defined settings. Note that if the `custom` mode is set, the RO instance weight configuration (RoWeightValues) parameter must be set.",
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

			"replication_delay_time": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Obsolete parameter, meaningless.",
			},
		},
	}
}

func resourceTencentCloudCdbRoGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_ro_group.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	var roGroupId string
	if v, ok := d.GetOk("ro_group_id"); ok {
		roGroupId = v.(string)
	}

	d.SetId(strings.Join([]string{instanceId, roGroupId}, FILED_SP))

	return resourceTencentCloudCdbRoGroupUpdate(d, meta)
}

func resourceTencentCloudCdbRoGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_ro_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	roGroupId := idSplit[1]

	roGroup, err := service.DescribeCdbRoGroupById(ctx, instanceId, roGroupId)
	if err != nil {
		return err
	}

	if roGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CdbRoGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if roGroup.RoGroupId != nil {
		_ = d.Set("ro_group_id", roGroup.RoGroupId)
	}

	if roGroup.RoGroupInfo != nil {
		roGroupInfoMap := map[string]interface{}{}

		if roGroup.RoGroupInfo.RoGroupName != nil {
			roGroupInfoMap["ro_group_name"] = roGroup.RoGroupInfo.RoGroupName
		}

		if roGroup.RoGroupInfo.RoMaxDelayTime != nil {
			roGroupInfoMap["ro_max_delay_time"] = roGroup.RoGroupInfo.RoMaxDelayTime
		}

		if roGroup.RoGroupInfo.RoOfflineDelay != nil {
			roGroupInfoMap["ro_offline_delay"] = roGroup.RoGroupInfo.RoOfflineDelay
		}

		if roGroup.RoGroupInfo.MinRoInGroup != nil {
			roGroupInfoMap["min_ro_in_group"] = roGroup.RoGroupInfo.MinRoInGroup
		}

		if roGroup.RoGroupInfo.WeightMode != nil {
			roGroupInfoMap["weight_mode"] = roGroup.RoGroupInfo.WeightMode
		}

		if roGroup.RoGroupInfo.ReplicationDelayTime != nil {
			roGroupInfoMap["replication_delay_time"] = roGroup.RoGroupInfo.ReplicationDelayTime
		}

		_ = d.Set("ro_group_info", []interface{}{roGroupInfoMap})
	}

	if roGroup.RoWeightValues != nil {
		roWeightValuesList := []interface{}{}
		for _, roWeightValues := range roGroup.RoWeightValues {
			roWeightValuesMap := map[string]interface{}{}

			if roGroup.RoWeightValues.InstanceId != nil {
				roWeightValuesMap["instance_id"] = roGroup.RoWeightValues.InstanceId
			}

			if roGroup.RoWeightValues.Weight != nil {
				roWeightValuesMap["weight"] = roGroup.RoWeightValues.Weight
			}

			roWeightValuesList = append(roWeightValuesList, roWeightValuesMap)
		}

		_ = d.Set("ro_weight_values", roWeightValuesList)

	}

	if roGroup.IsBalanceRoLoad != nil {
		_ = d.Set("is_balance_ro_load", roGroup.IsBalanceRoLoad)
	}

	if roGroup.ReplicationDelayTime != nil {
		_ = d.Set("replication_delay_time", roGroup.ReplicationDelayTime)
	}

	return nil
}

func resourceTencentCloudCdbRoGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_ro_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cdb.NewModifyRoGroupInfoRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	roGroupId := idSplit[1]

	request.InstanceId = &instanceId
	request.RoGroupId = &roGroupId

	immutableArgs := []string{"ro_group_id", "ro_group_info", "ro_weight_values", "is_balance_ro_load", "replication_delay_time"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().ModifyRoGroupInfo(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cdb roGroup failed, reason:%+v", logId, err)
		return err
	}

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCEED"}, 1*readRetryTimeout, time.Second, service.CdbRoGroupStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCdbRoGroupRead(d, meta)
}

func resourceTencentCloudCdbRoGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_ro_group.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
