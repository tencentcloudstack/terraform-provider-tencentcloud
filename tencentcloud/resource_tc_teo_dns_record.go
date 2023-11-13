/*
Provides a resource to create a teo dns_record

Example Usage

```hcl
resource "tencentcloud_teo_dns_record" "dns_record" {
  zone_id = &lt;nil&gt;
    type = &lt;nil&gt;
  name = &lt;nil&gt;
  content = &lt;nil&gt;
  mode = &lt;nil&gt;
  t_t_l = &lt;nil&gt;
  priority = &lt;nil&gt;
            }
```

Import

teo dns_record can be imported using the id, e.g.

```
terraform import tencentcloud_teo_dns_record.dns_record dns_record_id
```
*/
package tencentcloud

import (
"context"
"fmt"
"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
"log"
"strings"
"time"
)


func resourceTencentCloudTeoDnsRecord () * schema.Resource {
return & schema.Resource {
Create : resourceTencentCloudTeoDnsRecordCreate ,
Read : resourceTencentCloudTeoDnsRecordRead ,
Update : resourceTencentCloudTeoDnsRecordUpdate ,
Delete : resourceTencentCloudTeoDnsRecordDelete ,
Importer : & schema.ResourceImporter {
State : schema.ImportStatePassthrough ,
} ,
Schema : map[string] * schema.Schema {
"zone_id": {
  Required: true,
  Type: schema.TypeString,
  Description: "Site ID.",
},

"dns_record_id": {
  Computed: true,
  Type: schema.TypeString,
  Description: "DNS record ID.",
},

"type": {
  Required: true,
  Type: schema.TypeSet,
  Elem: &schema.Schema{
				Type: schema.TypeString,
	},
  Description: "DNS record Type. Valid values: `A`, `AAAA`, `CNAME`, `MX`, `TXT`, `NS`, `CAA`, `SRV`.",
},

"name": {
  Required: true,
  Type: schema.TypeString,
  Description: "DNS record Name.",
},

"content": {
  Required: true,
  Type: schema.TypeString,
  Description: "DNS record Content.",
},

"mode": {
  Required: true,
  Type: schema.TypeString,
  Description: "Proxy mode. Valid values:- `dns_only`: only DNS resolution of the subdomain is enabled.- `proxied`: subdomain is proxied and accelerated.",
},

"t_t_l": {
  Optional: true,
  Computed: true,
  Type: schema.TypeInt,
  Description: "Time to live of the DNS record cache in seconds.",
},

"priority": {
  Optional: true,
  Computed: true,
  Type: schema.TypeInt,
  Description: "Priority of the record. Valid value range: 1-50, the smaller value, the higher priority.",
},

"created_on": {
  Computed: true,
  Type: schema.TypeString,
  Description: "Creation date.",
},

"modified_on": {
  Computed: true,
  Type: schema.TypeString,
  Description: "Last modification date.",
},

"locked": {
  Computed: true,
  Type: schema.TypeBool,
  Description: "Whether the DNS record is locked.",
},

"status": {
  Computed: true,
  Type: schema.TypeString,
  Description: "Resolution status. Valid values: `active`, `pending`.",
},

"cname": {
  Computed: true,
  Type: schema.TypeString,
  Description: "CNAME address. Note: This field may return null, indicating that no valid value can be obtained.",
},

"domain_status": {
  Computed: true,
  Type: schema.TypeSet,
  Elem: &schema.Schema{
				Type: schema.TypeString,
	},
  Description: "Whether this domain enable load balancing, security, or l4 proxy capability. Valid values: `lb`, `security`, `l4`.",
},

} ,
}
} 

