function CloseModal (modal) {
    const modalContent = modal.querySelector('.modal-content')
    modalContent.classList.remove("opening");
    modalContent.classList.add("closing");

    setTimeout(() => {
        modal.classList.remove("open");
        modalContent.classList.remove("closing");
    }, MODAL_CLOSE_DELAY);
}

function OpenModalAt(modal, clientX, clientY) {
    const modalContent = modal.querySelector('.modal-content');
    console.log("open modals: ", document.querySelectorAll(".modal.open"))
    // Close other open modals that have a controller
    document.querySelectorAll(".modal.open").forEach(openModal => {
        if (openModal !== modal && openModal.__controller) {
            openModal.__controller.Close();
        }
    });

    const padding = 10;
    const maxLeft = window.innerWidth - modalContent.offsetWidth - padding;
    const maxTop = window.innerHeight - modalContent.offsetHeight - padding;

    modalContent.style.left = Math.min(clientX, maxLeft) + "px";
    modalContent.style.top = Math.min(clientY, maxTop) + "px";

    modalContent.classList.remove("opening", "closing");
    void modalContent.offsetWidth; // force reflow

    modal.classList.add("open");
    modalContent.classList.add("opening");
};



// class ModalController {
//     constructor(modal) {
//         this.modal = modal;
//         this.modalContent = modal.querySelector('.modal-content');
//     }

//     close() {
//         this.modalContent.classList.remove("opening");
//         this.modalContent.classList.add("closing");

//         setTimeout(() => {
//         this.modal.classList.remove("open");
//         this.modalContent.classList.remove("closing");
//         }, MODAL_CLOSE_DELAY);
//     }

//     openAt(clientX, clientY) {
//         document.querySelectorAll(".modal.open").forEach(openModal => {
//             if (openModal !== this.modal && openModal.__controller) {
//                 openModal.__controller.close();
//             }
//         });

//         const padding = 10;
//         const maxLeft = window.innerWidth - this.modalContent.offsetWidth - padding;
//         const maxTop = window.innerHeight - this.modalContent.offsetHeight - padding;

//         this.modalContent.style.left = Math.min(clientX, maxLeft) + "px";
//         this.modalContent.style.top = Math.min(clientY, maxTop) + "px";

//         this.modalContent.classList.remove("opening", "closing");
//         void this.modalContent.offsetWidth;

//         this.modal.classList.add("open");
//         this.modalContent.classList.add("opening");
//     }

//     configureButtons(buttonData = {}) {
//         const buttons = this.modalContent.querySelectorAll('[data-action]');

//         buttons.forEach(button => {
//             const actionName = button.dataset.action;
//             button = RemoveAllListeners(button)

//             button.addEventListener('click', (e) => {
//                 this[actionName](e, buttonData)
//             });
//         })
//     }
// }



// class ChatModalController extends ModalController {
//     constructor(modal) {
//         super(modal)
//         this.renderer = new Renderer()
//     }

//     editName(e, { ChatId }) {
//         const chat = document.querySelector(`[data-chatid="${ChatId}"]`);
//         const chatName = chat.querySelector('.chat-name');
//         const input = replaceWithInput(chatName, "Enter new name");
//         input.focus();
//         this.close();

//         input.addEventListener('keydown', (e) => {
//             if (e.key === "Enter") {
//                 EditChatNameRequest(input.value, ChatId).then(data => {
//                     if (data) {
//                         chatName.innerHTML = data.Name;
//                     }
//                 });
//             }
//         })
//     }

//     leaveChat(e, { ChatId }) {
//         this.close()
        
//         let isActive = false
//         if (e.target.classList.contains('active')) {
//             isActive = true
//         }

//         LeaveChatRequest(ChatId).then(data => {
//             if (!isActive) {
//                 return
//             }
//             this.renderer.render('chats', data.Chats, true)
//             this.renderer.render('messages', data.Messages, true)  
//             const chat = document.querySelector(`[data-chatid="${data.NewActiveChatId}"]`);
//             chat.classList.add('active')
//         })
//     }
// }

// class MessageModalController extends ModalController {
//     constructor(modal) {
//         super(modal)
//     }

//     _deleteMessage(e, { ChatId, MessageId, UserId }) {
//         return ChatId, MessageId, UserId
//     }
// }

// class ContactsModalController extends ModalController {
//     constructor(modal) {
//         super(modal)
//     }

//     configureButtons() {
//         return
//     }
// }




// CHAT CONTROLLERS ============================================================

// const chatControllerRegistry = {
//     Chat: {
//         containerId: 'chats-container',
//         elemClass: '.chat',
//         dataId: 'chatid',
//         elemType: 'Chat'
//     },
//     Contact: {
//         containerId: 'contacts-container',
//         elemClass: '.contact',
//         dataId: 'contactchatid',
//         elemType: 'Contact'
//     }
// }

// function AttachChatControllers() {
//     for (const config of Object.values(chatControllerRegistry)) {
//         const container = document.getElementById(config.containerId) 
//         container.__controller = new ChatController(config)
//         container.addEventListener('click', (e) => {
//             const chatElem = e.target.closest(container.__controller.config.elemClass)
//             if (chatElem) {
//                 container.__controller.handleChatActivation(chatElem)
//             }
//         })
//     }
// }

// class ChatController {
//     constructor(config) {
//         this.renderer = new Renderer()
//         this.config = config
//         this.container = document.getElementById(config.containerId)
//     }

//     handleChatActivation(target) {
//         const chatId = target.dataset[this.config.dataId]
//         if (!chatId) {return}

//         this._setActiveElemenById(chatId)
//         this._loadMessages(chatId)
//     }

//     _setActiveElemenById(targetId) {
//         const elems = this.container.querySelectorAll(this.config.elemClass)
//         for (const chat of elems) {
//             const isTarget = chat.dataset[this.config.dataId] === targetId
//             chat.classList.toggle('active', isTarget)
//         }
//     }

//     _loadMessages(targetId) {
//         SwitchChatRequest(this.config.elemType, targetId)
//         .then(data => {
//             this.renderer.render(data.Messages, true)
//         })
//         .catch(err => {
//             console.log("Failed to handle chat activation: ", err)
//         })
//     }
// }

// INPUT CONTROLLERS ===============================================================================



