import './UnitItem.css';
import React from 'react';
import uImage from '../../images/unit.png';
import gImage from '../../images/group.png';

import {
    Words
} from 'arwes';


class UnitItem extends React.Component {
    render() {
        return (
            <div className="unit-item-container">
                {this.props.type === "unit" ?
                <img className="unit-item-profile-image" alt="unit" src={uImage}></img>:
                <img className="group-item-profile-image" alt="group" src={gImage}></img>}

                <div className="unit-item-parent">
                    <div className="unit-item">
                        <Words className="unit-item-name"> {this.props.unit.name} </Words>
                        <Words className="unit-item-ip"> {this.props.unit.ip} </Words>
                    </div>
                </div>
            </div>
            
        );
    }
}

export default UnitItem;
