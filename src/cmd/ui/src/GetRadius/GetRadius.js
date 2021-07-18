import React, { useState, forwardRef, useImperativeHandle } from 'react';
import Modal from 'react-modal';
import './GetRadius.css'

Modal.setAppElement('#root');

const GetRadius = forwardRef(({ onGetRadius }, ref) => {
    const [isOpen, setIsOpen] = useState(false)
    const [radius, setRadius] = useState(0)

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
        setRadius(parseInt(window.$('input').val()))

        if (event.key === 'Enter') {
            onGetRadius(radius)
            closeModal()
        } 
    }

    return (
        <div>
            <Modal
                isOpen={isOpen}
                onRequestClose={closeModal}
                className="get-radius-modal">
                <button className="close" onClick={() => {
                    closeModal()
                }}>
                    &times;
                </button>
                <p className="get-radius-msg"> Enter the required radius: </p>
                <input type="text" className="radius-text-field" id="fname" name="fname"
                    onChange={handleClick} onKeyDown={handleClick}></input>
            </Modal>
        </div>
    );

}); export default GetRadius;
