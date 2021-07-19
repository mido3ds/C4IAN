import { NotificationContainer, NotificationManager } from 'react-notifications';
import React, { useState, useEffect, useRef } from 'react';
import Profile from './Profile/Profile'
import Home from './Home/Home'
import LogIn from './LogIn/LogIn'
import Menu from './Menu/Menu'
import PlayAudio from './PlayAudio/PlayAudio'
import Streams from './Streams/Streams'
import { receivedCodes } from './codes'
import { postMsg } from './Api/Api';

import 'react-notifications/lib/notifications.css';

import './index.css';
import './App.css';


function App() {
  const playAudioRef = useRef(null);

  const [audioModalName, setAudioModalName] = useState(null)
  const [audioModalData, setAudioModalData] = useState(null)

  const [selectedTab, setSelectedTab] = useState("Streams")

  const [streamsTime, setStreamsTimer] = useState(null)
  const [streams, setStreams] = useState([])

  var onPlayAudio = (name, data) => {
    setAudioModalName(name);
    setAudioModalData(data);
    playAudioRef.current.open()
  }

  var onEndStreamRequest = (data) => {
    console.log("Hello")
    setStreams(streams => {
      streams = streams.filter(stream => {
        return stream.id !== data.id
      })
      return streams
    })

    postMsg(data.src, { code: 3, })
  }

  var perodicStartStream = (data) => {
    // Resend start stream request
    postMsg(data.src, { code: 2, })
    if (streams.some(e => e.ID === data.ID)) {
      setTimeout(() => {
        perodicStartStream(data)
      }, 50 * 1000)
    }
  }

  var perodicCheckStreamsPage = () => {
    setStreams(streams => {
      streams.forEach(stream => {
        postMsg(stream.src, { code: 3, })
      })
      return []
    })
  }

  var onReceiveStream = (data) => {
    if (streams.some(e => e.ID === data.ID)) {
      //streams[streams.findIndex(stream => stream.id === data.ID)].body = data.body;
    } else {
      streams.push(data)
      setTimeout(() => {
        perodicStartStream(data)
      }, 50 * 1000)
    }
  }

  useEffect(() => {
    if (selectedTab === "Streams") {
      setStreamsTimer(null)
    } else {
      setStreamsTimer(setTimeout(() => {
        perodicCheckStreamsPage()
      }, 5 * 60 * 1000))
    }
  }, [selectedTab])

  useEffect(() => {

    window.$('.menu').css('visibility', 'visible')
    window.$('.menu .item span').each(function () { window.$(this).removeClass('selected') })

    window.$('.menu .item span')
      .filter(function (idx) { return this.innerHTML === selectedTab })
      .addClass('selected')

    var eventSource = new EventSource("http://localhost:3170/events")
    eventSource.addEventListener("msg", ev => {
      var data = JSON.parse(ev.data)
      NotificationManager.info(data.src + ": " + receivedCodes[data.code]);
    })

    eventSource.addEventListener("audio", ev => {
      var data = JSON.parse(ev.data)
      NotificationManager.info(data.src + " sends audio message, click here to play it!", '', 3000, () => onPlayAudio(data.src, data.body), true);
    })

    eventSource.addEventListener("video-fragment", ev => {
      onReceiveStream(JSON.parse(ev.data))
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
          <LogIn onLogIn={() => { onChange("Map") }} />
          : selectedTab === "Map" ?
            <Home />
            : selectedTab === "Units" ?
              <Profile />
              : selectedTab === "Streams" ?
                <Streams streams={streams} onEndStream={stream => onEndStreamRequest(stream)}/>
                : <> </>
      }
    </>
  );
}

export default App;
