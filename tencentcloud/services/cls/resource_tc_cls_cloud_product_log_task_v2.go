package cls

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clsv20201016 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsCloudProductLogTaskV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsCloudProductLogTaskV2Create,
		Read:   resourceTencentCloudClsCloudProductLogTaskV2Read,
		Update: resourceTencentCloudClsCloudProductLogTaskV2Update,
		Delete: resourceTencentCloudClsCloudProductLogTaskV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance ID.",
			},

			"assumer_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cloud product identification, Values: CDS, CWP, CDB, TDSQL-C, MongoDB, TDStore, DCDB, MariaDB, PostgreSQL, BH, APIS.",
			},

			"log_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Log type, Values: CDS-AUDIT, CDS-RISK, CDB-AUDIT, TDSQL-C-AUDIT, MongoDB-AUDIT, MongoDB-SlowLog, MongoDB-ErrorLog, TDMYSQL-SLOW, DCDB-AUDIT, DCDB-SLOW, DCDB-ERROR, MariaDB-AUDIT, MariaDB-SLOW, MariaDB-ERROR, PostgreSQL-SLOW, PostgreSQL-ERROR, PostgreSQL-AUDIT, BH-FILELOG, BH-COMMANDLOG, APIS-ACCESS.",
			},

			"cloud_product_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cloud product region. There are differences in the input format of different log types in different regions. Please refer to the following example:\n- CDS(all log type): ap-guangzhou\n- CDB-AUDIT: gz\n- TDSQL-C-AUDIT: gz\n- MongoDB-AUDIT: gz\n- MongoDB-SlowLog: ap-guangzhou\n- MongoDB-ErrorLog: ap-guangzhou\n- TDMYSQL-SLOW: gz\n- DCDB(all log type): gz\n- MariaDB(all log type): gz\n- PostgreSQL(all log type): gz\n- BH(all log type): overseas-polaris(Domestic sites overseas)/fsi-polaris(Domestic sites finance)/general-polaris(Domestic sites)/intl-sg-prod(International sites)\n- APIS(all log type): gz.",
			},

			"cls_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CLS target region.",
			},

			"logset_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Log set name, required if `logset_id` is not filled in. If the log set does not exist, it will be automatically created.",
			},

			"topic_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The name of the log topic is required when `topic_id` is not filled in. If the log theme does not exist, it will be automatically created.",
			},

			"extend": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Log configuration extension information, generally used to store additional log delivery configurations.",
			},

			"logset_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Log set ID.",
			},

			"topic_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Log theme ID.",
			},

			"force_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicate whether to forcibly delete the corresponding logset and topic. If set to true, it will be forcibly deleted. Default is false.",
			},
		},
	}
}

