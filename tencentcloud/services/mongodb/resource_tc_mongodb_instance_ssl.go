package mongodb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudMongodbInstanceSsl() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMongodbInstanceSslCreate,
		Read:   resourceTencentCloudMongodbInstanceSslRead,
		Update: resourceTencentCloudMongodbInstanceSslUpdate,
		Delete: resourceTencentCloudMongodbInstanceSslDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, for example: cmgo-p8vnipr5.",
			},

			"enable": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable SSL. Valid values: `true` - enable SSL, `false` - disable SSL.",
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "SSL status. Valid values: `0` - SSL is disabled, `1` - SSL is enabled.",
			},

			"expired_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Certificate expiration time, format: 2023-05-01 12:00:00. This field is only available when SSL is enabled.",
			},

			"cert_url": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Certificate download link. This field is only available when SSL is enabled.",
			},
		},
	}
}

func resourceTencentCloudMongodbInstanceSslCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_ssl.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMongodbInstanceSslUpdate(d, meta)
}

func resourceTencentCloudMongodbInstanceSslRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_ssl.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	sslStatus, err := service.DescribeMongodbInstanceSSLById(ctx, instanceId)
	if err != nil {
		return err
	}

	if sslStatus == nil || sslStatus.Response == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MongodbInstanceSsl` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if sslStatus.Response.Status != nil {
		_ = d.Set("status", sslStatus.Response.Status)
		if *sslStatus.Response.Status == 1 {
			_ = d.Set("enable", true)
		} else {
			_ = d.Set("enable", false)
		}
	}

	if sslStatus.Response.ExpiredTime != nil {
		_ = d.Set("expired_time", sslStatus.Response.ExpiredTime)
	}

	if sslStatus.Response.CertUrl != nil {
		_ = d.Set("cert_url", sslStatus.Response.CertUrl)
	}

	return nil
}

func resourceTencentCloudMongodbInstanceSslUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_ssl.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	if d.HasChange("enable") {
		enable := d.Get("enable").(bool)

		err := service.ModifyMongodbInstanceSSL(ctx, instanceId, enable)
		if err != nil {
			return err
		}

		// Wait for SSL configuration to take effect
		err = resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			sslStatus, e := service.DescribeMongodbInstanceSSLById(ctx, instanceId)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if sslStatus == nil || sslStatus.Response == nil || sslStatus.Response.Status == nil {
				return resource.RetryableError(fmt.Errorf("ssl status response is nil"))
			}

			expectedStatus := 0
			if enable {
				expectedStatus = 1
			}

			if int(*sslStatus.Response.Status) == expectedStatus {
				return nil
			}

			return resource.RetryableError(fmt.Errorf("ssl status is %d, waiting for %d", *sslStatus.Response.Status, expectedStatus))
		})

		if err != nil {
			log.Printf("[CRITAL]%s update mongodb instance ssl failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudMongodbInstanceSslRead(d, meta)
}

func resourceTencentCloudMongodbInstanceSslDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_ssl.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	// Disable SSL when deleting the resource
	if err := service.ModifyMongodbInstanceSSL(ctx, instanceId, false); err != nil {
		log.Printf("[CRITAL]%s disable mongodb instance ssl on delete failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
