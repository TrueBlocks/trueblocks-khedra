package app

import (
	"strings"

	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/coordinator"
	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/types"
	"github.com/TrueBlocks/trueblocks-sdk/v6/services"
	"github.com/TrueBlocks/trueblocks-sdk/v6/services/monitor"
)

// ServiceFactory creates and configures services based on configuration
type ServiceFactory struct {
	config      *types.Config
	logger      *types.CustomLogger
	coordinator *coordinator.ScraperMonitorCoordinator
}

// NewServiceFactory creates a new service factory
func NewServiceFactory(config *types.Config, logger *types.CustomLogger, coord *coordinator.ScraperMonitorCoordinator) *ServiceFactory {
	return &ServiceFactory{
		config:      config,
		logger:      logger,
		coordinator: coord,
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
			scraperSvc.SetScrapeCompleteCallback(sf.coordinator.HandleEvent)
			activeServices = append(activeServices, scraperSvc)
		case "monitor":
			monitorSvc := sf.createMonitorService(svc)
			if monitorSvc != nil {
				chains := strings.Split(strings.ReplaceAll(sf.config.EnabledChains(), " ", ""), ",")
				for _, chain := range chains {
					sf.coordinator.RegisterMonitor(chain, monitorSvc)
				}
				activeServices = append(activeServices, monitorSvc)
			}
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
func (sf *ServiceFactory) createMonitorService(svc types.Service) *monitor.MonitorService {
	chains := strings.Split(strings.ReplaceAll(sf.config.EnabledChains(), " ", ""), ",")
	config := monitor.MonitorConfig{
		WatchlistDir:    sf.config.General.MonitorsFolder,
		CommandsDir:     sf.config.General.MonitorsFolder,
		BatchSize:       svc.BatchSize,
		Concurrency:     svc.Concurrency,
		Sleep:           svc.Sleep,
		MaxBlocksPerRun: 0,
	}
	monitorSvc, err := monitor.NewMonitorService(sf.logger.GetLogger(), chains, config)
	if err != nil {
		sf.logger.Error("Failed to create monitor service", "error", err)
		return nil
	}

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
