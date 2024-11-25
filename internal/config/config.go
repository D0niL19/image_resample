package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"imageResample/pkg/utils"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string `yaml:"env" env-default:"dev"`
	HttpServer  `yaml:"http_server"`
	PathOrig    string
	PathRes     string
	ImageWidth  int
	ImageHeight int
}

type HttpServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"360s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatalln("CONFIG_PATH environment variable not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalln("CONFIG_PATH does not exist")
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalln("Can not read config file")
	}

	ParseFlags(&cfg)

	CreateDirectoriesIfNotExists(&cfg)

	return &cfg
}

func ParseFlags(cfg *Config) {
	flag.StringVar(&cfg.PathOrig, "path-orig", "/tmp/img_orig/", "Директория для исходных изображений")
	flag.StringVar(&cfg.PathRes, "path-res", "/tmp/img_res/", "Директория для обработанных изображений")
	flag.IntVar(&cfg.ImageWidth, "width", 200, "Ширина для ресэмплинга")
	flag.IntVar(&cfg.ImageHeight, "height", 200, "Высота для ресэмплинга")

	flag.Parse()
}

func CreateDirectoriesIfNotExists(cfg *Config) {
	utils.EnsureDir(cfg.PathOrig)
	utils.EnsureDir(cfg.PathRes)
}
