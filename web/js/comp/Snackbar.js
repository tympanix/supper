import React, { Component } from 'react';

import { EventEmitter} from 'events'

let snackEvent = new EventEmitter()

function push(type, title, message, tag) {
  snackEvent.emit("change", {
    type: type,
    title: capitalizeFirstLetter(title),
    message: capitalizeFirstLetter(message),
    tag: tag,
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
      created: new Date().getTime(),
      timer: this.removeSnackTimer(_id)
    })

    if (snack.tag !== undefined) {
      this.upsertSnack(snack)
    } else {
      this.insertSnack(snack)
    }

  }

  insertSnack(snack) {
    this.setState((prev) => {
      return {snacks: prev.snacks.concat([snack])}
    })
  }

  upsertSnack(snack) {
    this.setState((prev) => {
      let found = prev.snacks.find(s => s.tag === snack.tag)
      if (found) {
        return {snacks: prev.snacks.map(s => {
          if (s === found) {
            clearTimeout(found.timer)
            return Object.assign(found, snack)
          } else {
            return s
          }
        })}
      } else {
        return {snacks: prev.snacks.concat([snack])}
      }
    })
  }

  removeSnackTimer(id) {
    return setTimeout(() => {
      this.setState((prev) => {
        return {snacks: prev.snacks.filter((s) => s.id !== id)}
      })
    }, 5000)
  }

  componentWillMount() {
    snackEvent.on("change", this.handleSnack)
  }

  componentWillUnmount() {
    snackEvent.removeListener("change", this.handleSnack)
  }

  static notify(title, message, tag) {
    push("notify", title, message, tag)
  }

  static error(title, message, tag) {
    push("error", title, message, tag)
  }

  static success(title, message, tag) {
    push("success", title, message, tag)
  }

  static warning(title, message, tag) {
    push("warning", title, message, tag)
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
