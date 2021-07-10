import React, {useState} from 'react';
import UnitItem from "../UnitItem/UnitItem.js"
import './UnitList.css';
import anime from 'animejs'

function UnitList() {
    
    return (
        <div className="unit-list-wrap">
            <div className="unit-list-upper-arrow-area">
                <i className="fas fa-caret-up fa-lg unit-list-upper-arrow"></i>
            </div>
            <div id="card-slider" className="unit-list-area">
                <UnitItem name="One" />
                <UnitItem name="Two" />
                <UnitItem name="Three" />
                <UnitItem name="Four" />
                <UnitItem name="Five" />
                <UnitItem name="Six" />
                <UnitItem name="Seven" />
                <UnitItem name="Eight" />
            </div>
            <div className="unit-list-lower-arrow-area area-active">
                <i className="fas fa-caret-down fa-lg unit-list-lower-arrow arrow-active"></i>
            </div>
        </div>
    );
}
export default UnitList;
