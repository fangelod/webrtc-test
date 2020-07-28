import {
  getCallsErr,
  getCallsSuccess,
  iceCandidate,
  iceCandidateErr,
  iceCandidateSuccess,
  saveTrack,
  saveTrackErr,
  saveTrackSuccess,
  startCallErr,
  startCallSuccess
} from 'App/actions';
import {
  GET_CALLS,
  ICE_CANDIDATE,
  SAVE_TRACK,
  START_CALL
} from 'App/constants';
import { selectConnectionId, selectUser } from 'App/selectors';
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
    const user = yield select(selectUser);

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

export function* doStartCall() {
  try {
    const servers = yield call(request, '/iceservers', { method: 'GET' });
    const constraints = { audio: false, video: true };
    const user = yield select(selectUser);

    navigator.mediaDevices.getUserMedia(constraints).then(stream => {
      const pc = new RTCPeerConnection({ iceServers: servers });
      pc.onicecandidate = e => {
        if (e.candidate) {
          windowChannel.put(iceCandidate(pc, e.candidate));
        }
      };
      pc.oniceconnectionstatechange = e => console.log(`ICE Connection State: ${pc.iceConnectionState}`);
      pc.ontrack = e => {
        if (e.streams) {
          e.streams.forEach(s => windowChannel.put(saveTrack(pc, s)));
        }
      };
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
              let call = response.data.call;
              call.connection = pc;
              windowChannel.put(startCallSuccess(call));
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

export function* iceCandidateWatcher() {
  yield takeEvery(ICE_CANDIDATE, doIceCandidate);
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
  saveTrackWatcher,
  startCallWatcher,
];
