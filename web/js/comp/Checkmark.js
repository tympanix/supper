import React, { Component } from 'react';
import { CSSTransition } from 'react-transition-group'
import { EventEmitter } from 'events'

var checkmarkEvent = new EventEmitter()

class Checkmark extends Component {
  constructor() {
    super()
    this.state = {
      show: false
    }
    this.show = this.show.bind(this)
  }

  static show() {
    checkmarkEvent.emit('show')
  }

  show() {
    if (this.state.show === true) {
      return
    }
    this.setState({show: true})
    setTimeout(() => {
      this.setState({show: false})
    }, 2000)
  }

  componentWillMount() {
    checkmarkEvent.on('show', this.show)
  }

  componentWillUnmount() {
    checkmarkEvent.removeListener('show', this.show)
  }

  render() {
    if (!this.state.show) {
      return null
    }

    return (
      <div className="overlay transparent">
        <svg className="checkmark" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 52 52">
          <circle className="checkmark__circle" cx="26" cy="26" r="25" fill="none"/>
          <path className="checkmark__check" fill="none" d="M14.1 27.2l7.1 7.2 16.7-16.8"/>
        </svg>
      </div>
    )
  }
}

export default Checkmark
