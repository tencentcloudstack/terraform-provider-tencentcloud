package oceanus

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudOceanusResourceConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOceanusResourceConfigCreate,
		Read:   resourceTencentCloudOceanusResourceConfigRead,
		Update: resourceTencentCloudOceanusResourceConfigUpdate,
		Delete: resourceTencentCloudOceanusResourceConfigDelete,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Resource ID.",
			},
			"resource_loc": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Resource location.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The available storage types for resource location are currently limited to 1:COS.",
						},
						"param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "Json to describe resource location.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Resource bucket.",
									},
									"path": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Resource path.",
									},
									"region": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Resource region, if not set, use resource region, note: this field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
					},
				},
			},
			"remark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Resource description.",
			},
			"work_space_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Workspace SerialId.",
			},
			"version": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Resource Config Version.",
			},
		},
	}
}

func resourceTencentCloudOceanusResourceConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_resource_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = oceanus.NewCreateResourceConfigRequest()
		response   = oceanus.NewCreateResourceConfigResponse()
		resourceId string
		version    string
	)

	if v, ok := d.GetOk("resource_id"); ok {
		request.ResourceId = helper.String(v.(string))
		resourceId = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "resource_loc"); ok {
		resourceLoc := oceanus.ResourceLoc{}
		if v, ok := dMap["storage_type"]; ok {
			resourceLoc.StorageType = helper.IntInt64(v.(int))
		}

		if paramMap, ok := helper.InterfaceToMap(dMap, "param"); ok {
			resourceLocParam := oceanus.ResourceLocParam{}
			if v, ok := paramMap["bucket"]; ok {
				resourceLocParam.Bucket = helper.String(v.(string))
			}

			if v, ok := paramMap["path"]; ok {
				resourceLocParam.Path = helper.String(v.(string))
			}

			if v, ok := paramMap["region"]; ok {
				resourceLocParam.Region = helper.String(v.(string))
			}

			resourceLoc.Param = &resourceLocParam
		}

		request.ResourceLoc = &resourceLoc
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("work_space_id"); ok {
		request.WorkSpaceId = helper.String(v.(string))
	}

	request.AutoDelete = helper.IntInt64(0)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOceanusClient().CreateResourceConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("oceanus resourceConfig not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create oceanus resourceConfig failed, reason:%+v", logId, err)
		return err
	}

	versionInt := *response.Response.Version
	version = strconv.FormatInt(versionInt, 10)
	d.SetId(strings.Join([]string{resourceId, version}, tccommon.FILED_SP))

	return resourceTencentCloudOceanusResourceConfigRead(d, meta)
}

func resourceTencentCloudOceanusResourceConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_resource_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = OceanusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	resourceId := idSplit[0]
	version := idSplit[1]

	resourceConfig, err := service.DescribeOceanusResourceConfigById(ctx, resourceId, version)
	if err != nil {
		return err
	}

	if resourceConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `OceanusResourceConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if resourceConfig.ResourceId != nil {
		_ = d.Set("resource_id", resourceConfig.ResourceId)
	}

	if resourceConfig.ResourceLoc != nil {
		resourceLocMap := map[string]interface{}{}

		if resourceConfig.ResourceLoc.StorageType != nil {
			resourceLocMap["storage_type"] = resourceConfig.ResourceLoc.StorageType
		}

		if resourceConfig.ResourceLoc.Param != nil {
			paramMap := map[string]interface{}{}

			if resourceConfig.ResourceLoc.Param.Bucket != nil {
				paramMap["bucket"] = resourceConfig.ResourceLoc.Param.Bucket
			}

			if resourceConfig.ResourceLoc.Param.Path != nil {
				paramMap["path"] = resourceConfig.ResourceLoc.Param.Path
			}

			if resourceConfig.ResourceLoc.Param.Region != nil {
				paramMap["region"] = resourceConfig.ResourceLoc.Param.Region
			}

			resourceLocMap["param"] = []interface{}{paramMap}
		}

		_ = d.Set("resource_loc", []interface{}{resourceLocMap})
	}

	if resourceConfig.Remark != nil {
		_ = d.Set("remark", resourceConfig.Remark)
	}

	if resourceConfig.Version != nil {
		_ = d.Set("version", resourceConfig.Version)
	}

	return nil
}

func resourceTencentCloudOceanusResourceConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_resource_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	immutableArgs := []string{"resource_id", "resource_loc", "remark", "work_space_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	return resourceTencentCloudOceanusResourceConfigRead(d, meta)
}

func resourceTencentCloudOceanusResourceConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_resource_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = OceanusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	resourceId := idSplit[0]
	version := idSplit[1]

	if err := service.DeleteOceanusResourceConfigById(ctx, resourceId, version); err != nil {
		return err
	}

	return nil
}
