"use strict";

let button = document.querySelector("#send-button");
let textbox = document.querySelector("#textbox");

button.onclick = function () {
    let text = textbox.value;
    // if (text == "") {
    //     return
    // }
    fetch("./send", { method: "POST", body: text });
    textbox.value = "";
}

let messageTemplate =
    `<li class="message">
<div class="message-header">from {sender} at {time}</div>
<div class="message-text">{text}</div>
</li>`;

let messageBox = document.querySelector(".message-box")

function addMessage(obj) {
    let html =
        messageTemplate
            .replace("{sender}", obj.sender)
            .replace("{text}", obj.text)
            .replace("{time}", obj.time);
    
    messageBox.insertAdjacentHTML("beforeend", html)
}

function update() {
    fetch("./update", { method: "POST", body: String(messageBox.children.length)})
        .then(
            response => response.json(),
            error => alert(error)
        )
        .then(
            newElements => newElements.forEach(el => addMessage(el)),
            error => alert(error)
        );
}

setInterval(update, 500);
