/*
Provides a resource to create a pts project

Example Usage

```hcl
resource "tencentcloud_pts_project" "project" {
  name = "ptsObjectName"
  description = &lt;nil&gt;
  tags {
		tag_key = &lt;nil&gt;
		tag_value = &lt;nil&gt;

  }
            }
```

Import

pts project can be imported using the id, e.g.

```
terraform import tencentcloud_pts_project.project project_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	pts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/pts/v20210728"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudPtsProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPtsProjectCreate,
		Read:   resourceTencentCloudPtsProjectRead,
		Update: resourceTencentCloudPtsProjectUpdate,
		Delete: resourceTencentCloudPtsProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "PtsObjectName, which must be required.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Pts object description.",
			},

			"tags": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Tags List.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag value.",
						},
					},
				},
			},

			"app_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "App ID.",
			},

			"uin": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "User ID.",
			},

			"sub_account_uin": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Sub-user ID.",
			},

			"created_at": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Creation time.",
			},

			"updated_at": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Update time.",
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Project status.",
			},
		},
	}
}

func resourceTencentCloudPtsProjectCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_project.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = pts.NewCreateProjectRequest()
		response  = pts.NewCreateProjectResponse()
		projectId string
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			tagSpec := pts.TagSpec{}
			if v, ok := dMap["tag_key"]; ok {
				tagSpec.TagKey = helper.String(v.(string))
			}
			if v, ok := dMap["tag_value"]; ok {
				tagSpec.TagValue = helper.String(v.(string))
			}
			request.Tags = append(request.Tags, &tagSpec)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePtsClient().CreateProject(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create pts project failed, reason:%+v", logId, err)
		return err
	}

	projectId = *response.Response.ProjectId
	d.SetId(projectId)

	return resourceTencentCloudPtsProjectRead(d, meta)
}

func resourceTencentCloudPtsProjectRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_project.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	projectId := d.Id()

	project, err := service.DescribePtsProjectById(ctx, projectId)
	if err != nil {
		return err
	}

	if project == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PtsProject` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if project.Name != nil {
		_ = d.Set("name", project.Name)
	}

	if project.Description != nil {
		_ = d.Set("description", project.Description)
	}

	if project.Tags != nil {
		tagsList := []interface{}{}
		for _, tags := range project.Tags {
			tagsMap := map[string]interface{}{}

			if project.Tags.TagKey != nil {
				tagsMap["tag_key"] = project.Tags.TagKey
			}

			if project.Tags.TagValue != nil {
				tagsMap["tag_value"] = project.Tags.TagValue
			}

			tagsList = append(tagsList, tagsMap)
		}

		_ = d.Set("tags", tagsList)

	}

	if project.AppId != nil {
		_ = d.Set("app_id", project.AppId)
	}

	if project.Uin != nil {
		_ = d.Set("uin", project.Uin)
	}

	if project.SubAccountUin != nil {
		_ = d.Set("sub_account_uin", project.SubAccountUin)
	}

	if project.CreatedAt != nil {
		_ = d.Set("created_at", project.CreatedAt)
	}

	if project.UpdatedAt != nil {
		_ = d.Set("updated_at", project.UpdatedAt)
	}

	if project.Status != nil {
		_ = d.Set("status", project.Status)
	}

	return nil
}

func resourceTencentCloudPtsProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_project.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := pts.NewUpdateProjectRequest()

	projectId := d.Id()

	request.ProjectId = &projectId

	immutableArgs := []string{"name", "description", "tags", "app_id", "uin", "sub_account_uin", "created_at", "updated_at", "status"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("tags") {
		if v, ok := d.GetOk("tags"); ok {
			for _, item := range v.([]interface{}) {
				tagSpec := pts.TagSpec{}
				if v, ok := dMap["tag_key"]; ok {
					tagSpec.TagKey = helper.String(v.(string))
				}
				if v, ok := dMap["tag_value"]; ok {
					tagSpec.TagValue = helper.String(v.(string))
				}
				request.Tags = append(request.Tags, &tagSpec)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePtsClient().UpdateProject(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update pts project failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudPtsProjectRead(d, meta)
}

func resourceTencentCloudPtsProjectDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_project.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PtsService{client: meta.(*TencentCloudClient).apiV3Conn}
	projectId := d.Id()

	if err := service.DeletePtsProjectById(ctx, projectId); err != nil {
		return err
	}

	return nil
}
