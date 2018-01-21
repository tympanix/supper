import React, { Component } from 'react';
import { CSSTransition } from 'react-transition-group'

class Spinner extends Component {
  render() {
    if (this.props.visible === false) {
      return null
    }
    return (
        <div className="overlay">
          <div className="spinner"></div>
        </div>
    )
  }
}

export default Spinner
