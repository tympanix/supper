import React, { Component } from 'react';
import axios from 'axios';

import Search from './Search'
import Spinner from './Spinner'

import API from '../api'

class Details extends Component {
  constructor() {
    super()

    this.state = {
      media: undefined,
      folder: undefined,
      failed: false,
    }
  }

  componentWillMount() {
    let folder = this.getLocationState()
    console.log(folder)
    API.getMediaDetails(folder)
      .then((media) => this.setState({media: media}))
      .catch(() => this.setState({failed: true}))
  }

  getLocationState() {
    try {
      let folder = this.props.location.state.folder
      if (folder) {
        this.setState({folder: folder})
        return folder
      }
    } catch (e) {
      this.setState({failed: true})
    }
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
    if (this.state.failed) {
      return <h1>No media found</h1>
    }

    if (this.state.media) {
      return this.renderMediaList()
    } else if (this.state.folder) {
      return <Spinner/>
    }

  }

}

export default Details
