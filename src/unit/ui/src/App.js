import { NotificationContainer, NotificationManager } from 'react-notifications';
import React, { useState, useEffect, useRef } from 'react';
import Profile from './Profile/Profile'
import LogIn from './LogIn/LogIn'
import Menu from './Menu/Menu'
import PlayAudio from './PlayAudio/PlayAudio'
import { receivedCodes } from './codes'
import { baseURL } from './Api/Api';

import 'react-notifications/lib/notifications.css';

import './index.css';


function App() {
  const playAudioRef = useRef(null);

  const [audioModalName, setAudioModalName] = useState(null)
  const [audioModalData, setAudioModalData] = useState(null)

  const [selectedTab, setSelectedTab] = useState("Profile")

  var onPlayAudio = (name, data) => {
    setAudioModalName(name);
    setAudioModalData(data);
    playAudioRef.current.open()
  }

  var onReceiveAudio = (data) => {
    setSelectedTab(selectedTab => {
      if (selectedTab !== "Log Out")
        NotificationManager.info(data.src + " sends audio message, click here to play it!", '', 3000, () => onPlayAudio("Command Center", data.body), true);
      return selectedTab
    })
  }

  var onReceiveMessage = (data) => {
    setSelectedTab(selectedTab => {
      if (selectedTab !== "Log Out")
        NotificationManager.info("Command Center" + ": " + receivedCodes[data.code]);
      return selectedTab
    })
  }

  useEffect(() => {
    window.$('.menu').css('visibility', 'visible')
    window.$('.menu .item span').each(function () { window.$(this).removeClass('selected') })
    window.$('.menu .item span')
      .filter(function (idx) { return this.innerHTML === selectedTab })
      .addClass('selected')

    var eventSource = new EventSource(baseURL+"events")
    eventSource.addEventListener("msg", ev => {
      onReceiveMessage(JSON.parse(ev.data))
    })

    eventSource.addEventListener("audio", ev => {
      onReceiveAudio(JSON.parse(ev.data))
    })
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
          <Profile />
          : <> </>
      }
    </>
  );
}

export default App;
