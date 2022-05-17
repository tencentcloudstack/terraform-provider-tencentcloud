resource "tencentcloud_security_group" "emr-sg" {
  name        = "keep-tf-emr-sg"
  description = "EMR预留安全组"
  project_id  = 0
}