resource "tencentcloud_cls_Logset" "logset_basic" {
  logset_name    = var.logset_name

}
data "tencentcloud_cls_logsets" "logsets" {
     filters {
                key = "logsetId"
                value = [tencentcloud_cls_logset.logset_basic.id]
             }
     limit = 2
}