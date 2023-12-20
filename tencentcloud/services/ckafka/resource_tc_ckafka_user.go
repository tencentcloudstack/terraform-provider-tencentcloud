package ckafka

import (
	"context"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudCkafkaUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCkafkaUserCreate,
		Read:   resourceTencentCloudCkafkaUserRead,
		Update: resourceTencentCloudCkafkaUserUpdate,
		Delete: resourceTencentCloudCkafkaUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the ckafka instance.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Account name used to access to ckafka instance.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password of the account.",
			},
			// computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the account.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last update time of the account.",
			},
		},
	}
}

func resourceTencentCloudCkafkaUserCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_user.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Get("instance_id").(string)
	accountName := d.Get("account_name").(string)
	password := d.Get("password").(string)

	ckafkaService := CkafkaService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	if err := ckafkaService.CreateUser(ctx, instanceId, accountName, password); err != nil {
		return fmt.Errorf("[CRITAL]%s create ckafka user failed, reason:%+v", logId, err)
	}
	d.SetId(instanceId + tccommon.FILED_SP + accountName)

	return resourceTencentCloudCkafkaUserRead(d, meta)
}

func resourceTencentCloudCkafkaUserRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_user.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	ckafkaService := CkafkaService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	id := d.Id()
	info, has, err := ckafkaService.DescribeUserByUserId(ctx, id)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}
	items := strings.Split(id, tccommon.FILED_SP)
	_ = d.Set("instance_id", items[0])
	_ = d.Set("account_name", info.Name)
	_ = d.Set("create_time", info.CreateTime)
	_ = d.Set("update_time", info.UpdateTime)

	return nil
}

func resourceTencentCloudCkafkaUserUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_user.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	ckafkaService := CkafkaService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	instanceId := d.Get("instance_id").(string)
	user := d.Get("account_name").(string)
	if d.HasChange("password") {
		old, new := d.GetChange("password")
		if err := ckafkaService.ModifyPassword(ctx, instanceId, user, old.(string), new.(string)); err != nil {
			return err
		}

	}

	return resourceTencentCloudCkafkaUserRead(d, meta)
}

func resourceTencentCloudCkafkaUserDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_user.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	ckafkaService := CkafkaService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	if err := ckafkaService.DeleteUser(ctx, d.Id()); err != nil {
		return err
	}
	return nil
}
