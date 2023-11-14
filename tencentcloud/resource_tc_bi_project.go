/*
Provides a resource to create a bi project

Example Usage

```hcl
resource "tencentcloud_bi_project" "project" {
  name = "abc"
  color_code = "#066EFF"
  logo = &lt;nil&gt;
  mark = "abc"
  is_apply = true
  default_panel_type = 123
}
```

Import

bi project can be imported using the id, e.g.

```
terraform import tencentcloud_bi_project.project project_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bi "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bi/v20220105"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudBiProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBiProjectCreate,
		Read:   resourceTencentCloudBiProjectRead,
		Update: resourceTencentCloudBiProjectUpdate,
		Delete: resourceTencentCloudBiProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Project name.",
			},

			"color_code": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Logo background color.",
			},

			"logo": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Project logo.",
			},

			"mark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Remark.",
			},

			"is_apply": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to allow users to apply.",
			},

			"default_panel_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Default kanban.",
			},
		},
	}
}

func resourceTencentCloudBiProjectCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_project.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = bi.NewCreateProjectRequest()
		response  = bi.NewCreateProjectResponse()
		projectId string
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("color_code"); ok {
		request.ColorCode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("logo"); ok {
		request.Logo = helper.String(v.(string))
	}

	if v, ok := d.GetOk("mark"); ok {
		request.Mark = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("is_apply"); ok {
		request.IsApply = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("default_panel_type"); ok {
		request.DefaultPanelType = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseBiClient().CreateProject(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create bi project failed, reason:%+v", logId, err)
		return err
	}

	projectId = *response.Response.projectId
	d.SetId(projectId)

	return resourceTencentCloudBiProjectRead(d, meta)
}

func resourceTencentCloudBiProjectRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_project.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := BiService{client: meta.(*TencentCloudClient).apiV3Conn}

	projectId := d.Id()

	project, err := service.DescribeBiProjectById(ctx, projectId)
	if err != nil {
		return err
	}

	if project == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `BiProject` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if project.Name != nil {
		_ = d.Set("name", project.Name)
	}

	if project.ColorCode != nil {
		_ = d.Set("color_code", project.ColorCode)
	}

	if project.Logo != nil {
		_ = d.Set("logo", project.Logo)
	}

	if project.Mark != nil {
		_ = d.Set("mark", project.Mark)
	}

	if project.IsApply != nil {
		_ = d.Set("is_apply", project.IsApply)
	}

	if project.DefaultPanelType != nil {
		_ = d.Set("default_panel_type", project.DefaultPanelType)
	}

	return nil
}

func resourceTencentCloudBiProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_project.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := bi.NewModifyProjectRequest()

	projectId := d.Id()

	request.ProjectId = &projectId

	immutableArgs := []string{"name", "color_code", "logo", "mark", "is_apply", "default_panel_type"}

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

	if d.HasChange("color_code") {
		if v, ok := d.GetOk("color_code"); ok {
			request.ColorCode = helper.String(v.(string))
		}
	}

	if d.HasChange("logo") {
		if v, ok := d.GetOk("logo"); ok {
			request.Logo = helper.String(v.(string))
		}
	}

	if d.HasChange("mark") {
		if v, ok := d.GetOk("mark"); ok {
			request.Mark = helper.String(v.(string))
		}
	}

	if d.HasChange("is_apply") {
		if v, ok := d.GetOkExists("is_apply"); ok {
			request.IsApply = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("default_panel_type") {
		if v, ok := d.GetOkExists("default_panel_type"); ok {
			request.DefaultPanelType = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseBiClient().ModifyProject(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update bi project failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudBiProjectRead(d, meta)
}

func resourceTencentCloudBiProjectDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_project.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := BiService{client: meta.(*TencentCloudClient).apiV3Conn}
	projectId := d.Id()

	if err := service.DeleteBiProjectById(ctx, projectId); err != nil {
		return err
	}

	return nil
}
