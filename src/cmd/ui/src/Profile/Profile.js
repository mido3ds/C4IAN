import './Profile.css';
import React from 'react';
import ProfileList from './ProfileList/ProfileList';
import UnitList from './UnitList/UnitList'
import Gallery from './Gallery/Gallery'
import Control from './Control/Control'
import Map from './Map/Map'
import ChatBox from './ChatBox/ChatBox'
import HeartBeatChart from './HeartBeatChart/HearBeatChart'
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

const profileComponents = {
    "Control": <Control />,
    "Audios": <Gallery type="audio" />,
    "Videos": <Gallery type="video" />,
    "Messages": <ChatBox/>,
    "Locations": <Map/>,
    "Heartbeats": <HeartBeatChart/>
}

class Profile extends React.Component {
    constructor(props) {
        super(props);
        this.state = { activatedTab: "Control", activatedUnit: {name: ""} };
    }

    render() {
        const { classes } = this.props;
        return (
            <div>
                <Arwes>
                    <Content className={`profile-root ${classes.root}`}>
                        <UnitList onChange={activatedUnit => this.setState({ ...this.state, activatedUnit: activatedUnit })}></UnitList>
                        <div data-augmented-ui="tl-2-clip-x tr-clip r-clip-y br-clip-x br-clip border l-rect-y bl-clip-x " className="profile-frame">
                            <ProfileList onChange={activatedTab => this.setState({ ...this.state, activatedTab: activatedTab })}></ProfileList>
                            {React.cloneElement(
                                profileComponents[this.state.activatedTab],
                                { unit: this.state.activatedUnit }
                            )}
                        </div>
                    </Content>
                </Arwes>
            </div>
        )
    }
}
export default withTemplate(withStyles(styles)(Profile));
