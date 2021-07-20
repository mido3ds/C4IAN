import { NotificationContainer, NotificationManager } from 'react-notifications';
import React, { useState, useEffect, useRef } from 'react';
import Profile from './Profile/Profile'
import LogIn from './LogIn/LogIn'
import Menu from './Menu/Menu'
import PlayAudio from './PlayAudio/PlayAudio'
import { receivedCodes } from './codes'
import net from 'net';
import ipc from 'node-ipc';
import { decode } from '@msgpack/msgpack';

import 'react-notifications/lib/notifications.css';
import './index.css';


function App() {
  const playAudioRef = useRef(null);
  const AudioMsgType = 6 & 0xff
  const CodeMsgType = 7 & 0xff
  const [audioModalName, setAudioModalName] = useState(null)
  const [audioModalData, setAudioModalData] = useState(null)
  const [msgs, appendToMsgs] = useState([])
  const [audios, appendToAudios] = useState([])
  const [unixSocket, setSocket] = useState(null)

  const [selectedTab, setSelectedTab] = useState("Profile")

  var onPlayAudio = (name, data) => {
    setAudioModalName(name);
    setAudioModalData(data);
    playAudioRef.current.open()
  }

  var onReceiveMessage = (data) => {
    setSelectedTab(selectedTab => {
      if (selectedTab !== "Log Out") {
        NotificationManager.info("Command Center" + ": " + receivedCodes[data.Body]);
        data["sent"] = false;
        appendToMsgs(data)
      }
      return selectedTab
    })
  }

  var onReceiveAudio = (data) => {
    setSelectedTab(selectedTab => {
      if (selectedTab !== "Log Out") {
        NotificationManager.info("Command Center sends audio message, click here to play it!", '', 3000, () => onPlayAudio("Command Center", data.Body), true);
        appendToAudios(data)
      }
      return selectedTab
    })
  }

  useEffect(() => {
    window.$('.menu').css('visibility', 'visible')
    window.$('.menu .item span').each(function () { window.$(this).removeClass('selected') })
    window.$('.menu .item span')
      .filter(function (idx) { return this.innerHTML === selectedTab })
      .addClass('selected')

      var socket = net.connect('/tmp/unit.hal.sock');
      socket.on('data', (data) => {
        const parsedData = JSON.parse(decode(data).buffer);
        console.log("received msg: ", parsedData)
        if (parsedData.Type == CodeMsgType) {
          onReceiveMessage(parsedData)
        }
        else if (parsedData.Type == AudioMsgType) {
          onReceiveAudio(parsedData)
        }
      });
      setSocket(socket)
  }, [])

  var onChange = (selectedTab) => {
    setSelectedTab(selectedTab)

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
          <Profile socket={unixSocket} msgs={msgs} audios={audios}/>
          : <> </>
      }
    </>
  );
}

export default App;
