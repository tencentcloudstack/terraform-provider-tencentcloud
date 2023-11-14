/*
Provides a resource to create a dcdb db_instance

Example Usage

```hcl
resource "tencentcloud_dcdb_db_instance" "db_instance" {
  zones = &lt;nil&gt;
  period = &lt;nil&gt;
  shard_memory = &lt;nil&gt;
  shard_storage = &lt;nil&gt;
  shard_node_count = &lt;nil&gt;
  shard_count = &lt;nil&gt;
  count = &lt;nil&gt;
  project_id = &lt;nil&gt;
  vpc_id = &lt;nil&gt;
  subnet_id = &lt;nil&gt;
  db_version_id = "5.7.17"
  auto_voucher = &lt;nil&gt;
  voucher_ids = &lt;nil&gt;
  security_group_id = &lt;nil&gt;
  instance_name = &lt;nil&gt;
  ipv6_flag = &lt;nil&gt;
  resource_tags {
		tag_key = &lt;nil&gt;
		tag_value = &lt;nil&gt;

  }
  init_params {
		param = &lt;nil&gt;
		value = ""

  }
  dcn_region = &lt;nil&gt;
  dcn_instance_id = &lt;nil&gt;
  auto_renew_flag = &lt;nil&gt;
  security_group_ids = &lt;nil&gt;
}
```

Import

dcdb db_instance can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_db_instance.db_instance db_instance_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudDcdbDbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbDbInstanceCreate,
		Read:   resourceTencentCloudDcdbDbInstanceRead,
		Update: resourceTencentCloudDcdbDbInstanceUpdate,
		Delete: resourceTencentCloudDcdbDbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zones": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The availability zone distribution of shard nodes can be filled with up to two availability zones. When the shard specification is one master and two slaves, two of the nodes are in the first availability zone.Note that the current availability zone that can be sold needs to be pulled through the DescribeDCDBSaleInfo interface.",
			},

			"period": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The length of time you want to buy, unit: month.",
			},

			"shard_memory": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Shard memory size, unit: GB, can pass DescribeShardSpecQuery the instance specification to obtain.",
			},

			"shard_storage": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Shard storage size, unit: GB, can pass DescribeShardSpecQuery the instance specification to obtain.",
			},

			"shard_node_count": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Number of single shard nodes, can pass DescribeShardSpecQuery the instance specification to obtain.",
			},

			"shard_count": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The number of instance fragments, the optional range is 2-8, and new fragments can be added to a maximum of 64 fragments by upgrading the instance.",
			},

			"count": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The number of instances to be purchased.",
			},

			"project_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Project ID, which can be obtained by viewing the project list, if not passed, it will be associated with the default project.",
			},

			"vpc_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Virtual private network ID, if not passed or passed empty, it means that it is created as a basic network.",
			},

			"subnet_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Virtual private network subnet ID, required when VpcId is not empty.",
			},

			"db_version_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Database engine version, currently available: 8.0.18, 10.1.9, 5.7.17.8.0.18 - MySQL 8.0.18;10.1.9 - Mariadb 10.1.9;5.7.17 - Percona 5.7.17If not filled, the default is 5.7.17, which means Percona 5.7.17.",
			},

			"auto_voucher": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to automatically use vouchers for payment, not used by default.",
			},

			"voucher_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Voucher ID list, currently only supports specifying one voucher.",
			},

			"security_group_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The id of security group.",
			},

			"instance_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Instance name, you can set the name of the instance independently through this field.",
			},

			"ipv6_flag": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to support IPv6.",
			},

			"resource_tags": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Array of tag key-value pairs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The key of tag.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The value of tag.",
						},
					},
				},
			},

			"init_params": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Parameter list. The optional values of this interface are:character_set_server (character set, must be passed),lower_case_table_names (table name is case sensitive, must be passed, 0 - sensitive; 1 - insensitive),innodb_page_size (innodb data page, default 16K),sync_mode ( Synchronous mode: 0 - asynchronous; 1 - strong synchronous; 2 - strong synchronous degenerate. The default is strong synchronous degenerate).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of parameter.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The value of parameter.",
						},
					},
				},
			},

			"dcn_region": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "DCN source region.",
			},

			"dcn_instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "DCN source instance ID.",
			},

			"auto_renew_flag": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Automatic renewal flag, 0 means the default state (the user has not set it, that is, the initial state is manual renewal, and the user has activated the prepaid non-stop privilege and will also perform automatic renewal).1 means automatic renewal, 2 means no automatic renewal (user setting).if the business has no concept of renewal or automatic renewal is not required, it needs to be set to 0.",
			},

			"security_group_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security group ids, the security group can be passed in the form of an array, compatible with the previous SecurityGroupId parameter.",
			},
		},
	}
}

func resourceTencentCloudDcdbDbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_db_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		createDCDBInstanceRequest  = dcdb.NewCreateDCDBInstanceRequest()
		createDCDBInstanceResponse = dcdb.NewCreateDCDBInstanceResponse()
	)
	if v, ok := d.GetOk("zones"); ok {
		zonesSet := v.(*schema.Set).List()
		for i := range zonesSet {
			zones := zonesSet[i].(string)
			request.Zones = append(request.Zones, &zones)
		}
	}

	if v, ok := d.GetOkExists("period"); ok {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("shard_memory"); ok {
		request.ShardMemory = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("shard_storage"); ok {
		request.ShardStorage = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("shard_node_count"); ok {
		request.ShardNodeCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("shard_count"); ok {
		request.ShardCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("count"); ok {
		request.Count = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("project_id"); ok {
		request.ProjectId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_version_id"); ok {
		request.DbVersionId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("auto_voucher"); ok {
		request.AutoVoucher = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("voucher_ids"); ok {
		voucherIdsSet := v.(*schema.Set).List()
		for i := range voucherIdsSet {
			voucherIds := voucherIdsSet[i].(string)
			request.VoucherIds = append(request.VoucherIds, &voucherIds)
		}
	}

	if v, ok := d.GetOk("security_group_id"); ok {
		request.SecurityGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("ipv6_flag"); ok {
		request.Ipv6Flag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("resource_tags"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			resourceTag := dcdb.ResourceTag{}
			if v, ok := dMap["tag_key"]; ok {
				resourceTag.TagKey = helper.String(v.(string))
			}
			if v, ok := dMap["tag_value"]; ok {
				resourceTag.TagValue = helper.String(v.(string))
			}
			request.ResourceTags = append(request.ResourceTags, &resourceTag)
		}
	}

	if v, ok := d.GetOk("init_params"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			dBParamValue := dcdb.DBParamValue{}
			if v, ok := dMap["param"]; ok {
				dBParamValue.Param = helper.String(v.(string))
			}
			if v, ok := dMap["value"]; ok {
				dBParamValue.Value = helper.String(v.(string))
			}
			request.InitParams = append(request.InitParams, &dBParamValue)
		}
	}

	if v, ok := d.GetOk("dcn_region"); ok {
		request.DcnRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dcn_instance_id"); ok {
		request.DcnInstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("auto_renew_flag"); ok {
		request.AutoRenewFlag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		for i := range securityGroupIdsSet {
			securityGroupIds := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().CreateDCDBInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dcdb dbInstance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudDcdbDbInstanceRead(d, meta)
}

func resourceTencentCloudDcdbDbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_db_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	dbInstanceId := d.Id()

	dbInstance, err := service.DescribeDcdbDbInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if dbInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DcdbDbInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if dbInstance.Zones != nil {
		_ = d.Set("zones", dbInstance.Zones)
	}

	if dbInstance.Period != nil {
		_ = d.Set("period", dbInstance.Period)
	}

	if dbInstance.ShardMemory != nil {
		_ = d.Set("shard_memory", dbInstance.ShardMemory)
	}

	if dbInstance.ShardStorage != nil {
		_ = d.Set("shard_storage", dbInstance.ShardStorage)
	}

	if dbInstance.ShardNodeCount != nil {
		_ = d.Set("shard_node_count", dbInstance.ShardNodeCount)
	}

	if dbInstance.ShardCount != nil {
		_ = d.Set("shard_count", dbInstance.ShardCount)
	}

	if dbInstance.Count != nil {
		_ = d.Set("count", dbInstance.Count)
	}

	if dbInstance.ProjectId != nil {
		_ = d.Set("project_id", dbInstance.ProjectId)
	}

	if dbInstance.VpcId != nil {
		_ = d.Set("vpc_id", dbInstance.VpcId)
	}

	if dbInstance.SubnetId != nil {
		_ = d.Set("subnet_id", dbInstance.SubnetId)
	}

	if dbInstance.DbVersionId != nil {
		_ = d.Set("db_version_id", dbInstance.DbVersionId)
	}

	if dbInstance.AutoVoucher != nil {
		_ = d.Set("auto_voucher", dbInstance.AutoVoucher)
	}

	if dbInstance.VoucherIds != nil {
		_ = d.Set("voucher_ids", dbInstance.VoucherIds)
	}

	if dbInstance.SecurityGroupId != nil {
		_ = d.Set("security_group_id", dbInstance.SecurityGroupId)
	}

	if dbInstance.InstanceName != nil {
		_ = d.Set("instance_name", dbInstance.InstanceName)
	}

	if dbInstance.Ipv6Flag != nil {
		_ = d.Set("ipv6_flag", dbInstance.Ipv6Flag)
	}

	if dbInstance.ResourceTags != nil {
		resourceTagsList := []interface{}{}
		for _, resourceTags := range dbInstance.ResourceTags {
			resourceTagsMap := map[string]interface{}{}

			if dbInstance.ResourceTags.TagKey != nil {
				resourceTagsMap["tag_key"] = dbInstance.ResourceTags.TagKey
			}

			if dbInstance.ResourceTags.TagValue != nil {
				resourceTagsMap["tag_value"] = dbInstance.ResourceTags.TagValue
			}

			resourceTagsList = append(resourceTagsList, resourceTagsMap)
		}

		_ = d.Set("resource_tags", resourceTagsList)

	}

	if dbInstance.InitParams != nil {
		initParamsList := []interface{}{}
		for _, initParams := range dbInstance.InitParams {
			initParamsMap := map[string]interface{}{}

			if dbInstance.InitParams.Param != nil {
				initParamsMap["param"] = dbInstance.InitParams.Param
			}

			if dbInstance.InitParams.Value != nil {
				initParamsMap["value"] = dbInstance.InitParams.Value
			}

			initParamsList = append(initParamsList, initParamsMap)
		}

		_ = d.Set("init_params", initParamsList)

	}

	if dbInstance.DcnRegion != nil {
		_ = d.Set("dcn_region", dbInstance.DcnRegion)
	}

	if dbInstance.DcnInstanceId != nil {
		_ = d.Set("dcn_instance_id", dbInstance.DcnInstanceId)
	}

	if dbInstance.AutoRenewFlag != nil {
		_ = d.Set("auto_renew_flag", dbInstance.AutoRenewFlag)
	}

	if dbInstance.SecurityGroupIds != nil {
		_ = d.Set("security_group_ids", dbInstance.SecurityGroupIds)
	}

	return nil
}

func resourceTencentCloudDcdbDbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_db_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dcdb.NewModifyDBInstanceNameRequest()

	dbInstanceId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"zones", "period", "shard_memory", "shard_storage", "shard_node_count", "shard_count", "count", "project_id", "vpc_id", "subnet_id", "db_version_id", "auto_voucher", "voucher_ids", "security_group_id", "instance_name", "ipv6_flag", "resource_tags", "init_params", "dcn_region", "dcn_instance_id", "auto_renew_flag", "security_group_ids"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("instance_name") {
		if v, ok := d.GetOk("instance_name"); ok {
			request.InstanceName = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().ModifyDBInstanceName(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dcdb dbInstance failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDcdbDbInstanceRead(d, meta)
}

func resourceTencentCloudDcdbDbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_db_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	dbInstanceId := d.Id()

	if err := service.DeleteDcdbDbInstanceById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
