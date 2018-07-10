import { EventEmitter } from 'events'

import Snackbar from './comp/Snackbar'
import { configStore } from './stores'
import { join } from 'path'

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
    this.__listen()
    this.__connect()
    this.__id = 0
    this.__handlers = []
  }

  __create() {
    this.__proto = document.location.protocol === "http:" ? "ws://" : "wss://"
    this.__path = join(document.location.host, this.__proxypath || "", "/api/wss")
    this.ws = new WebSocket(this.__proto + this.__path)
    this.ws.onerror = function(evt) {
      Snackbar.error("Websocket", "Could not instantiate connection :(", "WS")
    }
    this.ws.onmessage = this.__handle.bind(this)
  }

  __connect() {
    this.__proxypath = configStore.getAll().proxypath
    if (this.__proxypath) {
      this.__create()
    }
  }

  __listen() {
    let self = this
    configStore.on("change", () => {
      self.__connect()
    })
  }

  __handle(event) {
    console.log(event.data)
    var data = this.__json(event.data)

    if (data === undefined) {
      return Snackbar.error("Websocket", "Could not read websocket message", "WS")
    }

    Object.assign(data, {wsid: this.__id})
    this.__id++
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
