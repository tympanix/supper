import React, { Component } from 'react'

import Subtitles from './Subtitles'
import FileTags from './FileTags'
import EpisodeTag from './EpisodeTag'

class FileList extends Component {
  constructor() {
    super()

    this.state = {
      files: []
    }
  }

  componentWillMount() {
    try {
      this.setState({files: this.props.files})
    } catch (e) {
      this.setState({failed: true})
    }
  }

  componentWillReceiveProps(props) {
    this.setState({files: props.files})
  }

  render() {
    let files = this.state.files.map((f) => {
      return (
        <li key={f.filepath}>
          <span className="name">{f.media.name}</span>
          <div className="content">
            <FileTags media={f}/>
            <span className="right">
              <Subtitles list={f.subtitles}/>
            </span>
          </div>
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
