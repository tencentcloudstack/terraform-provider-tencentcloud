/*
Provides a resource to create a cam service_linked_role

Example Usage

```hcl
resource "tencentcloud_cam_service_linked_role" "service_linked_role" {
  qcs_service_name = "postgreskms.postgres.cloud.tencent.com"
  custom_suffix = "x-1"
  description = "desc cam"
  tags = {
    "createdBy" = "terraform"
  }
}

```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCamServiceLinkedRole() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudCamServiceLinkedRoleRead,
		Create: resourceTencentCloudCamServiceLinkedRoleCreate,
		Update: resourceTencentCloudCamServiceLinkedRoleUpdate,
		Delete: resourceTencentCloudCamServiceLinkedRoleDelete,
		Schema: map[string]*schema.Schema{
			"qcs_service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Authorization service, the Tencent Cloud service principal with this role attached.",
			},

			"custom_suffix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The custom suffix, based on the string you provide, is combined with the prefix provided by the service to form the full role name.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "role description.",
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
		request        = cam.NewCreateServiceLinkedRoleRequest()
		response       *cam.CreateServiceLinkedRoleResponse
		roleId         string
		qcsServiceName = ""
		customSuffix   = ""
	)

	if v, ok := d.GetOk("qcs_service_name"); ok {
		qcsServiceName = v.(string)
		request.QCSServiceName = helper.Strings([]string{v.(string)})
	}

	if v, ok := d.GetOk("custom_suffix"); ok {
		customSuffix = v.(string)
		request.CustomSuffix = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		for k, v := range tags {
			key := k
			value := v
			request.Tags = append(request.Tags, &cam.RoleTags{
				Key:   &key,
				Value: &value,
			})
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().CreateServiceLinkedRole(request)
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
		log.Printf("[CRITAL]%s create cam serviceLinkedRole failed, reason:%+v", logId, err)
		return err
	}

	roleId = *response.Response.RoleId

	d.SetId(roleId + FILED_SP + qcsServiceName + FILED_SP + customSuffix)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::cam:%s:uin/:RoleId/%s", region, roleId)
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

	resourceId := d.Id()
	items := strings.Split(resourceId, FILED_SP)
	if len(items) != 3 {
		return fmt.Errorf("invalid ID %s", resourceId)
	}
	roleId := items[0]
	qcsServiceName := items[1]
	customSuffix := items[2]

	serviceLinkedRole, err := service.DescribeCamServiceLinkedRole(ctx, roleId)
	if err != nil {
		return err
	}

	if serviceLinkedRole == nil {
		d.SetId("")
		return fmt.Errorf("resource `serviceLinkedRole` %s does not exist", roleId)
	}

	if qcsServiceName != "" {
		_ = d.Set("qcs_service_name", qcsServiceName)
	}

	if customSuffix != "" {
		_ = d.Set("custom_suffix", customSuffix)
	}

	if serviceLinkedRole.Description != nil {
		_ = d.Set("description", serviceLinkedRole.Description)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "cam", "RoleId", tcClient.Region, roleId)
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
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	request := cam.NewUpdateRoleDescriptionRequest()

	resourceId := d.Id()
	items := strings.Split(resourceId, FILED_SP)
	if len(items) != 3 {
		return fmt.Errorf("invalid ID %s", resourceId)
	}
	roleId := items[0]
	// qcsServiceName := items[1]
	// customSuffix := items[2]

	request.RoleId = &roleId

	if d.HasChange("qcs_service_name") {
		return fmt.Errorf("`qcs_service_name` do not support change now.")
	}

	if d.HasChange("custom_suffix") {
		return fmt.Errorf("`custom_suffix` do not support change now.")
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().UpdateRoleDescription(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cam serviceLinkedRole failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
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

	resourceId := d.Id()
	items := strings.Split(resourceId, FILED_SP)
	if len(items) != 3 {
		return fmt.Errorf("invalid ID %s", resourceId)
	}
	roleId := items[0]
	// qcsServiceName := items[1]
	// customSuffix := items[2]

	serviceLinkedRole, err := service.DescribeCamServiceLinkedRole(ctx, roleId)
	if err != nil {
		return err
	}
	if serviceLinkedRole == nil || serviceLinkedRole.RoleName == nil {
		return fmt.Errorf("When querying serviceLinkedRole, an error occurs")
	}

	deletionTaskId, err := service.DeleteCamServiceLinkedRoleById(ctx, *serviceLinkedRole.RoleName)
	if err != nil {
		return err
	}

	err = resource.Retry(3*readRetryTimeout, func() *resource.RetryError {
		response, _ := service.DescribeCamServiceLinkedRoleDeleteStatus(ctx, deletionTaskId)
		// if errRet != nil {
		// 	return retryError(errRet, InternalError)
		// }
		if response == nil || response.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("When querying the deletion status, an error occurred"))
		}

		instance := response.Response
		if *instance.Status == "SUCCEEDED" {
			return nil
		}
		if *instance.Status == "FAILED" {
			return resource.NonRetryableError(fmt.Errorf("serviceLinkedRole status is %v, operate failed.", *instance.Status))
		}
		return resource.RetryableError(fmt.Errorf("serviceLinkedRole status is %v, retry...", *instance.Status))
	})
	if err != nil {
		return err
	}
	return nil
}
