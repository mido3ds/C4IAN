import { NotificationContainer, NotificationManager } from 'react-notifications';
import React, { useState, useEffect, useRef } from 'react';
import Profile from './Profile/Profile'
import PlayAudio from './PlayAudio/PlayAudio'
import { receivedCodes } from './codes'
import { baseURL } from './Api/Api'

import 'react-notifications/lib/notifications.css';
import './index.css';


function App() {
  const playAudioRef = useRef(null);
  const [audioModalName, setAudioModalName] = useState(null);
  const [audioModalData, setAudioModalData] = useState(null);
  const [msgs, setMsgs] = useState([]);
  const [audios, setAudios] = useState([]);

  var onPlayAudio = (name, data) => {
    setAudioModalName(name);
    setAudioModalData(data);
    playAudioRef.current.open();
  }

  var onReceiveMessage = (code) => {
    NotificationManager.info("Command Center" + ": " + receivedCodes[code]);
    const msg = {
      Type: 7 & 0xff,
      sent: false,
      Body: code,
    };
    setMsgs((msgs)=> {
      return [...msgs, msg]
    });
  }

  var onReceiveAudio = (audio) => {
    NotificationManager.info("Command Center sends audio message, click here to play it!", '', 3006, () => onPlayAudio("Command Center", audio), true);
    setAudios((audios)=> {
      return [...audios, audio]
    });
  }

  useEffect(() => {
      var eventSource = new EventSource(baseURL+'events')
      eventSource.addEventListener("CODE-EVENT", ev => {
        onReceiveMessage(ev.data)
      })
  
      eventSource.addEventListener("AUDIO-EVENT", ev => {
        onReceiveAudio(JSON.parse(ev.data))
      })
  }, [])

  return (
    <>
      <NotificationContainer />
      <PlayAudio name={audioModalName} audio={audioModalData} ref={playAudioRef}></PlayAudio>
      <Profile msgs={msgs} audios={audios} setMsgs={setMsgs}/>
    </>
  );
}

export default App;
