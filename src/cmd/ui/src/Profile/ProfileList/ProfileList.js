
import React from 'react';
import './ProfileList.css';

function ProfileList() {
    return (
        <div className="list-container">
            <div data-augmented-ui="bl-clip-x " className="upper-tap">
                <i className="fas fa-caret-up fa-lg upper-arrow"></i>
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
                <i className="fas fa-caret-down fa-lg down-arrow"></i>
            </div>
        </div>
    );
}
export default ProfileList;
