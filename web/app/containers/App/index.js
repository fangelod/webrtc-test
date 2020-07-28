import { getCalls, saveName, startCall } from 'App/actions';
import { selectCalls, selectTracks, selectUser } from 'App/selectors';

import {
  AppBar,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  IconButton,
  List,
  ListItem,
  ListItemSecondaryAction,
  ListItemText,
  Paper,
  TextField,
  Typography,
} from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { AddIcCall } from '@material-ui/icons';

const useStyles = makeStyles(theme => {
  return {
    mainDiv: {
      width: '100%',
      height: '100%',
      display: 'flex',
      flexDirection: 'column',
    },
    appBar: {
      height: '48px',
      alignItems: 'center',
      justifyContent: 'space-between',
      flexDirection: 'row',
    },
    title: {
      marginLeft: '10px',
    },
    content: {
      flexGrow: 1,
      display: 'flex',
      flexDirection: 'row',
      alignItems: 'stretch',
    },
    submit: {
      background: props => {
        if (props.name === "") {
          return theme.palette.text.disabled;
        }

        return theme.palette.primary.main;
      },
      color: props => {
        if (props.name === "") {
          return theme.palette.common.black;
        }
        return `${theme.palette.common.white} !important`
      }
    },
    calls: {
      background: 'red',
      width: '250px',
    },
    streams: {
      background: 'blue',
      flex: 1,
    }
  };
});

const App = () => {
  const [name, setName] = useState("");
  const [open, setOpen] = useState(true);
  const classes = useStyles({ name: name });
  const dispatch = useDispatch();
  const calls = useSelector(selectCalls);
  const streams = useSelector(selectTracks);
  const user = useSelector(selectUser);

  useEffect(() => {
    dispatch(getCalls());
  }, []);

  return (
    <div id={'content'}>
      <div id={'signalingContainer'} style={{ display: 'none' }}>
        Browser base64 Session Description<br />
        <textarea id={'localSessionDescription'} readOnly={true}></textarea> <br /> 
        Golang base64 Session Description<br />
        <textarea id={'remoteSessionDescription'}></textarea> <br />
        <button onClick={() => console.log('heh')}> Start Session </button> <br />
      </div> <br />
      Video<br />
      <video id={'video1'} width={'160'} height={'120'} autoplay muted></video><br />
      <button class={'createSessionButton'} onClick={() => dispatch(startCall())}> Publish a Broadcast </button>
      <button class={'createSessionButton'} onClick={() => dispatch(joinCall())}> Join a Broadcast </button><br /><br />
      Logs<br />
      <div id={'logs'}></div>
    </div>
  );
};
export default App;
