import {
  getCallsErr,
  getCallsSuccess,
  iceCandidate,
  iceCandidateErr,
  iceCandidateSuccess,
  joinCallSuccess,
  renegotiate,
  renegotiateSuccess,
  saveTrack,
  saveTrackErr,
  saveTrackSuccess,
  startCallErr,
  startCallSuccess
} from 'App/actions';
import {
  GET_CALLS,
  ICE_CANDIDATE,
  JOIN_CALL,
  RENEGOTIATE,
  SAVE_TRACK,
  START_CALL
} from 'App/constants';
import { selectCallId, selectPeerConnection, selectConnectionId, selectUser } from 'App/selectors';
import request from 'utils/request';

import axios from 'axios';
import {
  call,
  put,
  select,
  spawn,
  take,
  takeEvery
} from 'redux-saga/effects';
import delay from '@redux-saga/delay-p';
import { channel } from 'redux-saga';

const windowChannel = channel();

export function* doGetCalls() {
  try {
    const calls = yield call(request, "/calls", { method: 'GET' });

    yield put(getCallsSuccess(calls));
  } catch (err) {
    yield put(getCallsErr(err));
  }
}

export function* retryIceCandidate(payload) {
  yield delay(1000);
  yield put(iceCandidate(payload.connection, payload.candidate));
}

export function* doIceCandidate(action) {
  try {
    const id = yield select(selectConnectionId, action.payload.connection);
    //const user = yield select(selectUser);

    if (!id) {
      // Call hasn't been saved in the state yet requeue candidate
      yield spawn(retryIceCandidate, action.payload);
      return;
    }

    const requestURL = `/calls/${id}/icecandidate`;
    yield call(request, requestURL, {
      method: 'POST',
      data: {
        candidate: JSON.stringify(action.payload.candidate),
        user: user,
      }
    });

    yield put(iceCandidateSuccess());
  } catch (err) {
    console.error(err);
    yield put(iceCandidateErr(err));
  }
}

export function* retrySaveTrack(payload) {
  yield delay(1000);
  yield put(saveTrack(payload.connection, payload.stream));
}

export function* doSaveTrack(action) {
  try {
    const id = yield select(selectConnectionId, action.payload.connection);

    if (!id) {
      // Call hasn't been saved in the state yet requeue stream
      yield spawn(retrySaveTrack, action.payload);
      return;
    }

    yield put(saveTrackSuccess(id, action.payload.stream));
  } catch (err) {
    yield put(saveTrackErr(err));
  }
}

function* doJoinCall(action) {
  try {
    const servers = yield call(request, '/iceservers', { method: 'GET' });
    const pc = new RTCPeerConnection({ iceServers: servers });
    const user = yield select(selectUser);
    const constraints = { audio: false, video: true };

    pc.onicecandidate = e => {
      if (e.candidate) {
        //windowChannel.put(iceCandidate(pc, e.candidate));
      }
    };
    pc.oniceconnectionstatechange = e => console.log(`ICE Connection State: ${pc.iceConnectionState}`);
    pc.onsignalingstatechange = ev => console.log(`Signaling State: ${ev.target.signalingState}`);
    pc.onnegotiationneeded = ev => console.log('Negotiation Needed', pc);

    pc.ontrack = function (event) {
      console.log("Got a track event", event);

      //Generate new video player every time we see a track
      var newVideo = document.createElement("video")
      newVideo.srcObject = event.streams[0];
      newVideo.autoplay = true;
      newVideo.controls = true;
      newVideo.width = '160';
      newVideo.height = '120';

      document.getElementById("videoDiv").appendChild(newVideo);
      
      //var el = document.getElementById('video2');
      //el.srcObject = event.streams[0];
      //el.autoplay = true;
      //el.controls = true;
    };

    navigator.mediaDevices.getUserMedia(constraints).then(stream => {
      var el = document.getElementById("video1");
      el.srcObject = stream;
      el.autoplay = true;
      el.controls = true;

      stream.getTracks().forEach(track => pc.addTrack(track, stream));
    }).catch(e => console.error(e));
    
    pc.addTransceiver('video');

    console.log("Receivers");
    console.log(pc.getReceivers());
      
    console.log("Senders");
    console.log(pc.getSenders());
          
    pc.createOffer().then(desc => {
      pc.setLocalDescription(desc).then(() => {
        axios({
          url: `/calls/${action.payload}/join`,
          method: 'POST',
          data: {
            id: action.payload,
            offer: desc,
            user: user
          }
        }).then(response => {
          if (response.status !== 200) {
            throw new Error(response.statusText);
          }
          pc.setRemoteDescription(response.data.sdp).then(() => {
            windowChannel.put(joinCallSuccess(action.payload, pc));
          }).catch(e => console.error(e));
        }).catch(e => console.error(e));
      }).catch(e => console.error(e));
    }).catch(e => console.error(e));
  } catch (err) {
    console.error(err);
    yield put(joinCallErr(action.payload, err));
  }
}

