import React, { Component } from 'react';

class FileTags extends Component {
  render() {
    let meta = this.props.media.media.metadata

    let keys = ["quality", "source"]

    let tags = keys.map((k) => {
      return meta[k] && <span key={k} className="tag">{meta[k]}</span>
    })

    return (
      <span className="tags">
        {tags}
      </span>
    )
  }
}

export default FileTags