func resourceTencentCloudClsCloudProductLogTaskV2Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_cloud_product_log_task_v2.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId              = tccommon.GetLogId(tccommon.ContextNil)
		ctx                = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request            = clsv20201016.NewCreateCloudProductLogCollectionRequest()
		instanceId         string
		assumerName        string
		logType            string
		cloudProductRegion string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("assumer_name"); ok {
		request.AssumerName = helper.String(v.(string))
		assumerName = v.(string)
	}

	if v, ok := d.GetOk("log_type"); ok {
		request.LogType = helper.String(v.(string))
		logType = v.(string)
	}

	if v, ok := d.GetOk("cloud_product_region"); ok {
		request.CloudProductRegion = helper.String(v.(string))
		cloudProductRegion = v.(string)
	}

	if v, ok := d.GetOk("cls_region"); ok {
		request.ClsRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("logset_name"); ok {
		request.LogsetName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topic_name"); ok {
		request.TopicName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("extend"); ok {
		request.Extend = helper.String(v.(string))
	}

	if v, ok := d.GetOk("logset_id"); ok {
		request.LogsetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topic_id"); ok {
		request.TopicId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsV20201016Client().CreateCloudProductLogCollectionWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create cls cloud product log task failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cls cloud product log task failed, reason:%+v", logId, err)
		return err
	}

	// wait
	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"1"}, 10*tccommon.ReadRetryTimeout, time.Second, service.ClsCloudProductLogTaskStateRefreshFunc(ctx, instanceId, assumerName, logType, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	d.SetId(strings.Join([]string{instanceId, assumerName, logType, cloudProductRegion}, tccommon.FILED_SP))

	return resourceTencentCloudClsCloudProductLogTaskV2Read(d, meta)
}

func resourceTencentCloudClsCloudProductLogTaskV2Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_cloud_product_log_task_v2.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service     = ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		deleteForce bool
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	assumerName := idSplit[1]
	logType := idSplit[2]
	cloudProductRegion := idSplit[3]

	respData, err := service.DescribeClsCloudProductLogTaskById(ctx, instanceId, assumerName, logType)
	if err != nil {
		return err
	}

	if respData == nil || len(respData.Tasks) < 1 {
		log.Printf("[WARN]%s resource `tencentcloud_cls_cloud_product_log_task_v2` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("assumer_name", assumerName)
	_ = d.Set("log_type", logType)
	_ = d.Set("cloud_product_region", cloudProductRegion)

	if respData.Tasks[0].ClsRegion != nil {
		_ = d.Set("cls_region", respData.Tasks[0].ClsRegion)
	}

	if respData.Tasks[0].Extend != nil {
		_ = d.Set("extend", respData.Tasks[0].Extend)
	}

	if respData.Tasks[0].LogsetId != nil {
		_ = d.Set("logset_id", respData.Tasks[0].LogsetId)
		info, err := service.DescribeClsLogset(ctx, *respData.Tasks[0].LogsetId)
		if err != nil {
			return err
		}

		if info != nil {
			if info.LogsetName != nil {
				_ = d.Set("logset_name", info.LogsetName)
			}
		}
	}

	if respData.Tasks[0].TopicId != nil {
		_ = d.Set("topic_id", respData.Tasks[0].TopicId)
		info, err := service.DescribeClsTopicById(ctx, *respData.Tasks[0].TopicId)
		if err != nil {
			return err
		}

		if info != nil {
			if info.TopicName != nil {
				_ = d.Set("topic_name", info.TopicName)
			}
		}
	}

	if v, ok := d.GetOkExists("force_delete"); ok {
		deleteForce = v.(bool)
	}

	_ = d.Set("force_delete", deleteForce)

	return nil
}

func resourceTencentCloudClsCloudProductLogTaskV2Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_cloud_product_log_task_v2.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	immutableArgs := []string{"cls_region"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	assumerName := idSplit[1]
	logType := idSplit[2]
	cloudProductRegion := idSplit[3]

	needChange := false
	mutableArgs := []string{"extend"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := clsv20201016.NewModifyCloudProductLogCollectionRequest()
		request.InstanceId = helper.String(instanceId)
		request.AssumerName = helper.String(assumerName)
		request.LogType = helper.String(logType)
		request.CloudProductRegion = helper.String(cloudProductRegion)
		if v, ok := d.GetOk("extend"); ok {
			request.Extend = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsV20201016Client().ModifyCloudProductLogCollectionWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update cls cloud product log task failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudClsCloudProductLogTaskV2Read(d, meta)
}

func resourceTencentCloudClsCloudProductLogTaskV2Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_cloud_product_log_task_v2.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request     = clsv20201016.NewDeleteCloudProductLogCollectionRequest()
		deleteForce bool
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	assumerName := idSplit[1]
	logType := idSplit[2]
	cloudProductRegion := idSplit[3]

	request.InstanceId = helper.String(instanceId)
	request.AssumerName = helper.String(assumerName)
	request.LogType = helper.String(logType)
	request.CloudProductRegion = helper.String(cloudProductRegion)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsV20201016Client().DeleteCloudProductLogCollectionWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete cls cloud product log task failed, reason:%+v", logId, err)
		return err
	}

	// wait delete
	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"3"}, 10*tccommon.ReadRetryTimeout, time.Second, service.ClsCloudProductLogTaskStateRefreshFunc(ctx, instanceId, assumerName, logType, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	if v, ok := d.GetOkExists("force_delete"); ok {
		deleteForce = v.(bool)
	}

	if deleteForce {
		var (
			request1 = clsv20201016.NewDeleteTopicRequest()
			request2 = clsv20201016.NewDeleteLogsetRequest()
		)

		if v, ok := d.GetOk("topic_id"); ok {
			request1.TopicId = helper.String(v.(string))
		}

		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsV20201016Client().DeleteTopicWithContext(ctx, request1)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request1.GetAction(), request1.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s delete cls cloud product log task topic failed, reason:%+v", logId, err)
			return err
		}

		if v, ok := d.GetOk("logset_id"); ok {
			request2.LogsetId = helper.String(v.(string))
		}

		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsV20201016Client().DeleteLogsetWithContext(ctx, request2)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request2.GetAction(), request2.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s delete cls cloud product log task logset failed, reason:%+v", logId, err)
			return err
		}
	}

	return nil
}
