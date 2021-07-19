import React, { useEffect, useState, forwardRef, useImperativeHandle } from 'react';
import Modal from 'react-modal';
import ReactPlayer from 'react-player'
import './PlayVideo.css'

Modal.setAppElement('#root');

const baseURL = "http://localhost:3170/api/";
const M3U8Name   = "index.m3u8";

const HLSConfig = {
    maxLoadingDelay: 4,
    minAutoBitrate: 0,
    lowLatencyMode: true,
}

const PlayVideo = forwardRef(({ videoData }, ref) => {
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
                    <ReactHlsPlayer
                      src={baseURL + "stream/" + videoData.src + "/" + videoData.body + "/" + M3U8Name}
                      autoPlay={true}
                      controls={true}
                      hlsConfig={HLSConfig}  
                    />
                    : <> </>
                }
            </Modal>
        </div>
    );
}); export default PlayVideo;
