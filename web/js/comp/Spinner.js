import React, { Component } from 'react';
import { CSSTransition } from 'react-transition-group'

class Spinner extends Component {
  render() {
    if (this.props.visible === false) {
      return null
    }
    return (
      <CSSTransition
        classNames="fade"
        timeout={1000}>
        <div className="overlay">
          <div className="spinner"></div>
        </div>
      </CSSTransition>
    )
  }
}

export default Spinner
