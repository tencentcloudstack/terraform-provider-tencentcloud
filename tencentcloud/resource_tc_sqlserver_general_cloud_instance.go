/*
Provides a resource to create a sqlserver general_cloud_instance

Example Usage

```hcl
resource "tencentcloud_sqlserver_general_cloud_instance" "general_cloud_instance" {
  zone = "ap-guangzhou-1"
  memory =
  storage =
  cpu =
  machine_type = "CLOUD_SSD"
  instance_charge_type = "postpaid"
  project_id =
  goods_num = 1
  subnet_id = "subnet-bdoe83fa"
  vpc_id = "vpc-dsp338hz"
  period =
  auto_voucher =
  voucher_ids =
  d_b_version = ""
  auto_renew_flag =
  security_group_list =
  weekly =
  start_time = ""
  span =
  multi_zones =
  resource_tags {
		tag_key = ""
		tag_value = ""

  }
  collation = ""
  time_zone = ""
}
```

Import

sqlserver general_cloud_instance can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_general_cloud_instance.general_cloud_instance general_cloud_instance_id
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

func resourceTencentCloudSqlserverGeneralCloudInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverGeneralCloudInstanceCreate,
		Read:   resourceTencentCloudSqlserverGeneralCloudInstanceRead,
		Update: resourceTencentCloudSqlserverGeneralCloudInstanceUpdate,
		Delete: resourceTencentCloudSqlserverGeneralCloudInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance AZ, such as ap-guangzhou-1 (Guangzhou Zone 1). Purchasable AZs for an instance can be obtained through the DescribeZones API.",
			},

			"memory": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Memory, unit: GB.",
			},

			"storage": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Instance disk storage, unit: GB.",
			},

			"cpu": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Cpu，unit: CORE.",
			},

			"machine_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The host type of purchased instance. Valid values: CLOUD_PREMIUM (virtual machine with premium cloud disk), CLOUD_SSD (virtual machine with SSD).",
			},

			"instance_charge_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Payment mode, the value supports PREPAID (prepaid), POSTPAID (postpaid).",
			},

			"project_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Project ID.",
			},

			"goods_num": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Several instances are purchased this time, and the default value is 1. The value does not exceed 10.",
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

			"d_b_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sqlserver version, currently all supported versions are: 2008R2 (SQL Server 2008 R2 Enterprise), 2012SP3 (SQL Server 2012 Enterprise), 201202 (SQL Server 2012 Standard), 2014SP2 (SQL Server 2014 Enterprise), 201402 (SQL Server 2014 Standard) , 2016SP1 (SQL Server 2016 Enterprise), 201602 (SQL Server 2016 Standard), 2017 (SQL Server 2017 Enterprise), 201702 (SQL Server 2017 Standard), 2019 (SQL Server 2019 Enterprise), 201902 (SQL Server 2019 Standard). Each region supports different versions for sale, and the version information that can be sold in each region can be pulled through the DescribeProductConfig interface. If left blank, the default version is 2008R2.",
			},

			"auto_renew_flag": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Automatic renewal flag: 0-normal renewal 1-automatic renewal, the default is 1 automatic renewal. Valid only when purchasing a prepaid instance.",
			},

			"security_group_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security group list, fill in the security group ID in the form of sg-xxx.",
			},

			"weekly": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Maintainable time window configuration, in weeks, indicates the days of the week that allow maintenance, 1-7 represent Monday to weekend respectively.",
			},

			"start_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Maintainable time window configuration, daily maintainable start time.",
			},

			"span": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Maintainable time window configuration, duration, unit: hour.",
			},

			"multi_zones": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to deploy across availability zones, the default value is false.",
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

func resourceTencentCloudSqlserverGeneralCloudInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_cloud_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = sqlserver.NewCreateCloudDBInstancesRequest()
		response   = sqlserver.NewCreateCloudDBInstancesResponse()
		instanceId string
	)
	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
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

	if v, ok := d.GetOk("instance_charge_type"); ok {
		request.InstanceChargeType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("project_id"); ok {
		request.ProjectId = helper.IntInt64(v.(int))
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

	if v, ok := d.GetOk("d_b_version"); ok {
		request.DBVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("auto_renew_flag"); ok {
		request.AutoRenewFlag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("security_group_list"); ok {
		securityGroupListSet := v.(*schema.Set).List()
		for i := range securityGroupListSet {
			securityGroupList := securityGroupListSet[i].(string)
			request.SecurityGroupList = append(request.SecurityGroupList, &securityGroupList)
		}
	}

	if v, ok := d.GetOk("weekly"); ok {
		weeklySet := v.(*schema.Set).List()
		for i := range weeklySet {
			weekly := weeklySet[i].(int)
			request.Weekly = append(request.Weekly, helper.IntInt64(weekly))
		}
	}

	if v, ok := d.GetOk("start_time"); ok {
		request.StartTime = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("span"); ok {
		request.Span = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("multi_zones"); ok {
		request.MultiZones = helper.Bool(v.(bool))
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().CreateCloudDBInstances(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver generalCloudInstance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudSqlserverGeneralCloudInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverGeneralCloudInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_cloud_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	generalCloudInstanceId := d.Id()

	generalCloudInstance, err := service.DescribeSqlserverGeneralCloudInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if generalCloudInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverGeneralCloudInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if generalCloudInstance.Zone != nil {
		_ = d.Set("zone", generalCloudInstance.Zone)
	}

	if generalCloudInstance.Memory != nil {
		_ = d.Set("memory", generalCloudInstance.Memory)
	}

	if generalCloudInstance.Storage != nil {
		_ = d.Set("storage", generalCloudInstance.Storage)
	}

	if generalCloudInstance.Cpu != nil {
		_ = d.Set("cpu", generalCloudInstance.Cpu)
	}

	if generalCloudInstance.MachineType != nil {
		_ = d.Set("machine_type", generalCloudInstance.MachineType)
	}

	if generalCloudInstance.InstanceChargeType != nil {
		_ = d.Set("instance_charge_type", generalCloudInstance.InstanceChargeType)
	}

	if generalCloudInstance.ProjectId != nil {
		_ = d.Set("project_id", generalCloudInstance.ProjectId)
	}

	if generalCloudInstance.GoodsNum != nil {
		_ = d.Set("goods_num", generalCloudInstance.GoodsNum)
	}

	if generalCloudInstance.SubnetId != nil {
		_ = d.Set("subnet_id", generalCloudInstance.SubnetId)
	}

	if generalCloudInstance.VpcId != nil {
		_ = d.Set("vpc_id", generalCloudInstance.VpcId)
	}

	if generalCloudInstance.Period != nil {
		_ = d.Set("period", generalCloudInstance.Period)
	}

	if generalCloudInstance.AutoVoucher != nil {
		_ = d.Set("auto_voucher", generalCloudInstance.AutoVoucher)
	}

	if generalCloudInstance.VoucherIds != nil {
		_ = d.Set("voucher_ids", generalCloudInstance.VoucherIds)
	}

	if generalCloudInstance.DBVersion != nil {
		_ = d.Set("d_b_version", generalCloudInstance.DBVersion)
	}

	if generalCloudInstance.AutoRenewFlag != nil {
		_ = d.Set("auto_renew_flag", generalCloudInstance.AutoRenewFlag)
	}

	if generalCloudInstance.SecurityGroupList != nil {
		_ = d.Set("security_group_list", generalCloudInstance.SecurityGroupList)
	}

	if generalCloudInstance.Weekly != nil {
		_ = d.Set("weekly", generalCloudInstance.Weekly)
	}

	if generalCloudInstance.StartTime != nil {
		_ = d.Set("start_time", generalCloudInstance.StartTime)
	}

	if generalCloudInstance.Span != nil {
		_ = d.Set("span", generalCloudInstance.Span)
	}

	if generalCloudInstance.MultiZones != nil {
		_ = d.Set("multi_zones", generalCloudInstance.MultiZones)
	}

	if generalCloudInstance.ResourceTags != nil {
		resourceTagsList := []interface{}{}
		for _, resourceTags := range generalCloudInstance.ResourceTags {
			resourceTagsMap := map[string]interface{}{}

			if generalCloudInstance.ResourceTags.TagKey != nil {
				resourceTagsMap["tag_key"] = generalCloudInstance.ResourceTags.TagKey
			}

			if generalCloudInstance.ResourceTags.TagValue != nil {
				resourceTagsMap["tag_value"] = generalCloudInstance.ResourceTags.TagValue
			}

			resourceTagsList = append(resourceTagsList, resourceTagsMap)
		}

		_ = d.Set("resource_tags", resourceTagsList)

	}

	if generalCloudInstance.Collation != nil {
		_ = d.Set("collation", generalCloudInstance.Collation)
	}

	if generalCloudInstance.TimeZone != nil {
		_ = d.Set("time_zone", generalCloudInstance.TimeZone)
	}

	return nil
}

func resourceTencentCloudSqlserverGeneralCloudInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_cloud_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := sqlserver.NewUpgradeDBInstanceRequest()

	generalCloudInstanceId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"zone", "memory", "storage", "cpu", "machine_type", "instance_charge_type", "project_id", "goods_num", "subnet_id", "vpc_id", "period", "auto_voucher", "voucher_ids", "d_b_version", "auto_renew_flag", "security_group_list", "weekly", "start_time", "span", "multi_zones", "resource_tags", "collation", "time_zone"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
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

	if d.HasChange("d_b_version") {
		if v, ok := d.GetOk("d_b_version"); ok {
			request.DBVersion = helper.String(v.(string))
		}
	}

	if d.HasChange("multi_zones") {
		if v, ok := d.GetOkExists("multi_zones"); ok {
			request.MultiZones = helper.Bool(v.(bool))
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
		log.Printf("[CRITAL]%s update sqlserver generalCloudInstance failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverGeneralCloudInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverGeneralCloudInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_cloud_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	generalCloudInstanceId := d.Id()

	if err := service.DeleteSqlserverGeneralCloudInstanceById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
