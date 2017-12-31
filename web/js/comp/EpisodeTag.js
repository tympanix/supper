import React, { Component } from 'react';

class EpisodeTag extends Component {
  render() {
    let file = this.props.media

    if (file.type == "show") {
      return <span className="tag">S{file.media.season}E{file.media.episode}</span>
    } else {
      return null
    }
  }
}

export default EpisodeTag
