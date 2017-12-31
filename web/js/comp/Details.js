import React, { Component } from 'react';

import Search from './Search'
import Spinner from './Spinner'
import FileList from './FileList'
import Flag from './Flag'
import DownloadButtons from './DownloadButtons'

import API from '../api'

class Details extends Component {
  constructor() {
    super()

    this.state = {
      media: undefined,
      folder: undefined,
      busy: false,
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

  downloadSubtitles(lang) {
    this.setState({busy: true})
    let folder = this.state.folder
    API.downloadSubtitles(folder, lang).then((data) => {
      this.setState({media: data})
    }).finally(() => {
      this.setState({busy: false})
    })
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

          <section className="dark">
            <header>
              <h3 className="center">Download Subtitles</h3>
            </header>
            <DownloadButtons
              disabled={this.state.busy}
              onDownload={this.downloadSubtitles.bind(this)}/>
          </section>

          <section>
            <h3>Files</h3>
            <FileList files={this.state.media}/>
            <Spinner visible={this.state.busy}/>
          </section>
        </section>
      )
    }

  }

}

export default Details
