import './UnitItem.css';
import React from 'react';
import uImage from '../../images/unit.png';
import {
    Words
} from 'arwes';


class UnitItem extends React.Component {
    render() {
        return (
            <>
            {this.props.unit ? 
            <div className="unit-item-container">
                <img className="unit-item-profile-image" alt="unit" src={uImage}></img>
                <div className="unit-item-parent">
                    <div className="unit-item">
                        <Words className="unit-item-name"> {this.props.unit.name} </Words>
                        <Words className="unit-item-ip"> {this.props.unit.ip} </Words>
                    </div>
                </div>
            </div> : <> </>
            }
            </>
        );
    }
}

export default UnitItem;
