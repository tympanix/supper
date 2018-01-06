import React, { Component } from 'react';

import Flag from './Flag'

class Subtitles extends Component {
  click(event, lang) {
    this.props.onClick && this.props.onClick(event, lang)
  }

  render() {
    let subs = this.props.list.map((s) => {
      return (
        <span key={s.language} data-tooltip={s.language}>
          <Flag lang={s.code} onClick={this.click.bind(this)} />
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
