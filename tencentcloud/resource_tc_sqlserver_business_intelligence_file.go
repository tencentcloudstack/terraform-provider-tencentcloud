/*
Provides a resource to create a sqlserver business_intelligence_file

Example Usage

```hcl
resource "tencentcloud_sqlserver_business_intelligence_instance" "business_intelligence_instance" {
  zone = "ap-guangzhou-6"
  memory = 4
  storage = 20
  cpu = 2
  machine_type = "CLOUD_PREMIUM"
  project_id = 0
  subnet_id = "subnet-dwj7ipnc"
  vpc_id = "vpc-4owdpnwr"
  db_version = "201603"
  security_group_list = []
  weekly = [1, 2, 3, 4, 5, 6, 7]
  start_time = "00:00"
  span = 6
  instance_name = "create_db_name"
}

resource "tencentcloud_sqlserver_business_intelligence_file" "business_intelligence_file" {
  instance_id = tencentcloud_sqlserver_business_intelligence_instance.business_intelligence_instance.id
  file_url = "https://keep-sqlserver-1308919341.cos.ap-guangzhou.myqcloud.com/test.xlsx"
  file_type = "FLAT"
  remark = "test case."
}
```

Import

sqlserver business_intelligence_file can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_business_intelligence_file.business_intelligence_file business_intelligence_file_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSqlserverBusinessIntelligenceFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverBusinessIntelligenceFileCreate,
		Read:   resourceTencentCloudSqlserverBusinessIntelligenceFileRead,
		Delete: resourceTencentCloudSqlserverBusinessIntelligenceFileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},
			"file_url": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cos Url.",
			},
			"file_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "File Type FLAT - Flat File as Data Source, SSIS - ssis project package.",
			},
			"remark": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "remark.",
			},
		},
	}
}

func resourceTencentCloudSqlserverBusinessIntelligenceFileCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_business_intelligence_file.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId             = getLogId(contextNil)
		CreateBIFRequest  = sqlserver.NewCreateBusinessIntelligenceFileRequest()
		CreateBIFResponse = sqlserver.NewCreateBusinessIntelligenceFileResponse()
		instanceId        string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		CreateBIFRequest.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("file_url"); ok {
		CreateBIFRequest.FileURL = helper.String(v.(string))
	}

	if v, ok := d.GetOk("file_type"); ok {
		CreateBIFRequest.FileType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		CreateBIFRequest.Remark = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().CreateBusinessIntelligenceFile(CreateBIFRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, CreateBIFRequest.GetAction(), CreateBIFRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("sqlserver CreateBusinessIntelligenceFile not exists")
			return resource.NonRetryableError(e)
		}

		CreateBIFResponse = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver businessIntelligenceFile failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{instanceId, *CreateBIFResponse.Response.FileTaskId}, FILED_SP))

	return resourceTencentCloudSqlserverBusinessIntelligenceFileRead(d, meta)
}

func resourceTencentCloudSqlserverBusinessIntelligenceFileRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_business_intelligence_file.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	instanceId := idSplit[0]
	fileName := idSplit[1]

	businessIntelligenceFile, err := service.DescribeSqlserverBusinessIntelligenceFileById(ctx, instanceId, fileName)
	if err != nil {
		return err
	}

	if businessIntelligenceFile == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverBusinessIntelligenceFile` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if businessIntelligenceFile.InstanceId != nil {
		_ = d.Set("instance_id", businessIntelligenceFile.InstanceId)
	}

	if businessIntelligenceFile.FileURL != nil {
		_ = d.Set("file_url", businessIntelligenceFile.FileURL)
	}

	if businessIntelligenceFile.FileType != nil {
		_ = d.Set("file_type", businessIntelligenceFile.FileType)
	}

	if businessIntelligenceFile.Remark != nil {
		_ = d.Set("remark", businessIntelligenceFile.Remark)
	}

	return nil
}

func resourceTencentCloudSqlserverBusinessIntelligenceFileDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_business_intelligence_file.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	instanceId := idSplit[0]
	fileName := idSplit[1]

	if err := service.DeleteSqlserverBusinessIntelligenceFileById(ctx, instanceId, fileName); err != nil {
		return err
	}

	return nil
}
