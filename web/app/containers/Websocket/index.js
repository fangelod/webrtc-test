import websocket from 'utils/websocket';
import { useEffect } from 'react';

const Websockets = () => {
  useEffect(() => {
    websocket.callback = msg => {
      console.log(msg);
    };
  }, []);

  return <div />;
};

export default Websockets;
