package lighthouse

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
)

func ResourceTencentCloudLighthouseKeyPairAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseKeyPairAttachmentCreate,
		Read:   resourceTencentCloudLighthouseKeyPairAttachmentRead,
		Delete: resourceTencentCloudLighthouseKeyPairAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"key_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Key pair ID.",
			},

			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
		},
	}
}

func resourceTencentCloudLighthouseKeyPairAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_key_pair_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = lighthouse.NewAssociateInstancesKeyPairsRequest()
	)

	keyId := d.Get("key_id").(string)
	instanceId := d.Get("instance_id").(string)

	request.KeyIds = []*string{&keyId}
	request.InstanceIds = []*string{&instanceId}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().AssociateInstancesKeyPairs(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create lighthouse keyPairAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(keyId + tccommon.FILED_SP + instanceId)

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*tccommon.ReadRetryTimeout, time.Second, service.LighthouseInstanceStateRefreshFunc(instanceId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudLighthouseKeyPairAttachmentRead(d, meta)
}

func resourceTencentCloudLighthouseKeyPairAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_key_pair_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	keyId := idSplit[0]
	instanceId := idSplit[1]

	keyPairAttachment, err := service.DescribeLighthouseKeyPairAttachmentById(ctx, keyId)
	if err != nil {
		return err
	}

	if keyPairAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LighthouseKeyPairAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if keyPairAttachment.KeyId != nil {
		_ = d.Set("key_id", keyPairAttachment.KeyId)
	}

	if keyPairAttachment.AssociatedInstanceIds != nil {
		for _, v := range keyPairAttachment.AssociatedInstanceIds {
			if *v == instanceId {
				_ = d.Set("instance_id", *v)
				break
			}
		}
	}

	return nil
}

func resourceTencentCloudLighthouseKeyPairAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_key_pair_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	keyId := idSplit[0]
	instanceId := idSplit[1]

	if err := service.DeleteLighthouseKeyPairAttachmentById(ctx, keyId, instanceId); err != nil {
		return err
	}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*tccommon.ReadRetryTimeout, time.Second, service.LighthouseInstanceStateRefreshFunc(instanceId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
