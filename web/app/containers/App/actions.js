import {
  GET_CALLS,
  GET_CALLS_FAILURE,
  GET_CALLS_SUCCESS,
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

export const renegotiateSuccess = (id, answer) => ({
  type: RENEGOTIATE_SUCCESS,
  payload: {
    id: id,
    answer: answer
  }
});

export const saveName = name => ({ type: SAVE_NAME, payload: name });

export const startCall = () => ({ type: START_CALL });

export const startCallErr = err => ({ type: START_CALL_FAILURE, payload: err });

export const startCallSuccess = (call, answer) => ({
  type: START_CALL_SUCCESS,
  payload: {
    call: call,
    answer: answer
  }
});