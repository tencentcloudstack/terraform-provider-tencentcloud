/*
Provides a resource to create a dcdb db_instance

Example Usage

```hcl
resource "tencentcloud_dcdb_db_instance" "db_instance" {
  instance_name = "test_dcdb_db_instance"
  zones = ["ap-guangzhou-5"]
  period = 1
  shard_memory = "2"
  shard_storage = "10"
  shard_node_count = "2"
  shard_count = "2"
  vpc_id = local.vpc_id
  subnet_id = local.subnet_id
  db_version_id = "8.0"
  resource_tags {
	tag_key = "aaa"
	tag_value = "bbb"
  }
  init_params {
	 param = "character_set_server"
	 value = "utf8mb4"
  }
  init_params {
	param = "lower_case_table_names"
	value = "1"
  }
  init_params {
	param = "sync_mode"
	value = "2"
  }
  init_params {
	param = "innodb_page_size"
	value = "16384"
  }
  security_group_ids = [local.sg_id]
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
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
				Description: "&amp;quot;The availability zone distribution of shard nodes can be filled with up to two availability zones. When the shard specification is one master and two slaves, two of the nodes are in the first availability zone.&amp;quot;&amp;quot;Note that the current availability zone that can be sold needs to be pulled through the DescribeDCDBSaleInfo interface.&amp;quot;.",
			},

			"period": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The length of time you want to buy, unit: month.",
			},

			"shard_memory": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "&amp;quot;Shard memory size, unit: GB, can pass DescribeShardSpec&amp;quot;&amp;quot;Query the instance specification to obtain.&amp;quot;.",
			},

			"shard_storage": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "&amp;quot;Shard storage size, unit: GB, can pass DescribeShardSpec&amp;quot;&amp;quot;Query the instance specification to obtain.&amp;quot;.",
			},

			"shard_node_count": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "&amp;quot;Number of single shard nodes, can pass DescribeShardSpec&amp;quot;&amp;quot;Query the instance specification to obtain.&amp;quot;.",
			},

			"shard_count": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The number of instance fragments, the optional range is 2-8, and new fragments can be added to a maximum of 64 fragments by upgrading the instance.",
			},

			// "count": {
			// 	Optional:    true,
			// 	Type:        schema.TypeInt,
			// 	Description: "The number of instances to be purchased.",
			// },

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
				Description: "&amp;quot;Database engine version, currently available: 8.0.18, 10.1.9, 5.7.17.&amp;quot;&amp;quot;8.0.18 - MySQL 8.0.18;&amp;quot;&amp;quot;10.1.9 - Mariadb 10.1.9;&amp;quot;&amp;quot;5.7.17 - Percona 5.7.17&amp;quot;&amp;quot;If not filled, the default is 5.7.17, which means Percona 5.7.17.&amp;quot;.",
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
				Description: "&amp;quot;parameter list. The optional values of this interface are:&amp;quot;&amp;quot;character_set_server (character set, must be passed),&amp;quot;&amp;quot;lower_case_table_names (table name is case sensitive, must be passed, 0 - sensitive; 1 - insensitive),&amp;quot;&amp;quot;innodb_page_size (innodb data page, default 16K),&amp;quot;&amp;quot;sync_mode ( Synchronous mode: 0 - asynchronous; 1 - strong synchronous; 2 - strong synchronous degenerate. The default is strong synchronous degenerate)&amp;quot;.",
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
				Description: "&amp;quot;Automatic renewal flag, 0 means the default state (the user has not set it, that is, the initial state is manual renewal, and the user has activated the prepaid non-stop privilege and will also perform automatic renewal).&amp;quot;&amp;quot;1 means automatic renewal, 2 means no automatic renewal (user setting).&amp;quot;&amp;quot;if the business has no concept of renewal or automatic renewal is not required, it needs to be set to 0.&amp;quot;.",
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
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request    = dcdb.NewCreateDCDBInstanceRequest()
		response   = dcdb.NewCreateDCDBInstanceResponse()
		instanceId string
		service    = DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	)
	if v, ok := d.GetOk("zones"); ok {
		zonesSet := v.(*schema.Set).List()
		request.Zones = helper.InterfacesStringsPoint(zonesSet)
	}

	if v, _ := d.GetOk("period"); v != nil {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("shard_memory"); v != nil {
		request.ShardMemory = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("shard_storage"); v != nil {
		request.ShardStorage = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("shard_node_count"); v != nil {
		request.ShardNodeCount = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("shard_count"); v != nil {
		request.ShardCount = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("project_id"); v != nil {
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

	if v, _ := d.GetOk("auto_voucher"); v != nil {
		request.AutoVoucher = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("voucher_ids"); ok {
		voucherIdsSet := v.(*schema.Set).List()
		request.VoucherIds = helper.InterfacesStringsPoint(voucherIdsSet)
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, _ := d.GetOk("ipv6_flag"); v != nil {
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

	if v, _ := d.GetOk("auto_renew_flag"); v != nil {
		request.AutoRenewFlag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		request.SecurityGroupIds = helper.InterfacesStringsPoint(securityGroupIdsSet)
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

	if response == nil || len(response.Response.InstanceIds) < 1 {
		d.SetId("")
		return fmt.Errorf("[CRITAL]%s create dcdb dbInstance failed.", d.Id())
	}

	instanceId = *response.Response.InstanceIds[0]
	d.SetId(instanceId)

	if len(request.InitParams) < 1 {
		defaultInitParams := []*dcdb.DBParamValue{
			{
				Param: helper.String("character_set_server"),
				Value: helper.String("utf8mb4"),
			},
			{
				Param: helper.String("lower_case_table_names"),
				Value: helper.String("1"),
			},
			{
				Param: helper.String("sync_mode"),
				Value: helper.String("2"),
			},
			{
				Param: helper.String("innodb_page_size"),
				Value: helper.String("16384"),
			},
		}
		request.InitParams = defaultInitParams
	}

	initRet, flowId, e := service.InitDcdbDbInstance(ctx, instanceId, request.InitParams)
	if e != nil {
		return e
	}
	if !initRet {
		return fmt.Errorf("db instance init failed")
	}

	if flowId != nil {
		// need to wait init operation success
		// 0:success; 1:failed, 2:running
		conf := BuildStateChangeConf([]string{}, []string{"0"}, 3*readRetryTimeout, time.Second, service.DcdbDbInstanceStateRefreshFunc(helper.UInt64Int64(*flowId), []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}
	return resourceTencentCloudDcdbDbInstanceRead(d, meta)
}

func resourceTencentCloudDcdbDbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_db_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	ret, err := service.DescribeDcdbDbInstance(ctx, instanceId)
	if err != nil {
		return err
	}

	if ret == nil || len(ret.Instances) < 1 {
		d.SetId("")
		return fmt.Errorf("resource `DcdbDbInstance` %s does not exist", d.Id())
	}

	dbInstance := ret.Instances[0]

	if dbInstance.Zone != nil {
		v, _ := helper.StrToStrList(*dbInstance.Zone)
		_ = d.Set("zones", v)
	}

	// if dbInstance.Period != nil {
	// 	_ = d.Set("period", dbInstance.Period)
	// }

	if dbInstance.ShardDetail[0] != nil { // Memory and Storage is params for one shard
		shard := dbInstance.ShardDetail[0]

		if shard.Memory != nil {
			_ = d.Set("shard_memory", shard.Memory)
		}

		if shard.Storage != nil {
			_ = d.Set("shard_storage", shard.Storage)
		}
	}

	if dbInstance.NodeCount != nil {
		_ = d.Set("shard_node_count", dbInstance.NodeCount)
	}

	if dbInstance.ShardCount != nil {
		_ = d.Set("shard_count", dbInstance.ShardCount)
	}

	if dbInstance.ProjectId != nil {
		_ = d.Set("project_id", dbInstance.ProjectId)
	}

	if dbInstance.VpcId != nil {
		_ = d.Set("vpc_id", helper.Int64ToStrPoint(*dbInstance.VpcId))
	}

	if dbInstance.SubnetId != nil {
		_ = d.Set("subnet_id", helper.Int64ToStrPoint(*dbInstance.SubnetId))
	}

	if dbInstance.DbVersionId != nil {
		_ = d.Set("db_version_id", dbInstance.DbVersionId)
	}

	// if dbInstance.AutoVoucher != nil {
	// 	_ = d.Set("auto_voucher", dbInstance.AutoVoucher)
	// }

	// if dbInstance.VoucherIds != nil {
	// 	_ = d.Set("voucher_ids", dbInstance.VoucherIds)
	// }

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

			if resourceTags.TagKey != nil {
				resourceTagsMap["tag_key"] = resourceTags.TagKey
			}

			if resourceTags.TagValue != nil {
				resourceTagsMap["tag_value"] = resourceTags.TagValue
			}

			resourceTagsList = append(resourceTagsList, resourceTagsMap)
		}

		_ = d.Set("resource_tags", resourceTagsList)

	}

	// if dbInstance.InitParams != nil {
	// 	initParamsList := []interface{}{}
	// 	for _, initParams := range dbInstance.InitParams {
	// 		initParamsMap := map[string]interface{}{}

	// 		if dbInstance.InitParams.Param != nil {
	// 			initParamsMap["param"] = dbInstance.InitParams.Param
	// 		}

	// 		if dbInstance.InitParams.Value != nil {
	// 			initParamsMap["value"] = dbInstance.InitParams.Value
	// 		}

	// 		initParamsList = append(initParamsList, initParamsMap)
	// 	}

	// 	_ = d.Set("init_params", initParamsList)

	// }

	if dcn, err := service.DescribeDcnDetailById(ctx, instanceId); dcn != nil {
		_ = d.Set("dcn_region", dcn.Region)
		_ = d.Set("dcn_instance_id", dcn.InstanceId)
	} else {
		return err
	}

	if dbInstance.AutoRenewFlag != nil {
		_ = d.Set("auto_renew_flag", dbInstance.AutoRenewFlag)
	}

	if sg, err := service.DescribeDcdbSecurityGroup(ctx, instanceId); sg != nil {
		sgIds := make([]*string, 0, len(sg.Groups))
		for _, sg := range sg.Groups {
			sgIds = append(sgIds, sg.SecurityGroupId)
		}
		_ = d.Set("security_group_ids", sgIds)
	} else {
		return err
	}

	return nil
}

func resourceTencentCloudDcdbDbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_db_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dcdb.NewModifyDBInstanceNameRequest()

	instanceId := d.Id()

	request.InstanceId = helper.String(instanceId)
	if d.HasChange("zones") {
		return fmt.Errorf("`zones` do not support change now.")
	}
	if d.HasChange("period") {
		return fmt.Errorf("`period` do not support change now.")
	}
	if d.HasChange("shard_memory") {
		return fmt.Errorf("`shard_memory` do not support change now.")
	}
	if d.HasChange("shard_storage") {
		return fmt.Errorf("`shard_storage` do not support change now.")
	}
	if d.HasChange("shard_node_count") {
		return fmt.Errorf("`shard_node_count` do not support change now.")
	}
	if d.HasChange("shard_count") {
		return fmt.Errorf("`shard_count` do not support change now.")
	}

	if d.HasChange("project_id") {
		return fmt.Errorf("`project_id` do not support change now.")
	}
	if d.HasChange("vpc_id") {
		return fmt.Errorf("`vpc_id` do not support change now.")
	}
	if d.HasChange("subnet_id") {
		return fmt.Errorf("`subnet_id` do not support change now.")
	}
	if d.HasChange("db_version_id") {
		return fmt.Errorf("`db_version_id` do not support change now.")
	}
	if d.HasChange("auto_voucher") {
		return fmt.Errorf("`auto_voucher` do not support change now.")
	}
	if d.HasChange("voucher_ids") {
		return fmt.Errorf("`voucher_ids` do not support change now.")
	}

	if d.HasChange("ipv6_flag") {
		return fmt.Errorf("`ipv6_flag` do not support change now.")
	}
	if d.HasChange("resource_tags") {
		return fmt.Errorf("`resource_tags` do not support change now.")
	}
	if d.HasChange("init_params") {
		return fmt.Errorf("`init_params` do not support change now.")
	}
	if d.HasChange("dcn_region") {
		return fmt.Errorf("`dcn_region` do not support change now.")
	}
	if d.HasChange("dcn_instance_id") {
		return fmt.Errorf("`dcn_instance_id` do not support change now.")
	}
	if d.HasChange("auto_renew_flag") {
		return fmt.Errorf("`auto_renew_flag` do not support change now.")
	}
	if d.HasChange("security_group_ids") {
		return fmt.Errorf("`security_group_ids` do not support change now.")
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
	instanceId := d.Id()

	if err := service.DeleteDcdbDbInstanceById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
