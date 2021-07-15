import React, { useState } from 'react';
import Modal from 'react-modal';
import './PlayAudio.css'

import ReactAudioPlayer from 'react-audio-player'
import blobToBuffer from 'blob-to-buffer'

Modal.setAppElement('#root');

class PlayAudio extends React.Component {
    constructor(props) {
        super(props)

        this.state = {
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

    convertAudioToURL = (blob) => {
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

    componentDidMount() {
        this.convertAudioToURL(this.props.audioBolb)
    }


    render() {
        return (
            <div>
                <Modal
                    isOpen={this.state.isOpen}
                    onRequestClose={this.closeModal}
                    className="play-audio-modal">
                    <button className="close" onClick={this.closeModal}>
                        &times;
                    </button>
                    <p className="play-audio-msg"> {this.props.name + "'s audio message"}. </p>
                    {this.state.audioUrl ?
                        <ReactAudioPlayer
                            id='audio'
                            controls
                            className="audio-player"
                            src={this.state.audioUrl}
                        ></ReactAudioPlayer>

                        : <> </>
                    }
                </Modal>
            </div>
        );
    }

} export default PlayAudio;
