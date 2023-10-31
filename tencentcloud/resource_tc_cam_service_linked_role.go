/*
Provides a resource to create a cam service_linked_role

Example Usage

```hcl
resource "tencentcloud_cam_service_linked_role" "service_linked_role" {
  qcs_service_name = ["cvm.qcloud.com","ekslog.tke.cloud.tencent.com"]
  custom_suffix = "tf"
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
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				ForceNew:    true,
				Description: "Authorization service, the Tencent Cloud service principal with this role attached.",
			},

			"custom_suffix": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The custom suffix, based on the string you provide, is combined with the prefix provided by the service to form the full role name. This field is not allowed to contain the character `_`.",
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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceTencentCloudCamServiceLinkedRoleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_service_linked_role.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cam.NewCreateServiceLinkedRoleRequest()
		response *cam.CreateServiceLinkedRoleResponse
		roleId   string
	)

	if v, ok := d.GetOk("qcs_service_name"); ok {
		serviceName := v.(*schema.Set).List()
		serviceNameArr := make([]*string, 0, len(serviceName))
		for _, name := range serviceName {
			serviceNameArr = append(serviceNameArr, helper.String(name.(string)))
		}
		request.QCSServiceName = serviceNameArr
	}

	if v, ok := d.GetOk("custom_suffix"); ok {
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

	d.SetId(roleId)
	//ctx := context.WithValue(context.TODO(), logIdKey, logId)
	//if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
	//	tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
	//	resourceName := fmt.Sprintf("qcs::cam:%s:uin/:role/tencentcloudServiceRole/%s", "", roleId)
	//	if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
	//		return err
	//	}
	//}
	return resourceTencentCloudCamServiceLinkedRoleRead(d, meta)
}

func resourceTencentCloudCamServiceLinkedRoleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_service_linked_role.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}

	roleId := d.Id()

	serviceLinkedRole, err := service.DescribeCamServiceLinkedRole(ctx, roleId)
	if err != nil {
		return err
	}

	if serviceLinkedRole == nil {
		d.SetId("")
		return fmt.Errorf("resource `serviceLinkedRole` %s does not exist", roleId)
	}

	if serviceLinkedRole.PolicyDocument != nil {
		var documentJson Document
		err = json.Unmarshal([]byte(*serviceLinkedRole.PolicyDocument), &documentJson)
		if err != nil {
			return err
		}
		if documentJson.Statement != nil && len(documentJson.Statement) > 0 {
			principal := documentJson.Statement[0].Principal
			if principal.Service != nil && len(principal.Service) > 0 {
				_ = d.Set("qcs_service_name", principal.Service)
			}
		}
	}

	if serviceLinkedRole.RoleName != nil {
		roleName := strings.Split(*serviceLinkedRole.RoleName, "_")
		if len(roleName) > 0 {
			_ = d.Set("custom_suffix", roleName[len(roleName)-1])
		}
	}

	if serviceLinkedRole.Description != nil {
		_ = d.Set("description", serviceLinkedRole.Description)
	}

	if serviceLinkedRole.Tags != nil {
		tagsMap := map[string]interface{}{}
		for _, tag := range serviceLinkedRole.Tags {
			if tag.Key != nil && tag.Value != nil {
				tagsMap[*tag.Key] = tag.Value
			}
		}
		_ = d.Set("tags", tagsMap)
	}

	return nil
}

func resourceTencentCloudCamServiceLinkedRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_service_linked_role.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	roleId := d.Id()

	if d.HasChange("description") {
		request := cam.NewUpdateRoleDescriptionRequest()
		request.RoleId = &roleId

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
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
	}

	if d.HasChange("tags") {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("cam", "role/tencentcloudServiceRole", "", d.Id())
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

	roleId := d.Id()

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
