import React from 'react';
import Modal from 'react-modal';
import ReactPlayer from 'react-player'
import './PlayVideo.css'

Modal.setAppElement('#root');

class PlayVideo extends React.Component {
    constructor(props) {
        super(props)

        this.state = {
            isOpen: false,
            video: null
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

    componentDidMount() {
        import(this.props.videoUrl)
        .then(module => this.setState({ video: module.default }))  
    }

    render() {
        return (
            <div>
                <Modal
                    isOpen={this.state.isOpen}
                    onRequestClose={this.closeModal}
                    className="play-video-modal">
                    <button className="close" onClick={() => {
                        this.setState({
                            videoUrl: null
                        })
                        this.closeModal()
                        }}>
                        &times;
                    </button>
                    {this.props.videoUrl ?
                        <ReactPlayer controls url={this.state.video} />
                        : <> </>
                    }
                </Modal>
            </div>
        );
    }

} export default PlayVideo;
