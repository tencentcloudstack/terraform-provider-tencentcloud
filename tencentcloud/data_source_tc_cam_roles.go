/*
Use this data source to query detailed information of CAM roles

Example Usage

```hcl
# query by role_id
data "tencentcloud_cam_roles" "foo" {
  role_id     = tencentcloud_cam_role.foo.id
}

# query by name
data "tencentcloud_cam_roles" "bar" {
  name        = "cam-role-test"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCamRoles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCamRolesRead,

		Schema: map[string]*schema.Schema{
			"role_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of the CAM role to be queried.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the CAM role to be queried.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the CAM policy to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"role_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of CAM roles. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of CAM role.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of CAM role.",
						},
						"document": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy document of CAM role.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of CAM role.",
						},
						"console_login": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicate whether the CAM role can be login or not.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The create time of the CAM role.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The last update time of the CAM role.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCamRolesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cam_roles.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	params := make(map[string]interface{})
	if v, ok := d.GetOk("role_id"); ok {
		params["role_id"] = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		params["name"] = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		params["description"] = v.(string)
	}

	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var roles []*cam.RoleInfo
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := camService.DescribeRolesByFilter(ctx, params)
		if e != nil {
			return retryError(e)
		}
		roles = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM roles failed, reason:%s\n", logId, err.Error())
		return err
	}
	roleList := make([]map[string]interface{}, 0, len(roles))
	ids := make([]string, 0, len(roles))
	for _, role := range roles {
		mapping := map[string]interface{}{
			"role_id":     *role.RoleId,
			"name":        *role.RoleName,
			"document":    *role.PolicyDocument,
			"description": *role.Description,
			"create_time": *role.AddTime,
			"update_time": *role.UpdateTime,
		}
		if int(*role.ConsoleLogin) == 1 {
			mapping["console_login"] = true
		} else {
			mapping["console_login"] = false
		}
		roleList = append(roleList, mapping)
		ids = append(ids, *role.RoleId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("role_list", roleList); e != nil {
		log.Printf("[CRITAL]%s provider set CAM role list fail, reason:%s\n", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), roleList); e != nil {
			return e
		}
	}

	return nil
}
