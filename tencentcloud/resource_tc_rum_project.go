/*
Provides a resource to create a rum project

Example Usage

```hcl
resource "tencentcloud_rum_project" "project" {
  name = "projectName"
  instance_id = "rum-pasZKEI3RLgakj"
  rate = "100"
  enable_url_group = "0"
  type = "web"
  repo = ""
  url = "iac-tf.com"
  desc = "projectDesc-1"
}

```
Import

rum project can be imported using the id, e.g.
```
$ terraform import tencentcloud_rum_project.project project_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	rum "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/rum/v20210622"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudRumProject() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudRumProjectRead,
		Create: resourceTencentCloudRumProjectCreate,
		Update: resourceTencentCloudRumProjectUpdate,
		Delete: resourceTencentCloudRumProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the created project (required and up to 200 characters).",
			},

			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Business system ID.",
			},

			"rate": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project sampling rate (greater than or equal to 0).",
			},

			"enable_url_group": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Whether to enable aggregation.",
			},

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project type (valid values: `web`, `mp`, `android`, `ios`, `node`, `hippy`, `weex`, `viola`, `rn`).",
			},

			"repo": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Repository address of the project (optional and up to 256 characters).",
			},

			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Webpage address of the project (optional and up to 256 characters).",
			},

			"desc": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "	Description of the created project (optional and up to 1,000 characters).",
			},

			"creator": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creator ID.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creata Time.",
			},

			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique project key (12 characters).",
			},

			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance name.",
			},

			"instance_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance key.",
			},

			"is_star": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Starred status. `1`: yes; `0`: no.",
			},

			"project_status": {
				Type:        schema.TypeInt,
				Computed:    true,
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
		response *rum.CreateProjectResponse
		id       uint64
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceID = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rate"); ok {
		request.Rate = helper.String(v.(string))
	}

	if v := d.Get("enable_url_group"); v != nil {
		request.EnableURLGroup = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("repo"); ok {
		request.Repo = helper.String(v.(string))
	}

	if v, ok := d.GetOk("url"); ok {
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
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create rum project failed, reason:%+v", logId, err)
		return err
	}

	id = *response.Response.ID

	d.SetId(strconv.Itoa(int(id)))
	return resourceTencentCloudRumProjectRead(d, meta)
}

func resourceTencentCloudRumProjectRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_project.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	projectId := d.Id()

	project, err := service.DescribeRumProject(ctx, projectId)

	if err != nil {
		return err
	}

	if project == nil {
		d.SetId("")
		return fmt.Errorf("resource `project` %s does not exist", projectId)
	}

	if project.Name != nil {
		_ = d.Set("name", project.Name)
	}

	if project.InstanceID != nil {
		_ = d.Set("instance_id", project.InstanceID)
	}

	if project.Rate != nil {
		_ = d.Set("rate", project.Rate)
	}

	if project.EnableURLGroup != nil {
		_ = d.Set("enable_url_group", project.EnableURLGroup)
	}

	if project.Type != nil {
		_ = d.Set("type", project.Type)
	}

	if project.Repo != nil {
		_ = d.Set("repo", project.Repo)
	}

	if project.URL != nil {
		_ = d.Set("url", project.URL)
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

	id, e := strconv.Atoi(projectId)
	if e != nil {
		return fmt.Errorf("[ERROR]%s api[%s] sting to uint64 error, err [%s]", logId, request.GetAction(), e)
	}

	request.ID = helper.Uint64(uint64(id))

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("instance_id") {
		if v, ok := d.GetOk("instance_id"); ok {
			request.InstanceID = helper.String(v.(string))
		}

	}

	if d.HasChange("rate") {
		if v, ok := d.GetOk("rate"); ok {
			request.Rate = helper.String(v.(string))
		}
	}

	if d.HasChange("enable_url_group") {
		if v, ok := d.GetOk("enable_url_group"); ok {
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

	if d.HasChange("url") {
		if v, ok := d.GetOk("url"); ok {
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
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create rum project failed, reason:%+v", logId, err)
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

	if err := service.DeleteRumProjectById(ctx, projectId); err != nil {
		return err
	}

	return nil
}
