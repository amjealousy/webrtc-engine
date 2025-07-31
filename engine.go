package main

import (
	"fmt"
	"github.com/pion/ice/v4"
	"github.com/pion/stun/v3"
	"github.com/pion/transport/v3/stdnet"
	"github.com/pion/webrtc/v4"
	"log/slog"
)

var (
	staticEngine = EngineBuilder()
	peerAPI      = struct {
		config webrtc.Configuration
		api    *webrtc.API
	}{
		config: webrtc.Configuration{
			PeerIdentity: "go-signal",
			SDPSemantics: webrtc.SDPSemanticsUnifiedPlan,
		},
		api: GetPeerAPI(),
	}
)

func GetPeerAPI() *webrtc.API {
	api := webrtc.NewAPI(webrtc.WithSettingEngine(*staticEngine))
	return api
}
func GetPeerConnection() *webrtc.PeerConnection {
	connection, err := peerAPI.api.NewPeerConnection(peerAPI.config)
	if err != nil {
		slog.Error("Can`t create peer connection: %w", err)
	}
	return connection
}
func EngineBuilder() *webrtc.SettingEngine {
	engine := webrtc.SettingEngine{}
	engine.SetNetworkTypes([]webrtc.NetworkType{1})
	engine.SetICEBindingRequestHandler(func(m *stun.Message, local, remote ice.Candidate, pair *ice.CandidatePair) bool {
		n := fmt.Sprintf("ICEBinding request - %v , %v <> %v , pair - %v", m, local, remote, pair)
		slog.Info(n)
		return false
	})
	customNet, _ := stdnet.NewNet()
	engine.SetNet(customNet)

	return &engine
}
