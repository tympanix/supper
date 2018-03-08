import React, { Component } from 'react';

import EpisodeTag from './EpisodeTag'

const UNKW = "Unknown"

class FileTags extends Component {
  render() {
    let media = this.props.media
    let meta = media.media.metadata

    let keys = ["quality", "source", "codec", "group"]

    let tags = keys.map((k) => {
      return meta[k] && <span key={k} className="tag">{meta[k]}</span>
    })

    return (
      <span className="tags">
        <EpisodeTag media={media} />
        {tags}
      </span>
    )
  }
}

export default FileTags
