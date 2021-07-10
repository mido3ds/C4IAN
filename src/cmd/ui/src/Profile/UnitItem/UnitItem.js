import './UnitItem.css';
import React from 'react';
import uImage from '../../images/unit.png';
import {
    Words
} from 'arwes';


class UnitItem extends React.Component {

    render() {
        return (
            <div className="unit-item-container">
                <img className="unit-item-profile-image" alt="unit" src={uImage}></img>
                <div className="unit-item-parent">
                    <div className="unit-item">
                        <Words animate className="unit-item-name"> Ahmed Mahmoud</Words>
                        <Words animate className="unit-item-ip"> 192.168.1.1 </Words>
                    </div>
                </div>
            </div>
            
        );
    }
}

export default UnitItem;
