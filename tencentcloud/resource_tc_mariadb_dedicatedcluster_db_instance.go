/*
Provides a resource to create a mariadb dedicatedcluster_db_instance

Example Usage

```hcl
resource "tencentcloud_mariadb_dedicatedcluster_db_instance" "dedicatedcluster_db_instance" {
  goods_num 	= 1
  memory 		= 2
  storage 		= 10
  cluster_id 	= "dbdc-24odnuhr"
  vpc_id 		= "vpc-ii1jfbhl"
  subnet_id 	= "subnet-3ku415by"
  db_version_id = "8.0"
  instance_name = "cluster-mariadb-test-1"
}

```
Import

mariadb dedicatedcluster_db_instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_mariadb_dedicatedcluster_db_instance.dedicatedcluster_db_instance tdsql-050g3fmv
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMariadbDedicatedclusterDbInstance() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMariadbDedicatedclusterDbInstanceRead,
		Create: resourceTencentCloudMariadbDedicatedclusterDbInstanceCreate,
		Update: resourceTencentCloudMariadbDedicatedclusterDbInstanceUpdate,
		Delete: resourceTencentCloudMariadbDedicatedclusterDbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"goods_num": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "number of instance.",
			},

			"memory": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "instance memory.",
			},

			"storage": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "instance disk storage.",
			},

			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "dedicated cluster id.",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
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
				Computed:    true,
				Description: "db engine version, default to 0.",
			},

			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "name of this instance.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudMariadbDedicatedclusterDbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_dedicatedcluster_db_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request    = mariadb.NewCreateDedicatedClusterDBInstanceRequest()
		response   *mariadb.CreateDedicatedClusterDBInstanceResponse
		instanceId string
	)

	if v, ok := d.GetOk("goods_num"); ok {
		request.GoodsNum = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("memory"); ok {
		request.Memory = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("storage"); ok {
		request.Storage = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
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

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().CreateDedicatedClusterDBInstance(request)
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
		log.Printf("[CRITAL]%s create mariadb dedicatedclusterDbInstance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceIds[0]
	d.SetId(instanceId)

	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}
	initParams := []*mariadb.DBParamValue{
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

	initRet, err := service.InitDbInstance(ctx, instanceId, initParams)
	if err != nil {
		return err
	}
	if !initRet {
		return fmt.Errorf("db instance init failed")
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::mariadb:%s:uin/:mariadb-dedicatedcluster-instance/%s", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudMariadbDedicatedclusterDbInstanceRead(d, meta)
}

func resourceTencentCloudMariadbDedicatedclusterDbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_dedicatedcluster_db_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	paramMap := make(map[string]interface{})
	paramMap["instance_ids"] = []*string{&instanceId}

	dbInstances, err := service.DescribeMariadbDbInstancesByFilter(ctx, paramMap)

	if err != nil {
		return err
	}

	if len(dbInstances) < 1 {
		d.SetId("")
		return fmt.Errorf("resource `dedicatedclusterDbInstance` %s does not exist", instanceId)
	}

	_ = d.Set("goods_num", len(dbInstances))

	dbInstance := dbInstances[0]
	if dbInstance.Memory != nil {
		_ = d.Set("memory", dbInstance.Memory)
	}

	if dbInstance.Storage != nil {
		_ = d.Set("storage", dbInstance.Storage)
	}

	if dbInstance.ExclusterId != nil {
		_ = d.Set("cluster_id", dbInstance.ExclusterId)
	}

	if dbInstance.UniqueVpcId != nil {
		_ = d.Set("vpc_id", dbInstance.UniqueVpcId)
	}

	if dbInstance.UniqueSubnetId != nil {
		_ = d.Set("subnet_id", dbInstance.UniqueSubnetId)
	}

	if dbInstance.DbVersionId != nil {
		_ = d.Set("db_version_id", dbInstance.DbVersionId)
	}

	if dbInstance.InstanceName != nil {
		_ = d.Set("instance_name", dbInstance.InstanceName)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "mariadb", "mariadb-dedicatedcluster-instance", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudMariadbDedicatedclusterDbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_dedicatedcluster_db_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := mariadb.NewModifyDBInstanceNameRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if d.HasChange("goods_num") {

		return fmt.Errorf("`goods_num` do not support change now.")

	}

	if d.HasChange("memory") {

		return fmt.Errorf("`memory` do not support change now.")

	}

	if d.HasChange("storage") {

		return fmt.Errorf("`storage` do not support change now.")

	}

	if d.HasChange("cluster_id") {

		return fmt.Errorf("`cluster_id` do not support change now.")

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

	if d.HasChange("instance_name") {
		if v, ok := d.GetOk("instance_name"); ok {
			request.InstanceName = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().ModifyDBInstanceName(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mariadb dedicatedclusterDbInstance failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("mariadb", "mariadb-dedicatedcluster-instance", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudMariadbDedicatedclusterDbInstanceRead(d, meta)
}

func resourceTencentCloudMariadbDedicatedclusterDbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_dedicatedcluster_db_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	paramMap := make(map[string]interface{})
	paramMap["instance_ids"] = []*string{&instanceId}
	err := resource.Retry(3*readRetryTimeout, func() *resource.RetryError {
		dbInstances, errRet := service.DescribeMariadbDbInstancesByFilter(ctx, paramMap)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if dbInstances == nil || len(dbInstances) < 1 {
			return nil
		}

		if *dbInstances[0].Status == 2 {
			isolateRequest := mariadb.NewIsolateDedicatedDBInstanceRequest()
			isolateRequest.InstanceId = &instanceId
			errIsolate := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				if _, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().IsolateDedicatedDBInstance(isolateRequest); e != nil {
					return resource.NonRetryableError(fmt.Errorf("delete db instance failed, err: %v", e))
				}
				return nil
			})
			if errIsolate != nil {
				return resource.NonRetryableError(fmt.Errorf("db instance error %v, operate failed", errIsolate))
			}
			return resource.RetryableError(fmt.Errorf("db instance initializing, retry..."))
		}

		if *dbInstances[0].Status < 0 {
			ee := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				if e := service.DeleteMariadbDbInstance(ctx, instanceId); e != nil {
					return resource.NonRetryableError(fmt.Errorf("delete db instance failed, err: %v", e))
				}
				return nil
			})
			if ee != nil {
				return resource.NonRetryableError(fmt.Errorf("db instance error %v, operate failed", ee))
			}
			return resource.RetryableError(fmt.Errorf("db instance initializing, retry..."))
		}
		return resource.RetryableError(fmt.Errorf("db instance status is %v, retry...", *dbInstances[0].Status))
	})
	if err != nil {
		return err
	}

	return nil
}
