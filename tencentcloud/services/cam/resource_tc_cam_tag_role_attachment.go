package cam

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCamTagRoleAttachment() *schema.Resource {
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
	defer tccommon.LogElapsed("resource.tencentcloud_cam_tag_role_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCamClient().TagRole(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cam TagRole failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(roleName + tccommon.FILED_SP + roleId)

	return resourceTencentCloudCamTagRoleReadAttachment(d, meta)
}

func resourceTencentCloudCamTagRoleReadAttachment(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_tag_role_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CamService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
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
	defer tccommon.LogElapsed("resource.tencentcloud_cam_tag_role_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CamService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
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
