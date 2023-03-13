package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	maddr "github.com/multiformats/go-multiaddr"
	"github.com/pion/stun"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sourcegraph/jsonrpc2"
	websocketjsonrpc2 "github.com/sourcegraph/jsonrpc2/websocket"
	"github.com/spf13/viper"
	"golang.org/x/crypto/acme/autocert"
	sfuLog "main/pkg/logger"
	"main/pkg/middlewares/datachannel"
	"main/pkg/node"
	"main/pkg/sfu"
	"net"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

var (
	file     string
	conf     = sfu.Config{}
	nodePort *uint
	domain   string
)

func showHelp() {
	fmt.Printf("Usage:%s {params}\n", os.Args[0])
	fmt.Println("      -n {p2p listen port}")
	fmt.Println("      -c {config file}")
	fmt.Println("      -d (domain)")
	fmt.Println("      -h (show help info)")
}

func parse() bool {
	nodePort = flag.Uint("n", 6666, "node port")
	flag.StringVar(&file, "c", "config.toml", "config file")
	flag.StringVar(&domain, "d", "", "domain")
	help := flag.Bool("h", false, "help info")
	flag.Parse()

	if *help {
		return false
	}
	return true
}

func load() bool {
	_, err := os.Stat(file)
	if err != nil {
		return false
	}

	viper.SetConfigFile(file)
	viper.SetConfigType("toml")

	err = viper.ReadInConfig()
	if err != nil {
		return false
	}
	err = viper.GetViper().Unmarshal(&conf)
	if err != nil {
		return false
	}

	return true
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	if !parse() {
		showHelp()
		os.Exit(-1)
	}

	if !load() {
		showHelp()
		os.Exit(-1)
	}

	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.With().Caller().Logger()

	sfu.Logger = sfuLog.New()

	ip, err := GetExternalIP(context.Background(), []string{"stun.l.google.com:19302"})
	if err != nil {
		panic(err)
	}
	conf.WebRTC.Candidates.NAT1To1IPs = []string{ip}

	s := sfu.NewSFU(conf)
	dc := s.NewDatachannel(sfu.APIChannelLabel)
	dc.Use(datachannel.SubscriberAPI)

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain),
	}

	dir := cacheDir()
	if dir != "" {
		certManager.Cache = autocert.DirCache(dir)
	}

	server := &http.Server{
		Addr: ":https",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}
	log.Printf("Serving http/https for domains: %+v", domain)
	go func() {
		// serve HTTP, which will redirect automatically to HTTPS
		h := certManager.HTTPHandler(nil)
		err := http.ListenAndServe(":http", h)
		if err != nil {
			panic(err)
		}

	}()

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	ctx := context.Background()

	n := node.NewNode(":57000" + "libp2p-webrtc.privkey")
	if err := n.Start(ctx, uint16(*nodePort)); err != nil {
		panic(err)
	}

	var bootstrapNodes []maddr.Multiaddr
	addr1, _ := maddr.NewMultiaddr("/ip4/51.38.127.87/tcp/6666/p2p/12D3KooWSYhTs8ykMDMDwdVc1gKb6EuSvWYCwFQaP2HyWezvqwQe")
	addr2, _ := maddr.NewMultiaddr("/ip4/51.75.160.227/tcp/6666/p2p/12D3KooWJzHpvL2bzcruPenptg6BrrnHjfk79LABE7Ju5qYjEA8q")

	bootstrapNodes = append(bootstrapNodes, addr1)
	bootstrapNodes = append(bootstrapNodes, addr2)

	if err := n.Bootstrap(ctx, bootstrapNodes); err != nil {
		log.Error().Err(err).Msg("bootstrap")
	}

	http.Handle("/ws", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error().Err(err).Msg("upgrader")
			return
		}
		defer c.Close()

		p := NewParticipant(sfu.NewPeer(s), n)
		defer p.Close()

		jc := jsonrpc2.NewConn(r.Context(), websocketjsonrpc2.NewObjectStream(c), p)
		<-jc.DisconnectNotify()
	}))

	err = server.ListenAndServeTLS("", "")
	if err != nil {
		panic(err)
	}
}

func cacheDir() (dir string) {
	if u, _ := user.Current(); u != nil {
		dir = filepath.Join(os.TempDir(), "cache-golang-autocert-"+u.Username)
		if err := os.MkdirAll(dir, 0700); err == nil {
			return dir
		}
	}
	return ""
}

// GetExternalIP return external IP for localAddr from stun server.
func GetExternalIP(ctx context.Context, stunServers []string) (string, error) {
	if len(stunServers) == 0 {
		return "", errors.New("STUN servers are required but not defined")
	}
	dialer := &net.Dialer{}
	conn, err := dialer.Dial("udp4", stunServers[0])
	if err != nil {
		return "", err
	}
	c, err := stun.NewClient(conn)
	if err != nil {
		return "", err
	}
	defer c.Close()

	message, err := stun.Build(stun.TransactionID, stun.BindingRequest)
	if err != nil {
		return "", err
	}

	var stunErr error
	// sufficiently large buffer to not block it
	ipChan := make(chan string, 20)
	err = c.Start(message, func(res stun.Event) {
		if res.Error != nil {
			stunErr = res.Error
			return
		}

		var xorAddr stun.XORMappedAddress
		if err := xorAddr.GetFrom(res.Message); err != nil {
			stunErr = err
			return
		}
		ip := xorAddr.IP.To4()
		if ip != nil {
			ipChan <- ip.String()
		}
	})
	if err != nil {
		return "", err
	}

	ctx1, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	select {
	case nodeIP := <-ipChan:
		return nodeIP, nil
	case <-ctx1.Done():
		msg := "could not determine public IP"
		if stunErr != nil {
			return "", errors.Wrap(stunErr, msg)
		} else {
			return "", fmt.Errorf(msg)
		}
	}
}
