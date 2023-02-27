package main

import (
	"CustomerOrderMonoRepo/CustomerService/src/cmd"
	ordercmd "CustomerOrderMonoRepo/OrderService/src/cmd"
	"CustomerOrderMonoRepo/config"
	"github.com/erenkaratas99/COApiCore/pkg"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

func init() {
	pkg.InitLogrusConfig()

}

func main() {
	singleHTTPClient := pkg.NewSingletonClient()
	cfg, err := config.GetRepoConfig("dev")
	if err != nil {
		log.Fatal(err)
	}
	mongoClient, err := pkg.GetMongoClient(cfg.MongoDuration, cfg.MongoClientURI)
	if err != nil {
		log.Fatal(err)
	}
	bootMicroservices("both", singleHTTPClient, mongoClient, cfg)
}

func bootMicroservices(boot string, singleHTTPClient *pkg.RestClient, mc *mongo.Client, cfg *config.RepoConfigs) {
	if boot == "both" {
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			customercmd.Boot(singleHTTPClient, mc, &cfg.ServiceCfgs.Ccfg)
		}()
		go func() {
			defer wg.Done()
			ordercmd.Boot(singleHTTPClient, mc, &cfg.ServiceCfgs.Ocfg)
		}()
		wg.Wait()
	} else if boot == "customer" {
		customercmd.Boot(singleHTTPClient, mc, &cfg.ServiceCfgs.Ccfg)
	} else if boot == "order" {
		ordercmd.Boot(singleHTTPClient, mc, &cfg.ServiceCfgs.Ocfg)
	}
}
