package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTencentCloudClickhouseAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClickhouseAccountCreate,
		Read:   resourceTencentCloudClickhouseAccountRead,
		Update: resourceTencentCloudClickhouseAccountUpdate,
		Delete: resourceTencentCloudClickhouseAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance id.",
			},
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User name.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password.",
			},
			"describe": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description.",
			},
		},
	}
}

func resourceTencentCloudClickhouseAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clickhouse_account.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	instanceId := d.Get("instance_id").(string)
	userName := d.Get("user_name").(string)
	params := make(map[string]interface{})
	params["instance_id"] = instanceId
	params["user_name"] = userName
	params["password"] = d.Get("password").(string)
	if v, ok := d.GetOk("describe"); ok {
		params["describe"] = v.(string)
	}
	service := CdwchService{client: meta.(*TencentCloudClient).apiV3Conn}
	err := service.ActionAlterCkUser(ctx, ACTION_ALTER_CK_USER_ADD_SYSTEM_USER, params)
	if err != nil {
		return err
	}
	d.SetId(instanceId + FILED_SP + userName)

	return resourceTencentCloudClickhouseAccountRead(d, meta)
}

func resourceTencentCloudClickhouseAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clickhouse_account.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := CdwchService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("tencentcloud_clickhouse_account id is broken, id is %s", d.Id())
	}
	accounts, err := service.DescribeClickhouseAccountByUserName(ctx, idSplit[0], idSplit[1])
	if err != nil {
		return err
	}

	if len(accounts) == 0 {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClickhouseAccount` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	account := accounts[0]
	_ = d.Set("instance_id", account.InstanceId)
	_ = d.Set("user_name", account.UserName)
	_ = d.Set("describe", account.Describe)

	return nil
}

func resourceTencentCloudClickhouseAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clickhouse_account.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("tencentcloud_clickhouse_account id is broken, id is %s", d.Id())
	}

	immutableArgs := []string{"instance_id", "user_name", "describe"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("password") {
		params := make(map[string]interface{})
		params["instance_id"] = idSplit[0]
		params["user_name"] = idSplit[1]
		params["password"] = d.Get("password").(string)
		if v, ok := d.GetOk("describe"); ok {
			params["describe"] = v.(string)
		}
		service := CdwchService{client: meta.(*TencentCloudClient).apiV3Conn}
		err := service.ActionAlterCkUser(ctx, ACTION_ALTER_CK_USER_UPDATE_SYSTEM_USER, params)
		if err != nil {
			return err
		}
	}
	return resourceTencentCloudClickhouseAccountRead(d, meta)
}

func resourceTencentCloudClickhouseAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clickhouse_account.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := CdwchService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("tencentcloud_clickhouse_account id is broken, id is %s", d.Id())
	}

	err := service.DescribeCkSqlApis(ctx, idSplit[0], "", idSplit[1], DESCRIBE_CK_SQL_APIS_DELETE_SYSTEM_USER)
	if err != nil {
		return err
	}
	return nil
}
