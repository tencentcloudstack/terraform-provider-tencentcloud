package css

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCssBackupStream() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssBackupStreamCreate,
		Read:   resourceTencentCloudCssBackupStreamRead,
		Update: resourceTencentCloudCssBackupStreamUpdate,
		Delete: resourceTencentCloudCssBackupStreamDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"push_domain_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Push domain.",
			},

			"app_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "App name.",
			},

			"stream_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Stream id.",
			},

			"upstream_sequence": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Sequence.",
			},
		},
	}
}

func resourceTencentCloudCssBackupStreamCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_css_backup_stream.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var pushDomainName string
	if v, ok := d.GetOk("push_domain_name"); ok {
		pushDomainName = v.(string)
	}

	var appName string
	if v, ok := d.GetOk("app_name"); ok {
		appName = v.(string)
	}

	var streamName string
	if v, ok := d.GetOk("stream_name"); ok {
		streamName = v.(string)
	}

	d.SetId(strings.Join([]string{pushDomainName, appName, streamName}, tccommon.FILED_SP))

	return resourceTencentCloudCssBackupStreamUpdate(d, meta)
}

func resourceTencentCloudCssBackupStreamRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_css_backup_stream.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CssService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	pushDomainName := idSplit[0]
	appName := idSplit[1]
	streamName := idSplit[2]

	backupStream, err := service.DescribeCssBackupStreamById(ctx, pushDomainName, appName, streamName)
	if err != nil {
		return err
	}

	if backupStream == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CssBackupStream` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("push_domain_name", pushDomainName)
	_ = d.Set("app_name", appName)
	_ = d.Set("stream_name", streamName)

	if backupStream.UpstreamSequence != nil {
		_ = d.Set("upstream_sequence", backupStream.UpstreamSequence)
	}

	return nil
}

func resourceTencentCloudCssBackupStreamUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_css_backup_stream.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := css.NewSwitchBackupStreamRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	pushDomainName := idSplit[0]
	appName := idSplit[1]
	streamName := idSplit[2]

	request.PushDomainName = &pushDomainName
	request.AppName = &appName
	request.StreamName = &streamName

	if v, ok := d.GetOk("upstream_sequence"); ok {
		request.UpstreamSequence = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCssClient().SwitchBackupStream(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update css backupStream failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCssBackupStreamRead(d, meta)
}

func resourceTencentCloudCssBackupStreamDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_css_backup_stream.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
