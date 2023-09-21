package bridge_metamask

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/g45t345rt/g45w/utils"
)

//go:embed web/dist/*
var dist embed.FS

var PORT = 53943

type BridgeInData struct {
	WalletAddress string  `json:"walletAddress"`
	Symbol        string  `json:"symbol"`
	Amount        float64 `json:"amount"`
}

func createData(bridgeData BridgeInData) (base64Data string, err error) {
	jsonData, err := json.Marshal(bridgeData)
	if err != nil {
		return
	}

	base64Data = base64.URLEncoding.EncodeToString(jsonData)
	return
}

func Link(bridgeData BridgeInData) (url string, err error) {
	if utils.IsMobile() {
		return MetaMaskLink(bridgeData)
	}

	return HostLink(bridgeData)
}

func HostLink(bridgeData BridgeInData) (url string, err error) {
	base64Data, err := createData(bridgeData)
	if err != nil {
		return
	}

	url = fmt.Sprintf("http://127.0.0.1:%d/web/dist?data=%s", PORT, base64Data)
	return
}

func MetaMaskLink(bridgeData BridgeInData) (url string, err error) {
	base64Data, err := createData(bridgeData)
	if err != nil {
		return
	}

	url = fmt.Sprintf("https://metamask.app.link/dapp/bridge-in.deronfts.com?data=%s", base64Data)
	return
}

// Using http localhost server instead of https
// because browser don't accept self-signed certicate without warning

// For Metamask mobile I need to use https, so I have deployed a subdomain available at bridge-in.deronfts.com

func StartServer() error {
	if utils.IsMobile() {
		return nil
	}

check_port:
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", PORT))
	if err != nil {
		// desired port is in use
		PORT++
		goto check_port
	}
	listener.Close()

	handler := http.FileServer(http.FS(dist))
	http.Handle("/", handler)
	return http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", PORT), nil)
}

// keeping https code maybe for future
/*
func createCert() ([]byte, []byte, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	sn, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, err
	}

	tmpl := x509.Certificate{
		SerialNumber: sn,
		Subject: pkix.Name{
			Organization: []string{"G45W"},
			CommonName:   "G45W",
		},
		DNSNames:              []string{"127.0.0.1", "localhost"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	cert, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}

	certPem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert})
	keyPem, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, nil, err
	}

	keyPem = pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyPem})

	return certPem, keyPem, nil
}

func saveCert(certPem []byte, keyPem []byte) error {
	dataDir, err := app.DataDir()
	if err != nil {
		return err
	}

	certDir := filepath.Join(dataDir, "g45w", "cert")
	err = os.MkdirAll(certDir, os.ModePerm)
	if err != nil {
		return err
	}

	certPath := filepath.Join(certDir, "cert.crt")
	err = os.WriteFile(certPath, certPem, os.ModePerm)
	if err != nil {
		return err
	}

	keyPath := filepath.Join(certDir, "key.pem")
	err = os.WriteFile(keyPath, keyPem, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func loadCert() ([]byte, []byte, error) {
	dataDir, err := app.DataDir()
	if err != nil {
		return nil, nil, err
	}

	certDir := filepath.Join(dataDir, "g45w", "cert")
	err = os.MkdirAll(certDir, os.ModePerm)
	if err != nil {
		return nil, nil, err
	}

	certPath := filepath.Join(certDir, "cert.crt")
	certPem, err := os.ReadFile(certPath)
	if err != nil {
		return nil, nil, err
	}

	keyPath := filepath.Join(certDir, "key.pem")
	keyPem, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, nil, err
	}

	return certPem, keyPem, nil
}

func validateCertExpiration(certPem []byte) error {
	block, _ := pem.Decode(certPem)
	if block == nil {
		return fmt.Errorf("fail to decode cert")
	}

	cert, err := x509.ParseCertificate(certPem)
	if err != nil {
		return err
	}

	if time.Now().After(cert.NotAfter) {
		return fmt.Errorf("cert is expired")
	}

	return nil
}

{
	certPem, keyPem, err := loadCert()
	if err != nil {
		certPem, keyPem, err = createCert()
		if err != nil {
			return err
		}

		err = saveCert(certPem, keyPem)
		if err != nil {
			return err
		}
	}

	keyPair, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", PORT),
		Handler: http.FileServer(http.FS(dist)),
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
			Certificates:             []tls.Certificate{keyPair},
		},
	}

	return server.ListenAndServeTLS("", "")
}
*/
