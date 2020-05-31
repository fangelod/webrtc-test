import {
  SAVE_NAME
} from 'App/constants';

import { fromJS } from 'immutable';

const initialState = fromJS({
  calls: {},
  name: '',
});

function appReducer(state = initialState, action) {
  if (!action || !action.type) {
    return state;
  }

  switch (action.type) {
    case SAVE_NAME:
      return state.set('name', action.payload);
    default:
      return state;
  }
}
export default appReducer;