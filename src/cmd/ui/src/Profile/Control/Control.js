import './Control.css';
import React from 'react';

function Control() {
    return (
        <div className="control-container">
            <div data-augmented-ui="tr-clip br-clip bl-clip-y border" class="control-item"></div>

            <svg className="ctrl-svg" version="1.1" id="Layer_1" xmlns="http://www.w3.org/2000/svg" x="0px" y="0px"
                viewBox="0 0 1000 1000" styles="enable-background:new 0 0 1000 1000;" space="preserve">
                <circle class="st0" cx="500" cy="500" r="302.8">
                    <animateTransform attributeType="xml"
                        attributeName="transform"
                        type="rotate"
                        from="0 500 500"
                        to="360 500 500"
                        dur="100s"
                        repeatCount="indefinite" />
                </circle>
                <circle class="st1" cx="500" cy="500" r="237.7">
                    <animateTransform attributeType="xml"
                        attributeName="transform"
                        type="rotate"
                        from="0 500 500"
                        to="360 500 500"
                        dur="40s"
                        repeatCount="indefinite" />
                </circle>
                <circle class="st2" cx="500" cy="500" r="366.8" transform="rotate(0 500 500)">
                    <animateTransform attributeType="xml"
                        attributeName="transform"
                        type="rotate"
                        from="0 500 500"
                        to="-360 500 500"
                        dur="50s"
                        repeatCount="indefinite" />
                </circle>
                <circle class="st3" cx="500" cy="500" r="385.1" />
            </svg>

        </div>
    );
}
export default Control;
