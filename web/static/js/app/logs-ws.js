let loc = window.location;
let uri = 'ws:';

if (loc.protocol === 'https:') {
    uri = 'wss:';
}

uri += '//' + loc.host;
uri += '/logs/ws';

ws = new WebSocket(uri);

ws.onopen = function () {
    console.log('Websocket: Connecting to logging');
};

ws.onclose = function () {
    console.log('Websocket: Disconnected from logging');
};

ws.onmessage = function (evt) {
    let logEvent = JSON.parse(evt.data);
    let logTable = document.getElementById('logs');

    logTable.innerHTML += '<tr>' +
        '<td>' + logEvent.Time + '</td>' +
        '<td>' + logEvent.Level + '</td>' +
        '<td>' + logEvent.Component + '</td>' +
        '<td>' + logEvent.Message + '</td>' +
        '</tr>';

    scrollPageToBottom();
};

function scrollPageToBottom() {
    if ($(window).scrollTop() + $(window).height() > $(document).height() - 100) {
        let scrollingElement = (document.scrollingElement || document.body);
        scrollingElement.scrollTop = scrollingElement.scrollHeight;
    }
}