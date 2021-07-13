
import React, {useState} from 'react';
import './ProfileList.css';

const profileTabs = {
    0: "videos",
    1: "audios",
    2: "control",
}

function ProfileList({onChange}) {
    const[activeItem, setactiveItem] = useState(0);
    var down = () => {
        setactiveItem( () => {
            window.$(window.$('.list-item .middle').toArray()[activeItem]).removeClass('item-active')
            window.$(window.$('.list-item .middle .right-arrow').toArray()[activeItem]).removeClass('text-active')
            window.$(window.$('.list-item .middle .list-item-text').toArray()[activeItem]).removeClass('text-active')
            
            var newItem = activeItem < 2? activeItem + 1: 2;
            onChange(profileTabs[newItem])

            if(newItem === 0) {
                window.$('.upper-arrow').removeClass('active-arrow')
                window.$('.down-arrow').addClass('active-arrow')
            } else if (newItem === 1) {
                window.$('.down-arrow').addClass('active-arrow')
                window.$('.upper-arrow').addClass('active-arrow')
            } else if (newItem === 2) {
                window.$('.down-arrow').removeClass('active-arrow')
                window.$('.upper-arrow').addClass('active-arrow')
            }

            window.$(window.$('.list-item .middle').toArray()[newItem]).addClass('item-active')
            window.$(window.$('.list-item .middle .right-arrow').toArray()[newItem]).addClass('text-active')
            window.$(window.$('.list-item .middle .list-item-text').toArray()[newItem]).addClass('text-active')
            return activeItem < 2? activeItem + 1: 2
        })
    }

    var up = () => { 
        setactiveItem( () => {
            window.$(window.$('.list-item .middle').toArray()[activeItem]).removeClass('item-active')
            window.$(window.$('.list-item .middle .right-arrow').toArray()[activeItem]).removeClass('text-active')
            window.$(window.$('.list-item .middle .list-item-text').toArray()[activeItem]).removeClass('text-active')
            
            var newItem = activeItem > 0? activeItem - 1: 0
            onChange(profileTabs[newItem])

            if(newItem === 0) {
                window.$('.upper-arrow').removeClass('active-arrow')
                window.$('.down-arrow').addClass('active-arrow')
            } else if (newItem === 1) {
                window.$('.down-arrow').addClass('active-arrow')
                window.$('.upper-arrow').addClass('active-arrow')
            } else if (newItem === 2) {
                window.$('.down-arrow').removeClass('active-arrow')
                window.$('.upper-arrow').addClass('active-arrow')
            }

            window.$(window.$('.list-item .middle').toArray()[newItem]).addClass('item-active')
            window.$(window.$('.list-item .middle .right-arrow').toArray()[newItem]).addClass('text-active')
            window.$(window.$('.list-item .middle .list-item-text').toArray()[newItem]).addClass('text-active')
            return activeItem > 0? activeItem - 1: 0
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
                    <p className="list-item-text text-active"> Videos </p>
                </div>
                <div data-augmented-ui="br-clip-x tl-clip-x " className="lower"></div>
            </div>
            <div className="list-item"> 
                <div className="square"> </div>
                <div data-augmented-ui="bl-clip-x tr-clip-x " className="upper"></div>
                <div data-augmented-ui="tr-clip-x br-clip-x border" className="middle">
                    <i className="fas fa-caret-right fa-2x right-arrow"></i>
                    <p className="list-item-text"> Audios </p>
                </div>
                <div data-augmented-ui="tl-clip-x br-clip-x " className="lower"></div>
            </div>
            <div className="list-item">
                <div className="square"> </div>
                <div data-augmented-ui="bl-clip-x tr-clip-x " className="upper"></div>
                <div data-augmented-ui="tr-clip-x br-clip-x border" className="middle">
                    <i className="fas fa-caret-right fa-2x right-arrow"></i>
                    <p className="list-item-text"> Control </p>
                </div>
                <div data-augmented-ui="br-clip-x tl-clip-x " className="lower"></div>
            </div>
            <div data-augmented-ui="tl-clip-x " className="lower-tap">
                <i onClick={down}  className="fas fa-caret-down fa-lg down-arrow active-arrow"></i>
            </div>
        </div>
    );
}
export default ProfileList;
