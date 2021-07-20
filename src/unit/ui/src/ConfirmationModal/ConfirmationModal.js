import React from 'react';
import Modal from 'react-modal';
import './ConfirmationModal.css'
import { sentCodes } from '../codes'
Modal.setAppElement('#root');

class ConfirmationModal extends React.Component {
    constructor(props) {
        super(props)

        this.state = {
            isOpen: false,
        }
    }

    openModal = () => {
        this.setState({
            isOpen: true
        })
    }

    closeModal = () => {
        this.setState({
            isOpen: false
        })
    }

    render() {
        return (
            <div>
                <Modal
                    isOpen={this.state.isOpen}
                    onRequestClose={this.closeModal}
                    className="confirmation-modal">
                    <button className="close" onClick={this.closeModal}>
                        &times;
                    </button>
                    <p className="confirmation-msg"> {"Are you sure you want to send " + sentCodes[this.props.msgCode] + " message to " + this.props.name + "?"} </p>
                    <div className="confirmation-control">
                        <button onClick={() => {
                            this.props.onSend();
                            this.closeModal();
                        }} className="send-msg-button btn"> Send </button>
                        <button onClick={this.closeModal} className="close-msg-button btn"> Close </button>
                    </div>
                </Modal>
            </div>
        );
    }

} export default ConfirmationModal;
