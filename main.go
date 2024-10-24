package main

import (
	"fmt"
	fss "github.com/krackenservices/junit-api/fsservice"
	jp "github.com/krackenservices/junit-api/junitparser"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func SecureLogger() *slog.Logger {
	l := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				a.Value = slog.AnyValue(time.Now().UTC())
			}
			secrets := []string{"secret", "token", "password"}
			for _, secret := range secrets {
				if strings.Contains(strings.ToLower(a.Key), secret) && len(a.Value.String()) > 0 {
					a.Value = slog.StringValue("*****")
					break
				}
			}
			return a
		},
	})

	return slog.New(l)
}

var secLog *slog.Logger
var logLevel *slog.LevelVar
var testSuites []jp.TestSuites

type Config struct {
	LogLevel          string
	IP                string
	Port              string
	FileSystem        string
	ScanURI           string
	UpdateIntervalSec string
}

func getEnv(key, defaultstr string) string {
	if value, ok := os.LookupEnv(key); ok {
		secLog.Debug("getEnv:", "Value set", value)
		return value
	}
	secLog.Debug("getEnv:", "Value not set using default:", defaultstr)
	return defaultstr
}

func getData(cfg Config) {
	secLog.Info("Retrieving Data")
	fss := fss.FileSystemService{FsType: cfg.FileSystem}
	fss.Init()
	fsl, _ := fss.Fsi.ListFiles(cfg.ScanURI)

	var tmptestSuites []jp.TestSuites

	for _, suite := range fsl {
		contents, err := fss.Fsi.GetFileContents(fmt.Sprintf("%s/%s", cfg.ScanURI, suite))
		if err != nil {
			return
		}
		parsedSuites, err := jp.ParseJUnitXML(contents)
		if err != nil {
			fmt.Println("Error parsing XML:", err)
		}
		tmptestSuites = append(tmptestSuites, *parsedSuites)
	}
	testSuites = tmptestSuites
}

func (cfg *Config) ListCfg() {
	secLog.Debug("Listing Config")
	r := reflect.ValueOf(cfg).Elem()
	rt := r.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		rv := reflect.ValueOf(cfg)
		value := reflect.Indirect(rv).FieldByName(field.Name)
		secLog.Debug("Config:", field.Name, value.String())
	}

}

func setLogLeve(level string) {
	switch strings.ToLower(level) {
	case "":
	case "info":
		logLevel.Set(slog.LevelInfo)
	case "error":
		logLevel.Set(slog.LevelError)
	case "warn":
		logLevel.Set(slog.LevelWarn)
	case "debug":
		logLevel.Set(slog.LevelDebug)
	default:
		secLog.Error(fmt.Sprintf("Unknown LogLevel %s setting to Debug", level))
		logLevel.Set(slog.LevelDebug)
	}
}

func updateData(cfg Config) {
	getData(cfg)
	interval, err := strconv.Atoi(cfg.UpdateIntervalSec)
	if err != nil {
		secLog.Error(err.Error())
		interval = 30
	}
	tic := time.Second * time.Duration(interval)
	for range time.Tick(tic) {
		getData(cfg)
	}
}

func main() {
	logLevel = &slog.LevelVar{}
	secLog = SecureLogger()
	slog.SetLogLoggerLevel(slog.LevelDebug)
	secLog.Info("Starting...")

	cfg := Config{
		LogLevel:          getEnv("JAPI_LOGLEVEL", "INFO"),
		IP:                getEnv("JAPI_IP", "127.0.0.1"),
		Port:              getEnv("JAPI_PORT", "8080"),
		FileSystem:        getEnv("JAPI_FILESYSTEM", "local"),
		ScanURI:           getEnv("JAPI_SCANURI", "./testdata"),
		UpdateIntervalSec: getEnv("JAPI_UPDATESEC", "30"),
	}
	setLogLeve(cfg.LogLevel)
	cfg.ListCfg()

	go updateData(cfg)

	http.HandleFunc("/testsuites", testSuitesHandler)

	secLog.Info(fmt.Sprintf("Server is running on http://%s:%s/testsuites", cfg.IP, cfg.Port))
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", cfg.IP, cfg.Port), nil)
	if err != nil {
		secLog.Error(fmt.Sprintf("Error starting server: %v\n", err))
	}
}
