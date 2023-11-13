/*
Provides a resource to create a sqlserver general_cloud_ro_instance

Example Usage

```hcl
resource "tencentcloud_sqlserver_general_cloud_ro_instance" "general_cloud_ro_instance" {
  instance_id = ""
  zone = ""
  read_only_group_type =
  memory =
  storage =
  cpu =
  machine_type = ""
  read_only_group_forced_upgrade =
  read_only_group_id = ""
  read_only_group_name = ""
  read_only_group_is_offline_delay =
  read_only_group_max_delay_time =
  read_only_group_min_in_group =
  instance_charge_type = ""
  goods_num =
  subnet_id = ""
  vpc_id = ""
  period =
  security_group_list =
  auto_voucher =
  voucher_ids =
  resource_tags {
		tag_key = ""
		tag_value = ""

  }
  collation = ""
  time_zone = ""
}
```

Import

sqlserver general_cloud_ro_instance can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_general_cloud_ro_instance.general_cloud_ro_instance general_cloud_ro_instance_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudSqlserverGeneralCloudRoInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverGeneralCloudRoInstanceCreate,
		Read:   resourceTencentCloudSqlserverGeneralCloudRoInstanceRead,
		Update: resourceTencentCloudSqlserverGeneralCloudRoInstanceUpdate,
		Delete: resourceTencentCloudSqlserverGeneralCloudRoInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Read-only group type option, 1- Ship according to one instance and one read-only group, 2- Ship after creating a read-only group, all instances are under this read-only group, 3- All instances shipped are in the existing Some read-only groups below.",
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

			"read_only_group_forced_upgrade": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "0 - Default not to upgrade the master instance, 1 - Mandatory upgrade of the master instance to complete ro deployment; if the master instance is a non-cluster version, you need to fill in 1 to force the upgrade to a cluster version. Filling in 1 indicates that you have agreed to upgrade the master instance to a cluster instance.",
			},

			"read_only_group_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Required when ReadOnlyGroupType=3, existing read-only group ID.",
			},

			"read_only_group_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Required when ReadOnlyGroupType=2, the name of the newly created read-only group.",
			},

			"read_only_group_is_offline_delay": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Required when ReadOnlyGroupType=2, whether to enable the delayed elimination function for the newly created read-only group, 1-on, 0-off. When the delay between the read-only replica and the primary instance is greater than the threshold, it will be automatically removed.",
			},

			"read_only_group_max_delay_time": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Mandatory when ReadOnlyGroupType=2 and ReadOnlyGroupIsOfflineDelay=1, the threshold for delay culling of newly created read-only groups.",
			},

			"read_only_group_min_in_group": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Required when ReadOnlyGroupType=2 and ReadOnlyGroupIsOfflineDelay=1, the newly created read-only group retains at least the number of read-only replicas after delay elimination.",
			},

			"instance_charge_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Payment mode, the value supports PREPAID (prepaid), POSTPAID (postpaid).",
			},

			"goods_num": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Purchase several read-only instances this time, and the default value is 1.",
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
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Purchase instance period, the default value is 1, which means one month. The value cannot exceed 48.",
			},

			"security_group_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security group list, fill in the security group ID in the form of sg-xxx.",
			},

			"auto_voucher": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to use the voucher automatically; 1 - yes, 0 - no, not used by default.",
			},

			"voucher_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "An array of voucher IDs, currently only one voucher can be used for a single order.",
			},

			"resource_tags": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "A collection of tags bound to the new instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag value.",
						},
					},
				},
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
		},
	}
}

func resourceTencentCloudSqlserverGeneralCloudRoInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_cloud_ro_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = sqlserver.NewCreateCloudReadOnlyDBInstancesRequest()
		response   = sqlserver.NewCreateCloudReadOnlyDBInstancesResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("read_only_group_type"); ok {
		request.ReadOnlyGroupType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("memory"); ok {
		request.Memory = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("storage"); ok {
		request.Storage = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("cpu"); ok {
		request.Cpu = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("machine_type"); ok {
		request.MachineType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("read_only_group_forced_upgrade"); ok {
		request.ReadOnlyGroupForcedUpgrade = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("read_only_group_id"); ok {
		request.ReadOnlyGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("read_only_group_name"); ok {
		request.ReadOnlyGroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("read_only_group_is_offline_delay"); ok {
		request.ReadOnlyGroupIsOfflineDelay = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("read_only_group_max_delay_time"); ok {
		request.ReadOnlyGroupMaxDelayTime = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("read_only_group_min_in_group"); ok {
		request.ReadOnlyGroupMinInGroup = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("instance_charge_type"); ok {
		request.InstanceChargeType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("goods_num"); ok {
		request.GoodsNum = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("period"); ok {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("security_group_list"); ok {
		securityGroupListSet := v.(*schema.Set).List()
		for i := range securityGroupListSet {
			securityGroupList := securityGroupListSet[i].(string)
			request.SecurityGroupList = append(request.SecurityGroupList, &securityGroupList)
		}
	}

	if v, ok := d.GetOkExists("auto_voucher"); ok {
		request.AutoVoucher = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("voucher_ids"); ok {
		voucherIdsSet := v.(*schema.Set).List()
		for i := range voucherIdsSet {
			voucherIds := voucherIdsSet[i].(string)
			request.VoucherIds = append(request.VoucherIds, &voucherIds)
		}
	}

	if v, ok := d.GetOk("resource_tags"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			resourceTag := sqlserver.ResourceTag{}
			if v, ok := dMap["tag_key"]; ok {
				resourceTag.TagKey = helper.String(v.(string))
			}
			if v, ok := dMap["tag_value"]; ok {
				resourceTag.TagValue = helper.String(v.(string))
			}
			request.ResourceTags = append(request.ResourceTags, &resourceTag)
		}
	}

	if v, ok := d.GetOk("collation"); ok {
		request.Collation = helper.String(v.(string))
	}

	if v, ok := d.GetOk("time_zone"); ok {
		request.TimeZone = helper.String(v.(string))
	}

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

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudSqlserverGeneralCloudRoInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverGeneralCloudRoInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_cloud_ro_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	generalCloudRoInstanceId := d.Id()

	generalCloudRoInstance, err := service.DescribeSqlserverGeneralCloudRoInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if generalCloudRoInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverGeneralCloudRoInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if generalCloudRoInstance.InstanceId != nil {
		_ = d.Set("instance_id", generalCloudRoInstance.InstanceId)
	}

	if generalCloudRoInstance.Zone != nil {
		_ = d.Set("zone", generalCloudRoInstance.Zone)
	}

	if generalCloudRoInstance.ReadOnlyGroupType != nil {
		_ = d.Set("read_only_group_type", generalCloudRoInstance.ReadOnlyGroupType)
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

	if generalCloudRoInstance.MachineType != nil {
		_ = d.Set("machine_type", generalCloudRoInstance.MachineType)
	}

	if generalCloudRoInstance.ReadOnlyGroupForcedUpgrade != nil {
		_ = d.Set("read_only_group_forced_upgrade", generalCloudRoInstance.ReadOnlyGroupForcedUpgrade)
	}

	if generalCloudRoInstance.ReadOnlyGroupId != nil {
		_ = d.Set("read_only_group_id", generalCloudRoInstance.ReadOnlyGroupId)
	}

	if generalCloudRoInstance.ReadOnlyGroupName != nil {
		_ = d.Set("read_only_group_name", generalCloudRoInstance.ReadOnlyGroupName)
	}

	if generalCloudRoInstance.ReadOnlyGroupIsOfflineDelay != nil {
		_ = d.Set("read_only_group_is_offline_delay", generalCloudRoInstance.ReadOnlyGroupIsOfflineDelay)
	}

	if generalCloudRoInstance.ReadOnlyGroupMaxDelayTime != nil {
		_ = d.Set("read_only_group_max_delay_time", generalCloudRoInstance.ReadOnlyGroupMaxDelayTime)
	}

	if generalCloudRoInstance.ReadOnlyGroupMinInGroup != nil {
		_ = d.Set("read_only_group_min_in_group", generalCloudRoInstance.ReadOnlyGroupMinInGroup)
	}

	if generalCloudRoInstance.InstanceChargeType != nil {
		_ = d.Set("instance_charge_type", generalCloudRoInstance.InstanceChargeType)
	}

	if generalCloudRoInstance.GoodsNum != nil {
		_ = d.Set("goods_num", generalCloudRoInstance.GoodsNum)
	}

	if generalCloudRoInstance.SubnetId != nil {
		_ = d.Set("subnet_id", generalCloudRoInstance.SubnetId)
	}

	if generalCloudRoInstance.VpcId != nil {
		_ = d.Set("vpc_id", generalCloudRoInstance.VpcId)
	}

	if generalCloudRoInstance.Period != nil {
		_ = d.Set("period", generalCloudRoInstance.Period)
	}

	if generalCloudRoInstance.SecurityGroupList != nil {
		_ = d.Set("security_group_list", generalCloudRoInstance.SecurityGroupList)
	}

	if generalCloudRoInstance.AutoVoucher != nil {
		_ = d.Set("auto_voucher", generalCloudRoInstance.AutoVoucher)
	}

	if generalCloudRoInstance.VoucherIds != nil {
		_ = d.Set("voucher_ids", generalCloudRoInstance.VoucherIds)
	}

	if generalCloudRoInstance.ResourceTags != nil {
		resourceTagsList := []interface{}{}
		for _, resourceTags := range generalCloudRoInstance.ResourceTags {
			resourceTagsMap := map[string]interface{}{}

			if generalCloudRoInstance.ResourceTags.TagKey != nil {
				resourceTagsMap["tag_key"] = generalCloudRoInstance.ResourceTags.TagKey
			}

			if generalCloudRoInstance.ResourceTags.TagValue != nil {
				resourceTagsMap["tag_value"] = generalCloudRoInstance.ResourceTags.TagValue
			}

			resourceTagsList = append(resourceTagsList, resourceTagsMap)
		}

		_ = d.Set("resource_tags", resourceTagsList)

	}

	if generalCloudRoInstance.Collation != nil {
		_ = d.Set("collation", generalCloudRoInstance.Collation)
	}

	if generalCloudRoInstance.TimeZone != nil {
		_ = d.Set("time_zone", generalCloudRoInstance.TimeZone)
	}

	return nil
}

func resourceTencentCloudSqlserverGeneralCloudRoInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_cloud_ro_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := sqlserver.NewUpgradeDBInstanceRequest()

	generalCloudRoInstanceId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "zone", "read_only_group_type", "memory", "storage", "cpu", "machine_type", "read_only_group_forced_upgrade", "read_only_group_id", "read_only_group_name", "read_only_group_is_offline_delay", "read_only_group_max_delay_time", "read_only_group_min_in_group", "instance_charge_type", "goods_num", "subnet_id", "vpc_id", "period", "security_group_list", "auto_voucher", "voucher_ids", "resource_tags", "collation", "time_zone"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("instance_id") {
		if v, ok := d.GetOk("instance_id"); ok {
			request.InstanceId = helper.String(v.(string))
		}
	}

	if d.HasChange("memory") {
		if v, ok := d.GetOkExists("memory"); ok {
			request.Memory = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("storage") {
		if v, ok := d.GetOkExists("storage"); ok {
			request.Storage = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("cpu") {
		if v, ok := d.GetOkExists("cpu"); ok {
			request.Cpu = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("auto_voucher") {
		if v, ok := d.GetOkExists("auto_voucher"); ok {
			request.AutoVoucher = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("voucher_ids") {
		if v, ok := d.GetOk("voucher_ids"); ok {
			voucherIdsSet := v.(*schema.Set).List()
			for i := range voucherIdsSet {
				voucherIds := voucherIdsSet[i].(string)
				request.VoucherIds = append(request.VoucherIds, &voucherIds)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().UpgradeDBInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver generalCloudRoInstance failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverGeneralCloudRoInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverGeneralCloudRoInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_cloud_ro_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	generalCloudRoInstanceId := d.Id()

	if err := service.DeleteSqlserverGeneralCloudRoInstanceById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
