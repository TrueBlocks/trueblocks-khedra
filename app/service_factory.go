package app

import (
	"strings"

	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/types"
	"github.com/TrueBlocks/trueblocks-sdk/v6/services"
)

// ServiceFactory creates and configures services based on configuration
type ServiceFactory struct {
	config *types.Config
	logger *types.CustomLogger
}

// NewServiceFactory creates a new service factory
func NewServiceFactory(config *types.Config, logger *types.CustomLogger) *ServiceFactory {
	return &ServiceFactory{
		config: config,
		logger: logger,
	}
}

// CreateAllServices creates all configured services including control service
func (sf *ServiceFactory) CreateAllServices(controlSvc *services.ControlService) []services.Servicer {
	var activeServices []services.Servicer
	activeServices = append(activeServices, controlSvc)

	for _, svc := range sf.config.Services {
		switch svc.Name {
		case "scraper":
			scraperSvc := sf.createScraperService(svc)
			activeServices = append(activeServices, scraperSvc)
		case "monitor":
			monitorSvc := sf.createMonitorService(svc)
			activeServices = append(activeServices, monitorSvc)
		case "api":
			if svc.Enabled {
				apiSvc := sf.createApiService()
				activeServices = append(activeServices, apiSvc)
			}
		case "ipfs":
			if svc.Enabled {
				ipfsSvc := sf.createIpfsService()
				activeServices = append(activeServices, ipfsSvc)
			}
		}
	}

	return activeServices
}

// createScraperService creates and configures the scraper service
func (sf *ServiceFactory) createScraperService(svc types.Service) *services.ScrapeService {
	chains := strings.Split(strings.ReplaceAll(sf.config.EnabledChains(), " ", ""), ",")
	scraperSvc := services.NewScrapeService(
		sf.logger.GetLogger(),
		"all",
		chains,
		sf.config.Services["scraper"].Sleep,
		sf.config.Services["scraper"].BatchSize,
	)
	if !svc.Enabled {
		scraperSvc.Pause()
	}
	return scraperSvc
}

// createMonitorService creates and configures the monitor service
func (sf *ServiceFactory) createMonitorService(svc types.Service) *services.MonitorService {
	monitorSvc := services.NewMonitorService(nil)
	if !svc.Enabled {
		monitorSvc.Pause()
	}
	return monitorSvc
}

// createApiService creates the API service
func (sf *ServiceFactory) createApiService() *services.ApiService {
	return services.NewApiService(sf.logger.GetLogger())
}

// createIpfsService creates the IPFS service
func (sf *ServiceFactory) createIpfsService() *services.IpfsService {
	return services.NewIpfsService(sf.logger.GetLogger())
}