func resourceTencentCloudTeoDnsRecordCreate (d * schema.ResourceData,meta interface{}) error {
defer logElapsed ("resource.tencentcloud_teo_dns_record.create") ()
defer inconsistentCheck (d,meta) () 

logId := getLogId (contextNil) 

var (
request = teo.NewCreateDnsRecordRequest ()
response = teo.NewCreateDnsRecordResponse ()
zoneId string
dnsRecordId string
)
if v,ok := d . GetOk ("zone_id");ok {
zoneId = v .(string) 
 request . ZoneId = helper.String (v .(string))
} 

if v,ok := d . GetOk ("type");ok {
typeSet := v .(* schema.Set) . List () 
 for i := range typeSet {
type := typeSet [i] .(string)
request . Type = append(request . Type,& type)
}
} 

if v,ok := d . GetOk ("name");ok {
request . Name = helper.String (v .(string))
} 

if v,ok := d . GetOk ("content");ok {
request . Content = helper.String (v .(string))
} 

if v,ok := d . GetOk ("mode");ok {
request . Mode = helper.String (v .(string))
} 

if v,ok := d . GetOkExists ("t_t_l");ok {
request . TTL = helper.IntInt64 (v .(int))
} 

if v,ok := d . GetOkExists ("priority");ok {
request . Priority = helper.IntInt64 (v .(int))
} 

err := resource.Retry (writeRetryTimeout,func () * resource.RetryError {
result,e := meta .(* TencentCloudClient) . apiV3Conn . UseTeoClient () . CreateDnsRecord (request)
if e != nil {
return  retryError (e)
} else {
log.Printf ("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",logId,request . GetAction (),request . ToJsonString (),result . ToJsonString ())
}
response = result
return  nil
})
if err != nil {
log.Printf ("[CRITAL]%s create teo dnsRecord failed, reason:%+v",logId,err)
return  err
} 

zoneId = * response . Response . ZoneId
d . SetId (strings.Join ([] string {zoneId,dnsRecordId},FILED_SP)) 

service := TeoService {client:meta .(* TencentCloudClient) . apiV3Conn} 

conf := BuildStateChangeConf ([] string {},[] string {"active"},60 * readRetryTimeout,time.Second,service . TeoDnsRecordStateRefreshFunc (d . Id (),[] string {})) 

if _,e := conf . WaitForState ();e != nil {
return  e
} 

return resourceTencentCloudTeoDnsRecordRead (d,meta)
} 

func resourceTencentCloudTeoDnsRecordRead (d * schema.ResourceData,meta interface{}) error {
defer logElapsed ("resource.tencentcloud_teo_dns_record.read") ()
defer inconsistentCheck (d,meta) () 

logId := getLogId (contextNil) 

ctx := context.WithValue (context.TODO (),logIdKey,logId) 

service := TeoService {client:meta .(* TencentCloudClient) . apiV3Conn} 

idSplit := strings.Split (d . Id (),FILED_SP)
if len (idSplit) != 2 {
return fmt . Errorf ("id is broken,%s",d . Id ())
}
zoneId := idSplit [0]
dnsRecordId := idSplit [1]


dnsRecord,err := service . DescribeTeoDnsRecordById (ctx,zoneId,dnsRecordId)
if err != nil {
return err
} 

if dnsRecord == nil {
d . SetId ("")
log.Printf ("[WARN]%s resource `TeoDnsRecord` [%s] not found, please check if it has been deleted.\n",logId,d . Id ())
return nil
} 

if dnsRecord . ZoneId != nil {
_ = d . Set ("zone_id",dnsRecord . ZoneId)
} 

if dnsRecord . DnsRecordId != nil {
_ = d . Set ("dns_record_id",dnsRecord . DnsRecordId)
} 

if dnsRecord . Type != nil {
_ = d . Set ("type",dnsRecord . Type)
} 

if dnsRecord . Name != nil {
_ = d . Set ("name",dnsRecord . Name)
} 

if dnsRecord . Content != nil {
_ = d . Set ("content",dnsRecord . Content)
} 

if dnsRecord . Mode != nil {
_ = d . Set ("mode",dnsRecord . Mode)
} 

if dnsRecord . TTL != nil {
_ = d . Set ("t_t_l",dnsRecord . TTL)
} 

if dnsRecord . Priority != nil {
_ = d . Set ("priority",dnsRecord . Priority)
} 

if dnsRecord . CreatedOn != nil {
_ = d . Set ("created_on",dnsRecord . CreatedOn)
} 

if dnsRecord . ModifiedOn != nil {
_ = d . Set ("modified_on",dnsRecord . ModifiedOn)
} 

if dnsRecord . Locked != nil {
_ = d . Set ("locked",dnsRecord . Locked)
} 

if dnsRecord . Status != nil {
_ = d . Set ("status",dnsRecord . Status)
} 

if dnsRecord . Cname != nil {
_ = d . Set ("cname",dnsRecord . Cname)
} 

if dnsRecord . DomainStatus != nil {
_ = d . Set ("domain_status",dnsRecord . DomainStatus)
} 

return nil
} 

