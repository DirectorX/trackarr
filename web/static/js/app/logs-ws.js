// set socketUrl
let socketUrl = new URL('/ws', window.location.href);
socketUrl.protocol = socketUrl.protocol.replace('http', 'ws');

// create socket
let socket = new WebSocket(socketUrl.href);

socket.onopen = () => {
    console.log('Socket connected');
    socket.send(JSON.stringify({type: 'subscribe', 'data': 'logs'}))
};

socket.onclose = () => {
    console.log('Socket disconnected');
};

socket.onerror = error => {
    console.log(`Socket error: ${error}`);
};

socket.onmessage = e => {
    let event = JSON.parse(e.data);

    if (event.type !== 'log') {
        // ignore messages not of the type: log
        return;
    }

    let logTable = document.getElementById('logs');
    logTable.innerHTML += '<tr>' +
        '<td>' + event.data.time + '</td>' +
        '<td>' + event.data.level + '</td>' +
        '<td>' + event.data.component + '</td>' +
        '<td>' + event.data.message + '</td>' +
        '</tr>';

    scrollPageToBottom();
};


function scrollPageToBottom() {
    if ($(window).scrollTop() + $(window).height() > $(document).height() - 100) {
        let scrollingElement = (document.scrollingElement || document.body);
        scrollingElement.scrollTop = scrollingElement.scrollHeight;
    }
}