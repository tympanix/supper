var loc = window.location, new_uri;
if (loc.protocol === "https:") {
    new_uri = "wss:";
} else {
    new_uri = "ws:";
}
new_uri += "//" + loc.host;
new_uri += loc.pathname + "ws";

var ws = new WebSocket("ws://" + document.location.host + "/api/ws")

ws.onmessage = function(event) {
  console.log(event.data)
}
