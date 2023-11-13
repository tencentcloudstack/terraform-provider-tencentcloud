/*
Provides a resource to create a sqlserver config_instance_ro_group

Example Usage

```hcl
resource "tencentcloud_sqlserver_config_instance_ro_group" "config_instance_ro_group" {
  instance_id = "mssql-i1z41iwd"
  read_only_group_id = ""
  read_only_group_name = ""
  is_offline_delay =
  read_only_max_delay_time =
  min_read_only_in_group =
  weight_pairs {
		read_only_instance_id = ""
		read_only_weight =

  }
  auto_weight =
  balance_weight =
}
```

Import

sqlserver config_instance_ro_group can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_config_instance_ro_group.config_instance_ro_group config_instance_ro_group_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"log"
)

func resourceTencentCloudSqlserverConfigInstanceRoGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverConfigInstanceRoGroupCreate,
		Read:   resourceTencentCloudSqlserverConfigInstanceRoGroupRead,
		Update: resourceTencentCloudSqlserverConfigInstanceRoGroupUpdate,
		Delete: resourceTencentCloudSqlserverConfigInstanceRoGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"read_only_group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Read-only group ID.",
			},

			"read_only_group_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Read-only group name. If this parameter is not specified, it is not modified.",
			},

			"is_offline_delay": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to enable timeout culling function. 0- Disable the culling function. 1- Enable the culling function.",
			},

			"read_only_max_delay_time": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "After the timeout elimination function is enabled, the timeout threshold used, if this parameter is not filled, it will not be modified.",
			},

			"min_read_only_in_group": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "After the timeout removal function is enabled, the number of read-only copies retained by the read-only group at least, if this parameter is not filled, it will not be modified.",
			},

			"weight_pairs": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Read-only group instance weight modification set, if this parameter is not filled, it will not be modified.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"read_only_instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Read-only instance ID, in the format: mssqlro-3l3fgqn7.",
						},
						"read_only_weight": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Read-only instance weight, the range is 0-100.",
						},
					},
				},
			},

			"auto_weight": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "0-user-defined weight (adjusted according to WeightPairs), 1-system automatically assigns weight (WeightPairs is invalid), the default is 0.",
			},

			"balance_weight": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "0-do not rebalance the load, 1-rebalance the load, the default is 0.",
			},
		},
	}
}

func resourceTencentCloudSqlserverConfigInstanceRoGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_instance_ro_group.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverConfigInstanceRoGroupUpdate(d, meta)
}

func resourceTencentCloudSqlserverConfigInstanceRoGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_instance_ro_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	configInstanceRoGroupId := d.Id()

	configInstanceRoGroup, err := service.DescribeSqlserverConfigInstanceRoGroupById(ctx, instanceId)
	if err != nil {
		return err
	}

	if configInstanceRoGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverConfigInstanceRoGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if configInstanceRoGroup.InstanceId != nil {
		_ = d.Set("instance_id", configInstanceRoGroup.InstanceId)
	}

	if configInstanceRoGroup.ReadOnlyGroupId != nil {
		_ = d.Set("read_only_group_id", configInstanceRoGroup.ReadOnlyGroupId)
	}

	if configInstanceRoGroup.ReadOnlyGroupName != nil {
		_ = d.Set("read_only_group_name", configInstanceRoGroup.ReadOnlyGroupName)
	}

	if configInstanceRoGroup.IsOfflineDelay != nil {
		_ = d.Set("is_offline_delay", configInstanceRoGroup.IsOfflineDelay)
	}

	if configInstanceRoGroup.ReadOnlyMaxDelayTime != nil {
		_ = d.Set("read_only_max_delay_time", configInstanceRoGroup.ReadOnlyMaxDelayTime)
	}

	if configInstanceRoGroup.MinReadOnlyInGroup != nil {
		_ = d.Set("min_read_only_in_group", configInstanceRoGroup.MinReadOnlyInGroup)
	}

	if configInstanceRoGroup.WeightPairs != nil {
		weightPairsList := []interface{}{}
		for _, weightPairs := range configInstanceRoGroup.WeightPairs {
			weightPairsMap := map[string]interface{}{}

			if configInstanceRoGroup.WeightPairs.ReadOnlyInstanceId != nil {
				weightPairsMap["read_only_instance_id"] = configInstanceRoGroup.WeightPairs.ReadOnlyInstanceId
			}

			if configInstanceRoGroup.WeightPairs.ReadOnlyWeight != nil {
				weightPairsMap["read_only_weight"] = configInstanceRoGroup.WeightPairs.ReadOnlyWeight
			}

			weightPairsList = append(weightPairsList, weightPairsMap)
		}

		_ = d.Set("weight_pairs", weightPairsList)

	}

	if configInstanceRoGroup.AutoWeight != nil {
		_ = d.Set("auto_weight", configInstanceRoGroup.AutoWeight)
	}

	if configInstanceRoGroup.BalanceWeight != nil {
		_ = d.Set("balance_weight", configInstanceRoGroup.BalanceWeight)
	}

	return nil
}

func resourceTencentCloudSqlserverConfigInstanceRoGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_instance_ro_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := sqlserver.NewModifyReadOnlyGroupDetailsRequest()

	configInstanceRoGroupId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "read_only_group_id", "read_only_group_name", "is_offline_delay", "read_only_max_delay_time", "min_read_only_in_group", "weight_pairs", "auto_weight", "balance_weight"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().ModifyReadOnlyGroupDetails(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver configInstanceRoGroup failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverConfigInstanceRoGroupRead(d, meta)
}

func resourceTencentCloudSqlserverConfigInstanceRoGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_instance_ro_group.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
