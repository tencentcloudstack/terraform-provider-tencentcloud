package mongodb

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudMongodbInstanceSrvConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMongodbInstanceSrvConnectionCreate,
		Read:   resourceTencentCloudMongodbInstanceSrvConnectionRead,
		Update: resourceTencentCloudMongodbInstanceSrvConnectionUpdate,
		Delete: resourceTencentCloudMongodbInstanceSrvConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "MongoDB instance ID, for example: cmgo-p8vnipr5.",
			},

			"domain": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Custom domain for SRV connection. If not specified during creation, the system will use a default domain. After creation, this field will be populated with the actual domain. To set or modify a custom domain, use this field.",
			},
		},
	}
}

func resourceTencentCloudMongodbInstanceSrvConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_srv_connection.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Get("instance_id").(string)

	service := MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	// Step 1: Enable SRV connection (without custom domain)
	err := service.EnableSRVConnectionUrl(ctx, instanceId)
	if err != nil {
		log.Printf("[CRITAL]%s enable srv connection url failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	// Step 2: If custom domain is specified, modify it
	if v, ok := d.GetOk("domain"); ok {
		domain := v.(string)
		err = service.ModifySRVConnectionUrl(ctx, instanceId, domain)
		if err != nil {
			log.Printf("[CRITAL]%s modify srv connection url failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudMongodbInstanceSrvConnectionRead(d, meta)
}

func resourceTencentCloudMongodbInstanceSrvConnectionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_srv_connection.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Id()

	service := MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	domain, err := service.DescribeSRVConnectionDomain(ctx, instanceId)
	if err != nil {
		return err
	}

	_ = d.Set("instance_id", instanceId)

	if domain != nil {
		_ = d.Set("domain", domain)
	}

	return nil
}

func resourceTencentCloudMongodbInstanceSrvConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_srv_connection.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Id()

	if d.HasChange("domain") {
		service := MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

		domain := d.Get("domain").(string)

		err := service.ModifySRVConnectionUrl(ctx, instanceId, domain)
		if err != nil {
			log.Printf("[CRITAL]%s modify srv connection url failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudMongodbInstanceSrvConnectionRead(d, meta)
}

func resourceTencentCloudMongodbInstanceSrvConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_srv_connection.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Id()

	service := MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := service.DisableSRVConnectionUrl(ctx, instanceId)
	if err != nil {
		log.Printf("[CRITAL]%s disable srv connection url failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
