
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
    const [activeItem, setactiveItem] = useState(0);

    var down = () => {
        setactiveItem(() => {
            var newItem = activeItem < profileTabs.length - 1? activeItem + 1 : profileTabs.length - 1;
            onChange(profileTabs[newItem])

            if (newItem > 2) {
                setFirstTab(profileTabs[newItem - 2])
                setSecondTab(profileTabs[newItem - 1])
                setThirdTab(profileTabs[newItem])
            } else {
                window.$(window.$('.list-item .middle').toArray()[newItem]).addClass('item-active')
                window.$(window.$('.list-item .middle .right-arrow').toArray()[newItem]).addClass('text-active')
                window.$(window.$('.list-item .middle .list-item-text').toArray()[newItem]).addClass('text-active')
                window.$(window.$('.list-item .middle').toArray()[newItem - 1]).removeClass('item-active')
                window.$(window.$('.list-item .middle .right-arrow').toArray()[newItem - 1]).removeClass('text-active')
                window.$(window.$('.list-item .middle .list-item-text').toArray()[newItem - 1]).removeClass('text-active')
            }

            if (newItem >= 1 && newItem < profileTabs.length - 1) {
                window.$('.down-arrow').addClass('active-arrow')
                window.$('.upper-arrow').addClass('active-arrow')
            } else if (newItem === profileTabs.length - 1) {
                window.$('.down-arrow').removeClass('active-arrow')
                window.$('.upper-arrow').addClass('active-arrow')
            }

            return newItem
        })
    }

    var up = () => {
        setactiveItem(() => {
            var newItem = activeItem > 0 ? activeItem - 1 : 0
            onChange(profileTabs[newItem])

            if (newItem === 0) {
                window.$('.upper-arrow').removeClass('active-arrow')
                window.$('.down-arrow').addClass('active-arrow')
            } if (newItem >= 1 && newItem < profileTabs.length - 1) {
                window.$('.down-arrow').addClass('active-arrow')
                window.$('.upper-arrow').addClass('active-arrow')
            }

            if (newItem < profileTabs.length - 3) {
                setFirstTab(profileTabs[newItem])
                setSecondTab(profileTabs[newItem + 1])
                setThirdTab(profileTabs[newItem + 2])
            } else {
                window.$(window.$('.list-item .middle').toArray()[newItem - 2]).removeClass('item-active')
                window.$(window.$('.list-item .middle .right-arrow').toArray()[newItem - 2]).removeClass('text-active')
                window.$(window.$('.list-item .middle .list-item-text').toArray()[newItem - 2]).removeClass('text-active')
                window.$(window.$('.list-item .middle').toArray()[newItem - 3]).addClass('item-active')
                window.$(window.$('.list-item .middle .right-arrow').toArray()[newItem - 3]).addClass('text-active')
                window.$(window.$('.list-item .middle .list-item-text').toArray()[newItem - 3]).addClass('text-active')
            }

            return newItem
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
