package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cf-platform-eng/mongodb-on-demand-release/src/mongodb-service-adapter/adapter"
	jsonpatch "github.com/evanphx/json-patch"
)

var (
	configFilePath string
)

func main() {
	flag.StringVar(&configFilePath, "config", "", "Location of the config file")
	flag.Parse()

	config, err := LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("Error loading config file: %s", err)
	}

	logger := log.New(os.Stderr, "[mongodb-config-agent] ", log.LstdFlags)
	omClient := adapter.OMClient{Url: config.URL, Username: config.Username, ApiKey: config.APIKey}

	patchJSON, err := ioutil.ReadFile(config.PatchFile)
	if err != nil {
		log.Fatalf("Error reading JSON patch file: %s", err)
	}

	nodes := strings.Split(config.NodeAddresses, ",")
	ctx := &adapter.DocContext{
		ID:            config.ID,
		Key:           config.AuthKey,
		AdminPassword: config.AdminPassword,
		Nodes:         nodes,
		Version:       config.EngineVersion,
		RequireSSL:    config.RequireSSL,
	}

	if config.PlanID == adapter.PlanShardedCluster {
		var err error
		ctx.Cluster, err = adapter.NodesToCluster(nodes, config.Routers, config.ConfigServers, config.Replicas)
		if err != nil {
			logger.Fatal(err)
		}
	}

	logger.Printf("%+v", nodes)
	doc, err := omClient.LoadDoc(config.PlanID, ctx)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(doc)

	var modifiedDoc []byte
	if BytesToString(patchJSON) != "" {
		patch, err := jsonpatch.DecodePatch(patchJSON)
		if err != nil {
			logger.Fatalf("Error decoding patch file: %s", err)
		}

		modifiedDoc, err = patch.Apply([]byte(doc))
		if err != nil {
			logger.Fatalf("Error applying patch file: %s", err)
		}
	}

	automationAgentDoc := doc
	if BytesToString(patchJSON) != "" {
		automationAgentDoc = BytesToString(modifiedDoc)
	}

	monitoringAgentDoc, err := omClient.LoadDoc(adapter.MonitoringAgentConfiguration, ctx)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(monitoringAgentDoc)

	backupAgentDoc, err := omClient.LoadDoc(adapter.MonitoringAgentConfiguration, ctx)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(backupAgentDoc)

	for {
		logger.Printf("Checking group %s", config.GroupID)

		groupHosts, err := omClient.GetGroupHosts(config.GroupID)
		if err != nil {
			logger.Fatal(err)
		}

		//	logger.Printf("total number of hosts *** %v", groupHosts.TotalCount)
		if groupHosts.TotalCount == 0 {
			logger.Printf("Host count for %s is 0, configuring...", config.GroupID)

			err = omClient.ConfigureGroup(automationAgentDoc, config.GroupID)
			if err != nil {
				logger.Fatal(err)
			}

			err = omClient.ConfigureMonitoringAgent(monitoringAgentDoc, config.GroupID)
			if err != nil {
				logger.Fatal(err)
			}

			err = omClient.ConfigureBackupAgent(backupAgentDoc, config.GroupID)
			if err != nil {
				logger.Fatal(err)
			}

			logger.Printf("Configured group %s", config.GroupID)
		}

		time.Sleep(30 * time.Second)
	}
}

// BytesToString convert []byte to string
func BytesToString(data []byte) string {
	return string(data[:])
}
