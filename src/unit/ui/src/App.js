import { NotificationContainer, NotificationManager } from 'react-notifications';
import React, { useState, useEffect, useRef } from 'react';
import Profile from './Profile/Profile'
import LogIn from './LogIn/LogIn'
import Menu from './Menu/Menu'
import PlayAudio from './PlayAudio/PlayAudio'
import { receivedCodes } from './codes'
// import net from 'net';
// import ipc from 'node-ipc';
// import { decode } from '@msgpack/msgpack';
import { baseURL } from './Api/Api'

import 'react-notifications/lib/notifications.css';
import './index.css';


function App() {
  const playAudioRef = useRef(null);
  // const AudioMsgType = 6 & 0xff;
  // const CodeMsgType = 7 & 0xff;
  // const UnixSocketPath = '/tmp/unit.hal.sock';
  const [audioModalName, setAudioModalName] = useState(null);
  const [audioModalData, setAudioModalData] = useState(null);
  // const [unixSocket, setSocket] = useState(null);
  const [msgs, appendToMsgs] = useState([]);
  const [audios, appendToAudios] = useState([]);

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
          sent: false,
          Body: code,
        };
        appendToMsgs(msg);
      }
      return selectedTab;
    })
  }

  var onReceiveAudio = (audio) => {
    setSelectedTab(selectedTab => {
      if (selectedTab !== "Log Out") {
        NotificationManager.info("Command Center sends audio message, click here to play it!", '', 3000, () => onPlayAudio("Command Center", audio), true);
        appendToAudios(audio);
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

      // var socket = ipc.connectTo('unit-hal', UnixSocketPath);
      // socket.on('message', (data) => {
      //   const parsedData = JSON.parse(decode(data).buffer);
      //   console.log("received msg: ", parsedData);
      //   if (parsedData.Type == CodeMsgType) {
      //     onReceiveMessage(parsedData);
      //   }
      //   else if (parsedData.Type == AudioMsgType) {
      //     onReceiveAudio(parsedData);
      //   }
      // });
      // setSocket(socket);
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
          <Profile msgs={msgs} audios={audios}/>
          : <> </>
      }
    </>
  );
}

export default App;
