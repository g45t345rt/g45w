package bridge_metamask

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"gioui.org/x/browser"
)

//go:embed web/dist/*
var dist embed.FS

var PORT = 53943

type BridgeInData struct {
	WalletAddress string  `json:"walletAddress"`
	Symbol        string  `json:"symbol"`
	Amount        float64 `json:"amount"`
}

func Open(bridgeData BridgeInData) error {
	data, err := json.Marshal(bridgeData)
	if err != nil {
		return err
	}

	base64Data := base64.URLEncoding.EncodeToString(data)
	url := fmt.Sprintf("http://localhost:%d/web/dist?data=%s", PORT, base64Data)
	browser.OpenUrl(url)
	return nil
}

func StartServer() error {
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
	return http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
