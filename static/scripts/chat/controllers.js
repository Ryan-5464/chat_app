
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
        this.renderer = new Renderer()
    }

    editName(e, { ChatId }) {
        const chat = document.querySelector(`[data-chatid="${ChatId}"]`);
        const chatName = chat.querySelector('.chat-name');
        const input = replaceWithInput(chatName, "Enter new name");
        input.focus();
        this.close();

        input.addEventListener('keydown', (e) => {
            if (e.key === "Enter") {
                EditChatNameRequest(input.value, ChatId).then(data => {
                    if (data) {
                        chatName.innerHTML = data.Name;
                    }
                });
            }
        })
    }

    leaveChat(e, { ChatId }) {
        this.close()
        
        let isActive = false
        if (e.target.classList.contains('active')) {
            isActive = true
        }

        LeaveChatRequest(ChatId).then(data => {
            if (!isActive) {
                return
            }
            this.renderer.render('chats', data.Chats, true)
            this.renderer.render('messages', data.Messages, true)  
            const chat = document.querySelector(`[data-chatid="${data.NewActiveChatId}"]`);
            chat.classList.add('active')
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

// CHAT CONTROLLERS ============================================================

const chatControllerRegistry = {
    Chat: {
        containerId: 'chats-container',
        elemClass: '.chat',
        dataId: 'chatid',
        elemType: 'Chat'
    },
    Contact: {
        containerId: 'contacts-container',
        elemClass: '.contact',
        dataId: 'contactchatid',
        elemType: 'Contact'
    }
}

function InitializeChatControllers() {
    for (const config of Object.values(chatControllerRegistry)) {
        const controller = new ChatController(renderer, config)
        attachChatController(controller)
    }
}

function attachChatController(controller) {
    controller.container.addEventListener('click', (e) => {
        const chatElem = e.target.closest(controller.config.elemClass)
        if (chatElem) {
            controller.handleChatActivation(chatElem)
        }
    })
}

class ChatController {
    constructor(renderer, config) {
        this.renderer = renderer
        this.config = config
        this.container = document.getElementById(config.containerId)
    }

    handleChatActivation(target) {
        const chatId = target.dataset[this.config.dataId]
        if (!chatId) {return}

        this._setActiveElemenBytId(chatId)
        this._loadMessages(chatId)
    }

    _setActiveElemenBytId(targetId) {
        const elems = this.container.querySelectorAll(this.config.elemClass)
        for (const chat of elems) {
            isTarget = chat.dataset[this.config.dataId] === targetId
            chat.classList.toggle('active', isTarget)
        }
    }

    _loadMessages(targetId) {
        SwitchChatRequest(this.config.elemType, targetId)
        .then(data => {
            this.renderer.render(data.Messages, true)
        })
        .catch(err => {
            console.log("Failed to handle chat activation: ", err)
        })
    }
}

