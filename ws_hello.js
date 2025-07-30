window.onload = function() {
	var socket = new WebSocket("wss://localhost:1024/hello");
	socket.onmessage = m => {
	div = document.createElement("div");
	div.innerText = JSON.parse(m.data);
	document.body.append(div);
	};
}
