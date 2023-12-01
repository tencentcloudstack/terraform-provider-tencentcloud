/*
Provides a resource to create a sqlserver config_instance_ro_group

Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-example"
  description = "desc."
}

resource "tencentcloud_sqlserver_basic_instance" "example" {
  name                   = "tf-example"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = tencentcloud_vpc.vpc.id
  subnet_id              = tencentcloud_subnet.subnet.id
  project_id             = 0
  memory                 = 4
  storage                = 100
  cpu                    = 2
  machine_type           = "CLOUD_PREMIUM"
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  security_groups        = [tencentcloud_security_group.security_group.id]

  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_sqlserver_readonly_instance" "example" {
  name                     = "tf_example"
  availability_zone        = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type              = "POSTPAID_BY_HOUR"
  vpc_id                   = tencentcloud_vpc.vpc.id
  subnet_id                = tencentcloud_subnet.subnet.id
  memory                   = 4
  storage                  = 20
  master_instance_id       = tencentcloud_sqlserver_basic_instance.example.id
  readonly_group_type      = 2
  read_only_group_name     = "tf_example_ro"
  is_offline_delay         = 1
  read_only_max_delay_time = 10
  min_read_only_in_group   = 0
  force_upgrade            = true
}

resource "tencentcloud_sqlserver_config_instance_ro_group" "example" {
  instance_id              = tencentcloud_sqlserver_readonly_instance.example.master_instance_id
  read_only_group_id       = tencentcloud_sqlserver_readonly_instance.example.readonly_group_id
  read_only_group_name     = "tf_example_ro_update"
  is_offline_delay         = 1
  read_only_max_delay_time = 5
  min_read_only_in_group   = 1
}
```

Import

sqlserver config_instance_ro_group can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_config_instance_ro_group.example mssql-ds1xhnt9#mssqlro-o6dv2ugx#0#0
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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

	var (
		instanceId      string
		readOnlyGroupId string
		autoWeight      string
		balanceWeight   string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("read_only_group_id"); ok {
		readOnlyGroupId = v.(string)
	}

	if v, ok := d.GetOk("auto_weight"); ok {
		autoWeight = strconv.Itoa(v.(int))
	} else {
		autoWeight = "0"
	}

	if v, ok := d.GetOk("balance_weight"); ok {
		balanceWeight = strconv.Itoa(v.(int))
	} else {
		balanceWeight = "0"
	}

	d.SetId(strings.Join([]string{instanceId, readOnlyGroupId, autoWeight, balanceWeight}, FILED_SP))

	return resourceTencentCloudSqlserverConfigInstanceRoGroupUpdate(d, meta)
}

func resourceTencentCloudSqlserverConfigInstanceRoGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_instance_ro_group.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}
	instanceId := idSplit[0]
	readOnlyGroupId := idSplit[1]
	autoWeight := idSplit[2]
	balanceWeight := idSplit[3]

	configInstanceRoGroup, err := service.DescribeSqlserverConfigInstanceRoGroupById(ctx, instanceId, readOnlyGroupId)
	if err != nil {
		return err
	}

	if configInstanceRoGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverConfigInstanceRoGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if configInstanceRoGroup.MasterInstanceId != nil {
		_ = d.Set("instance_id", configInstanceRoGroup.MasterInstanceId)
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

	if autoWeight == "0" {
		if configInstanceRoGroup.ReadOnlyInstanceSet != nil {
			weightPairsList := []interface{}{}
			for _, weightPairs := range configInstanceRoGroup.ReadOnlyInstanceSet {
				weightPairsMap := map[string]interface{}{}

				if weightPairs.InstanceId != nil {
					weightPairsMap["read_only_instance_id"] = weightPairs.InstanceId
				}

				if weightPairs.Weight != nil {
					weightPairsMap["read_only_weight"] = weightPairs.Weight
				}

				weightPairsList = append(weightPairsList, weightPairsMap)
			}

			_ = d.Set("weight_pairs", weightPairsList)

		}
	}

	tmpAutoWeight, _ := strconv.Atoi(autoWeight)
	tmpBalanceWeight, _ := strconv.Atoi(balanceWeight)
	_ = d.Set("auto_weight", tmpAutoWeight)
	_ = d.Set("balance_weight", tmpBalanceWeight)

	return nil
}

func resourceTencentCloudSqlserverConfigInstanceRoGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_instance_ro_group.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		request = sqlserver.NewModifyReadOnlyGroupDetailsRequest()
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}
	instanceId := idSplit[0]
	readOnlyGroupId := idSplit[1]

	request.InstanceId = &instanceId
	request.ReadOnlyGroupId = &readOnlyGroupId

	if v, ok := d.GetOk("read_only_group_name"); ok {
		request.ReadOnlyGroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("is_offline_delay"); ok {
		request.IsOfflineDelay = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("read_only_max_delay_time"); ok {
		request.ReadOnlyMaxDelayTime = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("min_read_only_in_group"); ok {
		request.MinReadOnlyInGroup = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("weight_pairs"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			parameter := sqlserver.ReadOnlyInstanceWeightPair{}
			if v, ok := dMap["read_only_instance_id"]; ok {
				parameter.ReadOnlyInstanceId = helper.String(v.(string))
			}
			if v, ok := dMap["read_only_weight"]; ok {
				parameter.ReadOnlyWeight = helper.IntInt64(v.(int))
			}
			request.WeightPairs = append(request.WeightPairs, &parameter)
		}
	}

	if v, ok := d.GetOk("auto_weight"); ok {
		request.AutoWeight = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("balance_weight"); ok {
		request.BalanceWeight = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().ModifyReadOnlyGroupDetails(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("sqlserver configInstanceRoGroup not exists")
			return resource.NonRetryableError(e)
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
