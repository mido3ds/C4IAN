import './Profile.css';
import React from 'react';
import Gallery from './Gallery/Gallery'
import Control from './Control/Control'
import {
    withStyles,
    Arwes,
    Content
} from 'arwes';
import withTemplate from '../withTemplate';
import ProfileList from './ProfileList/ProfileList';
import UnitList from './UnitList/UnitList'

const styles = theme => ({
    root: {
        background: 'rgba(7, 43, 41, 0.05);',
    },
});

const profileComponents = {
    "videos": <Gallery type="video" />,
    "audios": <Gallery type="audio" />,
    "control": <Control />,
}

class Profile extends React.Component {
    constructor(props) {
        super(props);
        this.state = { activatedTab: "videos", activatedUnitOrGroup: {name: "hello"} };
    }

    render() {
        const { classes } = this.props;
        return (
            <div>
                <Arwes>
                    <Content className={`profile-root ${classes.root}`}>
                        <UnitList type={this.props.type} onChange={activatedUnitOrGroup => this.setState({ ...this.state, activatedUnitOrGroup: activatedUnitOrGroup })}></UnitList>
                        <div data-augmented-ui="tl-2-clip-x tr-clip r-clip-y br-clip-x br-clip border l-rect-y bl-clip-x " className="profile-frame">
                            <ProfileList onChange={activatedTab => this.setState({ ...this.state, activatedTab: activatedTab })}></ProfileList>
                            {React.cloneElement(
                                profileComponents[this.state.activatedTab],
                                { UnitOrGroup: this.state.activatedUnitOrGroup }
                            )}
                        </div>
                    </Content>
                </Arwes>
            </div>
        )
    }
}
export default withTemplate(withStyles(styles)(Profile));
