import { List } from "immutable";

const selectCalls = state => {
  return state.getIn(['app', 'calls']);
};

const selectConnectionId = (state, connection) => {
  return state.getIn(['app', 'calls']).findKey(call => {
    return connection.localDescription.sdp === call.getIn(['connection']).localDescription.sdp;
  });
};

const selectCurrentCall = state => state.getIn(['app', 'currentCall']);

const selectTracks = state => {
  if (selectCurrentCall(state) === '') {
    return new List();
  }

  if (state.getIn(['app', 'calls', selectCurrentCall(state), 'streams'])) {
    return state.getIn(['app', 'calls', selectCurrentCall(state), 'streams']);
  }

  return new List();
};

const selectUser = state => state.getIn(['app', 'name']);

export {
  selectCalls,
  selectConnectionId,
  selectCurrentCall,
  selectTracks,
  selectUser
};