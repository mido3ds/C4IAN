import React, { useState, useEffect } from 'react';
import UnitItem from "../UnitItem/UnitItem.js"
import './UnitList.css';
import anime from 'animejs'
import { units }  from '../../units'

function UnitList({onChange}) {
    const [firstUnit, setFirstUnit] = useState(units[units.length - 1])
    const [secondUnit, setSecondUnit] = useState(units[0])
    const [thirdUnit, setThirdUnit] = useState(units[1])
    const [activeUnit, setActiveUnit] = useState(0);

    useEffect(() => {
       onChange(units[0])
    },[])

    var circularAddition = (Augend, Addend, len) => {
        return (Augend + Addend) % len;
    }

    var circularSubtract = (Minuend, Subtrahend, len) => {
        return (Minuend - Subtrahend + len) % len
    }

    var down = () => {
        var cards = window.$('.unit-item-container').toArray()
        setActiveUnit(() => {
            anime({
                targets: cards[2],
                scaleX: 0.8,
                scaleY: 0.8,
                top: [-100, 50],
                opacity: '40%',
                duration: 3000,
            })
            setThirdUnit(units[circularSubtract(activeUnit, 2, units.length)])
            
            anime({
                targets: cards[0],
                scaleX: 1,
                scaleY: 1,
                top: [50, 165],
                opacity: '100%',
                duration: 3000,
            })
            setFirstUnit(units[circularSubtract(activeUnit, 1, units.length)])

            anime({
                targets: cards[1],
                scaleX: 0.8,
                scaleY: 0.8,
                top: [165, 295],
                opacity: '40%',
                duration: 3000,
            })
            setSecondUnit(units[activeUnit])

            onChange(units[circularSubtract(activeUnit, 1, units.length)])
            return circularSubtract(activeUnit, 1, units.length)
        })
    }


    var up = () => {
        var cards = window.$('.unit-item-container').toArray()

        setActiveUnit(() => {
            anime({
                targets: cards[1],
                scaleX: 0.8,
                scaleY: 0.8,
                top: [165, 50],
                opacity: '40%',
                duration: 3000,
            })
            setSecondUnit(units[activeUnit])

            anime({
                targets: cards[2],
                scaleX: 1,
                scaleY: 1,
                top: [295, 165],
                opacity: '100%',
                duration: 3000,
            })
            setThirdUnit(units[circularAddition(activeUnit, 1, units.length)])

            anime({
                targets: cards[0],
                scaleX: 0.8,
                scaleY: 0.8,
                opacity: '40%',
                top: [350, 295],
                duration: 3000,
            })
            setFirstUnit(units[circularAddition(activeUnit, 2, units.length)])

            onChange(units[circularAddition(activeUnit, 1, units.length)])
            return circularAddition(activeUnit, 1, units.length)
        })
    }


    return (
        <div className="unit-list-wrap">
            <div className="unit-list-upper-arrow-area area-active">
                <i onClick={up} className="fas fa-caret-up fa-lg unit-list-upper-arrow arrow-active"></i>
            </div>
            <div id="card-slider" className="unit-list-area">
                <UnitItem unit={firstUnit} />
                <UnitItem unit={secondUnit} />
                <UnitItem unit={thirdUnit} />
            </div>
            <div className="unit-list-lower-arrow-area area-active">
                <i onClick={down} className="fas fa-caret-down fa-lg unit-list-lower-arrow arrow-active"></i>
            </div>
        </div>
    );
}
export default UnitList;
