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
        rtcpPLIInterval = time.Second *1
)

var (
	ongoingCalls map[uuid.UUID]*Call
)

type RTCSessionDescription struct {
	Description webrtc.SessionDescription `json:"offer"`
	User        string                    `json:"user"`
	Type        string                    `json:"type"`
}

type Call struct {
	ID           uuid.UUID                         `json:"id"`
	Participants map[string]*User                  `json:"participants"`
        Api          *webrtc.API                       `json:"-"`
}

type User struct {
        Name         string                            `json:"name"`
        AudioTrack   *webrtc.Track                     `json:"-"`
        VideoTrack   *webrtc.Track                     `json:"-"`
        PC           *webrtc.PeerConnection            `json:"-"`
        //HasTrackFrom
        HTF          []string                          `json:"-"`
}

func init() {
	ongoingCalls = make(map[uuid.UUID]*Call)
}

func NewCall(offer webrtc.SessionDescription, username string) (webrtc.SessionDescription, uuid.UUID, error) {
        call := &Call{
		ID:           uuid.NewV4(),
                Participants: make(map[string]*User),
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

        outputTrack, err := pc.NewTrack(videoCodecs[0].PayloadType, rand.Uint32(), "video", username)
        if err != nil {
                fmt.Println("NewTrack creation err")
                fmt.Println(err)
                pc.Close()
                return webrtc.SessionDescription{}, call.ID, err
        }
        
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

        user := &User {
                Name:   username,
                VideoTrack:     outputTrack,
                PC:     pc,
                HTF:    []string{username},
        }

        call.Participants[username] = user
        ongoingCalls[call.ID] = call
        ongoingCalls[call.ID].Api = api
        
        return ans, call.ID, nil
}

func GetCalls() ([]*Call, error) {
	var calls = make([]*Call, 0)
	for _, c := range ongoingCalls {
		calls = append(calls, c)
	}
	return calls, nil
}

func JoinCall(offer webrtc.SessionDescription, id uuid.UUID, username string) (webrtc.SessionDescription, error) {
	// if _, exists := ongoingCalls[id]; !exists {
	//	return nil, fmt.Errorf("unable to join call %v, call does not exist", id)
	// }

	//pc, err := newPeerConnection()
	//if err != nil {
	//	return webrtc.SessionDescription{}, err
	//}
        
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
        
        var trackName string
        // Grab tracks from all call participants and add to connection !!Can't do this unless front end knows how many tracks are coming back so it added the appropriate number of recievers!!
        for _, participant := range ongoingCalls[id].Participants {
                _, err = pc.AddTrack(participant.VideoTrack)
                if participant.Name == "Brent" {
                        fmt.Println("Hello There")
                }
                if err != nil {
                        pc.Close()
                        return webrtc.SessionDescription{}, err
                }
                //For now just grabbing one track and adding to list of tracks already grabbed.
                trackName = participant.Name
                break
        }
        fmt.Println("Before Track Creation")
        // Take remote track and make it a local track 
        //96 is the videocodec payloadType
        outputTrack, err := pc.NewTrack(96, rand.Uint32(), "video", username)
        if err != nil {
                fmt.Println("NewTrack creation err")
                fmt.Println(err)
                pc.Close()
                return webrtc.SessionDescription{}, err
        }

        //pc.AddTrack(outputTrack)

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

        if err := pc.SetRemoteDescription(offer); err != nil {
		fmt.Println("Err SRD")
                pc.Close()
		return webrtc.SessionDescription{}, err
	}

	ans, err := pc.CreateAnswer(nil)
	if err != nil {
                fmt.Println("Err CA")
		pc.Close()
		return webrtc.SessionDescription{}, err
	}
        
	err = pc.SetLocalDescription(ans)
        if err != nil {
                fmt.Println("Err SLD")
		pc.Close()
		return webrtc.SessionDescription{}, err
	}


        user := &User {
                Name:   username,
                VideoTrack:     outputTrack,
                PC:     pc,
                HTF:    []string{username, trackName},
        }

	ongoingCalls[id].Participants[username] = user
        fmt.Println("Call Joined")
	return ans, nil
}

func LeaveCall(username string, id uuid.UUID) error {
	if _, exists := ongoingCalls[id]; !exists {
		return fmt.Errorf("unable to leave call %v, call does not exist", id)
	}

	return nil
}

func RenegotiateCall(offer webrtc.SessionDescription, id uuid.UUID, username string) (webrtc.SessionDescription, error) {
        //Needs exists check

        pc := ongoingCalls[id].Participants[username].PC
        user := ongoingCalls[id].Participants[username]
        //!!Can't do this unless front end knows how many tracks to expect back beforehand!!
        //Stop gap fix is a collection of tracks already grabbed and checking against list to get tracks user doesnt have, has to be one at a time due to above issue
        for _, participant := range ongoingCalls[id].Participants {
                //Need Unique key for participants to avoid same name problems
                // Skip adding the current user's tracks to their own connection
                //currently adds all the tracks but really just need the new tracks/ tracks user doesn't have
                if participant.Name == username {
                        fmt.Println("General Kenobi")
                        continue
                }
                var needsTrack = true
                for _, name := range user.HTF {
                        if name == participant.Name {
                                needsTrack = false
                        }
                }
                if needsTrack {                     
                        _, err := pc.AddTrack(participant.VideoTrack)
                        if err != nil {
                                pc.Close()
                                return webrtc.SessionDescription{}, err
                        }
                        ongoingCalls[id].Participants[username].HTF = append(ongoingCalls[id].Participants[username].HTF, participant.Name)
                }
        }
        
        if err := pc.SetRemoteDescription(offer); err != nil {
		fmt.Println("Err SRD in REN")
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
		pc.Close()
		return webrtc.SessionDescription{}, err
	}

        return ans, nil
}
