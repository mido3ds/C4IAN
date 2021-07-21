import { NotificationContainer, NotificationManager } from 'react-notifications';
import React, { useState, useEffect, useRef } from 'react';
import Profile from './Profile/Profile'
import LogIn from './LogIn/LogIn'
import Menu from './Menu/Menu'
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

  const [selectedTab, setSelectedTab] = useState("Profile");

  var onPlayAudio = (name, data) => {
    setAudioModalName(name);
    setAudioModalData(data);
    playAudioRef.current.open();
  }

  var onReceiveMessage = (code) => {
    setSelectedTab(selectedTab => {
      if (selectedTab !== "Log Out") {
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
      return selectedTab;
    })
  }

  var onReceiveAudio = (audio) => {
    setSelectedTab(selectedTab => {
      if (selectedTab !== "Log Out") {
        NotificationManager.info("Command Center sends audio message, click here to play it!", '', 3006, () => onPlayAudio("Command Center", audio), true);
        setAudios((audios)=> {
          return [...audios, audio]
        });
      }
      return selectedTab;
    })
  }

  useEffect(() => {
    window.$('.menu').css('visibility', 'visible')
    window.$('.menu .item span').each(function () { window.$(this).removeClass('selected') })
    window.$('.menu .item span')
      .filter(function (idx) { return this.innerHTML === selectedTab })
      .addClass('selected')

      var eventSource = new EventSource(baseURL+'events')
      eventSource.addEventListener("CODE-EVENT", ev => {
        onReceiveMessage(ev.data)
      })
  
      eventSource.addEventListener("AUDIO-EVENT", ev => {
        onReceiveAudio(ev.data)
      })
  }, [])

  var onChange = (selectedTab) => {
    setSelectedTab(selectedTab);

    if (selectedTab === "Log Out") {
      window.$('.menu').css('visibility', 'hidden')
    } else {
      window.$('.menu').css('visibility', 'visible')
      window.$('.menu .item span').each(function () { window.$(this).removeClass('selected') })

      window.$('.menu .item span')
        .filter(function (idx) { return this.innerHTML === selectedTab })
        .addClass('selected')
    }
  }

  return (
    <>
      <NotificationContainer />
      <PlayAudio name={audioModalName} audio={audioModalData} ref={playAudioRef}></PlayAudio>
      <Menu onChange={selectedTab => onChange(selectedTab)}> </Menu>
      {
        selectedTab === "Log Out" ?
          <LogIn />
          : selectedTab === "Profile" ?
          <Profile msgs={msgs} audios={audios} setMsgs={setMsgs}/>
          : <> </>
      }
    </>
  );
}

export default App;
