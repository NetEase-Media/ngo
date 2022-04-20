package service

var defaultService *Service

func SetDefaultService(service *Service) {
	defaultService = service
}

func GetAppName() string {
	return defaultService.AppName
}

func GetClusterName() string {
	return defaultService.ClusterName
}

func New(opt *Options) (*Service, error) {
	if err := checkOptions(opt); err != nil {
		return nil, err
	}
	return &Service{
		AppName:     opt.AppName,
		ClusterName: opt.ClusterName,
	}, nil
}

type Service struct {
	opt         *Options
	AppName     string
	ClusterName string
}
