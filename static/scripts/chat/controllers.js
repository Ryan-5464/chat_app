
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

function AttachChatControllers() {
    for (const config of Object.values(chatControllerRegistry)) {
        const container = document.getElementById(config.containerId) 
        container.__controller = new ChatController(config)
        container.addEventListener('click', (e) => {
            const chatElem = e.target.closest(container.__controller.config.elemClass)
            if (chatElem) {
                container.__controller.handleChatActivation(chatElem)
            }
        })
    }
}

class ChatController {
    constructor(config) {
        this.renderer = new Renderer()
        this.config = config
        this.container = document.getElementById(config.containerId)
    }

    handleChatActivation(target) {
        const chatId = target.dataset[this.config.dataId]
        if (!chatId) {return}

        this._setActiveElemenById(chatId)
        this._loadMessages(chatId)
    }

    _setActiveElemenById(targetId) {
        const elems = this.container.querySelectorAll(this.config.elemClass)
        for (const chat of elems) {
            const isTarget = chat.dataset[this.config.dataId] === targetId
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

// INPUT CONTROLLERS ===============================================================================

function AttachNewChatInputController() {
    const input = document.getElementById('chat-name-input')
    input.__controller = new NewChatInputController()
    input.addEventListener('keydown', () => {
        if ('keydown' !== 'Enter') {
            return
        }
        input.controller.createChat(input.value)
    })
}

class NewChatInputController {
    constructor(config, renderer = new Renderer()) {
        this.config = config
        this.renderer = renderer
    }

    createChat(name) {
        this._newChatRequest(name)
    }

    _newChatRequest(name) {
        NewChatRequest(name).then(data => {
            this.renderer.render('chats', data.Chats, false)
            this.renderer.render('messages', [], true)
        })
    }
}

function AttachNewContactInputController() {
    const input = document.getElementById('add-contact-input-container')
    input.__controller = new NewContactInputController()
    input.addEventListener('keydown', () => {
        if ('keydown' !== 'Enter') {
            return
        }
        input.__controller.addContact(input.value)
    })
}

class NewContactInputController {
    constructor(config, renderer = new Renderer()) {
        this.config = config
        this.renderer = renderer
    }

    createChat(name) {
        this._newChatRequest(name)
    }

    _newChatRequest(name) {
        NewChatRequest(name).then(data => {
            this.renderer.render('chats', data.Chats, false)
            this.renderer.render('messages', [], true)
        })
    }
}

const inputControllerRegistry = {
    newChat: {
        request: NewChatRequest,
        renderers: [
            {config: 'chats', overwrite: false},
            {config: 'messages', overwrite: true},
        ]
    },
    addContact: {
        request: AddContactRequest,
        renderers: [
            {config: 'contacts', overwrite: false},
            {config: 'messages', overwrite: true},
        ]
    },
    newMessage: {
        request: NewMessageRequest,
        renderers: [
            {config: 'messages', overwrite: false},
        ]
    }
}

const chatInputController = new InputController(inputControllerRegistry.chat)
chatInputController.send(input.value)

class InputController {
    constructor(config, renderer = new Renderer()) {
        this.config = config
        this.renderer = renderer
    }

    send(reqData) {
        this.config.request(reqData).then(resData => {
            RenderResponse(resData, this.renderer)
        })
    } 
}



function CreateNewChat(chatName) {
    const reqData = {
        Name: chatName,
    };

    const chatRenderer = Render.bind({
        containerId: 'chats-container',
        elemFactory: ChatElement,
        overwrite: true,
    });

    const messageRenderer = Render.bind({
        containerId: 'messages-container',
        elemFactory: MessageElement,
        overwrite: true,
    });

    const request = () => NewChatRequest(reqData);
    const responseCallbacks = {
        Chats: value => RenderResponse(chatRenderer, value),
        Messages: value => RenderResponse(messageRenderer, value),
    };

    HandleRequest(request, responseCallbacks);
}

function HandleRequest(req, callbacks) {
    req()
        .then(data => HandleResponse(data, callbacks))
        .catch(err => console.error("Request failed:", err));
}

function HandleResponse(data, callbacks) {
    Object.entries(data).forEach(([key, value]) => {
        if (callbacks[key]) {
            callbacks[key](value);
        } else {
            console.warn(`No callback found for key: ${key}`);
        }
    });
}


function Render(data) {
    const container = document.getElementById(this.containerId);
    if (this.overwrite) {
        container.innerHTML = '';
    }
    data.forEach(obj => {
        container.appendChild(this.elemFactory(obj))
    });
}

class NewChatHandler extends RequestHandler {
    constructor(chatElemFactory, messageElemFactory) {
        super('api/chat/new', 'POST')

        this.renderChats = Render.bind({
            containerId: 'chats-container',
            elemFactory: chatElemFactory,
            overwrite: true,
        });

        this.renderMessages = Render.bind({
            containerId: 'messages-container',
            elemFactory: messageElemFactory,
            overwrite: true,
        });

        this.responseHandlers = {
            Chats: value => this.renderChats(value),
            Messages: value => this.renderMessages(value),
        }
    }

    Create(chatName) {
        this._request({Name: chatName})
        .then(resData => {
            Object.entries(resData).forEach(([key, value]) => {
                const handler = this.responseHandlers[key];
                if (handler) {
                    handler(value);
                } else {
                    console.warn(`No handler found for key: ${key}`, value);
                }
            })
        })
        .catch(error => {
            console.log("Create failed => error: ", error)
        })
    }


}

class RequestHandler {
    constructor(endpoint, method) {
        this.endpoint = endpoint
        this.method = methods[method]
    }

    methods = {
        GET: () => this._GET(),
        DELETE: () => this._DELETE(),
        POST: json => this._POST(json)
    }

    _GET() {
        return {
            method: 'GET',
        }
    }

    _POST(json) {
        return {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(json),
        }
    }

    _DELETE() {
        return {
            method: 'DELETE',
        }
    }

    async _request(json={}) {
        try {
            let response
            if (json == {}) {
                response = await fetch(this.endpoint, this.method());
            } else {
                response = await fetch(this.endpoint, this.method(json));
            }
            if (!response.ok) throw new Error("Network response was not ok");

            const data = await response.json();
            console.log("new chat response data: ", data);
            return data;
        } catch(error) {
            console.log("new chat request failed => error: ", error)
            throw error
        }

    }

}

const newChatHandler = new NewChatHandler(POST, ChatElement, MessageElement)
newChatHandler.Create("Test Chat")
