import React, { useState, forwardRef, useImperativeHandle } from 'react';
import Modal from 'react-modal';
import './GetPort.css'

Modal.setAppElement('#root');

const GetPort = forwardRef(({ onGetPort }, ref) => {
    const [isOpen, setIsOpen] = useState(false)

    var openModal = () => {
        setIsOpen(true)
    }

    useImperativeHandle(ref, () => ({
        open() {
            openModal()
        }

    }));

    var closeModal = () => {
        setIsOpen(false)
    }

    var handleClick = (event) => {
        if (event.key === 'Enter') {
            onGetPort(window.$('input').val())
            closeModal()
        }
    }

    return (
        <div>
            <Modal
                isOpen={isOpen}
                onRequestClose={closeModal}
                className="get-port-modal">
                <button className="close" onClick={() => {
                    closeModal()
                }}>
                    &times;
                </button>
                <p className="get-port-msg"> Enter Connection Port: </p>
                <input type="text" className="port-text-field" id="fname" name="fname"
                    onChange={handleClick} onKeyDown={handleClick}></input>
            </Modal>
        </div>
    );

}); export default GetPort;
