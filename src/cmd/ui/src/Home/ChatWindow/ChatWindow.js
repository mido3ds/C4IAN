import React, { useEffect, useState, forwardRef, useImperativeHandle } from 'react';
import Modal from 'react-modal';
import { getAllMsgs, getNames } from '../../Api/Api'
import { receivedCodes } from '../../codes'
import './ChatWindow.css'

Modal.setAppElement('#root');

const ChatWindow = forwardRef(({port}, ref) => {
    const [messages, setMessages] = useState([])
    const [unitsNames, setUnitsNames] = useState([])

    var hangdleMessage = (msg) => {
        var msgData = receivedCodes.hasOwnProperty(msg.code) ? receivedCodes[msg.code] : (msg.code).toString(10)
        setMessages((messages) => { 
            return [...messages, { name: unitsNames[msg.src], msg: msgData }]
        })
    }

    useImperativeHandle(ref, () => ({
        onNewMessage(message) {
            hangdleMessage(message)
        }
    }));

    useEffect(() => {
        if(!port) return
        getAllMsgs(port).then(msgs =>
            getNames(port).then(unitNames => {
                setMessages(messages => {
                    msgs.forEach(msg => {
                        var msgData = receivedCodes.hasOwnProperty(msg.code) ? receivedCodes[msg.code] : (msg.code).toString(10)
                        messages.push({ name: unitNames[msg.src], msg: msgData })
                    })
                    return messages
                })
                setUnitsNames(unitNames)
                var element = document.querySelector(".chat-window-wrapper");
                if (element) element.scrollTop = element.scrollHeight;
            }))
    },[port]);

    return (
        <div className="chat-window-wrapper">
            {!messages || !messages.length ?
                <div className="no-msgs">
                    <p className="no-received-msgs"> No received messages </p>
                </div>
                : messages.map((value, index) => {
                    return <p className="chat-msg"> <span className="sender-name">  {value.name + ":"} </span> {value.msg} </p>
                })
            }
        </div>
    );

}); export default ChatWindow;
