package config

import "os"

const (
	EnvKeyPodName = "POD_NAME"
	EnvKeyPodIp   = "POD_IP"
)

var (
	PodName = os.Getenv(EnvKeyPodName)
	PodIp   = os.Getenv(EnvKeyPodIp)
)
