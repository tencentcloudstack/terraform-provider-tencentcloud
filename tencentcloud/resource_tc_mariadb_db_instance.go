/*
Provides a resource to create a mariadb db_instance

Example Usage

```hcl
resource "tencentcloud_mariadb_db_instance" "db_instance" {
  zones = ""
  node_count = ""
  memory = ""
  storage = ""
  vpc_id = ""
  subnet_id = ""
  db_version_id = ""
  period = ""
  instance_name = ""
  tags = {
    "createdBy" = "terraform"
  }
}

```
Import

mariadb db_instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_mariadb_db_instance.db_instance dbInstance_id
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

func resourceTencentCloudMariadbDbInstance() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMariadbDbInstanceRead,
		Create: resourceTencentCloudMariadbDbInstanceCreate,

		Update: resourceTencentCloudMariadbDbInstanceUpdate,
		Delete: resourceTencentCloudMariadbDbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zones": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "available of instance.",
			},

			"node_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "number of node for instance.",
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

			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "vpc id for this instance.",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "subnet id for this instance, it&amp;#39;s required when vpcId is set.",
			},

			"db_version_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "db engine version for this instance, default to Percona 5.7.17.",
			},

			"period": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "subscribes month of instance",
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

func resourceTencentCloudMariadbDbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_db_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request     = mariadb.NewCreateDBInstanceRequest()
		response    *mariadb.CreateDBInstanceResponse
		instanceIds string
	)

	if v, ok := d.GetOk("zones"); ok {
		zonesSet := v.(*schema.Set).List()
		for i := range zonesSet {
			zones := zonesSet[i].(string)
			request.Zones = append(request.Zones, &zones)
		}
	}

	if v, ok := d.GetOk("node_count"); ok {
		request.NodeCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("memory"); ok {
		request.Memory = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("storage"); ok {
		request.Storage = helper.IntInt64(v.(int))
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

	if v, ok := d.GetOk("period"); ok {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().CreateDBInstance(request)
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
		log.Printf("[CRITAL]%s create mariadb dbInstance failed, reason:%+v", logId, err)
		return err
	}

	instanceIds = *response.Response.InstanceIds[0]

	d.SetId(instanceIds)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::mariadb:%s:uin/:mariadb-instance/%s", region, instanceIds)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}
	return resourceTencentCloudMariadbDbInstanceRead(d, meta)
}

func resourceTencentCloudMariadbDbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_db_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

	dbInstanceId := d.Id()

	dbInstances, err := service.DescribeMariadbDbInstance(ctx, dbInstanceId)

	if err != nil {
		return err
	}

	if dbInstances == nil {
		d.SetId("")
		return fmt.Errorf("resource `dbInstance` %s does not exist", dbInstanceId)
	}

	if len(dbInstances.Instances) == 0 {
		d.SetId("")
		return fmt.Errorf("resource `dbInstance` %s does not exist", dbInstanceId)
	}
	dbInstance := dbInstances.Instances[0]

	if dbInstance.Zone != nil {
		_ = d.Set("zones", dbInstance.Zone)
	}

	if dbInstance.NodeCount != nil {
		_ = d.Set("node_count", dbInstance.NodeCount)
	}

	if dbInstance.Memory != nil {
		_ = d.Set("memory", dbInstance.Memory)
	}

	if dbInstance.Storage != nil {
		_ = d.Set("storage", dbInstance.Storage)
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

	if dbInstance.InstanceName != nil {
		_ = d.Set("instance_name", dbInstance.InstanceName)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "mariadb", "mariadb-instance", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudMariadbDbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_db_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := mariadb.NewModifyDBInstanceNameRequest()

	instanceIds := d.Id()

	request.InstanceId = &instanceIds

	if d.HasChange("zones") {

		return fmt.Errorf("`zones` do not support change now.")

	}

	if d.HasChange("node_count") {

		return fmt.Errorf("`node_count` do not support change now.")

	}

	if d.HasChange("memory") {

		return fmt.Errorf("`memory` do not support change now.")

	}

	if d.HasChange("storage") {

		return fmt.Errorf("`storage` do not support change now.")

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

	if d.HasChange("period") {

		return fmt.Errorf("`period` do not support change now.")

	}

	if d.HasChange("instance_name") {

		return fmt.Errorf("`instance_name` do not support change now.")

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
		log.Printf("[CRITAL]%s create mariadb dbInstance failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("mariadb", "mariadb-instance", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudMariadbDbInstanceRead(d, meta)
}

func resourceTencentCloudMariadbDbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_db_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceIds := d.Id()

	if err := service.DeleteMariadbDbInstanceById(ctx, instanceIds); err != nil {
		return err
	}

	return nil
}
