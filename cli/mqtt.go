package cli

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	MQTTDefaultConnectTimeout = "5s"
	MQTTDefaultClientId       = "plumber"
)

type MQTTOptions struct {
	// Shared
	Address  string
	Topic    string
	Timeout  time.Duration
	ClientID string
	QoSLevel int

	// TLS-related pieces
	TLSCAFile         string
	TLSClientCertFile string
	TLSClientKeyFile  string
	InsecureTLS       bool

	// Read
	ReadTimeout time.Duration

	// Write
	WriteTimeout time.Duration
}

func HandleMQTTFlags(readCmd, writeCmd *kingpin.CmdClause, opts *Options) {
	rc := readCmd.Command("mqtt", "MQTT message system")

	addSharedMQTTFlags(rc, opts)
	addReadMQTTFlags(rc, opts)

	wc := writeCmd.Command("mqtt", "MQTT message system")

	addSharedMQTTFlags(wc, opts)
	addWriteMQTTFlags(wc, opts)
}

func addSharedMQTTFlags(cmd *kingpin.CmdClause, opts *Options) {
	clientId := fmt.Sprintf("%s-%s", MQTTDefaultClientId, uuid.New().String()[0:3])

	cmd.Flag("address", "Destination host address").Default("tcp://localhost:1883").StringVar(&opts.MQTT.Address)
	cmd.Flag("topic", "Topic to read message(s) from").Required().StringVar(&opts.MQTT.Topic)
	cmd.Flag("timeout", "Connect timeout").Default(MQTTDefaultConnectTimeout).
		DurationVar(&opts.MQTT.Timeout)
	cmd.Flag("client-id", "Client id presented to MQTT broker").
		Default(clientId).StringVar(&opts.MQTT.ClientID)
	cmd.Flag("qos", "QoS level to use for pub/sub (0, 1, 2)").Default("0").IntVar(&opts.MQTT.QoSLevel)
	cmd.Flag("tls-ca-file", "CA file (only needed if addr is ssl://").ExistingFileVar(&opts.MQTT.TLSCAFile)
	cmd.Flag("tls-client-cert-file", "Client cert file (only needed if addr is ssl://").
		ExistingFileVar(&opts.MQTT.TLSClientCertFile)
	cmd.Flag("tls-client-key-file", "Client key file (only needed if addr is ssl://").
		ExistingFileVar(&opts.MQTT.TLSClientKeyFile)
	cmd.Flag("insecure-tls", "Whether to verify server certificate").Default("false").
		BoolVar(&opts.MQTT.InsecureTLS)
}

func addReadMQTTFlags(cmd *kingpin.CmdClause, opts *Options) {
	cmd.Flag("read-timeout", "How long to wait for a message (default: forever)").
		Default("0s").DurationVar(&opts.MQTT.ReadTimeout)
}

func addWriteMQTTFlags(cmd *kingpin.CmdClause, opts *Options) {
	cmd.Flag("write-timeout", "How long to attempt to publish a message").
		Default("5s").DurationVar(&opts.MQTT.WriteTimeout)
}
