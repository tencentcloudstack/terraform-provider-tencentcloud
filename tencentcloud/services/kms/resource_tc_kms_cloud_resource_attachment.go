package kms

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudKmsCloudResourceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKmsCloudResourceAttachmentCreate,
		Read:   resourceTencentCloudKmsCloudResourceAttachmentRead,
		Delete: resourceTencentCloudKmsCloudResourceAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"key_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "CMK unique identifier.",
			},
			"product_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "A unique identifier for the cloud product.",
			},
			"resource_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The resource/instance ID of the cloud product.",
			},
			// computed
			"alias": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Alias.",
			},
			"description": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Description.",
			},
			"key_state": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Key state.",
			},
			"key_usage": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Key usage.",
			},
			"owner": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "owner.",
			},
		},
	}
}

func resourceTencentCloudKmsCloudResourceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kms_cloud_resource_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = kms.NewBindCloudResourceRequest()
		keyId      string
		productId  string
		resourceId string
	)

	if v, ok := d.GetOk("key_id"); ok {
		request.KeyId = helper.String(v.(string))
		keyId = v.(string)
	}

	if v, ok := d.GetOk("product_id"); ok {
		request.ProductId = helper.String(v.(string))
		productId = v.(string)
	}

	if v, ok := d.GetOk("resource_id"); ok {
		request.ResourceId = helper.String(v.(string))
		resourceId = v.(string)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseKmsClient().BindCloudResource(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create kms cloudResourceAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{keyId, productId, resourceId}, tccommon.FILED_SP))
	return resourceTencentCloudKmsCloudResourceAttachmentRead(d, meta)
}

func resourceTencentCloudKmsCloudResourceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kms_cloud_resource_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = KmsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	keyId := idSplit[0]
	productId := idSplit[1]
	resourceId := idSplit[2]

	cloudResourceAttachment, err := service.DescribeKmsCloudResourceAttachmentById(ctx, keyId)
	if err != nil {
		return err
	}

	if cloudResourceAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `KmsCloudResourceAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("key_id", keyId)
	_ = d.Set("product_id", productId)
	_ = d.Set("resource_id", resourceId)

	if cloudResourceAttachment.Alias != nil {
		_ = d.Set("alias", cloudResourceAttachment.Alias)
	}

	if cloudResourceAttachment.Description != nil {
		_ = d.Set("description", cloudResourceAttachment.Description)
	}

	if cloudResourceAttachment.KeyState != nil {
		_ = d.Set("key_state", cloudResourceAttachment.KeyState)
	}

	if cloudResourceAttachment.KeyUsage != nil {
		_ = d.Set("key_usage", cloudResourceAttachment.KeyUsage)
	}

	if cloudResourceAttachment.Owner != nil {
		_ = d.Set("owner", cloudResourceAttachment.Owner)
	}

	return nil
}

func resourceTencentCloudKmsCloudResourceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kms_cloud_resource_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = KmsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	keyId := idSplit[0]
	productId := idSplit[1]
	resourceId := idSplit[2]

	if err := service.DeleteKmsCloudResourceAttachmentById(ctx, keyId, productId, resourceId); err != nil {
		return err
	}

	return nil
}
