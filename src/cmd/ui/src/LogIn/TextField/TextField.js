import React, { useEffect, useState } from 'react';
import anime from 'animejs';
import './TextField.css';
import Typewriter from 'typewriter-effect/dist/core';

const states = {
    USERNAME: "username",
    PASSWORD: "password",
    DONE: "done"
}

function TextField({onLogIn}) {
    const [logInState, setLoginState] = useState(states.USERNAME)

    useEffect(() => {
        window.$('.back-text').bind('click', backToUsername)
        const input = document.querySelector('input');


        anime({
            targets: input,
            width: '255px',
            easing: 'easeInOutSine',
            duration: 300,
        });
    })

    var backToUsername = () => {
        setLoginState(states.USERNAME)
        window.$('input').css({
            'background': 'linear-gradient(to bottom, rgba(0,0,0,0.7) 0%, rgba(0,0,0,0.2) 100%)',
            'border-color': 'rgba(43,237,230,0.7)'
        }).attr('type', 'text')

        window.$('.back-text').css({ 'visibility': 'hidden' })
        window.$('.identification-text').css({ 'color': '#30fffe' })
        window.$('.logIn-root').css({ 'background': 'radial-gradient(circle, rgba(14,63,87,0.9164040616246498) 0%, rgba(0,0,0,0.9472163865546218) 81%)' })

        window.$('.identification-text').animate({ 'opacity': 0 }, 100, function () {
            window.$(this).html('PLEASE IDENTIFY YOURSELF').animate({ 'opacity': 1, 'color': '#30FFFE' });
        });
        window.$('input').val('');
    }

    var wrongData = () => {
        window.$('input').css({
            'background': 'linear-gradient(to bottom, rgba(237,55,55,0.1) 0%, rgba(237,55,55,0.2) 100%)',
            'border-color': '#ed3737'
        })
        window.$('.identification-text').css({ 'color': '#FFF' })
        window.$('.logIn-root').css({ 'background': 'radial-gradient(circle, rgba(87, 14, 14, 0.916) 0%, rgba(0,0,0,0.9472163865546218) 81%)' })
    }

    var correctUsername = () => {
        window.$('input').css({
            'background': 'linear-gradient(to bottom, rgba(0,0,0,0.7) 0%, rgba(0,0,0,0.2) 100%)',
            'border-color': 'rgba(43,237,230,0.7)'
        }).attr('type', 'password')
        window.$('.back-text').css({ 'visibility': 'visible' })
        window.$('.identification-text').css({ 'color': '#30fffe' })
        window.$('.logIn-root').css({ 'background': 'radial-gradient(circle, rgba(14,63,87,0.9164040616246498) 0%, rgba(0,0,0,0.9472163865546218) 81%)' })

        window.$('.identification-text').animate({ 'opacity': 0 }, 100, function () {
            window.$(this).html('ENTER PASSWORD').animate({ 'opacity': 1, 'color': '#30FFFE' });
        });
        window.$('input').val('');
    }

    var accessGranted = () => {
        window.$('.back-text').css({ 'visibility': 'hidden' })
        window.$('input').css({
            'background': 'linear-gradient(to bottom, rgba(0,0,0,0.7) 0%, rgba(0,0,0,0.2) 100%)',
            'border-color': 'rgba(43,237,230,0.7)'
        })
        window.$('.identification-text').css({ 'color': '#30fffe' })
        window.$('.logIn-root').css({ 'background': 'radial-gradient(circle, rgba(14,63,87,0.9164040616246498) 0%, rgba(0,0,0,0.9472163865546218) 81%)' })

        window.$('.identification-text').animate({ 'opacity': 0 }).promise().then(function () {
            new Typewriter(document.querySelector('.access-text'), {
                loop: false,
                delay: 50,
                cursor: ''
            }).typeString('ACCESS GRANTED').start();
        });

        window.$('input').animate({ 'opacity': 0 }, 200);
    }

    var welcome = () => {
        window.$('input').animate({ 'opacity': 0 }).promise().then(function() {
            window.$('.identification-text').animate({ 'opacity': 0 }).promise().then(
                function () {
                    window.$('.home-unit-image').css('visibility', 'visible');
                    window.$('.home-unit-image').animate({ 'opacity': 1 }, 300);
                    new Typewriter(document.querySelector('.hello-text'), {
                        loop: false,
                        delay: 75,
                        cursor: ''
                    }).typeString('HELLO, AHMED').start();
                    new Typewriter(document.querySelector('.welcome-text'), {
                        loop: false,
                        delay: 75,
                        cursor: ''
                    }).typeString('WELCOME BACK').start();
                }
            )
        });
        
        setTimeout(function() {onLogIn()}, 4000)
    }

    var enterClick = (event) => {
        var name = "Ahmed"
        var pass = "12345"
        if (event.key === 'Enter') {
            setLoginState(() => {
                switch (logInState) {
                    case states.USERNAME:
                        if (window.$('input').val() === name) {
                            correctUsername();
                            return states.PASSWORD;
                        } else {
                            wrongData();
                            return states.USERNAME;
                        }
                    case states.PASSWORD:
                        if (window.$('input').val() === pass) {
                            accessGranted();
                            welcome();
                            return states.DONE;
                        } else {
                            wrongData();
                            return states.PASSWORD;
                        }
                    default:
                        break;
                }
            })
        }
    }

    return (
        <input type="text" className="text-field" id="fname" name="fname"
            onKeyDown={enterClick}></input>
    );

}

export default TextField;
