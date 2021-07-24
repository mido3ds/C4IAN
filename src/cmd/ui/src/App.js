import { NotificationContainer, NotificationManager } from 'react-notifications';
import React, { useState, useEffect, useRef } from 'react';
import Profile from './Profile/Profile'
import Home from './Home/Home'
import LogIn from './LogIn/LogIn'
import Menu from './Menu/Menu'
import GetPort from './GetPort/GetPort';
import PlayAudio from './PlayAudio/PlayAudio'
import Streams from './Streams/Streams'
import { postMsg, getNames } from './Api/Api';

import 'react-notifications/lib/notifications.css';

import './index.css';
import './App.css';

function App() {
  const playAudioRef = useRef(null);
  const homeRef = useRef(null);
  const getPortRef = useRef(null);

  const [audioModalName, setAudioModalName] = useState(null)
  const [audioModalData, setAudioModalData] = useState(null)

  const [selectedTab, setSelectedTab] = useState("Log Out")

  const [streamsTimer, setStreamsTimer] = useState(null)
  const [streams, setStreams] = useState([])

  const [unitNames, setUnitNames] = useState(null)

  const [port, setPort] = useState(null)

  var onPlayAudio = (name, data) => {
    setAudioModalName(name);
    setAudioModalData(data);
    playAudioRef.current.open()
  }

  var onEndStreamRequest = (data) => {
    setStreams(streams => {
      streams = streams.filter(stream => {
        return stream.id !== data.id
      })
      return streams
    })

    setPort(port => {
      postMsg(data.src, { code: 3, }, port)
      return port
    })
  }

  var perodicStartStream = (data) => {
    setStreams(streams => {
      if (streams.some(stream => stream.id === data.id && stream.src === data.src)) {
        // Resend start stream request
        setPort(port => {
          postMsg(data.src, { code: 2, }, port)
          return port
        })
        setTimeout(() => {
          perodicStartStream(data)
        }, 30 * 1000)
      }
      return streams
    })
  }

  var perodicCheckStreamsPage = () => {
    setStreams(streams => {
      setPort(port => {
        streams.forEach(stream => {
          postMsg(stream.src, { code: 3, }, port)
        })
        return port
      })
      return []
    })
  }

  var onReceiveAudio = (data) => {
    setSelectedTab(selectedTab => {
      setUnitNames(unitNames => {
        if (selectedTab !== "Log Out")
          NotificationManager.info(data.src + " sends audio message, click here to play it!", '', 3000, () => onPlayAudio(unitNames[data.src], data.body), true);
        return unitNames
      })
      return selectedTab

    })
  }

  var onReceiveMessage = (data) => {
    setSelectedTab(selectedTab => {
      if (selectedTab === "Map")
        homeRef.current.onNewMessage(data)
      return selectedTab
    })
  }

  var onReceiveStream = (data) => {
    if (streams.some(stream => stream.src === data.src)) {
      streams[streams.findIndex(stream => stream.src === data.src)].id = data.id;
    } else {
      setStreams((streams) => { return [...streams, data] })
      perodicStartStream(data)

      setSelectedTab(selectedTab => {
        setUnitNames(unitNames => {
          if (selectedTab !== "Log Out")
            NotificationManager.info(unitNames[data.src] + " start streaming, click here to open streaming page!", '', 3000, () => onChange("Streams"), true);
          return unitNames
        })
        return selectedTab
      })
    }
  }

  var onGetPort = (port) => {
    setPort(() => {
      var eventSource = new EventSource("http://localhost:" + port + "/events")
      eventSource.addEventListener("msg", ev => {
        onReceiveMessage(JSON.parse(ev.data))
      })

      eventSource.addEventListener("audio", ev => {
        onReceiveAudio(JSON.parse(ev.data))
      })

      eventSource.addEventListener("video", ev => {
        onReceiveStream(JSON.parse(ev.data))
      })

      getNames(port).then(unitsData => {
        setUnitNames(unitsData)
      })

      setSelectedTab(selectedTab => {
        if (selectedTab === "Map") {
          homeRef.current.onChangePort(port)
        }
        return selectedTab
      })

      return port
    })

  }

  useEffect(() => {
    if (selectedTab === "Streams") {
      setStreamsTimer(streamsTimer => {
        if (streamsTimer)
          clearTimeout(streamsTimer);
        return null
      })
    } else {
      setStreamsTimer(setTimeout(() => {
        perodicCheckStreamsPage()
      }, 5 * 60 * 1000))
    }
  }, [selectedTab])


  useEffect(() => {
    getPortRef.current.open();
  }, [])

  var onChange = (sTab) => {
    setSelectedTab(sTab)

    if (sTab === "Log Out") {
      window.$('.menu').css('visibility', 'hidden')
    } else {
      window.$('.menu').css('visibility', 'visible')
      window.$('.menu .item span').each(function () { window.$(this).removeClass('selected') })

      window.$('.menu .item span')
        .filter(function (idx) { return this.innerHTML === sTab })
        .addClass('selected')
    }
  }

  return (
    <>
      <GetPort onGetPort={onGetPort} ref={getPortRef}> </GetPort>
      <NotificationContainer />
      <PlayAudio name={audioModalName} audio={audioModalData} ref={playAudioRef}></PlayAudio>
      <Menu onChange={selectedTab => onChange(selectedTab)}> </Menu>
      {
        selectedTab === "Log Out" ?
          <LogIn port={port} onLogIn={() => { onChange("Map") }} />
          : selectedTab === "Map" ?
            <Home port={port} ref={homeRef} selectedTab={selectedTab} />
            : selectedTab === "Units" ?
              <Profile port={port} />
              : selectedTab === "Streams" ?
                <Streams port={port} streams={streams} onEndStream={stream => onEndStreamRequest(stream)} />
                : <> </>
      }
    </>
  );
}

export default App;