function* doRenegotiateCall() {
  try {
    console.log('Renegotiating');
    const pc = yield select(selectPeerConnection);
    const callId = yield select(selectCallId);
    const user = yield select(selectUser);
    pc.addTransceiver('video');
    console.log("Receivers");
    console.log(pc.getReceivers());
      
    console.log("Senders");
    console.log(pc.getSenders());
    // Create offer in order to change connection state to allow for the new answer      
    pc.createOffer().then(desc => {
      pc.setLocalDescription(desc).then(() => {
        axios({
          url: `/calls/${callId}/renegotiate`,
          method: 'POST',
          data: {
            id: callId,
            user: user,
            offer: desc,
          }
        }).then(response => {
          console.log("Got Remote offer back");
          pc.setRemoteDescription(response.data.sdp).then(() => {
            windowChannel.put(renegotiateSuccess(pc));
          }).catch(e => console.error(e));
        }).catch(e => console.error(e));
      }).catch(e => console.error(e));
    }).catch(e => console.error(e));
  } catch (err) {
    console.error(err);
  }
}


function* doStartCall() {
  try {
    const servers = yield call(request, '/iceservers', { method: 'GET' });
    const constraints = { audio: false, video: true };
    const user = yield select(selectUser);

    navigator.mediaDevices.getUserMedia(constraints).then(stream => {
      const pc = new RTCPeerConnection({ iceServers: servers });
      pc.onicecandidate = e => {
        console.log("Ice Candidate", e);
        if (e.candidate) {
         //windowChannel.put(iceCandidate(pc, e.candidate));
        }
      };
      pc.oniceconnectionstatechange = e => console.log(`ICE Connection State: ${pc.iceConnectionState}`);
      pc.onsignalingstatechange = ev => console.log(`Signaling State: ${ev.target.signalingState}`);
      pc.onnegotiationneeded = ev => console.log('Negotiation Needed', ev);

      //Add track to the connection with pc.addTrack() / .getVideoTracks()[]
      //pc.addTrack(stream.getVideoTracks()[0]);
      stream.getTracks().forEach(track => pc.addTrack(track, stream));        
      
      //Old
      //pc.addStream(document.getElementById('video1').srcObject = stream);
      
      //sets broadcaster stream to first video player
      document.getElementById('video1').srcObject = stream;
      
      //pc.addTransceiver('video');
      
      pc.ontrack = (event) => {
        console.log("Got a track event", event);
        //Generate new video player every time we see a track
        var newVideo = document.createElement("video")
        newVideo.srcObject = event.streams[0];
        newVideo.autoplay = true;
        newVideo.controls = true;
        newVideo.width = '160';
        newVideo.height = '120';

        document.getElementById("videoDiv").appendChild(newVideo);

        //sets stream from server to the second video player
        //var element = document.getElementById('video2');
        //element.srcObject = event.streams[0];
        //element.autoplay = true;
        //element.controls = true;
      };

      console.log("Receivers");
      console.log(pc.getReceivers());
      
      console.log("Senders");
      console.log(pc.getSenders());

      pc.createOffer().then(desc => {
        pc.setLocalDescription(desc).then(() => {
          axios({
            url: '/calls',
            method: 'POST',
            data: {
              body: constraints,
              offer: desc,
              user: user
            }
          }).then(response => {
            if (response.status !== 200) {
              throw new Error(response.statusText);
            }

            pc.setRemoteDescription(response.data.sdp).then(() => {
              let callId = response.data.callId;
              console.log(callId);
              console.log(pc);
              windowChannel.put(startCallSuccess(callId, pc));
            }).catch(e => console.error(e));
          }).catch(e => console.error(e));
        }).catch(e => console.error(e));
      }).catch(e => console.error(e));
    }).catch(e => console.error(e));
  } catch (err) {
    console.error(err);
    yield put(startCallErr(err));
  }
}



export function* channelWatcher() {
  while (windowChannel) {
    const action = yield take(windowChannel);

    yield put(action);
  }
}

export function* getCallsWatcher() {
  yield takeEvery(GET_CALLS, doGetCalls);
}

export function* iceCandidateWatcher() {
  yield takeEvery(ICE_CANDIDATE, doIceCandidate);
}

export function* joinCallWatcher() {
  yield takeEvery(JOIN_CALL, doJoinCall);
}

export function* renegotiateWatcher() {
  yield takeEvery(RENEGOTIATE, doRenegotiateCall);
}

export function* saveTrackWatcher() {
  yield takeEvery(SAVE_TRACK, doSaveTrack);
}

export function* startCallWatcher() {
  yield takeEvery(START_CALL, doStartCall);
}

export default [
  channelWatcher,
  getCallsWatcher,
  iceCandidateWatcher,
  joinCallWatcher,
  renegotiateWatcher,
  saveTrackWatcher,
  startCallWatcher,
];
