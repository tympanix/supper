import React, { Component } from 'react';
import axios from 'axios';

import Search from './Search'
import Spinner from './Spinner'
import FileList from './FileList'

import API from '../api'

class Details extends Component {
  constructor() {
    super()

    this.state = {
      media: undefined,
      folder: undefined,
      loading: true,
      failed: false,
    }
  }

  componentWillMount() {
    let folder = this.getLocationState()
    API.getMediaDetails(folder)
      .then((media) => this.setState({media: media}))
      .catch(() => this.setState({failed: true}))
      .finally(() => this.setState({loading: false}))
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

    if (this.state.loading) {
      return <Spinner/>
    }

    if (this.state.media) {
      return (
        <section>
          <header>
            <h1>{this.state.folder.name}</h1>
          </header>

          <h3>Files</h3>
          <div className="flex">
            <FileList files={this.state.media}/>
          </div>
        </section>
      )
    }

  }

}

export default Details
