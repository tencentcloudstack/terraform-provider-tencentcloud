/*
Provides a resource to create a sqlserver general_cloud_ro_instance

Example Usage

If `read_only_group_type` is 1:

resource "tencentcloud_sqlserver_general_cloud_ro_instance" "general_cloud_ro_instance" {
  instance_id          = "mssql-gyg9xycl"
  zone                 = "ap-guangzhou-6"
  read_only_group_type = 1
  memory               = 4
  storage              = 100
  cpu                  = 2
  machine_type         = "CLOUD_BSSD"
  instance_charge_type = "POSTPAID"
  subnet_id            = "subnet-dwj7ipnc"
  vpc_id               = "vpc-4owdpnwr"
  security_group_list  = ["sg-7kpsbxdb"]
  collation            = "Chinese_PRC_CI_AS"
  time_zone            = "China Standard Time"
  resource_tags        = {
    test-key1 = "test-value1"
    test-key2 = "test-value2"
  }
}

If `read_only_group_type` is 2:

```hcl
resource "tencentcloud_sqlserver_general_cloud_ro_instance" "general_cloud_ro_instance" {
  instance_id                      = "mssql-gyg9xycl"
  zone                             = "ap-guangzhou-6"
  read_only_group_type             = 2
  read_only_group_name             = "test-ro-group"
  read_only_group_is_offline_delay = 1
  read_only_group_max_delay_time   = 10
  read_only_group_min_in_group     = 1
  memory                           = 4
  storage                          = 100
  cpu                              = 2
  machine_type                     = "CLOUD_BSSD"
  instance_charge_type             = "POSTPAID"
  subnet_id                        = "subnet-dwj7ipnc"
  vpc_id                           = "vpc-4owdpnwr"
  security_group_list              = ["sg-7kpsbxdb"]
  collation                        = "Chinese_PRC_CI_AS"
  time_zone                        = "China Standard Time"
  resource_tags                    = {
    test-key1 = "test-value1"
    test-key2 = "test-value2"
  }
}
```

If `read_only_group_type` is 3:

```hcl
resource "tencentcloud_sqlserver_general_cloud_ro_instance" "general_cloud_ro_instance" {
  instance_id          = "mssql-gyg9xycl"
  zone                 = "ap-guangzhou-6"
  read_only_group_type = 3
  memory               = 4
  storage              = 100
  cpu                  = 2
  machine_type         = "CLOUD_BSSD"
  read_only_group_id   = "mssqlrg-clboghrj"
  instance_charge_type = "POSTPAID"
  subnet_id            = "subnet-dwj7ipnc"
  vpc_id               = "vpc-4owdpnwr"
  security_group_list  = ["sg-7kpsbxdb"]
  collation            = "Chinese_PRC_CI_AS"
  time_zone            = "China Standard Time"
  resource_tags        = {
    test-key1 = "test-value1"
    test-key2 = "test-value2"
  }
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSqlserverGeneralCloudRoInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverGeneralCloudRoInstanceCreate,
		Read:   resourceTencentCloudSqlserverGeneralCloudRoInstanceRead,
		Update: resourceTencentCloudSqlserverGeneralCloudRoInstanceUpdate,
		Delete: resourceTencentCloudSqlserverGeneralCloudRoInstanceDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Primary instance ID, in the format: mssql-3l3fgqn7.",
			},
			"zone": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance Availability Zone, similar to ap-guangzhou-1 (Guangzhou District 1); the instance sales area can be obtained through the interface DescribeZones.",
			},
			"read_only_group_type": {
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validateIntegerInRange(1, 3),
				Description:  "Read-only group type option, 1- Ship according to one instance and one read-only group, 2- Ship after creating a read-only group, all instances are under this read-only group, 3- All instances shipped are in the existing Some read-only groups below.",
			},
			"memory": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Instance memory size, in GB.",
			},
			"storage": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Instance disk size, in GB.",
			},
			"cpu": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Number of instance cores.",
			},
			"machine_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The host disk type of the purchased instance, CLOUD_HSSD-enhanced SSD cloud disk for virtual machines, CLOUD_TSSD-extremely fast SSD cloud disk for virtual machines, CLOUD_BSSD-universal SSD cloud disk for virtual machines.",
			},
			"read_only_group_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Required when ReadOnlyGroupType=3, existing read-only group ID.",
			},
			"read_only_group_name": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Required when ReadOnlyGroupType=2, the name of the newly created read-only group.",
			},
			"read_only_group_is_offline_delay": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Required when ReadOnlyGroupType=2, whether to enable the delayed elimination function for the newly created read-only group, 1-on, 0-off. When the delay between the read-only replica and the primary instance is greater than the threshold, it will be automatically removed.",
			},
			"read_only_group_max_delay_time": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Mandatory when ReadOnlyGroupType=2 and ReadOnlyGroupIsOfflineDelay=1, the threshold for delay culling of newly created read-only groups.",
			},
			"read_only_group_min_in_group": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Required when ReadOnlyGroupType=2 and ReadOnlyGroupIsOfflineDelay=1, the newly created read-only group retains at least the number of read-only replicas after delay elimination.",
			},
			"instance_charge_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Payment mode, the value supports PREPAID (prepaid), POSTPAID (postpaid).",
			},
			"subnet_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "VPC subnet ID, in the form of subnet-bdoe83fa; SubnetId and VpcId need to be set at the same time or not set at the same time.",
			},
			"vpc_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "VPC network ID, in the form of vpc-dsp338hz; SubnetId and VpcId need to be set at the same time or not set at the same time.",
			},
			"period": {
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validateIntegerInRange(1, 48),
				Description:  "Purchase instance period, the default value is 1, which means one month. The value cannot exceed 48.",
			},
			"security_group_list": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Security group list, fill in the security group ID in the form of sg-xxx.",
			},
			"resource_tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
			"collation": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "System character set collation, default: Chinese_PRC_CI_AS.",
			},
			"time_zone": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "System time zone, default: China Standard Time.",
			},
			"ro_instance_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Primary read only instance ID, in the format: mssqlro-lbljc5qd.",
			},
		},
	}
}

func resourceTencentCloudSqlserverGeneralCloudRoInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_cloud_ro_instance.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId        = getLogId(contextNil)
		ctx          = context.WithValue(context.TODO(), logIdKey, logId)
		service      = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		request      = sqlserver.NewCreateCloudReadOnlyDBInstancesRequest()
		response     = sqlserver.NewCreateCloudReadOnlyDBInstancesResponse()
		instanceId   string
		roInstanceId string
		dealNames    string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("read_only_group_type"); ok {
		roGroupType := v.(int)
		if roGroupType == 1 {
			request.ReadOnlyGroupForcedUpgrade = helper.IntInt64(1)
		}
		if roGroupType == 2 {
			if v, ok := d.GetOk("read_only_group_name"); ok {
				request.ReadOnlyGroupName = helper.String(v.(string))
			}

			if v, ok := d.GetOkExists("read_only_group_is_offline_delay"); ok {
				readOnlyGroupIsOfflineDelay := v.(int)
				if readOnlyGroupIsOfflineDelay == 1 {
					if v, ok := d.GetOk("read_only_group_max_delay_time"); ok {
						request.ReadOnlyGroupMaxDelayTime = helper.IntInt64(v.(int))
					}

					if v, ok := d.GetOk("read_only_group_min_in_group"); ok {
						request.ReadOnlyGroupMinInGroup = helper.IntInt64(v.(int))
					}
				}
				request.ReadOnlyGroupIsOfflineDelay = helper.IntInt64(readOnlyGroupIsOfflineDelay)
			}

		} else if roGroupType == 3 {
			if v, ok := d.GetOk("read_only_group_id"); ok {
				request.ReadOnlyGroupId = helper.String(v.(string))
			}
			request.ReadOnlyGroupForcedUpgrade = helper.IntInt64(1)
		}
		request.ReadOnlyGroupType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("memory"); ok {
		request.Memory = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("storage"); ok {
		request.Storage = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("cpu"); ok {
		request.Cpu = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("machine_type"); ok {
		request.MachineType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_charge_type"); ok {
		request.InstanceChargeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("period"); ok {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("security_group_list"); ok {
		securityGroupListSet := v.(*schema.Set).List()
		for i := range securityGroupListSet {
			securityGroupList := securityGroupListSet[i].(string)
			request.SecurityGroupList = append(request.SecurityGroupList, &securityGroupList)
		}
	}

	if v, ok := d.GetOk("collation"); ok {
		request.Collation = helper.String(v.(string))
	}

	if v, ok := d.GetOk("time_zone"); ok {
		request.TimeZone = helper.String(v.(string))
	}

	request.GoodsNum = helper.IntInt64(1)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().CreateCloudReadOnlyDBInstances(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver generalCloudRoInstance failed, reason:%+v", logId, err)
		return err
	}

	dealNames = *response.Response.DealNames[0]
	roInstanceId, err = service.GetInfoFromDeal(ctx, dealNames)
	if err != nil {
		return err
	}

	if tags := helper.GetTags(d, "resource_tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::sqlserver:%s:uin/:instance/%s", region, roInstanceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	d.SetId(strings.Join([]string{instanceId, roInstanceId}, FILED_SP))

	return resourceTencentCloudSqlserverGeneralCloudRoInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverGeneralCloudRoInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_cloud_ro_instance.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	instanceId := idSplit[0]
	roInstanceId := idSplit[1]

	generalCloudRoInstance, err := service.DescribeSqlserverGeneralCloudRoInstanceById(ctx, roInstanceId)
	if err != nil {
		return err
	}

	if generalCloudRoInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverGeneralCloudRoInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if generalCloudRoInstance.InstanceId != nil {
		_ = d.Set("ro_instance_id", generalCloudRoInstance.InstanceId)
	}

	if generalCloudRoInstance.Zone != nil {
		_ = d.Set("zone", generalCloudRoInstance.Zone)
	}

	if generalCloudRoInstance.Memory != nil {
		_ = d.Set("memory", generalCloudRoInstance.Memory)
	}

	if generalCloudRoInstance.Storage != nil {
		_ = d.Set("storage", generalCloudRoInstance.Storage)
	}

	if generalCloudRoInstance.Cpu != nil {
		_ = d.Set("cpu", generalCloudRoInstance.Cpu)
	}

	if generalCloudRoInstance.Type != nil {
		_ = d.Set("machine_type", generalCloudRoInstance.Type)
	}

	if generalCloudRoInstance.PayMode != nil {
		if *generalCloudRoInstance.PayMode == 0 {
			_ = d.Set("instance_charge_type", SQLSERVER_TYPE_POSTPAID)
		} else {
			_ = d.Set("instance_charge_type", SQLSERVER_TYPE_PREPAID)
		}
	}

	if generalCloudRoInstance.UniqSubnetId != nil {
		_ = d.Set("subnet_id", generalCloudRoInstance.UniqSubnetId)
	}

	if generalCloudRoInstance.UniqVpcId != nil {
		_ = d.Set("vpc_id", generalCloudRoInstance.UniqVpcId)
	}

	if generalCloudRoInstance.Collation != nil {
		_ = d.Set("collation", generalCloudRoInstance.Collation)
	}

	if generalCloudRoInstance.TimeZone != nil {
		_ = d.Set("time_zone", generalCloudRoInstance.TimeZone)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "sqlserver", "instance", tcClient.Region, roInstanceId)
	if err != nil {
		return err
	}

	_ = d.Set("resource_tags", tags)

	securityGroupList, err := service.DescribeInstanceSecurityGroups(ctx, roInstanceId)
	if err != nil {
		return err
	}

	if securityGroupList == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlservereSecurityGroups` [%s] not found, please check if it has been deleted.", logId, d.Id())
		return nil
	}

	if securityGroupList != nil {
		_ = d.Set("security_group_list", securityGroupList)
	}

	roGroupList, err := service.DescribeReadonlyGroupList(ctx, instanceId)
	if err != nil {
		return err
	}

	if roGroupList == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlservereReadonlyGroup` [%s] not found, please check if it has been deleted.", logId, d.Id())
		return nil
	}

	for _, v := range roGroupList {
		readOnlyInstanceSet := v.ReadOnlyInstanceSet
		for _, readOnlyInstance := range readOnlyInstanceSet {
			if roInstanceId == *readOnlyInstance.InstanceId {
				if v.ReadOnlyGroupId != nil {
					_ = d.Set("read_only_group_id", v.ReadOnlyGroupId)
				}

				if v.ReadOnlyGroupName != nil {
					_ = d.Set("read_only_group_name", v.ReadOnlyGroupName)
				}

				if v.IsOfflineDelay != nil {
					_ = d.Set("read_only_group_is_offline_delay", v.IsOfflineDelay)
				}

				if v.ReadOnlyMaxDelayTime != nil {
					_ = d.Set("read_only_group_max_delay_time", v.ReadOnlyMaxDelayTime)
				}

				if v.MinReadOnlyInGroup != nil {
					_ = d.Set("read_only_group_min_in_group", v.MinReadOnlyInGroup)
				}
			}
		}
	}

	return nil
}

func resourceTencentCloudSqlserverGeneralCloudRoInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_cloud_ro_instance.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId            = getLogId(contextNil)
		ctx              = context.WithValue(context.TODO(), logIdKey, logId)
		client           = meta.(*TencentCloudClient).apiV3Conn
		sqlserverService = SqlserverService{client: client}
		request          = sqlserver.NewUpgradeDBInstanceRequest()
		dealId           string
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	roInstanceId := idSplit[1]

	immutableArgs := []string{"instance_id", "zone", "read_only_group_type", "machine_type", "read_only_group_forced_upgrade", "read_only_group_id", "read_only_group_name", "read_only_group_is_offline_delay", "read_only_group_max_delay_time", "read_only_group_min_in_group", "instance_charge_type", "subnet_id", "vpc_id", "period", "security_group_list", "auto_voucher", "voucher_ids", "resource_tags", "collation", "time_zone"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	request.InstanceId = &roInstanceId
	request.WaitSwitch = helper.IntInt64(0)

	if d.HasChange("memory") || d.HasChange("storage") || d.HasChange("cpu") {
		if v, ok := d.GetOkExists("memory"); ok {
			request.Memory = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOkExists("storage"); ok {
			request.Storage = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOkExists("cpu"); ok {
			request.Cpu = helper.IntInt64(v.(int))
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().UpgradeDBInstance(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			dealId = *result.Response.DealName
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update sqlserver generalCloudRoInstance failed, reason:%+v", logId, err)
			return err
		}

		_, err = sqlserverService.GetInfoFromDeal(ctx, dealId)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudSqlserverGeneralCloudRoInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverGeneralCloudRoInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_cloud_ro_instance.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	roInstanceId := idSplit[1]

	if err := service.TerminateSqlserverInstance(ctx, roInstanceId); err != nil {
		return err
	}

	if err := service.DeleteSqlserverGeneralCloudRoInstanceById(ctx, roInstanceId); err != nil {
		return err
	}

	return nil
}
