const websocket = {
  // The uri for the websocket
  uri: `${window.location.host}/websocket`,
  // Event prefix to listen to, if you have more than one instance of this class in your application you will want
  // these to be unique
  eventPrefix: '',

  callback: null,

  onMessage: function (_, e) {
    let msg = null;
    try {
      msg = JSON.parse(e.data);
    } catch (ex) {
      console.log(`Bad websocket message: ${ex}`);
    }
    if (msg !== null && this.callback) {
      // fire event
      this.callback.call(null, msg);
    }
  },

  // Optional
  onOpen: null,
  onClose: null,
  onError: null,

  // Whether the websocket is currently in an opened state
  isOpen: false,
  https: document.location.protocol === 'https:',

  // This can be configured
  retryPolicy: {
    retryNum: 1,
    _currentDelay: 5,
    scale: function () {
      return this._currentDelay * this.retryNum++;
    },
    reset: function () {
      this._currentDelay = 5;
      this.retryNum = 1;
    },
    max: null,
  },

  // The websocket itself
  socket: null,

  // Queue messages when the websocket is down
  messageQueue: {
    _data: [],
    flush: function (scope) {
      if (scope.isOpen && this._data.length) {
        console.log(
          `[${scope.eventPrefix}Websocket] flushing ${this._data
            .length} record(s).`
        );
        while (this._data.length) {
          const payload = this._data.shift();
          scope.socket.send(payload);
        }
      }
    },
    add: function (data) {
      this._data.push(data);
    },
  },

  send: function (payload) {
    const scope = this;
    if (!scope.isOpen) {
      if (typeof payload === 'object') {
        payload = JSON.stringify(payload);
      }
      scope.messageQueue.add(payload);
    } else {
      if (typeof payload === 'object') {
        payload = JSON.stringify(payload);
      }
      scope.socket.send(payload);
    }
  },

  open: function () {
    const scope = this;
    scope.socket = new WebSocket(`ws${scope.https ? 's' : ''}://${scope.uri}`);
    scope.setWSHandlers();
  },

  // Close the websocket
  close: function () {
    this.socket.close();
  },

  // Retry handler for the websocket
  retryWebsocket: function (policy) {
    const scope = this;
    if (policy.max !== null && policy.retryNum >= policy.max) {
      console.log(
        `[${scope.eventPrefix}Websocket] giving up reconnection after ${policy.retryNum} attempts.`
      );
      return;
    }
    const delay = policy.scale();
    console.log(
      `[${scope.eventPrefix}Websocket] attempting reconnection in ${delay} seconds.`
    );
    setTimeout(() => {
      scope.socket = new WebSocket(
        `ws${scope.https ? 's' : ''}://${scope.uri}`
      );
      scope.setWSHandlers();
    }, delay * 1000);
  },

  // Add handlers for the websocket
  setWSHandlers: function () {
    const scope = this;
    const socket = scope.socket;
    socket.onerror = function (e) {
      console.log(`[${scope.eventPrefix}Websocket] websocket error.`);
      if (typeof scope.onError === 'function') {
        scope.onError();
      }
    };
    socket.onopen = function () {
      console.log(`[${scope.eventPrefix}Websocket] connection open.`);
      scope.isOpen = true;
      scope.retryPolicy.reset();
      scope.messageQueue.flush(scope);
      if (typeof scope.onOpen === 'function') {
        scope.onOpen();
      }
    };
    socket.onclose = function () {
      console.log(`[${scope.eventPrefix}Websocket] connection closed.`);
      scope.isOpen = false;
      let retry = true;
      if (typeof scope.onClose === 'function') {
        retry = scope.onClose();
      }
      if (retry) {
        scope.retryWebsocket(scope.retryPolicy);
      }
    };
    socket.onmessage = scope.onMessage.bind(this, scope);
  },
};

websocket.socket = new WebSocket(
  `ws${websocket.https ? 's' : ''}://${websocket.uri}`
);
console.log(`Initializing websocket at ${websocket.socket.url}`);
websocket.isOpen = false;
websocket.setWSHandlers();

export default websocket;
