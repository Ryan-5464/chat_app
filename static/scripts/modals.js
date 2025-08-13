// document.addEventListener("DOMContentLoaded", function () {
//     addMessageModelEventListener()
//     addChatModelEventListener()
//     addContactModelEventListener()
// });

// function addChatModelEventListener() {
//     const chats = document.querySelectorAll('.chat');
//     const modal = document.getElementById("chatModal");
//     const modalContent = modal.querySelector(".modal-content");

//     chats.forEach((chat) => {
//         chat.addEventListener("contextmenu", function (e) {
//             e.preventDefault();
//             openModalAt(e.clientX, e.clientY, modal, modalContent);
//         });
//     });
// }

// function addMessageModelEventListener() {
//     const messages = document.querySelectorAll('.message');
//     const modal = document.getElementById("messageModal");
//     const modalContent = modal.querySelector(".modal-content");

//     messages.forEach((message) => {
//         message.addEventListener("contextmenu", function (e) {
//             e.preventDefault();
//             openModalAt(e.clientX, e.clientY, modal, modalContent);
//         });
//     });
// }

// function addContactModelEventListener() {
//     const chats = document.querySelectorAll('.contact');
//     const modal = document.getElementById("contactModal");
//     const modalContent = modal.querySelector(".modal-content");

//     chats.forEach((chat) => {
//         chat.addEventListener("contextmenu", function (e) {
//             e.preventDefault();
//             openModalAt(e.clientX, e.clientY, modal, modalContent);
//         });
//     });
// }

// function openModalAt(x, y, modal, modalContent) {
//     document.querySelectorAll(".modal.open").forEach(openModal => {
//         const openContent = openModal.querySelector(".modal-content");
//         if (openModal !== modal) {
//             closeModal(openModal, openContent);
//         }
//     });

//     const padding = 10;
//     const maxLeft = window.innerWidth - modalContent.offsetWidth - padding;
//     const maxTop = window.innerHeight - modalContent.offsetHeight - padding;

//     const left = Math.min(x, maxLeft);
//     const top = Math.min(y, maxTop);

//     modalContent.style.left = left + "px";
//     modalContent.style.top = top + "px";

//     modalContent.classList.remove("opening", "closing");
//     void modalContent.offsetWidth;

//     modal.classList.add("open");
//     modalContent.classList.add("opening");
// }

// function closeModal(modal, modalContent) {
//     modalContent.classList.remove("opening");
//     modalContent.classList.add("closing");

//     setTimeout(() => {
//         modal.classList.remove("open");
//         modalContent.classList.remove("closing");
//     }, 300); 
// }

// window.addEventListener("click", function (e) {
//     document.querySelectorAll(".modal").forEach(modal => {
//         const modalContent = modal.querySelector(".modal-content");
//         if (modal.classList.contains("open") && !modalContent.contains(e.target)) {
//             closeModal(modal, modalContent);
//         }
//     });
// });

document.addEventListener("DOMContentLoaded", function () {
    addDelegatedEventListeners();
    AttachMsgModalToMsgContainer()
});

function AttachMsgModalToMsgContainer() {
    const messageContainer = document.querySelector("#messages-container");

    messageContainer.addEventListener("contextmenu", function (e) {

        const elem = e.target.closest('.message')

        if (elem) {
            e.preventDefault();
            ConfigureMessageModal()
        }
    
    });
}

function ConfigureMessageModal() {
    const messageModal = document.getElementById("messageModal");
    
    const data = {
        ChatId: elem.getAttribute("data-chatid"),
        UserId: elem.getAttribute("data-userid"),
        MessageId: elem.getAttribute("data-messageid"),
    }

    messageModal = ConfigureModal(messageModal, data)
    messageModal.openAt(e.clientX, e.clientY)
}

function ConfigureChatModal() {
    const chatContainer = document.querySelector("#chats-container");

    chatContainer.addEventListener("contextmenu", function (e) {
        
        const elem = e.target.closest(".chat")
    
        if (elem) {
            e.preventDefault();
            const chatModal = document.getElementById("chatModal");

            const data = {
                ChatId: elem.getAttribute("data-chatid")
            }

            chatModal = ConfigureModal(chatModal, data)
            chatModal.openAt(e.clientX, e.clientY)
            setupChatModalEditButton(chatModal.dataset.ChatId, chatModal, chatModal.modalContent)
        }

    });
}

function ConfigureContactModal() {
    const contactsContainer = document.querySelector('#contacts-container')

    contactsContainer.addEventListener("contextmenu", function (e) {

        const elem = e.target.closest(".contact")

        if (elem) {
            e.preventDefault();
            const contactModal = document.getElementById("contactModal");
            
            contactModal = ConfigureModal(contactModal)
            contactModal.openAt(e.clientX, e.clientY)
        }
    })
}

function ConfigureModal(modal, data={}) {
    const newModal = {
        modal: modal,
        modalContent: modal.querySelector('.modal-content'),
        data: data,

        close() {
            this.modalContent.classList.remove("opening");
            this.modalContent.classList.add("closing");

            setTimeout(() => {
                this.modal.classList.remove("open");
                this.modalContent.classList.remove("closing");
            }, 300);
        },

        openAt(clientX, clientY) {
            document.querySelectorAll(".modal.open").forEach(openModal => {
                if (openModal !== this.modal) {
                    ConfigureModal(openModal).close();
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
    };

    return newModal;
}

function setupChatModalEditButton(chatId, modal, modalContent) {
  const chatEditBtn = document.getElementById('chat-edit-btn');
  chatEditBtn.addEventListener('click', () => {
    const chat = document.querySelector(`[data-chatid="${chatId}"]`);
    if (!chat) return;

    const chatName = chat.querySelector('.chat-name');
    chatName.innerHTML = ''; // clear existing content

    const form = document.createElement('form');
    form.classList.add('edit-chat-name-form');

    const input = document.createElement('input');
    input.type = 'text';
    input.name = 'Name';
    input.placeholder = 'Enter new name';
    input.required = true;
    input.className = 'input-elem';
    form.appendChild(input);

    // hidden input to send chatId
    const hiddenChatId = document.createElement('input');
    hiddenChatId.type = 'hidden';
    hiddenChatId.name = 'ChatId';
    hiddenChatId.value = chatId;
    form.appendChild(hiddenChatId);

    chatName.appendChild(form);

    input.focus();

    closeModal(modal, modalContent);

    form.addEventListener('submit', async (e) => {
      e.preventDefault();

      const name = form.elements['Name'].value;
      const chatId = form.elements['ChatId'].value;

      try {
        const response = await fetch('/api/chat/edit', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'same-origin',
          body: JSON.stringify({ Name: name, ChatId: chatId }),
        });

        if (!response.ok) {
          alert('Failed to update chat name.');
          return;
        }

        const text = await response.text();
        chatName.innerHTML = text;

      } catch (err) {
        alert('Network error while updating chat name.');
      }
    });
  });
}

window.addEventListener("click", function (e) {
  document.querySelectorAll(".modal").forEach(modal => {
    const modalContent = modal.querySelector(".modal-content");
    if (modal.classList.contains("open") && !modalContent.contains(e.target)) {
      closeModal(modal, modalContent);
    }
  });
});
