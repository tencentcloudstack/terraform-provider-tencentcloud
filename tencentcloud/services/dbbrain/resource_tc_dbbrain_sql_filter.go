package dbbrain

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDbbrainSqlFilter() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudDbbrainSqlFilterRead,
		Create: resourceTencentCloudDbbrainSqlFilterCreate,
		Update: resourceTencentCloudDbbrainSqlFilterUpdate,
		Delete: resourceTencentCloudDbbrainSqlFilterDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance id.",
			},

			"session_token": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "session token.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "user name.",
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "password.",
						},
					},
				},
			},

			"sql_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "sql type, optional value is SELECT, UPDATE, DELETE, INSERT, REPLACE.",
			},

			"filter_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "filter key.",
			},

			"max_concurrency": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "maximum concurreny.",
			},

			"duration": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "filter duration.",
			},

			"product": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "product, optional value is &amp;#39;mysql&amp;#39;, &amp;#39;cynosdb&amp;#39;.",
			},

			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "filter status.",
			},

			"filter_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "filter id.",
			},
		},
	}
}
func resourceTencentCloudDbbrainSqlFilterCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbbrain_sql_filter.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request    = dbbrain.NewCreateSqlFilterRequest()
		response   *dbbrain.CreateSqlFilterResponse
		instanceId string
		filterId   string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sql_type"); ok {
		request.SqlType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filter_key"); ok {
		request.FilterKey = helper.String(v.(string))
	}

	if v, ok := d.GetOk("max_concurrency"); ok {
		request.MaxConcurrency = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("duration"); ok {
		request.Duration = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("product"); ok {
		request.Product = helper.String(v.(string))
	}

	sessionToken, err := getSessionToken(d, meta, ctx)
	if err != nil {
		return err
	}
	if sessionToken != nil {
		request.SessionToken = sessionToken
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDbbrainClient().CreateSqlFilter(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dbbrain sqlFilter failed, reason:%+v", logId, err)
		return err
	}

	filterId = helper.Int64ToStr(*response.Response.FilterId)

	d.SetId(instanceId + tccommon.FILED_SP + filterId)
	return resourceTencentCloudDbbrainSqlFilterRead(d, meta)
}

func getSessionToken(d *schema.ResourceData, meta interface{}, ctx context.Context) (sessionToken *string, errRet error) {
	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		instanceId *string
		product    *string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		product = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "session_token"); ok {
		var user *string
		var pw *string
		if v, ok := dMap["user"]; ok {
			user = helper.String(v.(string))
		}
		if v, ok := dMap["password"]; ok {
			pw = helper.String(v.(string))
		}

		service := DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		sessionToken, errRet = service.getSessionToken(ctx, instanceId, user, pw, product)

		if errRet != nil {
			return
		}

		log.Printf("[DEBUG]%s verify user account success, sessionToken [%s]\n", logId, *sessionToken)
	}
	return
}

func resourceTencentCloudDbbrainSqlFilterRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbbrain_sql_filter.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		idSplit []string
	)

	idSplit = strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	filterId := idSplit[1]

	sqlFilter, err := service.DescribeDbbrainSqlFilter(ctx, helper.String(instanceId), helper.String(filterId))
	if err != nil {
		return err
	}

	if sqlFilter == nil {
		d.SetId("")
		return fmt.Errorf("resource `sqlFilter` %s does not exist", d.Id())
	}

	_ = d.Set("instance_id", instanceId)

	if sqlFilter.SqlType != nil {
		_ = d.Set("sql_type", sqlFilter.SqlType)
	}

	if sqlFilter.OriginKeys != nil {
		_ = d.Set("filter_key", sqlFilter.OriginKeys)
	}

	if sqlFilter.MaxConcurrency != nil {
		_ = d.Set("max_concurrency", sqlFilter.MaxConcurrency)
	}

	if sqlFilter.Status != nil {
		_ = d.Set("status", sqlFilter.Status)
	}

	if sqlFilter.Id != nil {
		_ = d.Set("filter_id", sqlFilter.Id)
	}

	return nil
}

func resourceTencentCloudDbbrainSqlFilterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbbrain_sql_filter.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request = dbbrain.NewModifySqlFiltersRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := helper.String(idSplit[0])
	filterId := idSplit[1]
	request.InstanceId = instanceId
	request.FilterIds = []*int64{helper.StrToInt64Point(filterId)}

	if d.HasChange("instance_id") {
		return fmt.Errorf("`instance_id` do not support change now.")
	}

	if d.HasChange("sql_type") {
		return fmt.Errorf("`sql_type` do not support change now.")
	}

	if d.HasChange("filter_key") {
		return fmt.Errorf("`filter_key` do not support change now.")
	}

	if d.HasChange("max_concurrency") {
		return fmt.Errorf("`max_concurrency` do not support change now.")
	}

	if d.HasChange("duration") {
		return fmt.Errorf("`duration` do not support change now.")
	}

	if d.HasChange("status") {
		if v, ok := d.GetOk("status"); ok {
			request.Status = helper.String(v.(string))
		}
	}

	if d.HasChange("product") {
		if v, ok := d.GetOk("product"); ok {
			request.Product = helper.String(v.(string))
		}
	}

	sessionToken, err := getSessionToken(d, meta, ctx)
	if err != nil {
		return err
	}
	if sessionToken != nil {
		request.SessionToken = sessionToken
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDbbrainClient().ModifySqlFilters(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s modify dbbrain sqlFilter failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDbbrainSqlFilterRead(d, meta)
}

func resourceTencentCloudDbbrainSqlFilterDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbbrain_sql_filter.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := helper.String(idSplit[0])
	filterId := helper.String(idSplit[1])
	sessionToken, err := getSessionToken(d, meta, ctx)
	if err != nil {
		return err
	}

	if err := service.DeleteDbbrainSqlFilterById(ctx, instanceId, filterId, sessionToken); err != nil {
		return err
	}

	return nil
}
