
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
                EditChatName(input.value)
            }
        })
    }

    leaveChat(e, { ChatId }) {
        this.close()
        
        let isActive = false
        if (e.target.classList.contains('active')) {
            isActive = true
        }

        LeaveChat(ChatId, isActive)
    }
}

class MessageModalController extends ModalController {
    constructor(modal) {
        super(modal)
    }

    deleteMessage(e, { ChatId, MessageId, UserId }) {
        DeleteMessage(ChatId, MessageId, UserId)
    }

    editMessage() {}
}

class ContactsModalController extends ModalController {
    constructor(modal) {
        super(modal)
    }

    removeContact() {}
}




