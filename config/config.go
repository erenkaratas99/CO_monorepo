package config

import (
	"github.com/erenkaratas99/COApiCore/pkg/customErrors"
	"net/http"
	"time"
)

type RepoConfigs struct {
	MongoDuration  time.Duration
	MongoClientURI string
	BaseUrl        string
	ServiceURL     string
	ServiceCfgs    ServiceCfgs
}

type ServiceCfgs struct {
	Ocfg OrderConfigs
	Ccfg CustomerConfigs
}

type OrderConfigs struct {
	DBName     string
	ColName    string
	Port       string
	BaseUrl    string
	ServiceURL string
}
type CustomerConfigs struct {
	DBName     string
	ColName    string
	Port       string
	BaseUrl    string
	ServiceURL string
}

var cfgs = map[string]RepoConfigs{
	"prod": {
		MongoDuration:  time.Second * 10,
		MongoClientURI: "mongodb://root:root1234@mongodb_docker:27017",
		BaseUrl:        "0.0.0.0",
		ServiceURL:     "customer_app:8001",
		ServiceCfgs: ServiceCfgs{
			Ccfg: CustomerConfigs{
				DBName:     "CustomerDB",
				ColName:    "customers",
				Port:       ":8000",
				BaseUrl:    "0.0.0.0",
				ServiceURL: "customer_app:8001",
			},
			Ocfg: OrderConfigs{
				DBName:     "OrderDB",
				ColName:    "orders",
				Port:       ":8000",
				BaseUrl:    "0.0.0.0",
				ServiceURL: "order_app:8000",
			},
		},
	},
	"dev": {
		MongoDuration:  time.Second * 10,
		MongoClientURI: "mongodb://localhost:27017/",
		BaseUrl:        "127.0.0.1",
		ServiceURL:     "localhost:8001",
		ServiceCfgs: ServiceCfgs{
			Ccfg: CustomerConfigs{
				DBName:     "CustomerDB",
				ColName:    "customers",
				Port:       ":8000",
				BaseUrl:    "127.0.0.1",
				ServiceURL: "localhost:8001",
			},
			Ocfg: OrderConfigs{
				DBName:     "OrderDB",
				ColName:    "orders",
				Port:       ":8001",
				BaseUrl:    "127.0.0.1",
				ServiceURL: "localhost:8001",
			},
		},
	},
}

func GetRepoConfig(env string) (*RepoConfigs, error) {
	config, isExist := cfgs[env]
	if !isExist {
		return nil, customErrors.NewHTTPError(http.StatusInternalServerError,
			"ConfigErr",
			"Service configs could not have fetched correctly.")
	}
	return &config, nil
}
