/*
Provides a resource to create a cdb deploy_group

Example Usage

```hcl
resource "tencentcloud_cdb_deploy_group" "deploy_group" {
  deploy_group_name = &lt;nil&gt;
  description = &lt;nil&gt;
  limit_num = &lt;nil&gt;
  dev_class =
}
```

Import

cdb deploy_group can be imported using the id, e.g.

```
terraform import tencentcloud_cdb_deploy_group.deploy_group deploy_group_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCdbDeployGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdbDeployGroupCreate,
		Read:   resourceTencentCloudCdbDeployGroupRead,
		Update: resourceTencentCloudCdbDeployGroupUpdate,
		Delete: resourceTencentCloudCdbDeployGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"deploy_group_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The name of deploy group. the maximum length cannot exceed 60 characters.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The description of deploy group. the maximum length cannot exceed 200 characters.",
			},

			"limit_num": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The limit on the number of instances on the same physical machine in deploy group affinity policy 1.",
			},

			"dev_class": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The device class of deploy group. optional value is SH12+SH02, TS85.",
			},
		},
	}
}

func resourceTencentCloudCdbDeployGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_deploy_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = cdb.NewCreateDeployGroupRequest()
		response      = cdb.NewCreateDeployGroupResponse()
		deployGroupId string
	)
	if v, ok := d.GetOk("deploy_group_name"); ok {
		request.DeployGroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("limit_num"); ok {
		request.LimitNum = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("dev_class"); ok {
		devClassSet := v.(*schema.Set).List()
		for i := range devClassSet {
			devClass := devClassSet[i].(string)
			request.DevClass = append(request.DevClass, &devClass)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().CreateDeployGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cdb deployGroup failed, reason:%+v", logId, err)
		return err
	}

	deployGroupId = *response.Response.DeployGroupId
	d.SetId(deployGroupId)

	return resourceTencentCloudCdbDeployGroupRead(d, meta)
}

func resourceTencentCloudCdbDeployGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_deploy_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	deployGroupId := d.Id()

	deployGroup, err := service.DescribeCdbDeployGroupById(ctx, deployGroupId)
	if err != nil {
		return err
	}

	if deployGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CdbDeployGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if deployGroup.DeployGroupName != nil {
		_ = d.Set("deploy_group_name", deployGroup.DeployGroupName)
	}

	if deployGroup.Description != nil {
		_ = d.Set("description", deployGroup.Description)
	}

	if deployGroup.LimitNum != nil {
		_ = d.Set("limit_num", deployGroup.LimitNum)
	}

	if deployGroup.DevClass != nil {
		_ = d.Set("dev_class", deployGroup.DevClass)
	}

	return nil
}

func resourceTencentCloudCdbDeployGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_deploy_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cdb.NewModifyNameOrDescByDpIdRequest()

	deployGroupId := d.Id()

	request.DeployGroupId = &deployGroupId

	immutableArgs := []string{"deploy_group_name", "description", "limit_num", "dev_class"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().ModifyNameOrDescByDpId(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cdb deployGroup failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCdbDeployGroupRead(d, meta)
}

func resourceTencentCloudCdbDeployGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_deploy_group.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	deployGroupId := d.Id()

	if err := service.DeleteCdbDeployGroupById(ctx, deployGroupId); err != nil {
		return err
	}

	return nil
}
