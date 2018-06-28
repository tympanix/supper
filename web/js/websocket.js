import { EventEmitter } from 'events'

import Snackbar from './comp/Snackbar'

const loggers = {
  "debug": Snackbar.notify,
  "info": Snackbar.success,
  "warn": Snackbar.warning,
  "error": Snackbar.error,
  "fatal": Snackbar.error,
}

class Websocket extends EventEmitter {

  constructor() {
    super()
    this.ws = new WebSocket("ws://" + document.location.host + "/api/ws")
    this.ws.onmessage = this.__handle.bind(this)
    this.__handlers = []
  }

  __handle(event) {
    console.log(event.data)
    var data = this.__json(event.data)

    if (data === undefined) {
      return Snackbar.error("Websocket", "Could not read websocket message", "WS")
    }

    this.emit("ws", data)

    var log = loggers[data.level] || Snackbar.error
    log("Update", data.message, data.data.job)
  }

  __json(data) {
    try {
      return JSON.parse(data)
    } catch (e) {
      console.log(e)
      return undefined
    }
  }

  subscribe(fn) {
    this.on("ws", fn)
    return function() {
      this.removeListener("ws", fn)
    }.bind(this)
  }

  remove(fn) {
    this.removeListener("ws", fn)
  }
}

export default new Websocket()
