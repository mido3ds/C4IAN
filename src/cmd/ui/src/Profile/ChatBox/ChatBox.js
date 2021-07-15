import React, { useRef, useEffect, useState } from 'react';
import uImage from '../../images/unit.png';
import { getMsgs } from '../../Api/Api'
import moment from 'moment';
import './ChatBox.css'

function ChatBox({ unit }) {
    const [msgs, setMsgs] = useState(null)

    useEffect(() => {
        if (unit) {
            setMsgs(()=> {
                var msgsData = getMsgs(unit.ip)
                if(!msgsData || !msgsData.length) {
                    window.$('.chat-container').css('overflow-y', 'hidden')
                } else {
                    window.$('.chat-container').css('overflow-y', 'scroll')
                }
            })
            var element = document.querySelector(".chat-container");
            element.scrollTop = element.scrollHeight;
        }
    })

    return (
        <div className="content-wrapper">
            <div className="chat-container">
                {msgs ?
                    <ul className="chat-box chatContainerScroll">
                        {msgs.map((msg, index) => {
                            return <>
                                {msg.sent ?
                                    <li className="chat-right">
                                        <div className="chat-hour"> {moment.unix(msg.time).format('hh:mm')} </div>
                                        <div className="chat-text"> {msg.code} </div>
                                        <div className="chat-avatar">
                                            <img className="unit-item-profile-image" alt="unit" src={uImage}></img>
                                            <div className="chat-name"> Me </div>
                                        </div>
                                    </li> :
                                    <li className="chat-left">
                                        <div className="chat-avatar">
                                            <img className="unit-item-profile-image" alt="unit" src={uImage}></img>
                                            <div className="chat-name"> {unit.name.substr(0, unit.name.indexOf(' '))} </div>
                                        </div>
                                        <div className="chat-text"> {msg.code} </div>
                                        <div className="chat-hour"> {moment.unix(msg.time).format('hh:mm')} </div>
                                    </li>
                                }
                            </>
                        })}
                    </ul> :
                    <div className="no-data-chat-msg">
                        <p> No data to be previewed </p>
                    </div>}
            </div>
        </div>
    )

} export default ChatBox;