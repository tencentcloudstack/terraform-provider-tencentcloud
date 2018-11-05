resource "tencentcloud_lb" "classic-lb" {
  type = 2
  forward = 0
  name = "tf-test-classic-lb"
  project_id = 0
}

resource "tencentcloud_lb" "application-lb" {
  type = 2
  forward = 1
  name = "tf-test-application-lb"
  project_id = 0
}
