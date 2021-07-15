import {NotificationContainer, NotificationManager} from 'react-notifications';
import React, { useState, useEffect, useRef } from 'react';
import Profile from './Profile/Profile'
import Home from './Home/Home'
import LogIn from './LogIn/LogIn'
import Menu from './Menu/Menu'
import PlayAudio from './PlayAudio/PlayAudio'
import Streams from './Streams/Streams'
import 'react-notifications/lib/notifications.css';

import './index.css';
import './App.css';

const tabsComponents = {
  "Map": <Home />,
  "Units": <Profile type="unit"/>,
}

function App() {
  const playAudioRef = useRef(null);
  const [audioModalName, setAudioModalName] = useState(null)
  const [audioModalData, setAudioModalData] = useState(null)

  const [selectedTab, setSelectedTab] = useState("Log Out")
  const [eventSource, setEventSource] = useState(new EventSource("http://localhost:3170/events"))

  var onPlayAudio = (name, data) => {
    setAudioModalName(name);
    setAudioModalData(data);
    playAudioRef.current.openModal()
  }

  useEffect(() => {
    eventSource.addEventListener("msg", ev => {
        var data = JSON.parse(ev.data)
        NotificationManager.info(data.src + ": " + data.code);
    })
    eventSource.addEventListener("audio", ev => {
        var data = JSON.parse(ev.data)
        NotificationManager.info(data.src + " sends audio message, click here to play it!", '' , 3000, () => onPlayAudio(data.src, data.body), true);
    })
  })

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
      <NotificationContainer/>
      <PlayAudio name={audioModalName} audioBolb={audioModalData} ref={playAudioRef}></PlayAudio>
      <Menu onChange={selectedTab => onChange(selectedTab)}> </Menu>
      <Home> </Home>
      {/*{selectedTab === "Log Out" ?
        <LogIn onLogIn={() => { onChange("Map") }} />
        : tabsComponents[selectedTab]
      }*/}
    </>
  );
}

export default App;
