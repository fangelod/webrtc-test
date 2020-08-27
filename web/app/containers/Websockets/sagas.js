import { SAVE_NAME } from 'App/constants';
import websocket from 'utils/websocket';

import { takeEvery } from 'redux-saga/effects';

function* saveNameWatcher() {
  yield takeEvery(SAVE_NAME, ({ payload }) => {
    websocket.send({ action: SAVE_NAME, name: payload });
  });
}

export default [
  saveNameWatcher,
]