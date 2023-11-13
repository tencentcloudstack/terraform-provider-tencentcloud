/*
Provides a resource to create a sqlserver business_intelligence_instance

Example Usage

```hcl
resource "tencentcloud_sqlserver_business_intelligence_instance" "business_intelligence_instance" {
  zone = "ap-guangzhou-1"
  memory = 10
  storage = 100
  cpu = 2
  machine_type = "CLOUD_SSD"
  project_id = 0
  goods_num = 1
  subnet_id = "subnet-bdoe83fa"
  vpc_id = "vpc-dsp338hz"
  d_b_version = ""
  security_group_list =
  weekly =
  start_time = ""
  span =
  resource_tags {
		tag_key = ""
		tag_value = ""

  }
}
```

Import

sqlserver business_intelligence_instance can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_business_intelligence_instance.business_intelligence_instance business_intelligence_instance_id
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

func resourceTencentCloudSqlserverBusinessIntelligenceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverBusinessIntelligenceInstanceCreate,
		Read:   resourceTencentCloudSqlserverBusinessIntelligenceInstanceRead,
		Update: resourceTencentCloudSqlserverBusinessIntelligenceInstanceUpdate,
		Delete: resourceTencentCloudSqlserverBusinessIntelligenceInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance AZ, such as ap-guangzhou-1 (Guangzhou Zone 1). Purchasable AZs for an instance can be obtained through theDescribeZones API.",
			},

			"memory": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Instance memory size in GB.",
			},

			"storage": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Instance disk size in GB.",
			},

			"cpu": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The number of CPU cores of the instance you want to purchase.",
			},

			"machine_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The host type of purchased instance. Valid values: CLOUD_PREMIUM (virtual machine with premium cloud disk), CLOUD_SSD (virtual machine with SSD).",
			},

			"project_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Project ID.",
			},

			"goods_num": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Number of instances purchased this time. Default value: 1.",
			},

			"subnet_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "VPC subnet ID in the format of subnet-bdoe83fa. Both SubnetId and VpcId need to be set or unset at the same time.",
			},

			"vpc_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "VPC ID in the format of vpc-dsp338hz. Both SubnetId and VpcId need to be set or unset at the same time.",
			},

			"d_b_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Supported versions of business intelligence server. Valid values: 201603 (SQL Server 2016 Integration Services), 201703 (SQL Server 2017 Integration Services), 201903 (SQL Server 2019 Integration Services). Default value: 201903. As the purchasable versions are region-specific, you can use the DescribeProductConfig API to query the information of purchasable versions in each region.",
			},

			"security_group_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security group list, which contains security group IDs in the format of sg-xxx.",
			},

			"weekly": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Configuration of the maintenance window, which specifies the day of the week when maintenance can be performed. Valid values: 1 (Monday), 2 (Tuesday), 3 (Wednesday), 4 (Thursday), 5 (Friday), 6 (Saturday), 7 (Sunday).",
			},

			"start_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Configuration of the maintenance window, which specifies the start time of daily maintenance.",
			},

			"span": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Configuration of the maintenance window, which specifies the maintenance duration in hours.",
			},

			"resource_tags": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Tags associated with the instances to be created.",
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
		},
	}
}

func resourceTencentCloudSqlserverBusinessIntelligenceInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_business_intelligence_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = sqlserver.NewCreateBusinessDBInstancesRequest()
		response   = sqlserver.NewCreateBusinessDBInstancesResponse()
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
		request.Cpu = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("machine_type"); ok {
		request.MachineType = helper.String(v.(string))
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

	if v, ok := d.GetOk("d_b_version"); ok {
		request.DBVersion = helper.String(v.(string))
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().CreateBusinessDBInstances(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver businessIntelligenceInstance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudSqlserverBusinessIntelligenceInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverBusinessIntelligenceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_business_intelligence_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	businessIntelligenceInstanceId := d.Id()

	businessIntelligenceInstance, err := service.DescribeSqlserverBusinessIntelligenceInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if businessIntelligenceInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverBusinessIntelligenceInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if businessIntelligenceInstance.Zone != nil {
		_ = d.Set("zone", businessIntelligenceInstance.Zone)
	}

	if businessIntelligenceInstance.Memory != nil {
		_ = d.Set("memory", businessIntelligenceInstance.Memory)
	}

	if businessIntelligenceInstance.Storage != nil {
		_ = d.Set("storage", businessIntelligenceInstance.Storage)
	}

	if businessIntelligenceInstance.Cpu != nil {
		_ = d.Set("cpu", businessIntelligenceInstance.Cpu)
	}

	if businessIntelligenceInstance.MachineType != nil {
		_ = d.Set("machine_type", businessIntelligenceInstance.MachineType)
	}

	if businessIntelligenceInstance.ProjectId != nil {
		_ = d.Set("project_id", businessIntelligenceInstance.ProjectId)
	}

	if businessIntelligenceInstance.GoodsNum != nil {
		_ = d.Set("goods_num", businessIntelligenceInstance.GoodsNum)
	}

	if businessIntelligenceInstance.SubnetId != nil {
		_ = d.Set("subnet_id", businessIntelligenceInstance.SubnetId)
	}

	if businessIntelligenceInstance.VpcId != nil {
		_ = d.Set("vpc_id", businessIntelligenceInstance.VpcId)
	}

	if businessIntelligenceInstance.DBVersion != nil {
		_ = d.Set("d_b_version", businessIntelligenceInstance.DBVersion)
	}

	if businessIntelligenceInstance.SecurityGroupList != nil {
		_ = d.Set("security_group_list", businessIntelligenceInstance.SecurityGroupList)
	}

	if businessIntelligenceInstance.Weekly != nil {
		_ = d.Set("weekly", businessIntelligenceInstance.Weekly)
	}

	if businessIntelligenceInstance.StartTime != nil {
		_ = d.Set("start_time", businessIntelligenceInstance.StartTime)
	}

	if businessIntelligenceInstance.Span != nil {
		_ = d.Set("span", businessIntelligenceInstance.Span)
	}

	if businessIntelligenceInstance.ResourceTags != nil {
		resourceTagsList := []interface{}{}
		for _, resourceTags := range businessIntelligenceInstance.ResourceTags {
			resourceTagsMap := map[string]interface{}{}

			if businessIntelligenceInstance.ResourceTags.TagKey != nil {
				resourceTagsMap["tag_key"] = businessIntelligenceInstance.ResourceTags.TagKey
			}

			if businessIntelligenceInstance.ResourceTags.TagValue != nil {
				resourceTagsMap["tag_value"] = businessIntelligenceInstance.ResourceTags.TagValue
			}

			resourceTagsList = append(resourceTagsList, resourceTagsMap)
		}

		_ = d.Set("resource_tags", resourceTagsList)

	}

	return nil
}

func resourceTencentCloudSqlserverBusinessIntelligenceInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_business_intelligence_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := sqlserver.NewModifyDBInstanceNameRequest()

	businessIntelligenceInstanceId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"zone", "memory", "storage", "cpu", "machine_type", "project_id", "goods_num", "subnet_id", "vpc_id", "d_b_version", "security_group_list", "weekly", "start_time", "span", "resource_tags"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().ModifyDBInstanceName(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver businessIntelligenceInstance failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverBusinessIntelligenceInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverBusinessIntelligenceInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_business_intelligence_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	businessIntelligenceInstanceId := d.Id()

	if err := service.DeleteSqlserverBusinessIntelligenceInstanceById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
