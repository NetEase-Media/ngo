package httplib

const (
	HeaderKeyDowngrade          = "X-Downgrade-Status" // 降级key
	HeaderValueDowngradeStatic  = "STATIC"
	HeaderValueDowngradeDynamic = "DYNAMIC"
	HeaderKeyCircuitbreaking    = "X-Envoy-Circuitbreaking" // 熔断
)
