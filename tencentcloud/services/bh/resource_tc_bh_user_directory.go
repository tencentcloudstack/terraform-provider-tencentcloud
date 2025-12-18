package bh

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bhv20230418 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bh/v20230418"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudBhUserDirectory() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBhUserDirectoryCreate,
		Read:   resourceTencentCloudBhUserDirectoryRead,
		Update: resourceTencentCloudBhUserDirectoryUpdate,
		Delete: resourceTencentCloudBhUserDirectoryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"dir_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Directory ID.",
			},

			"dir_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Directory name.",
			},

			"user_org_set": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "IOA group information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"org_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "IOA user organization ID.",
						},
						"org_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IOA user organization name.",
						},
						"org_id_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IOA user organization ID path.",
						},
						"org_name_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IOA user organization name path.",
						},
						"user_total": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of users under the IOA user organization ID.",
						},
					},
				},
			},

			"source": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "IOA associated user source type.",
			},

			"source_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "IOA associated user source name.",
			},

			// computed
			"directory_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Directory ID.",
			},

			"user_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of users included in the directory.",
			},
		},
	}
}

func resourceTencentCloudBhUserDirectoryCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_user_directory.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request     = bhv20230418.NewCreateUserDirectoryRequest()
		response    = bhv20230418.NewCreateUserDirectoryResponse()
		directoryId string
	)

	if v, ok := d.GetOkExists("dir_id"); ok {
		request.DirId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("dir_name"); ok {
		request.DirName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_org_set"); ok {
		for _, item := range v.([]interface{}) {
			userOrgSetMap := item.(map[string]interface{})
			userOrg := bhv20230418.UserOrg{}
			if v, ok := userOrgSetMap["org_id"].(int); ok {
				userOrg.OrgId = helper.IntUint64(v)
			}

			if v, ok := userOrgSetMap["org_name"].(string); ok && v != "" {
				userOrg.OrgName = helper.String(v)
			}

			if v, ok := userOrgSetMap["org_id_path"].(string); ok && v != "" {
				userOrg.OrgIdPath = helper.String(v)
			}

			if v, ok := userOrgSetMap["org_name_path"].(string); ok && v != "" {
				userOrg.OrgNamePath = helper.String(v)
			}

			if v, ok := userOrgSetMap["user_total"].(int); ok {
				userOrg.UserTotal = helper.IntUint64(v)
			}

			request.UserOrgSet = append(request.UserOrgSet, &userOrg)
		}
	}

	if v, ok := d.GetOkExists("source"); ok {
		request.Source = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("source_name"); ok {
		request.SourceName = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().CreateUserDirectoryWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create bh user directory failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create bh user directory failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Id == nil {
		return fmt.Errorf("Id is nil.")
	}

	directoryId = helper.UInt64ToStr(*response.Response.Id)
	d.SetId(directoryId)
	return resourceTencentCloudBhUserDirectoryRead(d, meta)
}

func resourceTencentCloudBhUserDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_user_directory.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service     = BhService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		directoryId = d.Id()
	)

	respData, err := service.DescribeBhUserDirectoryById(ctx, directoryId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_bh_user_directory` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.DirId != nil {
		_ = d.Set("dir_id", respData.DirId)
	}

	if respData.DirName != nil {
		_ = d.Set("dir_name", respData.DirName)
	}

	if respData.UserOrgSet != nil {
		userOrgSetList := make([]map[string]interface{}, 0, len(respData.UserOrgSet))
		for _, userOrgSet := range respData.UserOrgSet {
			userOrgSetMap := map[string]interface{}{}
			if userOrgSet.OrgId != nil {
				userOrgSetMap["org_id"] = userOrgSet.OrgId
			}

			if userOrgSet.OrgName != nil {
				userOrgSetMap["org_name"] = userOrgSet.OrgName
			}

			if userOrgSet.OrgIdPath != nil {
				userOrgSetMap["org_id_path"] = userOrgSet.OrgIdPath
			}

			if userOrgSet.OrgNamePath != nil {
				userOrgSetMap["org_name_path"] = userOrgSet.OrgNamePath
			}

			if userOrgSet.UserTotal != nil {
				userOrgSetMap["user_total"] = userOrgSet.UserTotal
			}

			userOrgSetList = append(userOrgSetList, userOrgSetMap)
		}

		_ = d.Set("user_org_set", userOrgSetList)
	}

	if respData.Source != nil {
		_ = d.Set("source", respData.Source)
	}

	if respData.SourceName != nil {
		_ = d.Set("source_name", respData.SourceName)
	}

	if respData.UserTotal != nil {
		_ = d.Set("user_count", respData.UserTotal)
	}

	if respData.Id != nil {
		_ = d.Set("directory_id", respData.Id)
	}

	return nil
}

func resourceTencentCloudBhUserDirectoryUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_user_directory.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		directoryId = d.Id()
	)

	if d.HasChange("user_org_set") {
		request := bhv20230418.NewModifyUserDirectoryRequest()
		if v, ok := d.GetOk("user_org_set"); ok {
			for _, item := range v.([]interface{}) {
				userOrgSetMap := item.(map[string]interface{})
				userOrg := bhv20230418.UserOrg{}
				if v, ok := userOrgSetMap["org_id"].(int); ok {
					userOrg.OrgId = helper.IntUint64(v)
				}

				if v, ok := userOrgSetMap["org_name"].(string); ok && v != "" {
					userOrg.OrgName = helper.String(v)
				}

				if v, ok := userOrgSetMap["org_id_path"].(string); ok && v != "" {
					userOrg.OrgIdPath = helper.String(v)
				}

				if v, ok := userOrgSetMap["org_name_path"].(string); ok && v != "" {
					userOrg.OrgNamePath = helper.String(v)
				}

				if v, ok := userOrgSetMap["user_total"].(int); ok {
					userOrg.UserTotal = helper.IntUint64(v)
				}

				request.UserOrgSet = append(request.UserOrgSet, &userOrg)
			}
		}

		request.Id = helper.StrToUint64Point(directoryId)
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().ModifyUserDirectoryWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update bh user directory failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudBhUserDirectoryRead(d, meta)
}

func resourceTencentCloudBhUserDirectoryDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_user_directory.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request     = bhv20230418.NewDeleteUserDirectoryRequest()
		directoryId = d.Id()
	)

	request.IdSet = append(request.IdSet, helper.StrToUint64Point(directoryId))
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().DeleteUserDirectoryWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete bh user directory failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
