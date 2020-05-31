import App from 'App';
import appReducer from 'App/reducer';
import appSaga from 'App/sagas';

import { createMuiTheme, ThemeProvider } from '@material-ui/core/styles';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import { applyMiddleware, compose, createStore } from 'redux';
import { composeWithDevTools } from 'redux-devtools-extension';
import { combineReducers } from 'redux-immutable';
import createSagaMiddleware from 'redux-saga';
import { spawn } from 'redux-saga/effects';

const rootReducer = combineReducers({
  app: appReducer
});

const sagas = [
  ...appSaga
];

function* rootSaga() {
  yield sagas.map(saga => spawn(saga));
}

const composeEnhancers =
  process.env.NODE_ENV !== 'production'
    ? composeWithDevTools
    : compose;

const sagaMiddleware = createSagaMiddleware();

const store = createStore(
  rootReducer,
  composeEnhancers(applyMiddleware(sagaMiddleware))
);

sagaMiddleware.run(rootSaga);

const theme = createMuiTheme({});

export function render(config) {
  ReactDOM.render(
    <ThemeProvider theme={theme}>
      <Provider store={store}>
        <App />
      </Provider>
    </ThemeProvider>,
    document.querySelector(config.selector)
  );
}