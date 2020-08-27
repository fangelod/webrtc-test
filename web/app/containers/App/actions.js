import {
  GET_CALLS,
  GET_CALLS_FAILURE,
  GET_CALLS_SUCCESS,
  ICE_CANDIDATE,
  ICE_CANDIDATE_FAILURE,
  ICE_CANDIDATE_SUCCESS,
  JOIN_CALL,
  JOIN_CALL_FAILURE,
  JOIN_CALL_SUCCESS,
  LEAVE_CALL,
  LEAVE_CALL_FAILURE,
  LEAVE_CALL_SUCCESS,
  RENEGOTIATE,
  RENEGOTIATE_FAILURE,
  RENEGOTIATE_SUCCESS,
  SAVE_NAME,
  SAVE_TRACK,
  SAVE_TRACK_FAILURE,
  SAVE_TRACK_SUCCESS,
  START_CALL,
  START_CALL_FAILURE,
  START_CALL_SUCCESS
} from 'App/constants';

export const getCalls = () => ({ type: GET_CALLS });

export const getCallsErr = err => ({
  type: GET_CALLS_FAILURE,
  payload: err
});

export const getCallsSuccess = calls => ({
  type: GET_CALLS_SUCCESS,
  payload: calls
});

export const iceCandidate = (connection, candidate) => ({
  type: ICE_CANDIDATE,
  payload: {
    connection: connection,
    candidate: candidate
  }
});

export const iceCandidateErr = err => ({
  type: ICE_CANDIDATE_FAILURE,
  payload: err
});

export const iceCandidateSuccess = () => ({ type: ICE_CANDIDATE_SUCCESS });

export const joinCall = (id, user) => ({
  type: JOIN_CALL,
  payload: id
});

export const joinCallErr = (id, err) => ({
  type: JOIN_CALL_FAILURE,
  payload: err,
  meta: id
});

export const joinCallSuccess = (id, answer) => ({
  type: JOIN_CALL_SUCCESS,
  payload: {
    id: id,
    answer: answer
  }
});

export const leaveCall = id => ({ type: LEAVE_CALL, payload: id });

export const leaveCallErr = (id, err) => ({
  type: LEAVE_CALL_FAILURE,
  payload: err,
  meta: id
});

export const leaveCallSuccess = id => ({
  type: LEAVE_CALL_SUCCESS,
  payload: id
});

export const renegotiate = id => ({ type: RENEGOTIATE, payload: id });

export const renegotiateErr = (id, err) => ({
  type: RENEGOTIATE_FAILURE,
  payload: err,
  meta: id
});

export const renegotiateForce = () => ({
  type: RENEGOTIATE
});
        

export const renegotiateSuccess = (id, answer) => ({
  type: RENEGOTIATE_SUCCESS,
  payload: {
    id: id,
    answer: answer
  }
});

export const saveName = name => ({ type: SAVE_NAME, payload: name });

export const saveTrack = (connection, stream) => ({
  type: SAVE_TRACK,
  payload: {
    connection: connection,
    stream: stream
  }
});

export const saveTrackErr = err => ({
  type: SAVE_TRACK_FAILURE,
  payload: err
});

export const saveTrackSuccess = (id, stream) => ({
  type: SAVE_TRACK_SUCCESS,
  payload: {
    id: id,
    stream: stream
  }
});

export const startCall = () => ({ type: START_CALL });

export const startCallErr = err => ({ type: START_CALL_FAILURE, payload: err });

export const startCallSuccess = (callId, peerConnection) => ({
  type: START_CALL_SUCCESS,
  payload: {
    id: callId,
    peerConnection: peerConnection
  }
});
