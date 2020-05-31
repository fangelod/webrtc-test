import {
  getCallsErr,
  getCallsSuccess,
  startCallErr,
  startCallSuccess
} from 'App/actions';
import { GET_CALLS, START_CALL } from 'App/constants';
import { selectUser } from 'App/selectors';
import request from 'utils/request';

import axios from 'axios';
import { call, put, select, take, takeEvery } from 'redux-saga/effects';
import { channel } from 'redux-saga';

const windowChannel = channel();

export function* doGetCalls() {
  console.log('doGetCalls');
  try {
    const calls = yield call(request, "/calls", { method: 'GET' });

    yield put(getCallsSuccess(calls));
  } catch (err) {
    yield put(getCallsErr(err));
  }
}

export function* doStartCall() {
  try {
    const servers = yield call(request, '/iceservers', { method: 'GET' });
    const constraints = { audio: true, video: true };
    const user = yield select(selectUser);

    navigator.mediaDevices.getUserMedia(constraints).then(stream => {
      const pc = new RTCPeerConnection({ iceServers: servers });
      pc.onicecandidate = e => {

      };
      pc.oniceconnectionstatechange = e => console.log(`ICE Connection State: ${pc.iceConnectionState}`);
      pc.ontrack = e => console.log(e);
      pc.onsignalingstatechange = ev => console.log(`Signaling State: ${ev.target.signalingState}`);
      pc.onnegotiationneeded = ev => console.log('Negotiation Needed');

      const tracks = stream.getTracks();
      tracks.forEach(track => pc.addTrack(track));

      pc.createOffer().then(desc => {
        pc.setLocalDescription(desc).then(() => {
          axios({
            url: '/calls',
            method: 'POST',
            data: {
              body: constraints,
              name: "",
              offer: desc,
              user: user
            }
          }).then(response => {
            if (response.status !== 200) {
              throw new Error(response.statusText);
            }

            pc.setRemoteDescription(response.data.sdp).then(() => {
              windowChannel.put(startCallSuccess(response.data.call));
            }).catch(e => console.error(e));
          }).catch(e => console.error(e));
        }).catch(e => console.error(e));
      }).catch(e => console.error(e))
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

export function* startCallWatcher() {
  yield takeEvery(START_CALL, doStartCall);
}

export default [
  channelWatcher,
  getCallsWatcher,
  startCallWatcher,
];