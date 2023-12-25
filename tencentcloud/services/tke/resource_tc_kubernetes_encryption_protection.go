package tke

import (
	"context"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudKubernetesEncryptionProtection() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKubernetesEncryptionProtectionCreate,
		Read:   resourceTencentCloudKubernetesEncryptionProtectionRead,
		Delete: resourceTencentCloudKubernetesEncryptionProtectionDelete,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "cluster id.",
			},

			"kms_configuration": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "kms encryption configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "kms id.",
						},
						"kms_region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "kms region.",
						},
					},
				},
			},
			"status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "kms encryption status.",
			},
		},
	}
}

func resourceTencentCloudKubernetesEncryptionProtectionCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tke_encryption_protection.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request   = tke.NewEnableEncryptionProtectionRequest()
		clusterId string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
		clusterId = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "kms_configuration"); ok {
		kMSConfiguration := tke.KMSConfiguration{}
		if v, ok := dMap["key_id"]; ok {
			kMSConfiguration.KeyId = helper.String(v.(string))
		}
		if v, ok := dMap["kms_region"]; ok {
			kMSConfiguration.KmsRegion = helper.String(v.(string))
		}
		request.KMSConfiguration = &kMSConfiguration
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeClient().EnableEncryptionProtection(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tke encryptionProtection failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(clusterId)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"Opened"}, 3*tccommon.ReadRetryTimeout, time.Second, service.TkeEncryptionProtectionStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudKubernetesEncryptionProtectionRead(d, meta)
}

func resourceTencentCloudKubernetesEncryptionProtectionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tke_encryption_protection.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	encryptionProtectionId := d.Id()

	encryptionProtection, err := service.DescribeTkeEncryptionProtectionById(ctx, encryptionProtectionId)
	if err != nil {
		return err
	}

	if encryptionProtection == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TkeEncryptionProtection` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if encryptionProtection.Status != nil {
		_ = d.Set("status", encryptionProtection.Status)
	}

	return nil
}

func resourceTencentCloudKubernetesEncryptionProtectionDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tke_encryption_protection.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	encryptionProtectionId := d.Id()

	if err := service.DeleteTkeEncryptionProtectionById(ctx, encryptionProtectionId); err != nil {
		return err
	}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"Closed"}, 3*tccommon.ReadRetryTimeout, time.Second, service.TkeEncryptionProtectionStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
