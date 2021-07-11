import React, { useState } from 'react';
import UnitItem from "../UnitItem/UnitItem.js"
import './UnitList.css';
import anime from 'animejs'

function UnitList() {

    var circularAddition = (Augend, Addend, len) => {
        return (Augend + Addend) % len;
    }

    var circularSubtract = (Minuend, Subtrahend, len) => {
        return (Minuend - Subtrahend + len) % len
    }

    const [activeItem, setActiveItem] = useState(0);
    var down = () => {
        var cards = window.$('.unit-item-container').toArray()
        for (var i = 0; i < cards.length; i++) {
            window.$(cards[i]).css("visibility", "hidden");
        }
        setActiveItem(() => {
            anime({
                targets: cards[circularSubtract(activeItem, 2, cards.length)],
                scaleX: 0.8,
                scaleY: 0.8,
                top: [-100, 50],
                opacity: '40%',
                begin: function() {
                    window.$(cards[circularSubtract(activeItem, 2, cards.length)]).css("visibility", "visible");
                },
                duration: 3000,
            })
            anime({
                targets: cards[circularSubtract(activeItem, 1, cards.length)],
                scaleX: 1,
                scaleY: 1,
                top: [50, 165],
                opacity: '100%',
                begin: function() {
                    window.$(cards[circularSubtract(activeItem, 1, cards.length)]).css("visibility", "visible");
                },
                duration: 3000,
            })
            anime({
                targets: cards[activeItem],
                scaleX: 0.8,
                scaleY: 0.8,
                top: [165, 295],
                opacity: '40%',
                begin: function() {
                    window.$(cards[activeItem]).css("visibility", "visible");
                },
                duration: 3000,
            })
            return circularSubtract(activeItem, 1, cards.length)
        })
    }


    var up = () => {
        var cards = window.$('.unit-item-container').toArray()

        setActiveItem(() => {
            anime({
                targets: cards[activeItem],
                scaleX: 0.8,
                scaleY: 0.8,
                top: [165, 50],
                opacity: '40%',
                begin: function() {
                    window.$(cards[activeItem]).css("visibility", "visible");
                },
                duration: 3000,
            })
            anime({
                targets: cards[circularAddition(activeItem, 1, cards.length)],
                scaleX: 1,
                scaleY: 1,
                top: [295, 165],
                opacity: '100%',
                begin: function() {
                    window.$(cards[circularAddition(activeItem, 1, cards.length)]).css("visibility", "visible");
                },
                duration: 3000,
            })
            anime({
                targets: cards[circularAddition(activeItem, 2, cards.length)],
                scaleX: 0.8,
                scaleY: 0.8,
                opacity: '40%',
                top: [350, 295],
                begin: function() {
                    window.$(cards[circularAddition(activeItem, 2, cards.length)]).css("visibility", "visible");
                },
                duration: 3000,
            })
            
            return circularAddition(activeItem, 1, cards.length)
        })
    }


    return (
        <div className="unit-list-wrap">
            <div className="unit-list-upper-arrow-area area-active">
                <i onClick={up} className="fas fa-caret-up fa-lg unit-list-upper-arrow arrow-active"></i>
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
                <i onClick={down} className="fas fa-caret-down fa-lg unit-list-lower-arrow arrow-active"></i>
            </div>
        </div>
    );
}
export default UnitList;
