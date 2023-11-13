/*
Provides a resource to create a rum project

Example Usage

```hcl
resource "tencentcloud_rum_project" "project" {
  name = &lt;nil&gt;
  instance_i_d = &lt;nil&gt;
  rate = &lt;nil&gt;
  enable_u_r_l_group = &lt;nil&gt;
  type = &lt;nil&gt;
  repo = &lt;nil&gt;
  u_r_l = &lt;nil&gt;
  desc = &lt;nil&gt;
              }
```

Import

rum project can be imported using the id, e.g.

```
terraform import tencentcloud_rum_project.project project_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rum "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/rum/v20210622"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudRumProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRumProjectCreate,
		Read:   resourceTencentCloudRumProjectRead,
		Update: resourceTencentCloudRumProjectUpdate,
		Delete: resourceTencentCloudRumProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Name of the created project (required and up to 200 characters).",
			},

			"instance_i_d": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Business system ID.",
			},

			"rate": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Project sampling rate (greater than or equal to 0).",
			},

			"enable_u_r_l_group": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Whether to enable aggregation.",
			},

			"type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Project type (valid values: `web`, `mp`, `android`, `ios`, `node`, `hippy`, `weex`, `viola`, `rn`).",
			},

			"repo": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Repository address of the project (optional and up to 256 characters).",
			},

			"u_r_l": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Webpage address of the project (optional and up to 256 characters).",
			},

			"desc": {
				Optional: true,
				Type:     schema.TypeString,
				Description: "	Description of the created project (optional and up to 1,000 characters).",
			},

			"creator": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Creator ID.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Creata Time.",
			},

			"key": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Unique project key (12 characters).",
			},

			"instance_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Instance name.",
			},

			"instance_key": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Instance key.",
			},

			"is_star": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Starred status. `1`: yes; `0`: no.",
			},

			"project_status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Project status (`1`: Creating; `2`: Running; `3`: Abnormal; `4`: Restarting; `5`: Stopping; `6`: Stopped; `7`: Terminating; `8`: Terminated).",
			},
		},
	}
}

func resourceTencentCloudRumProjectCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_project.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = rum.NewCreateProjectRequest()
		response = rum.NewCreateProjectResponse()
		iD       int
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_i_d"); ok {
		request.InstanceID = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rate"); ok {
		request.Rate = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("enable_u_r_l_group"); ok {
		request.EnableURLGroup = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("repo"); ok {
		request.Repo = helper.String(v.(string))
	}

	if v, ok := d.GetOk("u_r_l"); ok {
		request.URL = helper.String(v.(string))
	}

	if v, ok := d.GetOk("desc"); ok {
		request.Desc = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRumClient().CreateProject(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create rum project failed, reason:%+v", logId, err)
		return err
	}

	iD = *response.Response.ID
	d.SetId(helper.Int64ToStr(int64(iD)))

	return resourceTencentCloudRumProjectRead(d, meta)
}

func resourceTencentCloudRumProjectRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_project.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	projectId := d.Id()

	project, err := service.DescribeRumProjectById(ctx, iD)
	if err != nil {
		return err
	}

	if project == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RumProject` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if project.Name != nil {
		_ = d.Set("name", project.Name)
	}

	if project.InstanceID != nil {
		_ = d.Set("instance_i_d", project.InstanceID)
	}

	if project.Rate != nil {
		_ = d.Set("rate", project.Rate)
	}

	if project.EnableURLGroup != nil {
		_ = d.Set("enable_u_r_l_group", project.EnableURLGroup)
	}

	if project.Type != nil {
		_ = d.Set("type", project.Type)
	}

	if project.Repo != nil {
		_ = d.Set("repo", project.Repo)
	}

	if project.URL != nil {
		_ = d.Set("u_r_l", project.URL)
	}

	if project.Desc != nil {
		_ = d.Set("desc", project.Desc)
	}

	if project.Creator != nil {
		_ = d.Set("creator", project.Creator)
	}

	if project.CreateTime != nil {
		_ = d.Set("create_time", project.CreateTime)
	}

	if project.Key != nil {
		_ = d.Set("key", project.Key)
	}

	if project.InstanceName != nil {
		_ = d.Set("instance_name", project.InstanceName)
	}

	if project.InstanceKey != nil {
		_ = d.Set("instance_key", project.InstanceKey)
	}

	if project.IsStar != nil {
		_ = d.Set("is_star", project.IsStar)
	}

	if project.ProjectStatus != nil {
		_ = d.Set("project_status", project.ProjectStatus)
	}

	return nil
}

func resourceTencentCloudRumProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_project.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := rum.NewModifyProjectRequest()

	projectId := d.Id()

	request.ID = &iD

	immutableArgs := []string{"name", "instance_i_d", "rate", "enable_u_r_l_group", "type", "repo", "u_r_l", "desc", "creator", "create_time", "key", "instance_name", "instance_key", "is_star", "project_status"}

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

	if d.HasChange("instance_i_d") {
		if v, ok := d.GetOk("instance_i_d"); ok {
			request.InstanceID = helper.String(v.(string))
		}
	}

	if d.HasChange("rate") {
		if v, ok := d.GetOk("rate"); ok {
			request.Rate = helper.String(v.(string))
		}
	}

	if d.HasChange("enable_u_r_l_group") {
		if v, ok := d.GetOkExists("enable_u_r_l_group"); ok {
			request.EnableURLGroup = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("type") {
		if v, ok := d.GetOk("type"); ok {
			request.Type = helper.String(v.(string))
		}
	}

	if d.HasChange("repo") {
		if v, ok := d.GetOk("repo"); ok {
			request.Repo = helper.String(v.(string))
		}
	}

	if d.HasChange("u_r_l") {
		if v, ok := d.GetOk("u_r_l"); ok {
			request.URL = helper.String(v.(string))
		}
	}

	if d.HasChange("desc") {
		if v, ok := d.GetOk("desc"); ok {
			request.Desc = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRumClient().ModifyProject(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update rum project failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudRumProjectRead(d, meta)
}

func resourceTencentCloudRumProjectDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_project.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}
	projectId := d.Id()

	if err := service.DeleteRumProjectById(ctx, iD); err != nil {
		return err
	}

	return nil
}
