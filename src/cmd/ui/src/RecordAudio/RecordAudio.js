import React from 'react';
import Modal from 'react-modal';
import './RecordAudio.css'

import Recorder from 'react-mp3-recorder'
import ReactAudioPlayer from 'react-audio-player'
import blobToBuffer from 'blob-to-buffer'

Modal.setAppElement('#root');

class RecordAudio extends React.Component {
    constructor(props) {
        super(props)

        this.state = {
            audioBlob: null,
            audioUrl: null,
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

    _onRecordingComplete = (blob) => {
        blobToBuffer(blob, (err, buffer) => {
            if (err) {
                console.error(err)
                return
            }

            this.setState({
                audioBolb: blob
            })

            if (this.state.audioUrl) {
                window.URL.revokeObjectURL(this.state.audioUrl)
            }

            this.setState({
                audioUrl: window.URL.createObjectURL(blob)
            })
        })
    }

    _onRecordingError = (err) => {
        console.log('recording error', err)
    }

    render() {
        return (
            <div>
                <Modal
                    isOpen={this.state.isOpen}
                    onRequestClose={this.closeModal}
                    className="record-audio-modal"
                >
                    <button className="close" onClick={this.closeModal}>
                        &times;
                    </button>
                    <Recorder
                        onRecordingComplete={this._onRecordingComplete}
                        onRecordingError={this._onRecordingError}
                        className="record-icon"
                    />
                    <p className="record-msg"> Click and hold to start recording. </p>
                    {this.state.audioUrl ?
                        <>
                            <ReactAudioPlayer
                                id='audio'
                                controls
                                className="record-audio-player"
                                src={this.state.audioUrl}
                            ></ReactAudioPlayer>
                            <div className="audio-control">
                                <button onClick={() => {
                                    this.props.onSend(this.state.audioBolb)
                                    this.setState({
                                        audioUrl: null
                                    })
                                    this.closeModal()
                                }} className="send-audio-button btn"> Send Audio </button>
                                <button onClick={() => {
                                    this.setState({
                                        audioUrl: null
                                    })
                                    this.closeModal()
                                }} className="close-audio-button btn"> Close </button>
                            </div>
                        </>
                        : <> </>
                    }
                </Modal>
            </div>
        );
    }

} export default RecordAudio;
