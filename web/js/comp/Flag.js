import React, { Component } from 'react';

import flags from '../flags'

class Flag extends Component {
  clickHandler(lang) {
    return (event) => {
      this.props.onClick && this.props.onClick(event, lang)
    }
  }

  render() {
    let code = this.props.lang
    let flag = flags[code] || 'unknown'

    return (
      <span
        className={`flag-icon flag-icon-${flag}`}
        onClick={this.clickHandler(this.props.lang).bind(this)}
      ></span>
    )
  }
}

export default Flag
