package apm

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apmv20210622 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apm/v20210622"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudApmAssociationConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudApmAssociationConfigCreate,
		Read:   resourceTencentCloudApmAssociationConfigRead,
		Update: resourceTencentCloudApmAssociationConfigUpdate,
		Delete: resourceTencentCloudApmAssociationConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"product_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Associated product name. currently only supports Prometheus.",
			},

			"status": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{1, 2}),
				Description:  "Status of the association relationship: // association status: 1 (enabled), 2 (disabled).",
			},

			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Business system ID.",
			},

			"peer_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Associated product instance ID.",
			},

			"topic": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the CKafka message topic.",
			},
		},
	}
}

func resourceTencentCloudApmAssociationConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_association_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		instanceId  string
		productName string
	)

	if v, ok := d.GetOk("product_name"); ok {
		productName = v.(string)
	}

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(strings.Join([]string{instanceId, productName}, tccommon.FILED_SP))
	return resourceTencentCloudApmAssociationConfigUpdate(d, meta)
}

func resourceTencentCloudApmAssociationConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_association_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ApmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	productName := idSplit[1]

	respData, err := service.DescribeApmAssociationById(ctx, instanceId, productName)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_apm_association_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("product_name", productName)

	if respData.PeerId != nil {
		_ = d.Set("peer_id", respData.PeerId)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.Topic != nil {
		_ = d.Set("topic", respData.Topic)
	}

	return nil
}

func resourceTencentCloudApmAssociationConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_association_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	productName := idSplit[1]

	request := apmv20210622.NewModifyApmAssociationRequest()
	if v, ok := d.GetOkExists("status"); ok {
		request.Status = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("peer_id"); ok {
		request.PeerId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topic"); ok {
		request.Topic = helper.String(v.(string))
	}

	request.InstanceId = &instanceId
	request.ProductName = &productName
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseApmClient().ModifyApmAssociationWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update apm association failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudApmAssociationConfigRead(d, meta)
}

func resourceTencentCloudApmAssociationConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_association_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
