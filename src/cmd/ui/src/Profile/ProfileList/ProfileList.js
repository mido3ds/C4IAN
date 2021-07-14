
import React, { useState } from 'react';
import './ProfileList.css';

const profileTabs = [
    "Control",
    "Audios",
    "Videos",
    "Messages",
    "Locations",
    "Heartbeats"
]

function ProfileList({ onChange }) {
    const [firstTab, setFirstTab] = useState(profileTabs[0])
    const [secondTab, setSecondTab] = useState(profileTabs[1])
    const [thirdTab, setThirdTab] = useState(profileTabs[2])
    const [selectedTab, setSelectedTab] = useState(0);
    const [activeItem, setActiveItem] = useState(0);

    var down = () => {
        setActiveItem(() => {
            var newActiveItem = activeItem < 2 ? activeItem + 1 : 2;
            setSelectedTab(() => {
                var newSelectedTab = selectedTab < profileTabs.length - 1 ? selectedTab + 1 : profileTabs.length - 1;
                onChange(profileTabs[newSelectedTab])

                if (newSelectedTab >= 1 && newSelectedTab < profileTabs.length - 1) {
                    window.$('.down-arrow').addClass('active-arrow')
                    window.$('.upper-arrow').addClass('active-arrow')
                } else if (newSelectedTab === profileTabs.length - 1) {
                    window.$('.down-arrow').removeClass('active-arrow')
                    window.$('.upper-arrow').addClass('active-arrow')
                }

                if (newActiveItem !== activeItem) {
                    window.$(window.$('.list-item .middle').toArray()[newActiveItem - 1]).removeClass('item-active')
                    window.$(window.$('.list-item .middle .right-arrow').toArray()[newActiveItem - 1]).removeClass('text-active')
                    window.$(window.$('.list-item .middle .list-item-text').toArray()[newActiveItem - 1]).removeClass('text-active')
                    window.$(window.$('.list-item .middle').toArray()[newActiveItem]).addClass('item-active')
                    window.$(window.$('.list-item .middle .right-arrow').toArray()[newActiveItem]).addClass('text-active')
                    window.$(window.$('.list-item .middle .list-item-text').toArray()[newActiveItem]).addClass('text-active')
                }


                if (newActiveItem === 2 && newSelectedTab !== selectedTab) {
                    setFirstTab(profileTabs[newSelectedTab - 2])
                    setSecondTab(profileTabs[newSelectedTab - 1])
                    setThirdTab(profileTabs[newSelectedTab])
                }

                
                return newSelectedTab;
            })
            return newActiveItem
        })
    }

    var up = () => {
        setActiveItem(() => {
            var newActiveItem = activeItem > 0 ? activeItem - 1 : 0;
            setSelectedTab(() => {
                var newSelectedTab = selectedTab > 0 ? selectedTab - 1 : 0;
                onChange(profileTabs[newSelectedTab])

                if (newSelectedTab === 0) {
                    window.$('.upper-arrow').removeClass('active-arrow')
                    window.$('.down-arrow').addClass('active-arrow')
                } if (newSelectedTab >= 1 && newSelectedTab < profileTabs.length - 1) {
                    window.$('.down-arrow').addClass('active-arrow')
                    window.$('.upper-arrow').addClass('active-arrow')
                }

                if (newActiveItem !== activeItem) {
                    window.$(window.$('.list-item .middle').toArray()[newActiveItem + 1]).removeClass('item-active')
                    window.$(window.$('.list-item .middle .right-arrow').toArray()[newActiveItem + 1]).removeClass('text-active')
                    window.$(window.$('.list-item .middle .list-item-text').toArray()[newActiveItem + 1]).removeClass('text-active')
                    window.$(window.$('.list-item .middle').toArray()[newActiveItem]).addClass('item-active')
                    window.$(window.$('.list-item .middle .right-arrow').toArray()[newActiveItem]).addClass('text-active')
                    window.$(window.$('.list-item .middle .list-item-text').toArray()[newActiveItem]).addClass('text-active')
                }

                if (newActiveItem === 0 && newSelectedTab !== selectedTab) {
                    setFirstTab(profileTabs[newSelectedTab])
                    setSecondTab(profileTabs[newSelectedTab + 1])
                    setThirdTab(profileTabs[newSelectedTab + 2])
                }

                return newSelectedTab;
            })
            return newActiveItem
        })
    }

    return (
        <div className="list-container">
            <div data-augmented-ui="bl-clip-x " className="upper-tap">
                <i onClick={up} className="fas fa-caret-up fa-lg upper-arrow"></i>
            </div>
            <div className="list-item">
                <div className="square"> </div>
                <div data-augmented-ui="bl-clip-x tr-clip-x " className="upper"></div>
                <div data-augmented-ui="tr-clip-x br-clip-x border" className="middle item-active">
                    <i className="fas fa-caret-right fa-2x right-arrow text-active"></i>
                    <p className="list-item-text text-active"> {firstTab} </p>
                </div>
                <div data-augmented-ui="br-clip-x tl-clip-x " className="lower"></div>
            </div>
            <div className="list-item">
                <div className="square"> </div>
                <div data-augmented-ui="bl-clip-x tr-clip-x " className="upper"></div>
                <div data-augmented-ui="tr-clip-x br-clip-x border" className="middle">
                    <i className="fas fa-caret-right fa-2x right-arrow"></i>
                    <p className="list-item-text"> {secondTab} </p>
                </div>
                <div data-augmented-ui="tl-clip-x br-clip-x " className="lower"></div>
            </div>
            <div className="list-item">
                <div className="square"> </div>
                <div data-augmented-ui="bl-clip-x tr-clip-x " className="upper"></div>
                <div data-augmented-ui="tr-clip-x br-clip-x border" className="middle">
                    <i className="fas fa-caret-right fa-2x right-arrow"></i>
                    <p className="list-item-text"> {thirdTab} </p>
                </div>
                <div data-augmented-ui="br-clip-x tl-clip-x " className="lower"></div>
            </div>
            <div data-augmented-ui="tl-clip-x " className="lower-tap">
                <i onClick={down} className="fas fa-caret-down fa-lg down-arrow active-arrow"></i>
            </div>
        </div>
    );
}
export default ProfileList;
