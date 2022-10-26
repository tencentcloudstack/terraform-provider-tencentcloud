/*
Provides a resource to create a mariadb dedicatedcluster_db_instance

Example Usage

```hcl
resource "tencentcloud_mariadb_dedicatedcluster_db_instance" "dedicatedcluster_db_instance" {
  goods_num = ""
  memory = ""
  storage = ""
  cluster_id = ""
  vpc_id = ""
  subnet_id = ""
  db_version_id = ""
  instance_name = ""
  tags = {
    "createdBy" = "terraform"
  }
}

```
Import

mariadb dedicatedcluster_db_instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_mariadb_dedicatedcluster_db_instance.dedicatedcluster_db_instance dedicatedClusterDBInstance_id
```
*/

package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMariadbDedicatedClusterDBInstance() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMariadbDedicatedClusterDBInstanceRead,
		Create: resourceTencentCloudMariadbDedicatedClusterDBInstanceCreate,

		Update: resourceTencentCloudMariadbDedicatedClusterDBInstanceUpdate,
		Delete: resourceTencentCloudMariadbDedicatedClusterDBInstanceDelete,
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
				Description: "vpc id for instance.",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "subnet id for instance, it&amp;#39;s required when vpcId is set.",
			},

			"db_version_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "db engine version for instance, default to 0.",
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

func resourceTencentCloudMariadbDedicatedClusterDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_dedicatedcluster_db_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

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
			if err, ok := e.(*errors.TencentCloudSDKError); ok {
				if err.Code == "FailedOperation.PayFailed" {
					return &resource.RetryError{e, false}
				}
			}
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mariadb dedicatedClusterDBInstance failed, reason:%+v", logId, err)
		return err
	}

	dedicatedClusterDBInstanceId := *response.Response.InstanceIds[0]

	d.SetId(instanceId)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::mariadb:%s:uin/:mariadb-dedicatedcluster-instance/%s", region, dedicatedClusterDBInstanceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}
	return resourceTencentCloudMariadbDedicatedClusterDBInstanceRead(d, meta)
}

func resourceTencentCloudMariadbDedicatedClusterDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_dedicatedcluster_db_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	dedicatedClusterDBInstances, err := service.DescribeMariadbDedicatedClusterDBInstance(ctx, instanceId)

	if err != nil {
		return err
	}

	if dedicatedClusterDBInstances == nil {
		d.SetId("")
		return fmt.Errorf("resource `dedicatedClusterDBInstance` %s does not exist", instanceId)
	}

	if len(dedicatedClusterDBInstances.Instances) == 0 {
		d.SetId("")
		return fmt.Errorf("resource `dedicatedClusterDBInstance` %s does not exist", instanceId)
	}
	dedicatedClusterDBInstance := dedicatedClusterDBInstances.Instances[0]

	if dedicatedClusterDBInstance.Memory != nil {
		_ = d.Set("memory", dedicatedClusterDBInstance.Memory)
	}

	if dedicatedClusterDBInstance.Storage != nil {
		_ = d.Set("storage", dedicatedClusterDBInstance.Storage)
	}

	if dedicatedClusterDBInstance.VpcId != nil {
		_ = d.Set("vpc_id", dedicatedClusterDBInstance.VpcId)
	}

	if dedicatedClusterDBInstance.SubnetId != nil {
		_ = d.Set("subnet_id", dedicatedClusterDBInstance.SubnetId)
	}

	if dedicatedClusterDBInstance.DbVersionId != nil {
		_ = d.Set("db_version_id", dedicatedClusterDBInstance.DbVersionId)
	}

	if dedicatedClusterDBInstance.InstanceName != nil {
		_ = d.Set("instance_name", dedicatedClusterDBInstance.InstanceName)
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

func resourceTencentCloudMariadbDedicatedClusterDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
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
		log.Printf("[CRITAL]%s create mariadb dedicatedClusterDBInstance failed, reason:%+v", logId, err)
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

	return resourceTencentCloudMariadbDedicatedClusterDBInstanceRead(d, meta)
}

func resourceTencentCloudMariadbDedicatedClusterDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_dedicatedcluster_db_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	if err := service.DeleteMariadbDedicatedClusterDBInstanceById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
