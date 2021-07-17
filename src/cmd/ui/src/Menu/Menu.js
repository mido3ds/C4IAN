import React from 'react';
import './Menu.css';


function Menu({onChange}) {
   
    return (
        <div className="menu">
            <div className="label">Menu</div>
            <div className="spacer"></div>
            <div className="item" onClick={() => {onChange("Map")}}><span className="selected">Map</span></div>
            <div className="item" onClick={() => {onChange("Units")}}><span>Units</span></div>
            <div className="item" onClick={() => {onChange("Streams")}}><span>Streams</span></div>
            <div className="item" onClick={() => {onChange("Log Out")}}><span>Log Out</span></div>
        </div>
    );
}

export default Menu;