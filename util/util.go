package util

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	SUCCESS = "SUCCESS"
	FAILED  = "FAILED"
	TRACE   = "TRACE"
	ERROR   = "ERROR"
	WARNING = "WARNING"
	INFO    = "INFO"
)

var (
	AppConfiguration *Configuration
	Error            *log.Logger
	Info             *log.Logger
	Warning          *log.Logger
	Trace            *log.Logger
	LogLevel         string = "ERROR"
)

//SMTPConfig : SMTPConfig Properties
type SMTPConfig struct {
	Host        string `json:"smtpServerHost"`
	Port        int    `json:"smtpServerPort"`
	Username    string `json:"emailAuthUserName"`
	Password    string `json:"emailAuthPassword"`
	EmailFrom   string `json:"emailFrom"`
	EmailSender string `json:"emailSender"`
	TLS         bool   `json:"tls"`
}

//SMSConfig : Twilio SMS Config Properties
type SMSConfig struct {
	UserName   string `json:"twilioAccountId"`
	Password   string `json:"twilioAuthToken"`
	SenderName string `json:"smsSender"`
}

//WebhookConfig : Webhook Config Properties
type WebhookConfig struct {
	AppURL string `json:"appURL"`
	AppKey string `json:"appKey"`
}

//KafkaConfig : Apache Kafka Config Properties
type KafkaConfig struct {
	BootstrapServers string `json:"bootstrapServers"`
	KafkaTopic       string `json:"kafkaTopic"`
	KafkaTopicConfig string `json:"kafkaTopicConfig"`
	KafkaGroupId     string `json:"kafkaGroupId"`
	KafkaTimeout     int    `json:"kafkaTimeout"`
}

//DBConfig : MongoDB Config Properties
type DBConfig struct {
	MongoHost       string `json:"mongoHost"`
	MongoPort       int    `json:"mongoPort"`
	MongoDBUsername string `json:"mongoDBUsername"`
	MongoDBPassword string `json:"mongoDBPassword"`
	MongoDB         string `json:"mongoDB"`
	Collection      string `json:"collection"`
}

//Configuration : Configuration File
type Configuration struct {
	Workers         int               `json:"workers"`
	LogFileLocation string            `json:"logFileLocation"`
	Log             bool              `json:"log"`
	KafkaConfig     KafkaConfig       `json:"kafkaConfig"`
	DbConfig        DBConfig          `json:"dbConfig"`
	SmsConfig       SMSConfig         `json:"smsConfig"`
	SmtpConfig      SMTPConfig        `json:"emailConfig"`
	WebhookConfig   WebhookConfig     `json:"webhookConfig"`
	Templates       map[string]string `json:"templates"`
}

//SetLogLevel : Set Logging Level
func SetLogLevel(logLevel string) {
	if AppConfiguration.Log {
		f, err := os.OpenFile(AppConfiguration.LogFileLocation, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Error opening log file: %s", err.Error())
		}
		if !AppConfiguration.Log {
			initLog(ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard,
				false)
			return
		}
		switch strings.ToUpper(logLevel) {
		case TRACE:
			initLog(f, f, f, f, true)
			return
		case INFO:
			initLog(ioutil.Discard, f, f, f, true)
			return
		case WARNING:
			initLog(ioutil.Discard, ioutil.Discard, f, f, true)
			return
		case ERROR:
			initLog(ioutil.Discard, ioutil.Discard, ioutil.Discard, f, true)
			return
		default:
			initLog(ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard,
				false)
			f.Close()
			return
		}
	}
	initLog(ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard,
		false)
}

//NewConfiguration : Loads App Config from File
func NewConfiguration() {
	var jsonConfig *os.File
	dir, _ := filepath.Abs("../")
	jsonConfig, err := os.Open(dir + "/configuration.json")
	if err != nil {
		dir, _ := filepath.Abs("./")
		jsonConfig, err = os.Open(dir + "/configuration.json")
		if err != nil {
			fmt.Println("Error reading configuration file " + err.Error())
			return
		}
	}
	defer jsonConfig.Close()
	byteValue, err := ioutil.ReadAll(jsonConfig)
	if err != nil {
		fmt.Println("Error reading configuration file " + err.Error())
		return
	}
	er := json.Unmarshal(byteValue, &AppConfiguration)
	if er != nil {
		fmt.Println("Error parsing json configuration file ")
		return
	}
	SetLogLevel(LogLevel)
	return
}

//GetTemplate : Gets Template From Config File
func (config *Configuration) GetTemplate(templateId string) string {
	return AppConfiguration.Templates[templateId]
}

func initLog(traceHandle io.Writer, infoHandle io.Writer,
	warningHandle io.Writer, errorHandle io.Writer, isFlag bool) {
	flag := 0
	if isFlag {
		flag = log.Ldate | log.Ltime | log.Lshortfile | log.LstdFlags
	}
	Trace = log.New(traceHandle, "TRACE: ", flag)
	Info = log.New(infoHandle, "INFO: ", flag)
	Warning = log.New(warningHandle, "WARNING: ", flag)
	Error = log.New(errorHandle, "ERROR: ", flag)
}
