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
});

function addDelegatedEventListeners() {
  const chatContainer = document.querySelector("#chats-container");
  const messageContainer = document.querySelector("#messages-container");
  const contactContainer = document.querySelector("#contacts-container");

  if (chatContainer) {
    const chatModal = document.getElementById("chatModal");
    const chatModalContent = chatModal.querySelector(".modal-content");

    chatContainer.addEventListener("contextmenu", function (e) {
        const elem = e.target.closest(".chat")
        if (e.target.closest(".chat")) {
            e.preventDefault();
            openModalAt(e.clientX, e.clientY, chatModal, chatModalContent);
            const chatId = elem.getAttribute("data-chatid")
            setChatModalData(chatId)
            setupChatModalEditButton(chatId, chatModal, chatModalContent)
        }
    });
  }

  if (messageContainer) {
    const messageModal = document.getElementById("messageModal");
    const messageModalContent = messageModal.querySelector(".modal-content");

    messageContainer.addEventListener("contextmenu", function (e) {
      const elem = e.target.closest('.message')
      if (e.target.closest(".message")) {
        e.preventDefault();
        openModalAt(e.clientX, e.clientY, messageModal, messageModalContent);
        const chatId = elem.getAttribute("data-chatid")
        const userId = elem.getAttribute("data-userid")
        const messageId = elem.getAttribute("data-messageid")
        setMessageModalData(messageId, chatId, userId)
      }
    });
  }

  if (contactContainer) {
    const contactModal = document.getElementById("contactModal");
    const contactModalContent = contactModal.querySelector(".modal-content");

    contactContainer.addEventListener("contextmenu", function (e) {
      if (e.target.closest(".contact")) {
        e.preventDefault();
        openModalAt(e.clientX, e.clientY, contactModal, contactModalContent);
      }
    });
  }
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


function setChatModalData(chatId) {
  document.getElementById('chat-leave-btn').dataset.chatid = chatId;
}

function setMessageModalData(messageId, chatId, userId) {
  document.getElementById('msg-del-btn').dataset.chatid = chatId;
  document.getElementById('msg-del-btn').dataset.userid = userId;
  document.getElementById('msg-del-btn').dataset.messageid = messageId;
}

function openModalAt(x, y, modal, modalContent) {
  document.querySelectorAll(".modal.open").forEach(openModal => {
    const openContent = openModal.querySelector(".modal-content");
    if (openModal !== modal) {
      closeModal(openModal, openContent);
    }
  });

  const padding = 10;
  const maxLeft = window.innerWidth - modalContent.offsetWidth - padding;
  const maxTop = window.innerHeight - modalContent.offsetHeight - padding;

  const left = Math.min(x, maxLeft);
  const top = Math.min(y, maxTop);

  modalContent.style.left = left + "px";
  modalContent.style.top = top + "px";

  modalContent.classList.remove("opening", "closing");
  void modalContent.offsetWidth;

  modal.classList.add("open");
  modalContent.classList.add("opening");
}

function closeModal(modal, modalContent) {
  modalContent.classList.remove("opening");
  modalContent.classList.add("closing");

  setTimeout(() => {
    modal.classList.remove("open");
    modalContent.classList.remove("closing");
  }, 300);
}

window.addEventListener("click", function (e) {
  document.querySelectorAll(".modal").forEach(modal => {
    const modalContent = modal.querySelector(".modal-content");
    if (modal.classList.contains("open") && !modalContent.contains(e.target)) {
      closeModal(modal, modalContent);
    }
  });
});
