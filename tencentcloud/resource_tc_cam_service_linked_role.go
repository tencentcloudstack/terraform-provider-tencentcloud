/*
Provides a resource to create a cam service_linked_role

Example Usage

```hcl
resource "tencentcloud_cam_service_linked_role" "service_linked_role" {
  qcs_service_name = "recordlistdnspod.cdn.cloud.tencent.com"
  custom_suffix = ""
  description = "desc cam"
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

cam service_linked_role can be imported using the id, e.g.

```
terraform import tencentcloud_cam_service_linked_role.service_linked_role service_linked_role_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudCamServiceLinkedRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamServiceLinkedRoleCreate,
		Read:   resourceTencentCloudCamServiceLinkedRoleRead,
		Update: resourceTencentCloudCamServiceLinkedRoleUpdate,
		Delete: resourceTencentCloudCamServiceLinkedRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"qcs_service_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Authorization service, the Tencent Cloud service principal with this role attached.",
			},

			"custom_suffix": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The custom suffix, based on the string you provide, is combined with the prefix provided by the service to form the full role name.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Role description.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudCamServiceLinkedRoleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_service_linked_role.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cam.NewCreateServiceLinkedRoleRequest()
		response = cam.NewCreateServiceLinkedRoleResponse()
		roleId   string
	)
	if v, ok := d.GetOk("qcs_service_name"); ok {
		request.QcsServiceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("custom_suffix"); ok {
		request.CustomSuffix = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().CreateServiceLinkedRole(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cam serviceLinkedRole failed, reason:%+v", logId, err)
		return err
	}

	roleId = *response.Response.RoleId
	d.SetId(roleId)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::cam:%s:uin/:RoleId/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudCamServiceLinkedRoleRead(d, meta)
}

func resourceTencentCloudCamServiceLinkedRoleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_service_linked_role.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}

	serviceLinkedRoleId := d.Id()

	serviceLinkedRole, err := service.DescribeCamServiceLinkedRoleById(ctx, roleId)
	if err != nil {
		return err
	}

	if serviceLinkedRole == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CamServiceLinkedRole` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if serviceLinkedRole.QcsServiceName != nil {
		_ = d.Set("qcs_service_name", serviceLinkedRole.QcsServiceName)
	}

	if serviceLinkedRole.CustomSuffix != nil {
		_ = d.Set("custom_suffix", serviceLinkedRole.CustomSuffix)
	}

	if serviceLinkedRole.Description != nil {
		_ = d.Set("description", serviceLinkedRole.Description)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "cam", "RoleId", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudCamServiceLinkedRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_service_linked_role.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cam.NewUpdateRoleDescriptionRequest()

	serviceLinkedRoleId := d.Id()

	request.RoleId = &roleId

	immutableArgs := []string{"qcs_service_name", "custom_suffix", "description"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().UpdateRoleDescription(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cam serviceLinkedRole failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("cam", "RoleId", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudCamServiceLinkedRoleRead(d, meta)
}

func resourceTencentCloudCamServiceLinkedRoleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_service_linked_role.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}
	serviceLinkedRoleId := d.Id()

	if err := service.DeleteCamServiceLinkedRoleById(ctx, roleId); err != nil {
		return err
	}

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCEEDED"}, 3*readRetryTimeout, time.Second, service.CamServiceLinkedRoleStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
