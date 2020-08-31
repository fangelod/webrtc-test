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
        Track2       *webrtc.Track                     `json:"-"`
        Track3       *webrtc.Track                     `json:"-"`
        Api          *webrtc.API                       `json:"-"`
        pc           *webrtc.PeerConnection            `json:"-"`
}

func init() {
	ongoingCalls = make(map[uuid.UUID]*Call)
}

func NewCall(offer webrtc.SessionDescription, name, user string) (webrtc.SessionDescription, uuid.UUID, error) {
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
                return webrtc.SessionDescription{}, call.ID, err
        }
        
        videoCodecs := mediaEngine.GetCodecsByKind(webrtc.RTPCodecTypeVideo)
        if len(videoCodecs) == 0 {
                fmt.Println("Offer had no video codecs");
        }

        api := webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine))

        config := webrtc.Configuration{ICEServers: rtcIceServers}
        pc, err := api.NewPeerConnection(config)
        if err != nil {
                return webrtc.SessionDescription{}, call.ID, err
        }

        if _, err = pc.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
                pc.Close()
                return webrtc.SessionDescription{}, call.ID, err
        }
        fmt.Println(videoCodecs[0].PayloadType)
        outputTrack, err := pc.NewTrack(videoCodecs[0].PayloadType, rand.Uint32(), "video", "pion")
        if err != nil {
                fmt.Println("NewTrack creation err")
                fmt.Println(err)
                pc.Close()
                return webrtc.SessionDescription{}, call.ID, err
        }
        
        //outputTrack3, err := pc.NewTrack(videoCodecs[0].PayloadType, rand.Uint32(), "video", "pion3")
        if err != nil {
                fmt.Println("outputTrack3 creation err")
                fmt.Println(err)
                pc.Close()
                return webrtc.SessionDescription{}, call.ID, err
        }
        //pc.AddTrack(outputTrack)
        //pc.AddTrack(outputTrack3)

        pc.OnTrack(func(remoteTrack *webrtc.Track, receiver *webrtc.RTPReceiver) {

                go func() {
                        ticker := time.NewTicker(rtcpPLIInterval)
                        for range ticker.C {
                                if err := pc.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: remoteTrack.SSRC()}}); err != nil {
                                        fmt.Println(err)
                                }
                        }
                }()

                for {
                        packet, err := remoteTrack.ReadRTP()
                        if err != nil {
                                fmt.Println("Remote Track Read err")
                                panic(err)
                        }

                        packet.SSRC = outputTrack.SSRC()

                        if err := outputTrack.WriteRTP(packet); err != nil && err != io.ErrClosedPipe {
                                fmt.Println("Local Track Write err")
                                panic(err)
                        }
                }
        })

	//attachHandlers(pc)


        //pc.OnNegotiationNeeded(func() {
        //        fmt.Println("NEGOTIATION NEEDED")
        //})

	err = pc.SetRemoteDescription(offer); 
        if err != nil {
		pc.Close()
		return webrtc.SessionDescription{}, call.ID, err
	}
        
        

	ans, err := pc.CreateAnswer(nil)
	if err != nil {
		pc.Close()
		return webrtc.SessionDescription{}, call.ID, err
	}

	err = pc.SetLocalDescription(ans);
        if err != nil {
		pc.Close()
		return webrtc.SessionDescription{}, call.ID, err
	}
        
        call.Connections[user] = pc
        ongoingCalls[call.ID] = call
        ongoingCalls[call.ID].Api = api
        ongoingCalls[call.ID].Track = outputTrack
        //ongoingCalls[call.ID].Track3 = outputTrack3
        ongoingCalls[call.ID].pc = pc
        fmt.Println("outputTrack put in ongoingCalls")
        
        return ans, call.ID, nil
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

	//pc, err := newPeerConnection()
	//if err != nil {
	//	return webrtc.SessionDescription{}, err
	//}

	// body := RTCSessionDescription{}
	// if err := json.Unmarshal([]byte(offer), &body); err != nil {
	//	return nil, fmt.Errorf("failed to unmarshal RTCSessionDescription")
	// }

	
        config := webrtc.Configuration{ICEServers: rtcIceServers}
        pc, err := ongoingCalls[id].Api.NewPeerConnection(config)
        if err != nil {
                return webrtc.SessionDescription{}, err
        }
        
        //_, err = pc.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo)
        if err != nil {
                fmt.Println("Err adding transceiver")
                pc.Close()
                return webrtc.SessionDescription{}, err
        }

        //need to grab local tracks and addTracks to connection
        serverTrack := ongoingCalls[id].Track
        
        _, err = pc.AddTrack(serverTrack)
        if err != nil {
                pc.Close()
                return webrtc.SessionDescription{}, err
        }

        //take remote track and make it a local track 
        //96 is the videocodec payloadType
        outputTrack2, err := pc.NewTrack(96, rand.Uint32(), "video", "pion2")
        if err != nil {
                fmt.Println("NewTrack creation err")
                fmt.Println(err)
                pc.Close()
                return webrtc.SessionDescription{}, err
        }
        
        //_, err = pc.AddTrack(outputTrack2)
        //_, err = pc.AddTrack(ongoingCalls[id].Track3)
        if err != nil {
                fmt.Println("outputTrack2 add err")
        }

        //_, err = ongoingCalls[id].pc.AddTrack(outputTrack2)
        if err != nil {
                fmt.Println("addTrack err")
                pc.Close()
                return webrtc.SessionDescription{}, err
        }

        //answer, err := ongoingCalls[id].pc.CreateAnswer(nil)
        if err != nil {
                fmt.Println("Err creating answer")
        }

        //err = ongoingCalls[id].pc.SetLocalDescription(answer)
        if err != nil {
                fmt.Println("Err setting local description for caller one")
        }

        pc.OnTrack(func(remoteTrack2 *webrtc.Track, receiver2 *webrtc.RTPReceiver) {

                go func() {
                        ticker2 := time.NewTicker(rtcpPLIInterval)
                        for range ticker2.C {
                                if err := pc.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: remoteTrack2.SSRC()}}); err != nil {
                                        fmt.Println(err)
                                }
                        }
                }()

                for {
                        packet2, err := remoteTrack2.ReadRTP()
                        if err != nil {
                                fmt.Println("Remote Track Read err")
                                panic(err)
                        }


                        //packet2.SSRC = ongoingCalls[id].Track3.SSRC()
                        packet2.SSRC = outputTrack2.SSRC()


                        //if err := ongoingCalls[id].Track3.WriteRTP(packet2); err != nil && err != io.ErrClosedPipe {
                        if err := outputTrack2.WriteRTP(packet2); err != nil && err != io.ErrClosedPipe {
                                fmt.Println("Local Track Write err")
                                panic(err)
                        }
                }
        })

        // if err := pc.SetRemoteDescription(body.Description); err != nil {
        if err := pc.SetRemoteDescription(offer); err != nil {
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

        ongoingCalls[id].Track2 = outputTrack2
	// ongoingCalls[id].Participants = append(ongoingCalls[id].Participants, user)
	// ongoingCalls[id].Connections[user] = pc
        fmt.Println("Join Call Success")

	return ans, nil
}

func LeaveCall(user string, id uuid.UUID) error {
	if _, exists := ongoingCalls[id]; !exists {
		return fmt.Errorf("unable to leave call %v, call does not exist", id)
	}

	return nil
}

func RenegotiateCall(id uuid.UUID) (webrtc.SessionDescription, error) {
        pc := ongoingCalls[id].pc

        _, err := pc.AddTrack(ongoingCalls[id].Track2)
        if err != nil {
                fmt.Println("Err adding track")
                pc.Close()
                return webrtc.SessionDescription{}, err
        }
        
	ans, err := pc.CreateAnswer(nil)
	if err != nil {
                fmt.Println("Err creating answer")
		pc.Close()
		return webrtc.SessionDescription{}, err
	}

	err = pc.SetLocalDescription(ans)
        if err != nil {
                fmt.Println("Err setting local description")
		pc.Close()
		return webrtc.SessionDescription{}, err
	}

        return ans, nil
}
