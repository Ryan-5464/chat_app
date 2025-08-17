window.addEventListener("click", function (e) {
    document.querySelectorAll(".modal").forEach(modal => {
        const modalContent = modal.querySelector(".modal-content");
        if (modal.classList.contains("open") && !modalContent.contains(e.target)) {
            modal.__controller.close()
        }
    });
});

// dataKeys are selectors for data needed to request required data from server.
const modalRegistry = {
    message: {
        modalId: 'messageModal',
        controller: MessageModalController,
        dataKeys: {ChatId: 'chatid', UserId: 'userid', MessageId: 'messageid'},
        modalOpenOn: '.message',
        attachListenerTo: '#messages-container',
    },
    chat: {
        modalId: 'chatModal',
        controller: ChatModalController,
        dataKeys: {ChatId: 'chatid'},
        modalOpenOn: '.chat',
        attachListenerTo: '#chats-container',
    },
    contact: {
        modalId: 'contactModal',
        controller: ContactsModalController,
        dataKeys: {},
        modalOpenOn: '.contact',
        attachListenerTo: '#contacts-container',
    }
}

document.addEventListener("DOMContentLoaded", function () {
    Object.values(modalRegistry).forEach(config => {
        AttachDelegatedContextModal(config)
    })
});

function AttachDelegatedContextModal(config) {
    const container = document.querySelector(config.attachListenerTo)
    container.addEventListener("contextmenu", function (e) {
        
        const elem = e.target.closest(config.modalOpenOn)

        if (elem) {
            e.preventDefault()
            const modal = ConfigureModal(elem, config)
            modal.openAt(e.clientX, e.clientY)
        }
    })
}

function ConfigureModal(elem, config) {
    const modal = document.getElementById(config.modalId)

    let buttonData = {}
    Object.entries(config.dataKeys).forEach(([key, value]) => {
        buttonData[key] = elem.getAttribute(`data-${value}`)
    });

    return AttachModalController(modal, config, buttonData)
}

function AttachModalController(modal, config, buttonData={}) {
    if (!modal.__controller) {
        modal.__controller = new config.controller(modal);
    }
    modal.__controller.configureButtons(buttonData)
    return modal.__controller;
    }

function replaceWithInput(elem, placeholder) {
    elem.innerHTML = ''
    const input = document.createElement('input');
    input.id = 'chat-name-input'
    input.type = 'text';
    input.name = 'Name';
    input.placeholder = placeholder;
    input.required = true;
    input.className = 'input-elem';
    elem.appendChild(input);
    return input
}

