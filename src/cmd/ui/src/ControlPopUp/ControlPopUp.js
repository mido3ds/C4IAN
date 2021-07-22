import React from 'react';
import Modal from 'react-modal';
import Control from '../Profile/Control/Control'
import './ControlPopUp.css'

Modal.setAppElement('#root');

class ControlPopUp extends React.Component {
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
                    className="get-msg-modal"
                >
                    <button className="close" onClick={this.closeModal}>
                        &times;
                    </button>
                    {!this.props.group || (!this.props.radius && this.props.group.id === "broadcast") ? <> </> :
                        <Control port={this.props.port} type="group" unit={{ ...this.props.group, name: this.props.group.id !== "broadcast" ? this.props.group.id + " group" : "all units on radius " + this.props.radius.toString(10)}}> </Control>
                    }
                </Modal>
            </div>
        );
    }

} export default ControlPopUp;
