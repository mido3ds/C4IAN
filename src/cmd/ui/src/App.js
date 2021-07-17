import { NotificationContainer, NotificationManager } from 'react-notifications';
import React, { useState, useEffect, useRef } from 'react';
import Profile from './Profile/Profile'
import Home from './Home/Home'
import LogIn from './LogIn/LogIn'
import Menu from './Menu/Menu'
import PlayAudio from './PlayAudio/PlayAudio'
import Streams from './Streams/Streams'
import { codes } from './codes'

import 'react-notifications/lib/notifications.css';

import './index.css';
import './App.css';

const tabsComponents = {
  "Map": <Home />,
  "Units": <Profile />,
  "Streams": <Streams />,
  "Log Out": <LogIn />,
}

class App extends React.Component {
  constructor(props) {
    super(props)

    this.state = {
      audioModalName: null,
      audioModalData: null,
      selectedTab: "Log Out",
      eventSource: new EventSource("http://localhost:3170/events")
    }

    this.playAudioRef = React.createRef()
    this.onChange = this.onChange.bind(this)
  }


  onPlayAudio(name, data) {
    this.setState({ audioModalName: name })
    this.setState({ audioModalData: data })

    this.playAudioRef.current.openModal()
  }

  componentDidMount() {
    this.state.eventSource.addEventListener("msg", ev => {
      var data = JSON.parse(ev.data)
      NotificationManager.info(data.src + ": " + codes[data.code]);
    })
    this.state.eventSource.addEventListener("audio", ev => {
      var data = JSON.parse(ev.data)
      NotificationManager.info(data.src + " sends audio message, click here to play it!", '', 3000, () => this.onPlayAudio(data.src, data.body), true);
    })
  }

  onChange = (selectedTab) => {
    console.log(selectedTab)
    this.setState({ selectedTab: selectedTab })

    if (this.state.selectedTab === "Log Out") {
      window.$('.menu').css('visibility', 'hidden')
    } else {
      window.$('.menu').css('visibility', 'visible')
      window.$('.menu .item span').each(function () { window.$(this).removeClass('selected') })

      window.$('.menu .item span')
        .filter(function (idx) { return this.innerHTML === selectedTab })
        .addClass('selected')
    }
  }

  render() {
    return (
      <>
        <NotificationContainer />
        <PlayAudio name={this.state.audioModalName} audioBolb={this.state.audioModalData} ref={this.playAudioRef}></PlayAudio>
        <Menu onChange={this.onChange}> </Menu>

        {
          React.cloneElement(
            tabsComponents[this.state.selectedTab],
            { onLogIn: this.setState({selectedTab: "Map"})}
        )
        }
      </>
    );
  }

}

export default App;
