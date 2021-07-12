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
                        <div className="first-info">
                            <Words animate className="first-info-header"> Power </Words>
                            <div className="first-info-bar">
                                <progress className="progress-bar" max="100" value="80"></progress>
                            </div>
                        </div>
                        <div className="first-info">
                            <Words animate className="first-info-header"> Speed </Words>
                            <div className="first-info-bar">
                                <progress className="progress-bar" max="100" value="65"></progress>
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