func resourceTencentCloudTeoDnsRecordUpdate (d * schema.ResourceData,meta interface{}) error {
defer logElapsed ("resource.tencentcloud_teo_dns_record.update") ()
defer inconsistentCheck (d,meta) () 

logId := getLogId (contextNil) 

request := teo.NewModifyDnsRecordRequest () 



idSplit := strings.Split (d . Id (),FILED_SP)
if len (idSplit) != 2 {
return fmt . Errorf ("id is broken,%s",d . Id ())
}
zoneId := idSplit [0]
dnsRecordId := idSplit [1]


request . ZoneId = & zoneId
request . DnsRecordId = & dnsRecordId


immutableArgs := [] string {"zone_id","dns_record_id","type","name","content","mode","t_t_l","priority","created_on","modified_on","locked","status","cname","domain_status"}


for _,v := range immutableArgs {
if d . HasChange (v) {
return fmt.Errorf ("argument `%s` cannot be changed",v)
}
}


if d . HasChange ("type") {
if v,ok := d . GetOk ("type");ok {
typeSet := v .(* schema.Set) . List () 
 for i := range typeSet {
type := typeSet [i] .(string)
request . Type = append(request . Type,& type)
}
}
} 

if d . HasChange ("name") {
if v,ok := d . GetOk ("name");ok {
request . Name = helper.String (v .(string))
}
} 

if d . HasChange ("content") {
if v,ok := d . GetOk ("content");ok {
request . Content = helper.String (v .(string))
}
} 

if d . HasChange ("mode") {
if v,ok := d . GetOk ("mode");ok {
request . Mode = helper.String (v .(string))
}
} 

if d . HasChange ("t_t_l") {
if v,ok := d . GetOkExists ("t_t_l");ok {
request . TTL = helper.IntInt64 (v .(int))
}
} 

if d . HasChange ("priority") {
if v,ok := d . GetOkExists ("priority");ok {
request . Priority = helper.IntInt64 (v .(int))
}
} 

err := resource.Retry (writeRetryTimeout,func () * resource.RetryError {
result,e := meta .(* TencentCloudClient) . apiV3Conn . UseTeoClient () . ModifyDnsRecord (request)
if e != nil {
return  retryError (e)
} else {
log.Printf ("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",logId,request . GetAction (),request . ToJsonString (),result . ToJsonString ())
}
return  nil
})
if err != nil {
log.Printf ("[CRITAL]%s update teo dnsRecord failed, reason:%+v",logId,err)
return  err
} 

return resourceTencentCloudTeoDnsRecordRead (d,meta)
} 

func resourceTencentCloudTeoDnsRecordDelete (d * schema.ResourceData,meta interface{}) error {
defer logElapsed ("resource.tencentcloud_teo_dns_record.delete") ()
defer inconsistentCheck (d,meta) () 

logId := getLogId (contextNil)
ctx := context.WithValue (context.TODO (),logIdKey,logId) 

service := TeoService {client:meta .(* TencentCloudClient) . apiV3Conn}
idSplit := strings.Split (d . Id (),FILED_SP)
if len (idSplit) != 2 {
return fmt . Errorf ("id is broken,%s",d . Id ())
}
zoneId := idSplit [0]
dnsRecordId := idSplit [1]


if err := service . DeleteTeoDnsRecordById (ctx,zoneId,dnsRecordId) ; err != nil {
return err
} 

return nil
} 
