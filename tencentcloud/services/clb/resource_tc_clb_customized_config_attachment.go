package clb

import (
	"context"
	"fmt"
	"log"

	"github.com/pkg/errors"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	clbintl "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/clb/v20180317"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudClbCustomizedConfigAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbCustomizedConfigAttachmentCreate,
		Read:   resourceTencentCloudClbCustomizedConfigAttachmentRead,
		Update: resourceTencentCloudClbCustomizedConfigAttachmentUpdate,
		Delete: resourceTencentCloudClbCustomizedConfigAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of Customized Config.",
			},
			"bind_list": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Associated server or location.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_balancer_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Clb ID.",
						},
						"listener_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Listener ID.",
						},
						"domain": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Domain.",
						},
						"location_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Location ID.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudClbCustomizedConfigAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_customized_config_attachment.create")()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		request  = clbintl.NewAssociateCustomizedConfigRequest()
		configId string
	)

	if v, ok := d.GetOk("config_id"); ok {
		request.UconfigId = helper.String(v.(string))
		configId = v.(string)
	}

	if v, ok := d.GetOk("bind_list"); ok {
		for _, item := range v.(*schema.Set).List() {
			bindItem := clbintl.BindItem{}
			dMap := item.(map[string]interface{})
			if v, ok := dMap["load_balancer_id"]; ok && v.(string) != "" {
				bindItem.LoadBalancerId = helper.String(v.(string))
			}

			if v, ok := dMap["listener_id"]; ok && v.(string) != "" {
				bindItem.ListenerId = helper.String(v.(string))
			}

			if v, ok := dMap["domain"]; ok && v.(string) != "" {
				bindItem.Domain = helper.String(v.(string))
			}

			if v, ok := dMap["location_id"]; ok && v.(string) != "" {
				bindItem.LocationId = helper.String(v.(string))
			}

			request.BindList = append(request.BindList, &bindItem)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbIntlClient().AssociateCustomizedConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			if result == nil || result.Response == nil || result.Response.RequestId == nil {
				return resource.NonRetryableError(fmt.Errorf("Associate CLB Customized Config Failed, Response is nil."))
			}

			requestId := *result.Response.RequestId
			retryErr := waitForTaskFinishIntl(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbIntlClient())
			if retryErr != nil {
				return tccommon.RetryError(errors.WithStack(retryErr))
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s Associate CLB Customized Config Failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(configId)

	return resourceTencentCloudClbCustomizedConfigAttachmentRead(d, meta)
}

func resourceTencentCloudClbCustomizedConfigAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_customized_config_attachment.read")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		configId   = d.Id()
	)

	bindList, err := clbService.DescribeDescribeCustomizedConfigAssociateListById(ctx, configId)
	if err != nil {
		return err
	}

	if bindList == nil || len(bindList) == 0 {
		d.SetId("")
		return fmt.Errorf("resource `tencentcloud_clb_customized_config_attachment` %s does not exist", configId)
	}

	_ = d.Set("config_id", configId)

	tmpList := make([]map[string]interface{}, 0, len(bindList))
	for _, item := range bindList {
		dMap := make(map[string]interface{})
		if item.LoadBalancerId != nil {
			dMap["load_balancer_id"] = *item.LoadBalancerId
		}

		if item.ListenerId != nil {
			dMap["listener_id"] = *item.ListenerId
		}

		if item.Domain != nil {
			dMap["domain"] = *item.Domain
		}

		if item.LocationId != nil {
			dMap["location_id"] = *item.LocationId
		}

		tmpList = append(tmpList, dMap)
	}

	_ = d.Set("bind_list", tmpList)

	return nil
}

func resourceTencentCloudClbCustomizedConfigAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_customized_config_attachment.delete")()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		id    = d.Id()
	)

	if d.HasChange("bind_list") {
		oldInterface, newInterface := d.GetChange("bind_list")
		olds := oldInterface.(*schema.Set)
		news := newInterface.(*schema.Set)
		remove := olds.Difference(news).List()
		add := news.Difference(olds).List()
		if len(remove) > 0 {
			request := clbintl.NewDisassociateCustomizedConfigRequest()
			for _, item := range remove {
				bindItem := clbintl.BindItem{}
				dMap := item.(map[string]interface{})
				if v, ok := dMap["load_balancer_id"]; ok && v.(string) != "" {
					bindItem.LoadBalancerId = helper.String(v.(string))
				}

				if v, ok := dMap["listener_id"]; ok && v.(string) != "" {
					bindItem.ListenerId = helper.String(v.(string))
				}

				if v, ok := dMap["domain"]; ok && v.(string) != "" {
					bindItem.Domain = helper.String(v.(string))
				}

				if v, ok := dMap["location_id"]; ok && v.(string) != "" {
					bindItem.LocationId = helper.String(v.(string))
				}

				request.BindList = append(request.BindList, &bindItem)
			}

			request.UconfigId = helper.String(id)
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbIntlClient().DisassociateCustomizedConfig(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					if result == nil || result.Response == nil || result.Response.RequestId == nil {
						return resource.NonRetryableError(fmt.Errorf("Disassociate CLB Customized Config Failed, Response is nil."))
					}

					requestId := *result.Response.RequestId
					retryErr := waitForTaskFinishIntl(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbIntlClient())
					if retryErr != nil {
						return tccommon.RetryError(errors.WithStack(retryErr))
					}
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s Disassociate CLB Customized Config Failed, reason:%+v", logId, err)
				return err
			}
		}

		if len(add) > 0 {
			request := clbintl.NewAssociateCustomizedConfigRequest()
			request.UconfigId = helper.String(id)
			for _, item := range add {
				bindItem := clbintl.BindItem{}
				dMap := item.(map[string]interface{})
				if v, ok := dMap["load_balancer_id"]; ok && v.(string) != "" {
					bindItem.LoadBalancerId = helper.String(v.(string))
				}

				if v, ok := dMap["listener_id"]; ok && v.(string) != "" {
					bindItem.ListenerId = helper.String(v.(string))
				}

				if v, ok := dMap["domain"]; ok && v.(string) != "" {
					bindItem.Domain = helper.String(v.(string))
				}

				if v, ok := dMap["location_id"]; ok && v.(string) != "" {
					bindItem.LocationId = helper.String(v.(string))
				}

				request.BindList = append(request.BindList, &bindItem)
			}

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbIntlClient().AssociateCustomizedConfig(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					if result == nil || result.Response == nil || result.Response.RequestId == nil {
						return resource.NonRetryableError(fmt.Errorf("Associate CLB Customized Config Failed, Response is nil."))
					}

					requestId := *result.Response.RequestId
					retryErr := waitForTaskFinishIntl(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbIntlClient())
					if retryErr != nil {
						return tccommon.RetryError(errors.WithStack(retryErr))
					}
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s Associate CLB Customized Config Failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	return resourceTencentCloudClbCustomizedConfigAttachmentRead(d, meta)
}

func resourceTencentCloudClbCustomizedConfigAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_customized_config_attachment.delete")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = clbintl.NewDisassociateCustomizedConfigRequest()
		id      = d.Id()
	)

	request.UconfigId = helper.String(id)
	if v, ok := d.GetOk("bind_list"); ok {
		for _, item := range v.(*schema.Set).List() {
			bindItem := clbintl.BindItem{}
			dMap := item.(map[string]interface{})
			if v, ok := dMap["load_balancer_id"]; ok {
				bindItem.LoadBalancerId = helper.String(v.(string))
			}

			if v, ok := dMap["listener_id"]; ok {
				bindItem.ListenerId = helper.String(v.(string))
			}

			if v, ok := dMap["domain"]; ok {
				bindItem.Domain = helper.String(v.(string))
			}

			if v, ok := dMap["location_id"]; ok {
				bindItem.LocationId = helper.String(v.(string))
			}

			request.BindList = append(request.BindList, &bindItem)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbIntlClient().DisassociateCustomizedConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			if result == nil || result.Response == nil || result.Response.RequestId == nil {
				return resource.NonRetryableError(fmt.Errorf("Disassociate CLB Customized Config Failed, Response is nil."))
			}

			requestId := *result.Response.RequestId
			retryErr := waitForTaskFinishIntl(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbIntlClient())
			if retryErr != nil {
				return tccommon.RetryError(errors.WithStack(retryErr))
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s Disassociate CLB Customized Config Failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
