package utils

import "time"

const (
  ANSIBLE_LOG = "./devops/log/ansible.log"
)

func Now() string {
  //return time.Now().Format("2006-01-02 15:04:05")
  return time.Now().Format("20060102150405")
}