import { getCalls, saveName, startCall, joinCall, renegotiateForce, leaveCall } from 'App/actions';
import { selectCalls, selectTracks, selectUser } from 'App/selectors';
import Websockets from 'containers/Websocket';

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
    <React.Fragment>
      <div id={'webcamDiv'}>
        Myself<br />
        <video id={'video1'} width={'160'} height={'120'} autoPlay muted></video><br />
      </div>
      <div id={'videoDiv'}>
        Others<br />
      </div>
      <div id={'buttons'}>
        <button className={'createSessionButton'} onClick={() => dispatch(startCall())}> Start Call </button><br />
        {calls.valueSeq().map(call => {
          return (
            <button className={'createSessionButton'} onClick={() => dispatch(joinCall(call.get('id'), ''))}> {call.get('id')} </button>
          );
        })} 
        <br /><br />
        <button className={'checkMediaStream'} onClick={() => console.log(document.getElementById('video2').srcObject)}> Check Stream </button>
        <button className={'renegotiate'} onClick={() => dispatch(renegotiateForce())}> Renegotiate </button>
        <button className={'leave call'} onClick={() => dispatch(leaveCall())}> Leave Call </button>
      </div>
      <Dialog disableBackdropClick={true} disableEscapeKeyDown={true} open={open}>
        <DialogContent>
          <DialogContentText>
            Choose a name to continue
          </DialogContentText>
          <TextField
            autoFocus={true}
            onChange={event => setName(event.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button
            className={classes.submit}
            disabled={name === ""}
            onClick={() => {
              dispatch(saveName(name));
              setOpen(false);
            }}
          >
            Submit
          </Button>
        </DialogActions>
      </Dialog>
    </React.Fragment>
  );
};
export default App;
