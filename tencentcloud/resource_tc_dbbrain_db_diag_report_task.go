package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDbbrainDbDiagReportTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDbbrainDbDiagReportTaskCreate,
		Read:   resourceTencentCloudDbbrainDbDiagReportTaskRead,
		Delete: resourceTencentCloudDbbrainDbDiagReportTaskDelete,
		// contact_group, contact_person, send_mail_flag and product fileds can not query by read api
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"product": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values include: mysql - cloud database MySQL, cynosdb - cloud database CynosDB for MySQL.",
			},

			"start_time": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Start time, such as 2020-11-08T14:00:00+08:00.",
			},

			"end_time": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "End time, such as 2020-11-09T14:00:00+08:00.",
			},

			"send_mail_flag": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Whether to send mail: 0 - no, 1 - yes.",
			},

			"contact_person": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "An array of contact IDs to receive emails from.",
			},

			"contact_group": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "An array of contact group IDs to receive mail from.",
			},
		},
	}
}

func resourceTencentCloudDbbrainDbDiagReportTaskCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dbbrain_db_diag_report_task.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = dbbrain.NewCreateDBDiagReportTaskRequest()
		response   = dbbrain.NewCreateDBDiagReportTaskResponse()
		instanceId string
		product    string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		request.StartTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		request.EndTime = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("send_mail_flag"); ok {
		request.SendMailFlag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("contact_person"); ok {
		contactPersonSet := v.(*schema.Set).List()
		for i := range contactPersonSet {
			contactPerson := contactPersonSet[i].(int)
			request.ContactPerson = append(request.ContactPerson, helper.IntInt64(contactPerson))
		}
	}

	if v, ok := d.GetOk("contact_group"); ok {
		contactGroupSet := v.(*schema.Set).List()
		for i := range contactGroupSet {
			contactGroup := contactGroupSet[i].(int)
			request.ContactGroup = append(request.ContactGroup, helper.IntInt64(contactGroup))
		}
	}

	if v, ok := d.GetOk("product"); ok {
		product = v.(string)
		request.Product = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDbbrainClient().CreateDBDiagReportTask(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dbbrain dbDiagReportTask failed, reason:%+v", logId, err)
		return err
	}
	if response == nil || response.Response.AsyncRequestId == nil {
		return fmt.Errorf("[CRITAL]%s The dbbrain dbDiagReportTask id not found after creation", logId)
	}

	asyncRequestId := response.Response.AsyncRequestId
	d.SetId(helper.Int64ToStr(*asyncRequestId) + FILED_SP + instanceId + FILED_SP + product)

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"100"}, 3*readRetryTimeout, time.Second, service.DbbrainDbDiagReportTaskStateRefreshFunc(asyncRequestId, instanceId, product, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudDbbrainDbDiagReportTaskRead(d, meta)
}

func resourceTencentCloudDbbrainDbDiagReportTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dbbrain_db_diag_report_task.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	asyncRequestId := idSplit[0]
	instanceId := idSplit[1]
	product := idSplit[2]

	dbDiagReportTask, err := service.DescribeDbbrainDbDiagReportTaskById(ctx, helper.StrToInt64Point(asyncRequestId), instanceId, product)
	if err != nil {
		return err
	}

	if dbDiagReportTask == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DbbrainDbDiagReportTask` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if dbDiagReportTask.StartTime != nil {
		_ = d.Set("start_time", dbDiagReportTask.StartTime)
	}

	if dbDiagReportTask.EndTime != nil {
		_ = d.Set("end_time", dbDiagReportTask.EndTime)
	}

	// if dbDiagReportTask.SendMailFlag != nil {
	// 	_ = d.Set("send_mail_flag", dbDiagReportTask.SendMailFlag)
	// }

	// if dbDiagReportTask.ContactPerson != nil {
	// 	_ = d.Set("contact_person", dbDiagReportTask.ContactPerson)
	// }

	// if dbDiagReportTask.ContactGroup != nil {
	// 	_ = d.Set("contact_group", dbDiagReportTask.ContactGroup)
	// }

	// if dbDiagReportTask.Product != nil {
	// 	_ = d.Set("product", dbDiagReportTask.Product)
	// }

	return nil
}

func resourceTencentCloudDbbrainDbDiagReportTaskDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dbbrain_db_diag_report_task.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	asyncRequestId := idSplit[0]
	instanceId := idSplit[1]
	product := idSplit[2]

	if err := service.DeleteDbbrainDbDiagReportTaskById(ctx, *helper.StrToInt64Point(asyncRequestId), instanceId, product); err != nil {
		return err
	}

	return nil
}
