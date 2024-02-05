package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/p0t4t0sandwich/ampapi-go"
	"github.com/p0t4t0sandwich/ampapi-go/modules"
)

// -------------- Global Variables --------------

// Server manager
var serverManager ServerManager = ServerManager{}

func main() {
	// Get settings from settings.json
	settingsFile, err := os.ReadFile("./settings.json")
	if err != nil {
		fmt.Println("Error reading settings.json")
	}
	var settings Settings
	_ = json.Unmarshal(settingsFile, &settings)

	// Override settings from env
	ENV_IP_ADDRESS := os.Getenv("IP_ADDRESS")
	if ENV_IP_ADDRESS != "" {
		settings.IP_ADDRESS = ENV_IP_ADDRESS
	}
	ENV_PORT := os.Getenv("PORT")
	if ENV_PORT != "" {
		settings.PORT = ENV_PORT
	}
	ENV_AMP_API_URL := os.Getenv("AMP_API_URL")
	if ENV_AMP_API_URL != "" {
		settings.AMP_API_URL = ENV_AMP_API_URL
	}
	ENV_AMP_API_USERNAME := os.Getenv("AMP_API_USERNAME")
	if ENV_AMP_API_USERNAME != "" {
		settings.AMP_API_USERNAME = ENV_AMP_API_USERNAME
	}
	ENV_AMP_API_PASSWORD := os.Getenv("AMP_API_PASSWORD")
	if ENV_AMP_API_PASSWORD != "" {
		settings.AMP_API_PASSWORD = ENV_AMP_API_PASSWORD
	}

	controller, _ := modules.NewADS(
		settings.AMP_API_URL,
		settings.AMP_API_USERNAME,
		settings.AMP_API_PASSWORD,
	)
	serverManager = ServerManager{
		controller:   *controller,
		targetData:   make(map[string]Data[ampapi.IADSInstance, modules.ADS]),
		instanceData: make(map[string]Data[ampapi.Instance, interface{}]),
	}
	serverManager.InitInstnaceData()
	fmt.Println("Initialized instance data")

	gin.SetMode(gin.ReleaseMode)
	var router *gin.Engine = gin.Default()
	fmt.Println("break 1")

	// Static files
	router.StaticFile("/openapi.json", "./openapi.json")

	// Docs
	router.GET("/docs", getRoot)
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/docs")
	})
	fmt.Println("break 2")
	// API routes
	router.GET("/target/status/:targetName", getTargetStatus)
	router.GET("/instance/status/simple/:instanceName", getInstanceStatusSimple)
	router.GET("/server/status/:serverName", getServerStatus)
	router.GET("/server/status/simple/:serverName", getServerStatusSimple)
	router.Run(settings.IP_ADDRESS + ":" + settings.PORT)
	fmt.Println("Started server on " + settings.IP_ADDRESS + ":" + settings.PORT)

	// Update the server manager
	go func() {
		for {
			time.Sleep(20 * time.Second)
			serverManager.InitInstnaceData()
			fmt.Println("Updated instance data")
		}
	}()
}

// -------------- Structs --------------

// Settings - Settings struct
type Settings struct {
	IP_ADDRESS       string `json:"IP_ADDRESS"`
	PORT             string `json:"PORT"`
	AMP_API_URL      string `json:"AMP_API_URL"`
	AMP_API_USERNAME string `json:"AMP_API_USERNAME"`
	AMP_API_PASSWORD string `json:"AMP_API_PASSWORD"`
}

// Data - Generic data struct
type Data[T any, R any] struct {
	Data T
	API  R
}

// Instance handler
type ServerManager struct {
	controller   modules.ADS
	targetData   map[string]Data[ampapi.IADSInstance, modules.ADS]
	instanceData map[string]Data[ampapi.Instance, interface{}]
}

// -------------- Methods --------------

