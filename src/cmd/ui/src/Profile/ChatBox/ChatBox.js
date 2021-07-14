import React, { useRef, useEffect, useState } from 'react';
import uImage from '../../images/unit.png';
import './ChatBox.css'

function ChatBox() {

    useEffect(() => {
        var element = document.querySelector(".chat-container");
        element.scrollTop = element.scrollHeight;
    }, [])


    return (
        <div class="content-wrapper">
            <div class="chat-container">
                <ul class="chat-box chatContainerScroll">
                    <li class="chat-left">
                        <div class="chat-avatar">
                            <img className="unit-item-profile-image" alt="unit" src={uImage}></img>
                            <div class="chat-name"> Russell</div>
                        </div>
                        <div class="chat-text"> Attack </div>
                        <div class="chat-hour"> 08:55 </div>
                    </li>
                    <li class="chat-right">
                        <div class="chat-hour"> 08:56 </div>
                        <div class="chat-text"> Defense </div>
                        <div class="chat-avatar">
                            <img className="unit-item-profile-image" alt="unit" src={uImage}></img>
                            <div class="chat-name">Sam</div>
                        </div>
                    </li>
                    <li class="chat-left">
                        <div class="chat-avatar">
                            <img className="unit-item-profile-image" alt="unit" src={uImage}></img>
                            <div class="chat-name"> Russell</div>
                        </div>
                        <div class="chat-text"> Attack </div>
                        <div class="chat-hour"> 08:57 </div>
                    </li>
                    <li class="chat-left">
                        <div class="chat-avatar">
                            <img className="unit-item-profile-image" alt="unit" src={uImage}></img>
                            <div class="chat-name"> Russell</div>
                        </div>
                        <div class="chat-text"> Attack </div>
                        <div class="chat-hour"> 08:57 </div>
                    </li>
                    <li class="chat-right">
                        <div class="chat-hour"> 08:59 </div>
                        <div class="chat-text"> Defense </div>
                        <div class="chat-avatar">
                            <img className="unit-item-profile-image" alt="unit" src={uImage}></img>
                            <div class="chat-name"> Joyse </div>
                        </div>
                    </li>
                    <li class="chat-right">
                        <div class="chat-hour"> 08:59 </div>
                        <div class="chat-text"> Defense </div>
                        <div class="chat-avatar">
                            <img className="unit-item-profile-image" alt="unit" src={uImage}></img>
                            <div class="chat-name"> Joyse </div>
                        </div>
                    </li>
                    <li class="chat-left">
                        <div class="chat-avatar">
                            <img className="unit-item-profile-image" alt="unit" src={uImage}></img>
                            <div class="chat-name"> Russell</div>
                        </div>
                        <div class="chat-text"> Attack </div>
                        <div class="chat-hour"> 08:57 </div>
                    </li>
                </ul>
            </div>
        </div>
    )

} export default ChatBox;