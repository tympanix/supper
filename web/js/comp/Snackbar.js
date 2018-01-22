import React, { Component } from 'react';

import { EventEmitter} from 'events'

let snackEvent = new EventEmitter()

function push(type, title, message, delay) {
  snackEvent.emit("change", {
    type: type,
    title: capitalizeFirstLetter(title),
    message: capitalizeFirstLetter(message),
    delay: delay || 8000,
  })
}

function capitalizeFirstLetter(string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
}

let id = 0

class Snackbar extends Component {
  constructor() {
    super()

    this.handleSnack = this.handleSnack.bind(this)

    this.state = {
      snacks: []
    }
  }

  handleSnack(snack) {
    let _id = id++;

    snack = Object.assign({}, snack, {
      id: _id,
      created: new Date().getTime()
    })

    this.setState((prev) => {
      return {snacks: prev.snacks.concat([snack])}
    })

    setTimeout(() => {
      this.setState((prev) => {
        return {snacks: prev.snacks.filter((s) => s.id !== _id)}
      })
    }, 5000)
  }

  componentWillMount() {
    snackEvent.on("change", this.handleSnack)
  }

  componentWillUnmount() {
    snackEvent.removeListener("change", this.handleSnack)
  }

  static notify(title, message, delay) {
    push("notify", title, message, delay)
  }

  static error(title, message, delay) {
    push("error", title, message, delay)
  }

  static success(title, message, delay) {
    push("success", title, message, delay)
  }

  static warning(title, message, delay) {
    push("warning", title, message, delay)
  }

  render() {
    let snacks = this.state.snacks.map((s) => {
      return (
        <li className={s.type} key={s.id}>
          <span className="title">{s.title}</span>
          <span className="message">{s.message}</span>
        </li>
      )
    })

    return (
      <ul className="snackbar">
        {snacks}
      </ul>
    )
  }
}

export default Snackbar
