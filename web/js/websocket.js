import Snackbar from './comp/Snackbar'

var ws = new WebSocket("ws://" + document.location.host + "/api/ws")

const loggers = {
  "debug": Snackbar.notify,
  "info": Snackbar.success,
  "warn": Snackbar.warning,
  "error": Snackbar.error,
  "fatal": Snackbar.error,
}

ws.onmessage = function(event) {
  console.log(event.data)
  var data
  try {
    data = JSON.parse(event.data)
  } catch (e) {
    console.log(e)
    return Snackbar.error("Websocket", "Could not read websocket message", "WS")
  }

  var log = loggers[data.level] || Snackbar.error
  log("Update", data.message, data.data.job)
}
