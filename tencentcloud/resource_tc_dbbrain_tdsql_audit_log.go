/*
Provides a resource to create a dbbrain tdsql_audit_log

Example Usage

```hcl
resource "tencentcloud_dbbrain_tdsql_audit_log" "tdsql_audit_log" {
  product = ""
  node_request_type = ""
  instance_id = ""
  start_time = ""
  end_time = ""
  filter {
		host = 
		d_b_name = 
		user = 
		sent_rows = 
		affect_rows = 
		exec_time = 

  }
}
```

Import

dbbrain tdsql_audit_log can be imported using the id, e.g.

```
terraform import tencentcloud_dbbrain_tdsql_audit_log.tdsql_audit_log tdsql_audit_log_id
```
*/
package tencentcloud

import (
"context"
"fmt"
"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
"log"
"strings"
)


func resourceTencentCloudDbbrainTdsqlAuditLog () * schema.Resource {
return & schema.Resource {
Create : resourceTencentCloudDbbrainTdsqlAuditLogCreate ,
Read : resourceTencentCloudDbbrainTdsqlAuditLogRead ,
Update : resourceTencentCloudDbbrainTdsqlAuditLogUpdate ,
Delete : resourceTencentCloudDbbrainTdsqlAuditLogDelete ,
Importer : & schema.ResourceImporter {
State : schema.ImportStatePassthrough ,
} ,
Schema : map[string] * schema.Schema {
"product": {
  Required: true,
  Type: schema.TypeString,
  Description: "���务产品类型，支持值包括： dcdb - 云数据库 Tdsql， mariadb - 云数据库 MariaDB for MariaDB。.",
},

"node_request_type": {
  Required: true,
  Type: schema.TypeString,
  Description: "���Product保持一致。如：dcdb ,mariadb.",
},

"instance_id": {
  Required: true,
  Type: schema.TypeString,
  Description: "���例 ID 。.",
},

"start_time": {
  Required: true,
  Type: schema.TypeString,
  Description: "���始时间，如“2019-09-10 12:13:14”。	.",
},

"end_time": {
  Required: true,
  Type: schema.TypeString,
  Description: "���止时间，如“2019-09-11 10:13:14”。.",
},

"filter": {
  Optional: true,
  Type: schema.TypeList,
MaxItems: 1,
  Description: "���滤条件。可按设置的过滤条件过滤日志。.",
Elem: &schema.Resource{
    Schema: map[string]*schema.Schema{
      "host": {
  Type: schema.TypeSet,
  Elem: &schema.Schema{
				Type: schema.TypeString,
	},
  Optional: true,
  Description: "���户端地址。.",
},
"d_b_name": {
  Type: schema.TypeSet,
  Elem: &schema.Schema{
				Type: schema.TypeString,
	},
  Optional: true,
  Description: "���据库名称。.",
},
"user": {
  Type: schema.TypeSet,
  Elem: &schema.Schema{
				Type: schema.TypeString,
	},
  Optional: true,
  Description: "���户名。.",
},
"sent_rows": {
  Type: schema.TypeInt,
  Optional: true,
  Description: "���回行数。表示筛选返回行数大于该值的审计日志。.",
},
"affect_rows": {
  Type: schema.TypeInt,
  Optional: true,
  Description: "���响行数。表示筛选影响行数大于该值的审计日志。.",
},
"exec_time": {
  Type: schema.TypeInt,
  Optional: true,
  Description: "���行时间。单位为：µs。表示筛选执行时间大于该值的审计日志。.",
},

    },
  },
},

} ,
}
} 

func resourceTencentCloudDbbrainTdsqlAuditLogCreate (d * schema.ResourceData,meta interface{}) error {
defer logElapsed ("resource.tencentcloud_dbbrain_tdsql_audit_log.create") ()
defer inconsistentCheck (d,meta) () 

logId := getLogId (contextNil) 

var (
request = dbbrain.NewCreateAuditLogFileRequest ()
response = dbbrain.NewCreateAuditLogFileResponse ()
asyncRequestId int
instanceId int
product string
)
if v,ok := d . GetOk ("product");ok {
product = v .(string) 
 request . Product = helper.String (v .(string))
} 

if v,ok := d . GetOk ("node_request_type");ok {
request . NodeRequestType = helper.String (v .(string))
} 

if v,ok := d . GetOk ("instance_id");ok {
instanceId = v .(string) 
 request . InstanceId = helper.String (v .(string))
} 

if v,ok := d . GetOk ("start_time");ok {
request . StartTime = helper.String (v .(string))
} 

if v,ok := d . GetOk ("end_time");ok {
request . EndTime = helper.String (v .(string))
} 

if dMap,ok := helper . InterfacesHeadMap (d,"filter");ok {
auditLogFilter := dbbrain . AuditLogFilter {}
if v,ok := dMap ["host"];ok {
hostSet := v .(* schema.Set) . List () 
 for i := range hostSet {
host := hostSet [i] .(string)
auditLogFilter . Host = append(auditLogFilter . Host,& host)
}
}
if v,ok := dMap ["d_b_name"];ok {
dBNameSet := v .(* schema.Set) . List () 
 for i := range dBNameSet {
dBName := dBNameSet [i] .(string)
auditLogFilter . DBName = append(auditLogFilter . DBName,& dBName)
}
}
if v,ok := dMap ["user"];ok {
userSet := v .(* schema.Set) . List () 
 for i := range userSet {
user := userSet [i] .(string)
auditLogFilter . User = append(auditLogFilter . User,& user)
}
}
if v,ok := dMap ["sent_rows"];ok {
auditLogFilter . SentRows = helper.IntInt64 (v .(int))
}
if v,ok := dMap ["affect_rows"];ok {
auditLogFilter . AffectRows = helper.IntInt64 (v .(int))
}
if v,ok := dMap ["exec_time"];ok {
auditLogFilter . ExecTime = helper.IntInt64 (v .(int))
}
request . Filter = & auditLogFilter
} 

err := resource.Retry (writeRetryTimeout,func () * resource.RetryError {
result,e := meta .(* TencentCloudClient) . apiV3Conn . UseDbbrainClient () . CreateAuditLogFile (request)
if e != nil {
return  retryError (e)
} else {
log.Printf ("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",logId,request . GetAction (),request . ToJsonString (),result . ToJsonString ())
}
response = result
return  nil
})
if err != nil {
log.Printf ("[CRITAL]%s create dbbrain tdsqlAuditLog failed, reason:%+v",logId,err)
return  err
} 

asyncRequestId = * response . Response . AsyncRequestId
d . SetId (strings.Join ([] string {helper.Int64ToStr (asyncRequestId),helper.Int64ToStr (instanceId),product},FILED_SP)) 

return resourceTencentCloudDbbrainTdsqlAuditLogRead (d,meta)
} 

