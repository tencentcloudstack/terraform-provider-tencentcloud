/*
Provides a resource to create a dcdb hourdb_instance

Example Usage

```hcl
resource "tencentcloud_dcdb_hourdb_instance" "hourdb_instance" {
  instance_name = "test_dcdc_dc_instance"
  zones = ["ap-guangzhou-5", "ap-guangzhou-6"]
  shard_memory = "2"
  shard_storage = "10"
  shard_node_count = "2"
  shard_count = "2"
  vpc_id = local.vpc_id
  subnet_id = local.subnet_id
  security_group_id = local.sg_id
  db_version_id = "8.0"
  resource_tags {
	tag_key = "aaa"
	tag_value = "bbb"
  }
}

```
Import

dcdb hourdb_instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_dcdb_hourdb_instance.hourdb_instance hourdbInstance_id
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

func resourceTencentCloudDcdbHourdbInstance() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudDcdbHourdbInstanceRead,
		Create: resourceTencentCloudDcdbHourdbInstanceCreate,
		Update: resourceTencentCloudDcdbHourdbInstanceUpdate,
		Delete: resourceTencentCloudDcdbHourdbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zones": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "available zone.",
			},

			"shard_memory": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "memory(GB) for each shard. It can be obtained by querying api DescribeShardSpec.",
			},

			"shard_storage": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "storage(GB) for each shard. It can be obtained by querying api DescribeShardSpec.",
			},

			"shard_node_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "node count for each shard. It can be obtained by querying api DescribeShardSpec.",
			},

			"shard_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "instance shard count.",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "vpc id.",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "subnet id, it&amp;#39;s required when vpcId is set.",
			},

			"db_version_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "db engine version, default to Percona 5.7.17.",
			},

			"security_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "security group id.",
			},

			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "project id.",
			},

			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "name of this instance.",
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

			"ipv6_flag": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to support IPv6.",
			},

			"extranet_access": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to open the extranet access.",
			},

			"rs_access_strategy": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "RS nearest access mode, 0-no policy, 1-nearest access.",
			},

			"resource_tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "resource tags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "tag value.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudDcdbHourdbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_hourdb_instance.create")()
	defer inconsistentCheck(d, meta)()

	var (
		request       = dcdb.NewCreateHourDCDBInstanceRequest()
		response      *dcdb.CreateHourDCDBInstanceResponse
		instanceId    string
		dcnInstanceId string
		ipv6Flag      int
		logId         = getLogId(contextNil)
		ctx           = context.WithValue(context.TODO(), logIdKey, logId)
		service       = DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	if v, ok := d.GetOk("zones"); ok {
		zonesSet := v.(*schema.Set).List()
		for i := range zonesSet {
			zones := zonesSet[i].(string)
			request.Zones = append(request.Zones, &zones)
		}
	}

	if v, ok := d.GetOk("shard_memory"); ok {
		request.ShardMemory = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("shard_storage"); ok {
		request.ShardStorage = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("shard_node_count"); ok {
		request.ShardNodeCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("shard_count"); ok {
		request.ShardCount = helper.IntInt64(v.(int))
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

	if v, ok := d.GetOk("security_group_id"); ok {
		request.SecurityGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, _ := d.GetOk("ipv6_flag"); v != nil {
		request.Ipv6Flag = helper.IntInt64(v.(int))
		ipv6Flag = v.(int)
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

	if v, ok := d.GetOk("dcn_region"); ok {
		request.DcnRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dcn_instance_id"); ok {
		request.DcnInstanceId = helper.String(v.(string))
		dcnInstanceId = v.(string)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().CreateHourDCDBInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dcdb hourdbInstance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceIds[0]

	d.SetId(instanceId)

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

	initRet, flowId, err := service.InitDcdbDbInstance(ctx, instanceId, defaultInitParams)
	if err != nil {
		return err
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

	if dcnInstanceId != "" {
		// need to wait dcn init processing complete
		// 0:none; 1:creating, 2:running
		conf := BuildStateChangeConf([]string{}, []string{"2"}, 3*readRetryTimeout, time.Second, service.DcdbDcnStateRefreshFunc(instanceId, []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	if v, ok := d.GetOkExists("extranet_access"); ok && v != nil {
		flag := v.(bool)
		err := service.SetDcdbExtranetAccess(ctx, instanceId, ipv6Flag, flag)
		if err != nil {
			return err
		}
	}

	if v, ok := d.GetOkExists("rs_access_strategy"); ok && v != nil {
		rsStrategy := v.(int)
		err := service.SetRealServerAccessStrategy(ctx, instanceId, rsStrategy)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudDcdbHourdbInstanceRead(d, meta)
}

func resourceTencentCloudDcdbHourdbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_hourdb_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	hourdbInstances, err := service.DescribeDcdbHourdbInstance(ctx, instanceId)

	if err != nil {
		return err
	}

	if hourdbInstances == nil {
		d.SetId("")
		return fmt.Errorf("resource `hourdbInstance` %s does not exist", instanceId)
	}

	if *hourdbInstances.TotalCount > 1 || len(hourdbInstances.Instances) > 1 {
		d.SetId("")
		return fmt.Errorf("the count of resource `hourdbInstance` shoud not beyond one. count: %v\n", hourdbInstances.TotalCount)
	}

	hourdbInstance := hourdbInstances.Instances[0]
	if hourdbInstance.Zone != nil {
		_ = d.Set("zones", []*string{hourdbInstance.Zone})
	}

	if hourdbInstance.ShardDetail[0] != nil { // Memory and Storage is params for one shard
		shard := hourdbInstance.ShardDetail[0]

		if shard.Memory != nil {
			_ = d.Set("shard_memory", shard.Memory)
		}

		if shard.Storage != nil {
			_ = d.Set("shard_storage", shard.Storage)
		}
	}

	if hourdbInstance.NodeCount != nil {
		_ = d.Set("shard_node_count", hourdbInstance.NodeCount)
	}

	if hourdbInstance.ShardCount != nil {
		_ = d.Set("shard_count", hourdbInstance.ShardCount)
	}

	if hourdbInstance.UniqueVpcId != nil {
		_ = d.Set("vpc_id", hourdbInstance.UniqueVpcId)
	}

	if hourdbInstance.UniqueSubnetId != nil {
		_ = d.Set("subnet_id", hourdbInstance.UniqueSubnetId)
	}

	if hourdbInstance.DbVersionId != nil {
		_ = d.Set("db_version_id", hourdbInstance.DbVersionId)
	}

	if hourdbInstance.ProjectId != nil {
		_ = d.Set("project_id", hourdbInstance.ProjectId)
	}

	if hourdbInstance.InstanceName != nil {
		_ = d.Set("instance_name", hourdbInstance.InstanceName)
	}

	if hourdbInstance.Ipv6Flag != nil {
		_ = d.Set("ipv6_flag", hourdbInstance.Ipv6Flag)
	}

	if hourdbInstance.WanStatus != nil {
		//0-未开通；1-已开通；2-关闭；3-开通中
		if *hourdbInstance.WanStatus == DCDB_WAN_STATUS_UNOPEN || *hourdbInstance.WanStatus == DCDB_WAN_STATUS_CLOSED {
			_ = d.Set("extranet_access", false)
		}

		if *hourdbInstance.WanStatus == DCDB_WAN_STATUS_OPENED {
			_ = d.Set("extranet_access", true)
		}
	}

	if hourdbInstance.ResourceTags != nil {
		resourceTagsList := []interface{}{}
		for _, resourceTags := range hourdbInstance.ResourceTags {
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

	if sg, err := service.DescribeDcdbSecurityGroup(ctx, instanceId); sg != nil {
		sgId := ""
		if len(sg.Groups) > 0 {
			sgId = *sg.Groups[0].SecurityGroupId
		}
		_ = d.Set("security_group_id", sgId)
	} else {
		return err
	}

	// set dcn id and region
	if dcns, err := service.DescribeDcnDetailById(ctx, instanceId); dcns != nil {
		for _, dcn := range dcns {
			var master *dcdb.DcnDetailItem
			if *dcn.DcnFlag == DCDB_DCN_FLAG_MASTER {
				master = dcn
				_ = d.Set("dcn_region", master.Region)
				_ = d.Set("dcn_instance_id", master.InstanceId)
			}
		}
	} else {
		return err
	}

	return nil
}

func resourceTencentCloudDcdbHourdbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_hourdb_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request = dcdb.NewModifyDBInstanceNameRequest()
		service = DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if d.HasChange("zones") {
		return fmt.Errorf("`zones` do not support change now.")
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

	if d.HasChange("vpc_id") {
		return fmt.Errorf("`vpc_id` do not support change now.")
	}

	if d.HasChange("subnet_id") {
		return fmt.Errorf("`subnet_id` do not support change now.")
	}

	if d.HasChange("db_version_id") {
		return fmt.Errorf("`db_version_id` do not support change now.")
	}

	// wait for ModifyDBInstanceSecurityGroups
	// if d.HasChange("security_group_id") {
	// 	return fmt.Errorf("`security_group_id` do not support change now.")
	// }

	if v, ok := d.GetOkExists("extranet_access"); ok && v != nil {
		flag := v.(bool)
		var ipv6Flag int
		if v, _ := d.GetOk("ipv6_flag"); v != nil {
			ipv6Flag = v.(int)
		}
		err := service.SetDcdbExtranetAccess(ctx, instanceId, ipv6Flag, flag)
		if err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
	}

	if v, ok := d.GetOkExists("rs_access_strategy"); ok && v != nil {
		rsStrategy := v.(int)
		err := service.SetRealServerAccessStrategy(ctx, instanceId, rsStrategy)
		if err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
	}

	if d.HasChange("project_id") {
		if projectId, ok := d.GetOk("project_id"); ok {
			request := dcdb.NewModifyDBInstancesProjectRequest()

			request.InstanceIds = []*string{&instanceId}
			request.ProjectId = helper.IntInt64(projectId.(int))

			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().ModifyDBInstancesProject(request)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s operate dcdb modifyInstanceProjectOperation failed, reason:%+v", logId, err)
				return err
			}
		}
		time.Sleep(2 * time.Second)
	}

	if d.HasChange("dcn_region") {
		return fmt.Errorf("`dcn_region` do not support change now.")
	}
	if d.HasChange("dcn_instance_id") {
		return fmt.Errorf("`dcn_instance_id` do not support change now.")
	}

	if d.HasChange("dcn_region") {
		return fmt.Errorf("`dcn_region` do not support change now.")
	}
	if d.HasChange("dcn_instance_id") {
		return fmt.Errorf("`dcn_instance_id` do not support change now.")
	}

	if d.HasChange("resource_tags") {
		return fmt.Errorf("`resource_tags` do not support change now.")
	}

	if d.HasChange("instance_name") {
		if v, ok := d.GetOk("instance_name"); ok {
			request.InstanceName = helper.String(v.(string))
		}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().ModifyDBInstanceName(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s create dcdb hourdbInstance failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudDcdbHourdbInstanceRead(d, meta)
}

func resourceTencentCloudDcdbHourdbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_hourdb_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	if err := service.DeleteDcdbHourdbInstanceById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
