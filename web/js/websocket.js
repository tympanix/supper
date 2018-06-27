import Snackbar from './comp/Snackbar'

var ws = new WebSocket("ws://" + document.location.host + "/api/ws")

ws.onmessage = function(event) {
  console.log(event.data)
  var data
  try {
    data = JSON.parse(event.data)
  } catch (e) {
    console.log(e)
    return Snackbar.error("Websocket", "Could not read websocket message")
  }
  if (data.level == "info") {
    Snackbar.notify("Update", data.message)
  } else if (data.level == "error") {
    Snackbar.error("Update", data.message)
  } else if (data.level == "warn") {
    Snackbar.warning("Update", data.message)
  } else {
    Snackbar.error("Update", data.message)
  }
}
