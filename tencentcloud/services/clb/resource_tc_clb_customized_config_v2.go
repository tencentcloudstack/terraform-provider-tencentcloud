package clb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	clbintl "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClbCustomizedConfigV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbCustomizedConfigV2Create,
		Read:   resourceTencentCloudClbCustomizedConfigV2Read,
		Update: resourceTencentCloudClbCustomizedConfigV2Update,
		Delete: resourceTencentCloudClbCustomizedConfigV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"config_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of Customized Config.",
			},
			"config_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"CLB", "SERVER", "LOCATION"}),
				Description:  "Type of Customized Config. Valid values: `SERVER` and `LOCATION`.",
			},
			"config_content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Content of Customized Config.",
			},

			//computed
			"config_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of Customized Config.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of Customized Config.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Update time of Customized Config.",
			},
		},
	}
}

func resourceTencentCloudClbCustomizedConfigV2Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_customized_config_v2.create")()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		request  = clbintl.NewAddCustomizedConfigRequest()
		response = clbintl.NewAddCustomizedConfigResponse()
	)

	configType := d.Get("config_type").(string)

	request.ConfigName = helper.String(d.Get("config_name").(string))
	request.ConfigType = helper.String(configType)
	request.ConfigContent = helper.String(d.Get("config_content").(string))

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbIntlClient().AddCustomizedConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			if result == nil || result.Response == nil || result.Response.RequestId == nil {
				return resource.NonRetryableError(fmt.Errorf("Create CLB Customized Config Failed, Response is nil."))
			}

			requestId := *result.Response.RequestId
			retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
			if retryErr != nil {
				return tccommon.RetryError(errors.WithStack(retryErr))
			}
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s Create CLB Customized Config Failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.ConfigId == nil {
		return fmt.Errorf("ConfigId is nil.")
	}

	configId := *response.Response.ConfigId

	d.SetId(strings.Join([]string{configId, configType}, tccommon.FILED_SP))

	return resourceTencentCloudClbCustomizedConfigV2Read(d, meta)
}

func resourceTencentCloudClbCustomizedConfigV2Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_customized_config_v2.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	configId := idSplit[0]
	configType := idSplit[1]

	var config *clbintl.ConfigListItem
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := clbService.DescribeLbIntlCustomizedConfigById(ctx, configId, configType)
		if e != nil {
			return tccommon.RetryError(e)
		}

		config = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read CLB customized config failed, reason:%+v", logId, err)
		return err
	}

	if config == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("config_id", configId)
	_ = d.Set("config_type", configType)
	_ = d.Set("config_name", config.ConfigName)
	_ = d.Set("config_content", config.ConfigContent)
	_ = d.Set("create_time", config.CreateTimestamp)
	_ = d.Set("update_time", config.UpdateTimestamp)

	return nil
}

func resourceTencentCloudClbCustomizedConfigV2Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_customized_config_v2.update")()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	configId := idSplit[0]

	d.Partial(true)

	if d.HasChange("config_name") || d.HasChange("config_content") {
		request := clbintl.NewModifyCustomizedConfigRequest()
		request.UconfigId = &configId
		configName := d.Get("config_name").(string)
		configContent := d.Get("config_content").(string)
		request.ConfigName = &configName
		request.ConfigContent = &configContent

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbIntlClient().ModifyCustomizedConfig(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				if result == nil || result.Response == nil || result.Response.RequestId == nil {
					return resource.NonRetryableError(fmt.Errorf("Update CLB Customized Config Failed, Response is nil."))
				}

				requestId := *result.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
				if retryErr != nil {
					return tccommon.RetryError(errors.WithStack(retryErr))
				}
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s Update CLB Customized Config Failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudClbCustomizedConfigV2Read(d, meta)
}

func resourceTencentCloudClbCustomizedConfigV2Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_customized_config_v2.delete")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = clbintl.NewDeleteCustomizedConfigRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	configId := idSplit[0]

	request.UconfigIdList = []*string{&configId}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbIntlClient().DeleteCustomizedConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			if result == nil || result.Response == nil || result.Response.RequestId == nil {
				return resource.NonRetryableError(fmt.Errorf("Delete CLB Customized Config Failed, Response is nil."))
			}

			requestId := *result.Response.RequestId
			retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
			if retryErr != nil {
				return tccommon.RetryError(errors.WithStack(retryErr))
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s Delete CLB Customized Config Failed, reason:%+v", logId, err)
		return err
	}
	return nil
}