// InitInstnaceData - Initialize the instance data
func (s *ServerManager) InitInstnaceData() {
	// Get all instances
	var targets []ampapi.IADSInstance = nil
	targets, _ = s.controller.ADSModule.GetInstances()

	for i := 0; i <= len(targets)-1; i++ {
		var instances []ampapi.Instance = targets[i].AvailableInstances

		// Check if the target exists
		if s.TargetExists(targets[i].FriendlyName) {
			// Update the target data
			s.targetData[targets[i].FriendlyName] = Data[ampapi.IADSInstance, modules.ADS]{
				Data: targets[i],
				API:  s.targetData[targets[i].FriendlyName].API,
				fmt.Println[]
			}
		} else {
			// Add the target data
			s.targetData[targets[i].FriendlyName] = Data[ampapi.IADSInstance, modules.ADS]{
				Data: targets[i],
				API:  modules.ADS{},
			}
		}

		for j := 0; j <= len(instances)-1; j++ {
			// Check if the instance exists
			if s.InstanceExists(instances[j].FriendlyName) {
				// Update the instance data
				s.instanceData[instances[j].FriendlyName] = Data[ampapi.Instance, interface{}]{
					Data: instances[j],
					API:  s.instanceData[instances[j].FriendlyName].API,
				}
			} else {
				// Add the instance data
				s.instanceData[instances[j].FriendlyName] = Data[ampapi.Instance, interface{}]{
					Data: instances[j],
					API:  nil,
				}
			}
		}
	}

	// Remove targets that are not present in targets[i]
	for i := 0; i <= len(s.targetData)-1; i++ {
		var found bool = false
		for j := 0; j <= len(targets)-1; j++ {
			if s.targetData[targets[j].FriendlyName].Data.FriendlyName == targets[j].FriendlyName {
				found = true
				break
			}
		}
		if !found {
			delete(s.targetData, s.targetData[targets[i].FriendlyName].Data.FriendlyName)
		}

		var instances []ampapi.Instance = targets[i].AvailableInstances
		// Remove instances that are not present in instances[j]
		for j := 0; j <= len(s.instanceData)-1; j++ {
			var found bool = false
			for k := 0; k <= len(instances)-1; k++ {
				if s.instanceData[instances[k].FriendlyName].Data.FriendlyName == instances[k].FriendlyName {
					found = true
					break
				}
			}
			if !found {
				delete(s.instanceData, s.instanceData[instances[j].FriendlyName].Data.FriendlyName)
			}
		}
	}
}

// InitTargetAPI - Initialize a target's API
func (s *ServerManager) InitTargetAPI(targetName string) {
	// Get the target data
	var targetData Data[ampapi.IADSInstance, modules.ADS] = s.targetData[targetName]

	// Initialize the target's API
	var api interface{}
	api, err := s.controller.InstanceLogin(targetData.Data.InstanceId, "ADS")
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error logging in to target " + targetName)
	}
	targetData.API = *api.(*modules.ADS)

	// Store the target data
	s.targetData[targetName] = targetData
}

// InitInstanceAPI - Initialize an instance's API
func (s *ServerManager) InitInstanceAPI(instanceName string) {
	// Get the instance data
	var instanceData Data[ampapi.Instance, interface{}] = s.instanceData[instanceName]

	// Initialize the instance's API
	instanceData.API, _ = s.controller.InstanceLogin(instanceData.Data.InstanceID, instanceData.Data.Module)

	// Store the instance data
	s.instanceData[instanceName] = instanceData
}

// TargetExists - Check if a target exists
func (s *ServerManager) TargetExists(targetName string) bool {
	_, ok := s.targetData[targetName]
	return ok || targetName == "ADS01"
}

// InstanceExists - Check if an instance exists
func (s *ServerManager) InstanceExists(instanceName string) bool {
	_, ok := s.instanceData[instanceName]
	return ok || instanceName == "ADS01"
}

// GetTargetAPI - Get a target's API
func (s *ServerManager) GetTargetAPI(targetName string) modules.ADS {
	if targetName == "ADS01" {
		return s.controller
	}
	if s.targetData[targetName].API.Username == "" {
		s.InitTargetAPI(targetName)
	}
	return s.targetData[targetName].API
}

// GetTargetData - Get a target's data
func (s *ServerManager) GetTargetData(targetName string) ampapi.IADSInstance {
	return s.targetData[targetName].Data
}

