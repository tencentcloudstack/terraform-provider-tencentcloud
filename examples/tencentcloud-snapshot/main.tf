resource "tencentcloud_snapshot" "my-snapshot" {
  storage_id    = "${var.storage_id}"
  snapshot_name = "${var.snapshot_name}"
}
