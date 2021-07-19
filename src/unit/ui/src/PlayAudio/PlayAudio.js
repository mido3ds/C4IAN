import React, { useEffect, useState, forwardRef, useImperativeHandle } from 'react';
import Modal from 'react-modal';
import './PlayAudio.css'

import ReactAudioPlayer from 'react-audio-player'
import blobToBuffer from 'blob-to-buffer'

Modal.setAppElement('#root');


const PlayAudio = forwardRef(({ audio, name }, ref) => {
    const [isOpen, setIsOpen] = useState(false)
    const [audioUrl, setAudioUrl] = useState(null)
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

    var convertAudioToURL = (blob) => {
        blobToBuffer(blob, (err, buffer) => {
            if (err) {
                console.error(err)
                return
            }

            if (audioUrl) {
                window.URL.revokeObjectURL(audioUrl)
            }

            setAudioUrl(window.URL.createObjectURL(blob))
        })
    }

    var b64toBlob = (b64Data, contentType = '', sliceSize = 512) => {
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

        const blob = new Blob(byteArrays, { type: contentType });
        return blob;
    }

    useEffect(() => {
        if (!audio) return
        convertAudioToURL(b64toBlob(audio, 'audio/mpeg'))
    }, [audio])

    return (
        <div>
            <Modal
                isOpen={isOpen}
                onRequestClose={closeModal}
                className="play-audio-modal">
                <button className="close" onClick={() => {
                    setAudioUrl(null)
                    closeModal()
                }}>
                    &times;
                </button>
                <p className="play-audio-msg"> {name + "'s audio message"}. </p>
                {audioUrl ?
                    <ReactAudioPlayer
                        id='audio'
                        controls
                        className="audio-player"
                        src={audioUrl}
                    ></ReactAudioPlayer>
                    : <> </>
                }
            </Modal>
        </div>
    );
}); export default PlayAudio;
