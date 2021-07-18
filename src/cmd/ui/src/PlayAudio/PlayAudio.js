import React from 'react';
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
                audioBlob: blob
            })

            if (this.state.audioUrl) {
                window.URL.revokeObjectURL(this.state.audioUrl)
            }

            this.setState({
                audioUrl: window.URL.createObjectURL(blob)
            })
        })
    }

    b64toBlob(b64Data, contentType='', sliceSize=512) {
        const byteCharacters = atob(b64Data);
        const byteArrays = [];
      
        for (let offset = 0; offset < byteCharacters.length; offset += sliceSize) {
          const slice = byteCharacters.slice(offset, offset + sliceSize);
      
          const byteNumbers = new Array(slice.length);
          for (let i = 0; i < slice.length; i++) {
            byteNumbers[i] = slice.charCodeAt(i);
          }
      
          const byteArray = new Uint8Array(byteNumbers);
          byteArrays.push(byteArray);
        }
      
        const blob = new Blob(byteArrays, {type: contentType});
        return blob;
    }

    componentDidMount() {
        if(!this.props.audioBlob) return

        var blob = this.b64toBlob(this.props.audioBlob, 'audio/mpeg')
        this.convertAudioToURL(blob)
    }


    render() {
        return (
            <div>
                <Modal
                    isOpen={this.state.isOpen}
                    onRequestClose={this.closeModal}
                    className="play-audio-modal">
                    <button className="close" onClick={()=> {
                        this.setState({
                            audioUrl: null
                        })
                        this.closeModal()
                        }}>
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
