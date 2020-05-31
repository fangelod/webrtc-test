package rtc

import (
	"encoding/json"
	"fmt"

	webrtc "github.com/pion/webrtc/v2"
	uuid "github.com/satori/go.uuid"
)

var (
	ongoingCalls map[uuid.UUID]*Call
)

type Constraints struct {
	Audio bool `json:"audio"`
	Video bool `json:"video"`
}

type RTCSessionDescription struct {
	Constraints Constraints               `json:"body"`
	Description webrtc.SessionDescription `json:"offer"`
	Name        string                    `json:"name"`
	User        string                    `json:"user"`
	Type        string                    `json:"type"`
}

type Call struct {
	ID           uuid.UUID                         `json:"id"`
	Name         string                            `json:"name"`
	Participants []string                          `json:"participants"`
	Connections  map[string]*webrtc.PeerConnection `json:"-"`
	Tracks       map[string]*webrtc.Track          `json:"-"`
}

func init() {
	ongoingCalls = make(map[uuid.UUID]*Call)
}

func NewCall(offer webrtc.SessionDescription, name, user string) (webrtc.SessionDescription, *Call, error) {
	call := &Call{
		ID:           uuid.NewV4(),
		Name:         name,
		Participants: []string{user},
		Connections:  make(map[string]*webrtc.PeerConnection),
		Tracks:       make(map[string]*webrtc.Track),
	}

	pc, err := newPeerConnection()
	if err != nil {
		return webrtc.SessionDescription{}, nil, err
	}

	if err := pc.SetRemoteDescription(offer); err != nil {
		pc.Close()
		return webrtc.SessionDescription{}, nil, err
	}

	ans, err := pc.CreateAnswer(nil)
	if err != nil {
		pc.Close()
		return webrtc.SessionDescription{}, nil, err
	}

	if err := pc.SetLocalDescription(ans); err != nil {
		pc.Close()
		return webrtc.SessionDescription{}, nil, err
	}

	call.Connections[user] = pc
	ongoingCalls[call.ID] = call

	return ans, call, nil
}

func GetCalls() ([]*Call, error) {
	var calls = make([]*Call, 0)
	for _, c := range ongoingCalls {
		calls = append(calls, c)
	}
	return calls, nil
}

func JoinCall(user, offer string, id uuid.UUID) (*Call, error) {
	if _, exists := ongoingCalls[id]; !exists {
		return nil, fmt.Errorf("unable to join call %v, call does not exist", id)
	}

	pc, err := newPeerConnection()
	if err != nil {
		return nil, err
	}

	body := RTCSessionDescription{}
	if err := json.Unmarshal([]byte(offer), &body); err != nil {
		return nil, fmt.Errorf("failed to unmarshal RTCSessionDescription")
	}

	if err := pc.SetRemoteDescription(body.Description); err != nil {
		pc.Close()
		return nil, err
	}

	ans, err := pc.CreateAnswer(nil)
	if err != nil {
		pc.Close()
		return nil, err
	}

	if err := pc.SetLocalDescription(ans); err != nil {
		pc.Close()
		return nil, err
	}

	ongoingCalls[id].Participants = append(ongoingCalls[id].Participants, user)
	ongoingCalls[id].Connections[user] = pc

	return ongoingCalls[id], nil
}

func LeaveCall(user string, id uuid.UUID) error {
	if _, exists := ongoingCalls[id]; !exists {
		return fmt.Errorf("unable to leave call %v, call does not exist", id)
	}

	return nil
}
