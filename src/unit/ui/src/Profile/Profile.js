import './Profile.css';
import React from 'react';
import Gallery from './Gallery/Gallery'
import Control from './Control/Control'
import ChatBox from './ChatBox/ChatBox'
import ProfileList from './ProfileList/ProfileList'

import {
    withStyles,
    Arwes,
    Content
} from 'arwes';
import withTemplate from '../withTemplate';

const styles = theme => ({
    root: {
        background: 'rgba(7, 43, 41, 0.05);',
    },
});

class Profile extends React.Component {
    constructor(props) {
        super(props);
        this.state = { activatedTab: "Control", activatedUnit: { name: "Me" } };
    }

    updateActiveUnit(activatedUnit) {
        this.setState({ ...this.state, activatedUnit: activatedUnit })
    }

    updateActiveTab(activatedTab) {
        this.setState({ ...this.state, activatedTab: activatedTab })
    }

    render() {
        const { classes, msgs, audios, setMsgs } = this.props;
        const profileComponents = {
            "Control": <Control msgs={msgs} setMsgs={setMsgs} />,
            "Audios": <Gallery type="audio" audios={audios} />,
            "Messages": <ChatBox msgs={msgs} />,
        }

        return (
            <div>
                <Arwes>
                    <Content className={`profile-root ${classes.root}`}>
                        <div className="profile-container">
                            <div className="unit-name">
                                <p> {this.props.name} </p>
                            </div>
                            <ProfileList onChange={activatedTab => this.updateActiveTab(activatedTab)}></ProfileList>
                            <div data-augmented-ui="tl-2-clip-x tr-clip r-clip-y br-clip-x br-clip border l-rect-y bl-clip-x " className="profile-frame">
                                {React.cloneElement(
                                    profileComponents[this.state.activatedTab],
                                    { unit: this.state.activatedUnit, port: this.props.port }
                                )}
                            </div>
                        </div>
                    </Content>
                </Arwes>
            </div>
        )
    }
}
export default withTemplate(withStyles(styles)(Profile));
