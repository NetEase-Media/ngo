package env

import "os"

const (
	appNameKey     = "appname"
	clusterNameKey = "clustername"
)

func GetAppName() string {
	return os.Getenv(appNameKey)
}

func GetClusterName() string {
	return os.Getenv(clusterNameKey)
}

func Get(name string) string {
	return os.Getenv(name)
}
