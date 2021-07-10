import React from 'react';

export default class Template extends React.Component {

  componentDidMount () {
    this.removeServerStyles();
  }

  render () {
    return this.props.children;
  }

  removeServerStyles () {
    const pagesStyles = document.querySelector('#pages-styles');
    if (pagesStyles) pagesStyles.remove();
  }
}
