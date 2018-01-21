import React, { Component } from 'react';

import subtitleStore from '../stores/subtitle_store'

class ScoreTag extends Component {
  constructor() {
    super()
  }

  render() {
    let percent = this.props.score
    let color = 120*percent

    let styles = {
      backgroundColor: `hsl(${color}, 80%, 45%)`
    }

    return <span className="tag" style={styles}>{Math.round(percent*100)}%</span>
  }
}

export default ScoreTag