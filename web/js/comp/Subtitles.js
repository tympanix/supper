import React, { Component } from 'react';

import Flag from './Flag'

class Subtitles extends Component {
  render() {
    let subs = this.props.list.map((s) => {
      return (
        <span key={s.language} data-tooltip={s.language}>
          <Flag lang={s.code}/>
        </span>
      )
    })
    if (!subs || subs.length === 0) {
      return <span className="tag">No subtitles</span>
    } else {
      return (
        <span className="tags">
          {subs}
        </span>
      )
    }
  }
}

export default Subtitles
