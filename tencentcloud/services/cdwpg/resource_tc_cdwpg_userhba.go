package cdwpg

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdwpgv20201230 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwpg/v20201230"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCdwpgUserhba() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdwpgUserhbaCreate,
		Read:   resourceTencentCloudCdwpgUserhbaRead,
		Update: resourceTencentCloudCdwpgUserhbaUpdate,
		Delete: resourceTencentCloudCdwpgUserhbaDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance id.",
			},

			"hba_configs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "HBA configuration array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type.",
						},
						"database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database.",
						},
						"user": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "User.",
						},
						"address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP address.",
						},
						"mask": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Mask.",
						},
						"method": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Method.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCdwpgUserhbaCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdwpg_userhba.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	instanceId := d.Get("instance_id").(string)

	d.SetId(instanceId)

	return resourceTencentCloudCdwpgUserhbaUpdate(d, meta)
}

func resourceTencentCloudCdwpgUserhbaRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdwpg_userhba.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := CdwpgService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()
	_ = d.Set("instance_id", instanceId)
	respData, err := service.DescribeCdwpgUserhbaById(ctx, instanceId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `cdwpg_userhba` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	hbaConfigsList := make([]map[string]interface{}, 0, len(respData.HbaConfigs))
	if respData.HbaConfigs != nil {
		for _, hbaConfigs := range respData.HbaConfigs {
			hbaConfigsMap := map[string]interface{}{}

			if hbaConfigs.Type != nil {
				hbaConfigsMap["type"] = hbaConfigs.Type
			}

			if hbaConfigs.Database != nil {
				hbaConfigsMap["database"] = hbaConfigs.Database
			}

			if hbaConfigs.User != nil {
				hbaConfigsMap["user"] = hbaConfigs.User
			}

			if hbaConfigs.Address != nil {
				hbaConfigsMap["address"] = hbaConfigs.Address
			}

			if hbaConfigs.Mask != nil {
				hbaConfigsMap["mask"] = hbaConfigs.Mask
			}

			if hbaConfigs.Method != nil {
				hbaConfigsMap["method"] = hbaConfigs.Method
			}

			hbaConfigsList = append(hbaConfigsList, hbaConfigsMap)
		}

		_ = d.Set("hba_configs", hbaConfigsList)
	}

	return nil
}

func resourceTencentCloudCdwpgUserhbaUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdwpg_userhba.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	instanceId := d.Id()

	needChange := false
	mutableArgs := []string{"hba_configs"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := cdwpgv20201230.NewModifyUserHbaRequest()

		if v, ok := d.GetOk("instance_id"); ok {
			request.InstanceId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("hba_configs"); ok {
			for _, item := range v.([]interface{}) {
				hbaConfigsMap := item.(map[string]interface{})
				hbaConfig := cdwpgv20201230.HbaConfig{}
				if v, ok := hbaConfigsMap["type"]; ok {
					hbaConfig.Type = helper.String(v.(string))
				}
				if v, ok := hbaConfigsMap["database"]; ok {
					hbaConfig.Database = helper.String(v.(string))
				}
				if v, ok := hbaConfigsMap["user"]; ok {
					hbaConfig.User = helper.String(v.(string))
				}
				if v, ok := hbaConfigsMap["address"]; ok {
					hbaConfig.Address = helper.String(v.(string))
				}
				if v, ok := hbaConfigsMap["mask"]; ok {
					hbaConfig.Mask = helper.String(v.(string))
				}
				if v, ok := hbaConfigsMap["method"]; ok {
					hbaConfig.Method = helper.String(v.(string))
				}
				request.HbaConfigs = append(request.HbaConfigs, &hbaConfig)
			}
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCdwpgV20201230Client().ModifyUserHbaWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cdwpg userhba failed, reason:%+v", logId, err)
			return err
		}
	}

	_ = instanceId

	service := CdwpgService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"Serving"}, 10*tccommon.ReadRetryTimeout, time.Second, service.InstanceStateRefreshFunc(instanceId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCdwpgUserhbaRead(d, meta)
}

func resourceTencentCloudCdwpgUserhbaDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdwpg_userhba.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
