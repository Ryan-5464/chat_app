window.addEventListener("click", function (e) {
  document.querySelectorAll(".modal").forEach(modal => {
    const modalContent = modal.querySelector(".modal-content");
    if (modal.classList.contains("open") && !modalContent.contains(e.target)) {
      closeModal(modal, modalContent);
    }
  });
});

document.addEventListener("DOMContentLoaded", function () {
    const messageContainer = document.querySelector("#messages-container");
    AttachDelegatedContextModal(messageContainer, ConfigureMessageModal, '.message')

    const chatsContainer = document.querySelector("#chats-container");
    AttachDelegatedContextModal(chatsContainer, ConfigureChatsModal, '.chats')

    const contactsContainer = document.querySelector("#contacts-container");
    AttachDelegatedContextModal(contactsContainer, ConfigureContactsModal, '.contacts')

});

function AttachDelegatedContextModal(container, modalConfig, targetSelector) {
    container.addEventListener("contextmenu", function (e) {
        
        const elem = e.target.closest(targetSelector)

        if (elem) {
            e.preventDefault()
            const modal = modalConfig(elem)
            modal.openAt(e.clientX, e.clientY)
        }
    })
}

function ConfigureMessageModal(elem) {
    const messageModal = document.getElementById("messageModal");
    
    const data = {
        ChatId: elem.getAttribute("data-chatid"),
        UserId: elem.getAttribute("data-userid"),
        MessageId: elem.getAttribute("data-messageid"),
    }

    return ConfigureModal(messageModal, data)
}

function ConfigureChatsModal(elem) {
    const chatModal = document.getElementById("chatModal");

    const data = {
        ChatId: elem.getAttribute("data-chatid")
    }

    let buttonConfig = function(modal) {
        editNameButton(modal, this.ChatId)
    }
    buttonConfig = buttonConfig.bind(data)

    return ConfigureModal(chatModal, buttonConfig)
}

function ConfigureContactsModal() {
    const contactModal = document.getElementById("contactModal");
    return ConfigureModal(contactModal)
}

function ConfigureModal(modal, buttonConfig) {
  if (!modal.__controller) {
    modal.__controller = new ModalController(modal);
  }

  if (typeof buttonConfig === "function") {
    buttonConfig(modal);
  }

  return modal.__controller;
}

function editNameButton(modal, chatId) {
    const chatEditBtn = document.getElementById('chat-edit-btn');
    const chat = document.querySelector(`[data-chatid="${chatId}"]`);
    
    chatEditBtn.addEventListener('click', () => {

        const chatName = chat.querySelector('.chat-name');
        const replaceNameWithInput = replaceWithInput.bind(chatName)
        const input = replaceNameWithInput("Enter new name")
        input.focus();
        modal.close()

        input.addEventListener('keydown', (e) => {
            e.preventDefault();
            if (e.key === ENTER) {
                chatName.innerHTML = EditChatNameRequest(input.value, chatId)
            }
        }, { once: true })   
    })
}

function replaceWithInput(placeholder) {
    this.innerHTML = ''
    const input = document.createElement('input');
    input.type = 'text';
    input.name = 'Name';
    input.placeholder = placeholder;
    input.required = true;
    input.className = 'input-elem';
    this.appendChild(input);
    return input
}

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
    }, 300);
  }

  openAt(clientX, clientY) {
    document.querySelectorAll(".modal.open").forEach(openModal => {
      if (openModal !== this.modal) {
        if (openModal.__controller) {
          openModal.__controller.close();
        }
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
}

class ChatModalController extends ModalController {
    constructor(modal, chatId) {
        super(modal)
        buttonConfig(chatId)
    }

    buttonConfig(chatId) {
        editNameButton(chatId)
    }

    editNameButton(chatId) {
        const chatEditBtn = document.getElementById('chat-edit-btn');
        const chat = document.querySelector(`[data-chatid="${chatId}"]`);
        
        chatEditBtn.addEventListener('click', () => {

            const chatName = chat.querySelector('.chat-name');
            const replaceNameWithInput = replaceWithInput.bind(chatName)
            const input = replaceNameWithInput("Enter new name")
            input.focus();
            this.modal.close()

            input.addEventListener('keydown', (e) => {
                e.preventDefault();
                if (e.key === ENTER) {
                    chatName.innerHTML = EditChatNameRequest(input.value, chatId)
                }
            }, { once: true })   
        })
    }
}