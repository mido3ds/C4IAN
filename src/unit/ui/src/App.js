import { NotificationContainer, NotificationManager } from 'react-notifications';
import React, { useState, useEffect, useRef } from 'react';
import Profile from './Profile/Profile'
import PlayAudio from './PlayAudio/PlayAudio'
import GetPort from './GetPort/GetPort';
import { receivedCodes } from './codes'

import 'react-notifications/lib/notifications.css';
import './index.css';
import { getName } from './Api/Api';


function App() {
  const playAudioRef = useRef(null);
  const getPortRef = useRef(null);

  const [audioModalName, setAudioModalName] = useState(null);
  const [audioModalData, setAudioModalData] = useState(null);
  const [msgs, setMsgs] = useState([]);
  const [audios, setAudios] = useState([]);
  const [port, setPort] = useState(null)
  const [name, setName] = useState(null)

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
    setMsgs((msgs) => {
      return [...msgs, msg]
    });
  }

  var onReceiveAudio = (audio) => {
    NotificationManager.info("Command Center sends audio message, click here to play it!", '', 3006, () => onPlayAudio("Command Center", audio.body), true);
    setAudios((audios) => {
      return [...audios, audio]
    });
  }

  var onGetPort = (port) => {
    setPort(() => {
      var eventSource = new EventSource("http://localhost:" + port + "/events")
      eventSource.addEventListener("CODE-EVENT", ev => {
        onReceiveMessage(ev.data)
      })

      eventSource.addEventListener("AUDIO-EVENT", ev => {
        onReceiveAudio(JSON.parse(ev.data))
      })
      return port
    })
  }

  useEffect(() => {
    if (!port) getPortRef.current.open();
    else {
      getName(port).then(name => {
        setName(name.name);
        document.title = name.name;
      })
    }
  }, [port])

  return (
    <>
      <NotificationContainer />
      <GetPort onGetPort={onGetPort} ref={getPortRef}> </GetPort>
      <PlayAudio name={audioModalName} audio={audioModalData} ref={playAudioRef}></PlayAudio>

      {!port ? <> </> : <Profile name={name} port={port} msgs={msgs} audios={audios} setMsgs={setMsgs} />}
    </>
  );
}

export default App;
