package stats

import (
	"errors"
	"log"
	"time"

	"github.com/p0t4t0sandwich/ampapi-go"
	"github.com/p0t4t0sandwich/ampapi-go/modules"
)

// Service - service struct
type Service struct {
	settings   *Settings
	Controller *modules.ADS
	Targets    map[string]*TargetData
	Instances  map[string]*InstanceData
}

// NewService - Create a new service
func NewService(settings *Settings) *Service {
	s := &Service{
		settings:  settings,
		Targets:   make(map[string]*TargetData),
		Instances: make(map[string]*InstanceData),
	}
	s.InitData()
	return s
}

// InitData - Init the instance data
func (s *Service) InitData() {
	ads, err := modules.NewADS(
		s.settings.AMP_API_URL,
		s.settings.AMP_API_USERNAME,
		s.settings.AMP_API_PASSWORD,
	)
	if err != nil {
		panic(err)
	}
	s.Controller = ads

	targets, err := s.Controller.ADSModule.GetInstances()
	if err != nil {
		panic(err)
	}
	for _, target := range targets {
		s.Targets[target.FriendlyName] = &TargetData{
			TargetID: target.InstanceId,
			Target:   nil,
		}
		for _, instance := range target.AvailableInstances {
			s.Instances[instance.InstanceName] = &InstanceData{
				InstanceID: instance.InstanceID,
				Instance:   nil,
			}
		}
	}
	log.Println("Initialized instance data")

	go func() {
		for {
			time.Sleep(5 * time.Minute)
			s.RefreshData()
			log.Println("Updated instance data")
		}
	}()
}

// RefreshData - Refresh the data -- basically InitData without the nils
func (s *Service) RefreshData() {
	targets, err := s.Controller.ADSModule.GetInstances()
	if err != nil {
		err := s.ReauthController()
		if err != nil {
			panic(err)
		}
	}
	for _, target := range targets {
		if t, ok := s.Targets[target.FriendlyName]; ok {
			s.Targets[target.FriendlyName] = t
		} else {
			s.Targets[target.FriendlyName] = &TargetData{
				TargetID: target.InstanceId,
			}
		}
		for _, instance := range target.AvailableInstances {
			if i, ok := s.Instances[instance.InstanceName]; ok {
				s.Instances[instance.InstanceName] = i
			} else {
				s.Instances[instance.InstanceName] = &InstanceData{
					InstanceID: instance.InstanceID,
				}
			}
		}
	}
}

// ReauthController - Reauthenticate the controller
func (s *Service) ReauthController() error {
	controller, err := modules.NewADS(
		s.settings.AMP_API_URL,
		s.settings.AMP_API_USERNAME,
		s.settings.AMP_API_PASSWORD,
	)
	if err != nil {
		return err
	}
	s.Controller = controller
	return nil
}

// ReauthTarget - Reauthenticate a target
func (s *Service) ReauthTarget(targetName string) error {
	target := s.Targets[targetName]
	if target.Target == nil {
		t, err := s.Controller.InstanceLogin(target.TargetID, "ADS")
		if err != nil {
			return err
		}
		target.Target = t.(*modules.ADS)
	}
	return nil
}

// ReauthInstance - Reauthenticate an instance
func (s *Service) ReauthInstance(instanceName string) error {
	instance := s.Instances[instanceName]
	if instance.Instance == nil {
		i, err := s.Controller.InstanceLogin(instance.InstanceID, "CommonAPI")
		if err != nil {
			return err
		}
		instance.Instance = i.(*modules.CommonAPI)
	}
	return nil
}

// GetTargetStatus - Get the target status
func (s *Service) GetTargetStatus(targetName string) (*ampapi.Status, error) {
	if targetName == "ADS01" {
		status, err := s.Controller.Core.GetStatus()
		if err != nil {
			s.ReauthController()
			return nil, err
		}
		return &status, nil
	}

	target := s.Targets[targetName]
	if target.Target == nil {
		err := s.ReauthTarget(targetName)
		if err != nil {
			return nil, err
		}
	}

	status, err := target.Target.Core.GetStatus()
	if err != nil {
		s.Targets[targetName].Target = nil
		s.ReauthTarget(targetName)
		return nil, err
	}
	return &status, nil
}

// GetServerStatus - Get the instance status
func (s *Service) GetServerStatus(instanceName string) (*ampapi.Status, error) {
	instance := s.Instances[instanceName]
	if instance.Instance == nil {
		err := s.ReauthInstance(instanceName)
		if err != nil {
			return nil, err
		}
	}

	status, err := instance.Instance.Core.GetStatus()
	if status.State == 0 && status.Uptime == "" && status.Metrics == nil {
		return nil, errors.New("instance is not running")
	}
	if err != nil {
		s.Instances[instanceName].Instance = nil
		s.ReauthInstance(instanceName)
		return nil, err
	}
	return &status, nil
}

// InstanceStatusSimple - Get the instance status in a simple format
func (s *Service) InstanceStatusSimple(instanceName string) (string, error) {
	targets, err := s.Controller.ADSModule.GetInstances()
	if err != nil {
		return "", err
	}
	for _, target := range targets {
		for _, instance := range target.AvailableInstances {
			if instance.InstanceName == instanceName {
				instanceID := instance.InstanceID
				instanceStatuses, _ := s.Controller.ADSModule.GetInstanceStatuses()
				for _, instanceStatus := range instanceStatuses {
					if instanceStatus.InstanceID == instanceID {
						if instanceStatus.Running {
							return "Running", nil
						}
						return "Offline", nil
					}
				}
			}
		}
	}
	return "", err
}

// ServerStatusSimple - Get the server status in a simple format
func (s *Service) ServerStatusSimple(serverName string) (string, error) {
	instance := s.Instances[serverName]
	if instance.Instance == nil {
		err := s.ReauthInstance(serverName)
		if err != nil {
			return "", err
		}
	}

	status, err := instance.Instance.Core.GetStatus()
	if status.State == 0 && status.Uptime == "" && status.Metrics == nil {
		return "InstanceOffline", nil
	}
	if err != nil {
		s.Instances[serverName].Instance = nil
		s.ReauthInstance(serverName)
		return "", err
	}
	return status.State.String(), nil
}
