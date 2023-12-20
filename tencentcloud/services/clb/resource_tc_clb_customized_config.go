package clb

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClbCustomizedConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbCustomizedConfigCreate,
		Read:   resourceTencentCloudClbCustomizedConfigRead,
		Update: resourceTencentCloudClbCustomizedConfigUpdate,
		Delete: resourceTencentCloudClbCustomizedConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"config_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of Customized Config.",
			},
			"config_content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Content of Customized Config.",
			},
			"load_balancer_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of LoadBalancer Ids.",
			},
			//computed
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

func resourceTencentCloudClbCustomizedConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_customized_config.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	configName := d.Get("config_name").(string)
	configContent := d.Get("config_content").(string)

	request := clb.NewSetCustomizedConfigForLoadBalancerRequest()
	request.OperationType = helper.String("ADD")
	request.ConfigName = helper.String(configName)
	request.ConfigContent = helper.String(configContent)

	var response *clb.SetCustomizedConfigForLoadBalancerResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().SetCustomizedConfigForLoadBalancer(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
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
	d.SetId(*response.Response.ConfigId)

	if v, ok := d.GetOk("load_balancer_ids"); ok {
		loadBalancerIds := v.(*schema.Set).List()
		clbService := ClbService{
			client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
		}
		err := clbService.BindOrUnBindCustomizedConfigWithLbId(ctx,
			"BIND", *response.Response.ConfigId, loadBalancerIds)
		if err != nil {
			log.Printf("[CRITAL]%s Binding LB Customized Config Failed, reason:%+v", logId, err)
			return err
		}
	}
	return resourceTencentCloudClbCustomizedConfigRead(d, meta)
}

func resourceTencentCloudClbCustomizedConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_customized_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	configId := d.Id()
	clbService := ClbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var config *clb.ConfigListItem
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := clbService.DescribeLbCustomizedConfigById(ctx, configId)
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

	_ = d.Set("config_name", config.ConfigName)
	_ = d.Set("config_content", config.ConfigContent)
	_ = d.Set("create_time", config.CreateTimestamp)
	_ = d.Set("update_time", config.UpdateTimestamp)

	request := clb.NewDescribeCustomizedConfigAssociateListRequest()
	request.UconfigId = &configId
	var response *clb.DescribeCustomizedConfigAssociateListResponse
	assErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().DescribeCustomizedConfigAssociateList(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if assErr != nil {
		log.Printf("[CRITAL]%s Describe CLB Customized Config Associate List Failed, reason:%+v", logId, assErr)
		return err
	}
	_ = d.Set("load_balancer_ids", extractBindClbList(response.Response.BindList))

	return nil
}

func resourceTencentCloudClbCustomizedConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_customized_config.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	d.Partial(true)

	configId := d.Id()
	request := clb.NewSetCustomizedConfigForLoadBalancerRequest()
	request.UconfigId = &configId
	request.OperationType = helper.String("UPDATE")

	if d.HasChange("config_name") {
		configName := d.Get("config_name").(string)
		request.ConfigName = &configName
	}

	if d.HasChange("config_content") {
		configContent := d.Get("config_content").(string)
		request.ConfigContent = &configContent
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().SetCustomizedConfigForLoadBalancer(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
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

	if d.HasChange("load_balancer_ids") {
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService := ClbService{
			client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
		}
		old, new := d.GetChange("load_balancer_ids")
		olds := old.(*schema.Set)
		news := new.(*schema.Set)
		add := news.Difference(olds).List()
		remove := olds.Difference(news).List()
		if len(remove) > 0 {
			err := clbService.BindOrUnBindCustomizedConfigWithLbId(ctx,
				"UNBIND", configId, remove)
			if err != nil {
				log.Printf("[CRITAL]%s UnBinding LB Customized Config Failed, reason:%+v", logId, err)
				return err
			}
		}
		if len(add) > 0 {
			err := clbService.BindOrUnBindCustomizedConfigWithLbId(ctx,
				"BIND", configId, add)
			if err != nil {
				log.Printf("[CRITAL]%s Binding LB Customized Config Failed, reason:%+v", logId, err)
				return err
			}
		}
	}
	return nil
}

func resourceTencentCloudClbCustomizedConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_customized_config.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	configId := d.Id()
	clbService := ClbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := clbService.DeleteLbCustomizedConfigById(ctx, configId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete CLB customized config failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
