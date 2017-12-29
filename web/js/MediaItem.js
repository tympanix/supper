import React, { Component } from 'react';
import axios from 'axios';

export default class MediaItem extends Component {
  render() {

    return (
      // Add your component markup and other subcomponent references here.
      <li>{this.props.item.name}</li>
    );
  }
}
