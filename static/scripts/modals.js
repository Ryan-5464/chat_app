window.addEventListener("click", function (e) {
    document.querySelectorAll(".modal").forEach(modal => {
        const modalContent = modal.querySelector(".modal-content");
        if (modal.classList.contains("open") && !modalContent.contains(e.target)) {
            modal.__controller.close()
        }
    });
});


class ModalController {
    constructor(modal) {
        this.modal = modal;
        this.modalContent = modal.querySelector('.modal-content');
    }

    close() {
        this.modalContent.classList.remove("opening");
        this.modalContent.classList.add("closing");

        setTimeout(() => {
        this.modal.classList.remove("open");
        this.modalContent.classList.remove("closing");
        }, MODAL_CLOSE_DELAY);
    }

    openAt(clientX, clientY) {
        document.querySelectorAll(".modal.open").forEach(openModal => {
            if (openModal !== this.modal && openModal.__controller) {
                openModal.__controller.close();
            }
        });

        const padding = 10;
        const maxLeft = window.innerWidth - this.modalContent.offsetWidth - padding;
        const maxTop = window.innerHeight - this.modalContent.offsetHeight - padding;

        this.modalContent.style.left = Math.min(clientX, maxLeft) + "px";
        this.modalContent.style.top = Math.min(clientY, maxTop) + "px";

        this.modalContent.classList.remove("opening", "closing");
        void this.modalContent.offsetWidth;

        this.modal.classList.add("open");
        this.modalContent.classList.add("opening");
    }

    configureButtons(buttonData = {}) {
        const buttons = this.modalContent.querySelectorAll('[data-action]');

        buttons.forEach(button => {
            const actionName = button.dataset.action;
            button = RemoveAllListeners(button)

            button.addEventListener('click', (e) => {
                this[actionName](e, buttonData)
            });
        })
    }
}

class ChatModalController extends ModalController {
    constructor(modal) {
        super(modal)
    }

    editName(e, { ChatId }) {
        const chat = document.querySelector(`[data-chatid="${ChatId}"]`);
        const chatName = chat.querySelector('.chat-name');
        const input = replaceWithInput(chatName, "Enter new name");
        input.focus();
        this.close();

        input.addEventListener('keydown', (e) => {
            if (e.key === "Enter") {
                EditChatNameRequest(input.value, ChatId).then(newName => {
                    if (newName) {
                        chatName.innerHTML = newName;
                    }
                });
            }
        })
    }
}

class MessageModalController extends ModalController {
    constructor(modal) {
        super(modal)
    }

    _deleteMessage(e, { ChatId, MessageId, UserId }) {
        return ChatId, MessageId, UserId
    }
}

class ContactsModalController extends ModalController {
    constructor(modal) {
        super(modal)
    }

    configureButtons() {
        return
    }
}

const modalRegistry = {
    message: {
        modalId: 'messageModal',
        controller: MessageModalController,
        dataKeys: {ChatId: 'chatid', UserId: 'userid', MessageId: 'messageid'},
        modalOpenOn: '.message',
        attachListenerTo: '#messages-container'
    },
    chat: {
        modalId: 'chatModal',
        controller: ChatModalController,
        dataKeys: {ChatId: 'chatid'},
        modalOpenOn: '.chat',
        attachListenerTo: '#chats-container'
    },
    contact: {
        modalId: 'contactModal',
        controller: ContactsModalController,
        dataKeys: {},
        modalOpenOn: '.contact',
        attachListenerTo: '#contacts-container'
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

    return AttachModalController(modal, config.controller, buttonData)
}

function AttachModalController(modal, controllerClass, buttonData={}) {
    if (!modal.__controller) {
        modal.__controller = new controllerClass(modal);
    }
    modal.__controller.configureButtons(buttonData)
    return modal.__controller;
    }

function replaceWithInput(elem, placeholder) {
    elem.innerHTML = ''
    const input = document.createElement('input');
    input.type = 'text';
    input.name = 'Name';
    input.placeholder = placeholder;
    input.required = true;
    input.className = 'input-elem';
    elem.appendChild(input);
    return input
}

