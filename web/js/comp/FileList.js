import React, { Component } from 'react'

import Subtitles from './Subtitles'
import FileTags from './FileTags'
import EpisodeTag from './EpisodeTag'

class FileList extends Component {
  constructor() {
    super()

    this.state = {

    }
  }

  componentWillMount() {
    try {
      this.setState({files: this.props.files})
    } catch (e) {
      this.setState({failed: true})
    }
  }

  render() {
    let files = this.state.files.map((f) => {
      return (
        <li key={f.filepath}>
          <span className="name">{f.media.name}</span>
          <FileTags media={f}/>
          <span className="right">
            <Subtitles list={f.subtitles}/>
          </span>
        </li>
      )
    })

    return (
      <ul className="file-list">
        {files}
      </ul>
    )
  }
}

export default FileList
