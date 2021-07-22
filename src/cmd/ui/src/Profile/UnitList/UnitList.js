import React, { useState, useEffect } from 'react';
import UnitItem from "../UnitItem/UnitItem.js"
import { getUnits } from '../../Api/Api.js';
import anime from 'animejs'

import './UnitList.css';

function UnitList({ onChange, port }) {
    const [firstUnit, setFirstUnit] = useState(null)
    const [secondUnit, setSecondUnit] = useState(null)
    const [thirdUnit, setThirdUnit] = useState(null)
    const [activeUnit, setActiveUnit] = useState(0);
    const [units, setUnits] = useState([]);

    useEffect(() => {
        if ( units.length || !port ) return
        getUnits(port).then(initialData => {
            setUnits(units => {
                initialData.forEach(unit => {
                    units.push({ name: unit.name, ip: unit.ip})
                });
                
                onChange(units[0])
                setFirstUnit(units[units.length - 1])
                setSecondUnit(units[0])
                setThirdUnit(units[1])
                return units
            })
        })
    })

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

            anime({
                targets: cards[0],
                scaleX: 1,
                scaleY: 1,
                top: [50, 165],
                opacity: '100%',
                duration: 3000,
            })

            anime({
                targets: cards[1],
                scaleX: 0.8,
                scaleY: 0.8,
                top: [165, 295],
                opacity: '40%',
                duration: 3000,
            })

            setUnits(units => {
                setFirstUnit(units[circularSubtract(activeUnit, 1, units.length)])
                setSecondUnit(units[activeUnit])
                setThirdUnit(units[circularSubtract(activeUnit, 2, units.length)])
                onChange(units[circularSubtract(activeUnit, 1, units.length)])
                return units
            })

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

            anime({
                targets: cards[2],
                scaleX: 1,
                scaleY: 1,
                top: [295, 165],
                opacity: '100%',
                duration: 3000,
            })

            anime({
                targets: cards[0],
                scaleX: 0.8,
                scaleY: 0.8,
                opacity: '40%',
                top: [350, 295],
                duration: 3000,
            })

            setUnits(units => {
                setFirstUnit(units[circularAddition(activeUnit, 2, units.length)])
                setSecondUnit(units[activeUnit])
                setThirdUnit(units[circularAddition(activeUnit, 1, units.length)])
                onChange(units[circularAddition(activeUnit, 1, units.length)])
                return units
            })
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
