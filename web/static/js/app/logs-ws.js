let socket = glue();

socket.on("connected", function () {
    console.log('Socket connected!');
});

socket.on("disconnected", function () {
    console.log('Socket disconnected...');
});

socket.on("error", function (e, msg) {
    console.log('Socket error: ' + msg);
});

socket.onMessage(function (data) {
    let event = JSON.parse(data);

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
});

function scrollPageToBottom() {
    if ($(window).scrollTop() + $(window).height() > $(document).height() - 100) {
        let scrollingElement = (document.scrollingElement || document.body);
        scrollingElement.scrollTop = scrollingElement.scrollHeight;
    }
}