import './MapPopUp.css';
import React from 'react';
import uImage from '../../images/unit.png';
import withTemplate from '../../withTemplate';

import {
    Words,
    withStyles
} from 'arwes';

const styles = theme => ({
    root: {
        background: 'rgba(7, 43, 41, 0.05);',
    },
});

class MapPopUp extends React.Component {
    render() {
        return (
            <div>
                <div className="map-unit-container">
                   {this.props.selectedUnit ?
                    <>
                    <img className="map-unit-profile-image" alt="unit" src={uImage}></img>
                    <div className="map-unit-info">
                        <Words className="map-unit-name"> {this.props.selectedUnit.name} </Words>
                        <Words className="map-unit-ip"> {this.props.selectedUnit.ip} </Words>
                        <div className="info">
                            <Words animate className="info-header"> Battery </Words>
                            <div className="info-bar">
                                <progress className="progress-bar" max="100" value="80"></progress>
                            </div>
                        </div>
                        <div className="info">
                            <Words animate className="info-header"> Health </Words>
                            <div className="info-bar">
                                <progress className="progress-bar" max="100" value={!this.props.selectedUnit.hasOwnProperty("heartbeat") ? 0: this.props.selectedUnit.heartbeat <= 100 ? this.props.selectedUnit.heartbeat : this.props.selectedUnit.heartbeat > 100 ? 100 : 0}></progress>
                            </div>
                        </div>
                    </div> </> : 
                    <div className="no-unit"> 
                        <Words animate className="no-unit-msg"> No selected unit </Words> 
                   </div>} 
                </div>
            </div>
        );
    }
}

export default withTemplate(withStyles(styles)(MapPopUp));