package rtc

import (
	"fmt"

	webrtc "github.com/pion/webrtc/v2"
)

func newPeerConnection() (*webrtc.PeerConnection, error) {
	config := webrtc.Configuration{ICEServers: rtcIceServers}

	pc, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create new RTCPeerConnection")
	}

	if _, err := pc.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio); err != nil {
		pc.Close()
		return nil, err
	}

	if _, err := pc.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		pc.Close()
		return nil, err
	}

	return pc, nil
}
