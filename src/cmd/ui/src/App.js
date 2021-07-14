import React, { useState } from 'react';
import Profile from './Profile/Profile'
import Home from './Home/Home'
import LogIn from './LogIn/LogIn'
import Menu from './Menu/Menu'
import './index.css';
import './App.css';

const tabsComponents = {
  "Map": <Home />,
  "Units": <Profile type="unit"/>,
}

function App() {
  const [selectedTab, setSelectedTab] = useState("Log Out")

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
      <Menu onChange={selectedTab => onChange(selectedTab)}> </Menu>
      <Profile> </Profile>
      {/*{selectedTab === "Log Out" ?
        <LogIn onLogIn={() => { onChange("Map") }} />
        : tabsComponents[selectedTab]
      }*/}
    </>
  );
}

export default App;