func resourceTencentCloudDbbrainTdsqlAuditLogRead (d * schema.ResourceData,meta interface{}) error {
defer logElapsed ("resource.tencentcloud_dbbrain_tdsql_audit_log.read") ()
defer inconsistentCheck (d,meta) () 

logId := getLogId (contextNil) 

ctx := context.WithValue (context.TODO (),logIdKey,logId) 

service := DbbrainService {client:meta .(* TencentCloudClient) . apiV3Conn} 

idSplit := strings.Split (d . Id (),FILED_SP)
if len (idSplit) != 3 {
return fmt . Errorf ("id is broken,%s",d . Id ())
}
asyncRequestId := idSplit [0]
instanceId := idSplit [1]
product := idSplit [2]


tdsqlAuditLog,err := service . DescribeDbbrainTdsqlAuditLogById (ctx,asyncRequestId,instanceId,product)
if err != nil {
return err
} 

if tdsqlAuditLog == nil {
d . SetId ("")
log.Printf ("[WARN]%s resource `DbbrainTdsqlAuditLog` [%s] not found, please check if it has been deleted.\n",logId,d . Id ())
return nil
} 

if tdsqlAuditLog . Product != nil {
_ = d . Set ("product",tdsqlAuditLog . Product)
} 

if tdsqlAuditLog . NodeRequestType != nil {
_ = d . Set ("node_request_type",tdsqlAuditLog . NodeRequestType)
} 

if tdsqlAuditLog . InstanceId != nil {
_ = d . Set ("instance_id",tdsqlAuditLog . InstanceId)
} 

if tdsqlAuditLog . StartTime != nil {
_ = d . Set ("start_time",tdsqlAuditLog . StartTime)
} 

if tdsqlAuditLog . EndTime != nil {
_ = d . Set ("end_time",tdsqlAuditLog . EndTime)
} 

if tdsqlAuditLog . Filter != nil {
filterMap := map[string] interface{} {} 

if tdsqlAuditLog.Filter . Host != nil {
filterMap ["host"] = tdsqlAuditLog.Filter . Host
} 

if tdsqlAuditLog.Filter . DBName != nil {
filterMap ["d_b_name"] = tdsqlAuditLog.Filter . DBName
} 

if tdsqlAuditLog.Filter . User != nil {
filterMap ["user"] = tdsqlAuditLog.Filter . User
} 

if tdsqlAuditLog.Filter . SentRows != nil {
filterMap ["sent_rows"] = tdsqlAuditLog.Filter . SentRows
} 

if tdsqlAuditLog.Filter . AffectRows != nil {
filterMap ["affect_rows"] = tdsqlAuditLog.Filter . AffectRows
} 

if tdsqlAuditLog.Filter . ExecTime != nil {
filterMap ["exec_time"] = tdsqlAuditLog.Filter . ExecTime
} 

_ = d . Set ("filter",[] interface{} {filterMap})
} 

return nil
} 

func resourceTencentCloudDbbrainTdsqlAuditLogUpdate (d * schema.ResourceData,meta interface{}) error {
defer logElapsed ("resource.tencentcloud_dbbrain_tdsql_audit_log.update") ()
defer inconsistentCheck (d,meta) () 

logId := getLogId (contextNil) 

immutableArgs := [] string {"product","node_request_type","instance_id","start_time","end_time","filter"}


for _,v := range immutableArgs {
if d . HasChange (v) {
return fmt.Errorf ("argument `%s` cannot be changed",v)
}
}
return resourceTencentCloudDbbrainTdsqlAuditLogRead (d,meta)
} 

func resourceTencentCloudDbbrainTdsqlAuditLogDelete (d * schema.ResourceData,meta interface{}) error {
defer logElapsed ("resource.tencentcloud_dbbrain_tdsql_audit_log.delete") ()
defer inconsistentCheck (d,meta) () 

logId := getLogId (contextNil)
ctx := context.WithValue (context.TODO (),logIdKey,logId) 

service := DbbrainService {client:meta .(* TencentCloudClient) . apiV3Conn}
idSplit := strings.Split (d . Id (),FILED_SP)
if len (idSplit) != 3 {
return fmt . Errorf ("id is broken,%s",d . Id ())
}
asyncRequestId := idSplit [0]
instanceId := idSplit [1]
product := idSplit [2]


if err := service . DeleteDbbrainTdsqlAuditLogById (ctx,asyncRequestId,instanceId,product) ; err != nil {
return err
} 

return nil
} 
