package tpulsar

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctdmq "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmq"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
)

func ResourceTencentCloudTdmqNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqNamespaceCreate,
		Read:   resourceTencentCloudTdmqNamespaceRead,
		Update: resourceTencentCloudTdmqNamespaceUpdate,
		Delete: resourceTencentCloudTdmqNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"environ_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of namespace to be created.",
			},
			"msg_ttl": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The expiration time of unconsumed message.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Dedicated Cluster Id.",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the namespace.",
			},
			"retention_policy": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "The Policy of message to retain. Format like: `{time_in_minutes: Int, size_in_mb: Int}`. `time_in_minutes`: the time of message to retain; `size_in_mb`: the size of message to retain.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time_in_minutes": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "the time of message to retain.",
						},
						"size_in_mb": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "the size of message to retain.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTdmqNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_namespace.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tdmqService := svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	var (
		environ_name    string
		msg_ttl         uint64
		remark          string
		clusterId       string
		retentionPolicy tdmq.RetentionPolicy
	)
	if temp, ok := d.GetOk("environ_name"); ok {
		environ_name = temp.(string)
		if len(environ_name) < 1 {
			return fmt.Errorf("environ_name should be not empty string")
		}
	}

	msg_ttl = uint64(d.Get("msg_ttl").(int))

	if temp, ok := d.GetOk("cluster_id"); ok {
		clusterId = temp.(string)
	}

	if temp, ok := d.GetOk("remark"); ok {
		remark = temp.(string)
	}

	if temp, ok := d.GetOk("retention_policy"); ok {
		policy := temp.([]interface{})
		for _, item := range policy {
			value := item.(map[string]interface{})
			timeInMinutes := int64(value["time_in_minutes"].(int))
			sizeInMB := int64(value["size_in_mb"].(int))
			retentionPolicy.TimeInMinutes = &timeInMinutes
			retentionPolicy.SizeInMB = &sizeInMB
		}
	}
	environId, err := tdmqService.CreateTdmqNamespace(ctx, environ_name, msg_ttl, clusterId, remark, retentionPolicy)
	if err != nil {
		return err
	}

	d.SetId(strings.Join([]string{environId, clusterId}, tccommon.FILED_SP))
	return resourceTencentCloudTdmqNamespaceRead(d, meta)
}

func resourceTencentCloudTdmqNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_namespace.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	environId := idSplit[0]
	clusterId := idSplit[1]

	tdmqService := svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		info, has, e := tdmqService.DescribeTdmqNamespaceById(ctx, environId, clusterId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if !has {
			d.SetId("")
			return nil
		}

		_ = d.Set("environ_name", info.EnvironmentId)
		_ = d.Set("cluster_id", clusterId)
		_ = d.Set("msg_ttl", info.MsgTTL)
		_ = d.Set("remark", info.Remark)

		tmpList := make([]map[string]interface{}, 0)
		retentionPolicy := make(map[string]interface{}, 2)
		retentionPolicy["time_in_minutes"] = info.RetentionPolicy.TimeInMinutes
		retentionPolicy["size_in_mb"] = info.RetentionPolicy.SizeInMB
		tmpList = append(tmpList, retentionPolicy)
		_ = d.Set("retention_policy", tmpList)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudTdmqNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_instance.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	environId := idSplit[0]
	clusterId := idSplit[1]

	service := svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	var (
		msgTtl          uint64
		remark          string
		retentionPolicy = new(tdmq.RetentionPolicy)
	)

	old, now := d.GetChange("msg_ttl")
	if d.HasChange("msg_ttl") {
		msgTtl = uint64(now.(int))
	} else {
		msgTtl = uint64(old.(int))
	}

	old, now = d.GetChange("remark")
	if d.HasChange("remark") {
		remark = now.(string)
	} else {
		remark = old.(string)
	}

	_, now = d.GetChange("retention_policy")
	if d.HasChange("retention_policy") {
		policy := now.([]interface{})

		for _, item := range policy {
			value := item.(map[string]interface{})
			timeInMinutes := int64(value["time_in_minutes"].(int))
			sizeInMB := int64(value["size_in_mb"].(int))
			retentionPolicy.TimeInMinutes = &timeInMinutes
			retentionPolicy.SizeInMB = &sizeInMB
		}
	}

	d.Partial(true)
	if err := service.ModifyTdmqNamespaceAttribute(ctx, environId, msgTtl, remark, clusterId, retentionPolicy); err != nil {
		return err
	}

	d.Partial(false)
	return resourceTencentCloudTdmqNamespaceRead(d, meta)
}

func resourceTencentCloudTdmqNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_instance.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	environId := idSplit[0]
	clusterId := idSplit[1]

	service := svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		if err := service.DeleteTdmqNamespace(ctx, environId, clusterId); err != nil {
			if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
				if sdkErr.Code == svcvpc.VPCNotFound {
					return nil
				}
			}
			return resource.RetryableError(err)
		}
		return nil
	})

	return err
}
