import React from 'react';
import {
  ThemeProvider,
  createTheme,
  SoundsProvider,
  createSounds
} from 'arwes';

import createAppTheme from './createAppTheme';
import Template from './Template';

import small from './static/img/background-small.jpg';
import medium from './static/img/background-medium.jpg';
import large from './static/img/background-large.jpg';
import xlarge from './static/img/background-xlarge.jpg';
import pattern from './static/img/glow.png'

import click from './static/sound/click.mp3'
import typing from './static/sound/typing.mp3'
import deploy from './static/sound/deploy.mp3'

const resources = {
  background: {
    small: small,
    medium: medium,
    large: large,
    xlarge: xlarge
  },
  pattern: pattern,
};

const sounds = {
  shared: {
    volume: 0.6,
  },
  players: {
    click: {
      sound: { src: [click] },
      settings: { oneAtATime: true }
    },
    typing: {
      sound: { src: [typing] },
      settings: { oneAtATime: true }
    },
    deploy: {
      sound: { src: [deploy] },
      settings: { oneAtATime: true }
    },
  }
};

const GlobalTemplate = (App) => {
  return (props) => (
    <ThemeProvider theme={createTheme(createAppTheme())}>
      <SoundsProvider sounds={createSounds(sounds)}>
        <Template>
          <App resources={resources} {...props} />
        </Template>
      </SoundsProvider>
    </ThemeProvider>
  );
};

export default GlobalTemplate;
