import React, { useEffect, useState, forwardRef, useImperativeHandle } from 'react';
import Modal from 'react-modal';
import ReactPlayer from 'react-player'
import './PlayVideo.css'

Modal.setAppElement('#root');

const PlayVideo = forwardRef(({ videoUrl }, ref) => {
    const [isOpen, setIsOpen] = useState(false)
    const [video, setVideo] = useState(null)

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

    useEffect(() => {
        if (!videoUrl) return
        import(videoUrl)
            .then(module => setVideo(module.default))
    }, [videoUrl])

    return (
        <div>
            <Modal
                isOpen={isOpen}
                onRequestClose={closeModal}
                className="play-video-modal">
                <button className="close" onClick={() => {
                    setVideo(null)
                    closeModal()
                }}>
                    &times;
                </button>
                {videoUrl ?
                    <ReactPlayer controls url={video} />
                    : <> </>
                }
            </Modal>
        </div>
    );
}); export default PlayVideo;
