import React, { Component } from 'react';

import flags from '../flags'

class Flag extends Component {

  render() {
    let code = this.props.lang
    let flag = flags[code] || ''
    
    return <span className={`flag-icon flag-icon-${flag}`}></span>
  }
}

export default Flag
