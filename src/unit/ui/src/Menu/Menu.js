import React from 'react';
import './Menu.css';

function Menu({onChange}) {
    return (
        <div className="menu">
            <div className="label">Menu</div>
            <div className="spacer"></div>
            <div className="item" onClick={() => {onChange("Profile")}}><span>Profile</span></div>
            <div className="item" onClick={() => {onChange("Log Out")}}><span>Log Out</span></div>
        </div>
    );
}

export default Menu;