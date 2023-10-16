/*
Provides a resource to create a cam tag_role

Example Usage

```hcl
resource "tencentcloud_cam_tag_role_attachment" "tag_role" {
  tags {
    key = "test1"
    value = "test1"
  }
  role_id = "test-cam-tag"
}
```

Import

cam tag_role can be imported using the id, e.g.

```
terraform import tencentcloud_cam_tag_role_attachment.tag_role tag_role_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCamTagRoleAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamTagRoleCreateAttachment,
		Read:   resourceTencentCloudCamTagRoleReadAttachment,
		Delete: resourceTencentCloudCamTagRoleDeleteAttachment,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"tags": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Label.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Label.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Label.",
						},
					},
				},
			},

			"role_name": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Character name, at least one input with the character ID.",
			},

			"role_id": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Character ID, at least one input with the character name.",
			},
		},
	}
}

func resourceTencentCloudCamTagRoleCreateAttachment(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_tag_role_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cam.NewTagRoleRequest()
		roleName string
		roleId   string
	)
	if v, ok := d.GetOk("tags"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			roleTags := cam.RoleTags{}
			if v, ok := dMap["key"]; ok {
				roleTags.Key = helper.String(v.(string))
			}
			if v, ok := dMap["value"]; ok {
				roleTags.Value = helper.String(v.(string))
			}
			request.Tags = append(request.Tags, &roleTags)
		}
	}

	if v, ok := d.GetOk("role_name"); ok {
		roleName = v.(string)
		request.RoleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("role_id"); ok {
		roleId = v.(string)
		request.RoleId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().TagRole(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cam TagRole failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(roleName + FILED_SP + roleId)

	return resourceTencentCloudCamTagRoleReadAttachment(d, meta)
}

func resourceTencentCloudCamTagRoleReadAttachment(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_tag_role_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	roleName := idSplit[0]
	roleId := idSplit[1]
	TagRole, err := service.DescribeCamTagRoleById(ctx, roleName, roleId)
	if err != nil {
		return err
	}

	if TagRole == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CamTagRole` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if TagRole.Tags != nil {
		tagsList := []interface{}{}
		for _, tags := range TagRole.Tags {
			tagsMap := map[string]interface{}{}

			if tags.Key != nil {
				tagsMap["key"] = tags.Key
			}

			if tags.Value != nil {
				tagsMap["value"] = tags.Value
			}

			tagsList = append(tagsList, tagsMap)
		}

		_ = d.Set("tags", tagsList)

	}

	if TagRole.RoleName != nil {
		_ = d.Set("role_name", TagRole.RoleName)
	}

	if TagRole.RoleId != nil {
		_ = d.Set("role_id", TagRole.RoleId)
	}

	return nil
}

func resourceTencentCloudCamTagRoleDeleteAttachment(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_tag_role_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	roleName := idSplit[0]
	roleId := idSplit[1]

	keys := make([]*string, 0)
	if v, ok := d.GetOk("tags"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			var key string
			if v, ok := dMap["key"]; ok {
				key = v.(string)
			}
			keys = append(keys, &key)
		}
	}
	if err := service.DeleteCamTagRoleById(ctx, roleName, roleId, keys); err != nil {
		return err
	}

	return nil
}
