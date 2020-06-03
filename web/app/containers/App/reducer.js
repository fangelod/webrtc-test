import {
  GET_CALLS_SUCCESS,
  SAVE_NAME,
  SAVE_TRACK_SUCCESS,
  START_CALL_SUCCESS
} from 'App/constants';

import { fromJS, List } from 'immutable';

const initialState = fromJS({
  calls: {},
  currentCall: '',
  name: '',
});

function appReducer(state = initialState, action) {
  if (!action || !action.type) {
    return state;
  }

  switch (action.type) {
    case GET_CALLS_SUCCESS:
      return state.set('calls', fromJS(action.payload));

    case SAVE_NAME:
      return state.set('name', action.payload);

    case SAVE_TRACK_SUCCESS:
      if (!state.getIn(['calls', action.payload.id, 'streams'])) {
        const newList = new List();
        return state.setIn(['calls', action.payload.id, 'streams'], newList.push(action.payload.stream));
      }
      return state.setIn(
        ['calls', action.payload.id, 'streams'],
        state.getIn(['calls', action.payload.id, 'streams']).push(action.payload.stream)
      );

    case START_CALL_SUCCESS:
      return state
        .setIn(['calls', action.payload.id], fromJS(action.payload))
        .set('currentCall', action.payload.id);

    default:
      return state;
  }
}
export default appReducer;