// GetInstanceAPI - Get an instance's API
func (s *ServerManager) GetInstanceAPI(instanceName string) interface{} {
	if instanceName == "ADS01" {
		return s.controller
	}
	if s.instanceData[instanceName].API == nil {
		s.InitInstanceAPI(instanceName)
	}
	return s.instanceData[instanceName].API
}

// GetInstanceData - Get an instance's data
func (s *ServerManager) GetInstanceData(instanceName string) ampapi.Instance {
	return s.instanceData[instanceName].Data
}

// -------------- Functions --------------

// targetStatus - Get a target's status
func targetStatus(targetName string) (ampapi.Status, error) {
	return serverManager.GetTargetAPI(targetName).Core.GetStatus()
}

// instanceStatusSimple - Get the simple status of an instance
func instanceStatusSimple(instanceName string) string {
	for targetName := range serverManager.targetData {
		for _, instance := range serverManager.targetData[targetName].Data.AvailableInstances {
			if instance.FriendlyName == instanceName {
				instaceId := instance.InstanceID
				API := serverManager.GetTargetAPI(targetName)
				instanceStatuses, _ := API.ADSModule.GetInstanceStatuses()
				for i := 0; i <= len(instanceStatuses)-1; i++ {
					if instanceStatuses[i].InstanceID == instaceId {
						running := instanceStatuses[i].Running
						if running {
							return "Running"
						} else {
							return "Offline"
						}
					}
				}
			}
		}
	}

	return "Not implemented"
}

// serverStatus - Get a server's status
func serverStatus(serverName string) (ampapi.Status, error) {
	var instanceData ampapi.Instance = serverManager.GetInstanceData(serverName)
	switch instanceData.Module {
	case "ADS":
		return serverManager.GetTargetAPI(instanceData.FriendlyName).Core.GetStatus()
	case "Minecraft":
		return serverManager.GetInstanceAPI(instanceData.FriendlyName).(*modules.Minecraft).Core.GetStatus()
	case "GenericModule":
		return serverManager.GetInstanceAPI(instanceData.FriendlyName).(*modules.GenericModule).Core.GetStatus()
	default:
		return serverManager.GetInstanceAPI(instanceData.FriendlyName).(*modules.CommonAPI).Core.GetStatus()
	}
}

// serverStatusSimple - Get the simple status of a server
func serverStatusSimple(serverName string) (string, error) {
	status, err := serverStatus(serverName)
	return status.State.String(), err
}

// -------------- Handlers --------------

// Get root route
func getRoot(c *gin.Context) {
	// Read the html file
	html, err := os.ReadFile("./index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Return the html
	c.Data(http.StatusOK, "text/html", html)
}

// Get target status
func getTargetStatus(c *gin.Context) {
	// Get the target name
	targetName := c.Param("targetName")

	// Check if the target exists
	if !serverManager.TargetExists(targetName) {
		c.String(http.StatusNotFound, "Target not found")
		return
	}

	// Return the status
	status, err := targetStatus(targetName)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, status)
}

// Get instance status simple
func getInstanceStatusSimple(c *gin.Context) {
	// Get the instance name
	instanceName := c.Param("instanceName")

	// Check if the instance exists
	if !serverManager.InstanceExists(instanceName) {
		c.String(http.StatusNotFound, "Instance not found")
		return
	}

	// Return the status
	c.String(http.StatusOK, instanceStatusSimple(instanceName))
}5

// Get server status
func getServerStatus(c *gin.Context) {
	// Get the server name
	serverName := c.Param("serverName")

	// Check if the server exists
	if !serverManager.InstanceExists(serverName) {
		c.String(http.StatusNotFound, "Server not found")
		return
	}

	// Return the status
	status, err := serverStatus(serverName)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, status)
}

// Get server status simple
func getServerStatusSimple(c *gin.Context) {
	// Get the server name
	serverName := c.Param("serverName")

	// Check if the server exists
	if !serverManager.InstanceExists(serverName) {
		c.String(http.StatusNotFound, "Server not found")
		return
	}

	// Return the status
	status, err := serverStatusSimple(serverName)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, status)
}
