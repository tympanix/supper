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

  subtitleClick(event, lang) {
    this.props.languageClicked &&
    this.props.languageClicked(event, lang)
  }

  render() {
    let files = this.state.files.map((f) => {
      return (
        <li className="flex" key={f.filepath}>
          <div className="col">
            <div className="name">{f.media.name}</div>
            <FileTags media={f}/>
          </div>
          <div className="flex center subtitles">
            <Subtitles list={f.subtitles}
              onClick={this.subtitleClick.bind(this)} />
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
