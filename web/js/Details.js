import React, { Component } from 'react';
import axios from 'axios';

import Search from './Search'


class Details extends Component {
  constructor() {
    super()

    this.state = {
      media: undefined
    }
  }

  componentWillMount() {
    try {
      let media = this.props.location.state.media
      if (media) {
        this.setState({media})
      }
    } catch (e) {}
  }

  renderMediaList() {
    let media = this.state.media.map((item) => {
      return <h3>{item.media.name}</h3>
    })

    return (
      <div>{media}</div>
    )
  }

  render() {
    if (this.state.media) {
      return this.renderMediaList()
    }

    return (
      <h1>No media found</h1>
    );
  }

}

export default Details
