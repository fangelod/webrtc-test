package rtc

import (
	// "encoding/json"
	"fmt"
        "io"
        "time"
        "math/rand"

        "github.com/pion/rtcp"
	webrtc "github.com/pion/webrtc/v2"
	uuid "github.com/satori/go.uuid"
)

const (
        rtcpPLIInterval = time.Second *3
)

var (
	ongoingCalls map[uuid.UUID]*Call
        TestTrack *webrtc.Track
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
        Track        *webrtc.Track                     `json:"-"`
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
	//pc, err := newPeerConnection()
	//if err != nil {
	//	return webrtc.SessionDescription{}, nil, err
	//}
        
        mediaEngine := webrtc.MediaEngine{}
        err := mediaEngine.PopulateFromSDP(offer)
        if err != nil {
                return webrtc.SessionDescription{}, nil, err
        }
        
        videoCodecs := mediaEngine.GetCodecsByKind(webrtc.RTPCodecTypeVideo)
        if len(videoCodecs) == 0 {
                fmt.Println("Offer had no video codecs");
        }

        api := webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine))

        config := webrtc.Configuration{ICEServers: rtcIceServers}
        pc, err := api.NewPeerConnection(config)
        if err != nil {
                return webrtc.SessionDescription{}, nil, err
        }

        if _, err = pc.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
                pc.Close()
                return webrtc.SessionDescription{}, nil, err
        }
        
        outputTrack, err := pc.NewTrack(videoCodecs[0].PayloadType, rand.Uint32(), "video", "pion")
        if err != nil {
                fmt.Println("NewTrack creation err")
                fmt.Println(err)
                pc.Close()
                return webrtc.SessionDescription{}, nil, err
        }
        
        if _, err = pc.AddTrack(outputTrack); err !=nil {
                fmt.Println("AddTrack Err")
                fmt.Println(err)
                pc.Close()
                return webrtc.SessionDescription{}, nil, err
        }

        //localTrackChan := make(chan *webrtc.Track)
        pc.OnTrack(func(remoteTrack *webrtc.Track, receiver *webrtc.RTPReceiver) {

                go func() {
                        ticker := time.NewTicker(rtcpPLIInterval)
                        for range ticker.C {
                                if err := pc.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: remoteTrack.SSRC()}}); err != nil {
                                        fmt.Println(err)
                                }
                        }
                }()
                //localTrack, err := pc.NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), "video", "pion")
                //if err != nil {
                //        fmt.Println("New Track err")
                //        panic(err)
                //}
                //localTrackChan <- localTrack
                //ongoingCalls[call.ID].Track = localTrack

                //TestTrack = localTrack

                //rtpBuf := make([]byte, 1400)

                for {
                        //i, err := remoteTrack.Read(rtpBuf)
                        packet, err := remoteTrack.ReadRTP()
                        if err != nil {
                                fmt.Println("Remote Track Read err")
                                panic(err)
                        }

                        packet.SSRC = outputTrack.SSRC()

                        //if _, err = outputTrack.Write(rtpBuf[:i]); err != nil && err != io.ErrClosedPipe {
                        if err := outputTrack.WriteRTP(packet); err != nil && err != io.ErrClosedPipe {
                                fmt.Println("Local Track Write err")
                                panic(err)
                        }
                }
                fmt.Println("after for loop")
        })

	//attachHandlers(pc)

	err = pc.SetRemoteDescription(offer); 
        if err != nil {
		pc.Close()
		return webrtc.SessionDescription{}, nil, err
	}
        
        

	ans, err := pc.CreateAnswer(nil)
	if err != nil {
		pc.Close()
		return webrtc.SessionDescription{}, nil, err
	}

	err = pc.SetLocalDescription(ans);
        if err != nil {
		pc.Close()
		return webrtc.SessionDescription{}, nil, err
	}
        
        fmt.Println("Adding call to list of calls")
        call.Connections[user] = pc
        ongoingCalls[call.ID] = call
        //ongoingCalls[call.ID].Track = <-localTrackChan
        fmt.Println("After remote and local set, after channel thing")
        
        return ans, call, nil
}

func GetCalls() ([]*Call, error) {
	var calls = make([]*Call, 0)
	for _, c := range ongoingCalls {
		calls = append(calls, c)
	}
	return calls, nil
}

func JoinCall(user string, offer webrtc.SessionDescription, id uuid.UUID) (webrtc.SessionDescription, error) {
	// if _, exists := ongoingCalls[id]; !exists {
	//	return nil, fmt.Errorf("unable to join call %v, call does not exist", id)
	// }

	pc, err := newPeerConnection()
	if err != nil {
		return webrtc.SessionDescription{}, err
	}

	// body := RTCSessionDescription{}
	// if err := json.Unmarshal([]byte(offer), &body); err != nil {
	//	return nil, fmt.Errorf("failed to unmarshal RTCSessionDescription")
	// }

	// if err := pc.SetRemoteDescription(body.Description); err != nil {
        if err := pc.SetRemoteDescription(offer); err != nil {
		pc.Close()
		return webrtc.SessionDescription{}, err
	}
        
        //need to grab local tracks and addTrack
        //localTrack := ongoingCalls[id].Track

        _, err = pc.AddTrack(TestTrack)
        if err != nil {
                pc.Close()
                return webrtc.SessionDescription{}, err
        }

	ans, err := pc.CreateAnswer(nil)
	if err != nil {
		pc.Close()
		return webrtc.SessionDescription{}, err
	}

	err = pc.SetLocalDescription(ans)
        if err != nil {
		pc.Close()
		return webrtc.SessionDescription{}, err
	}

	// ongoingCalls[id].Participants = append(ongoingCalls[id].Participants, user)
	// ongoingCalls[id].Connections[user] = pc

	return ans, nil
}

func LeaveCall(user string, id uuid.UUID) error {
	if _, exists := ongoingCalls[id]; !exists {
		return fmt.Errorf("unable to leave call %v, call does not exist", id)
	}

	return nil
}
