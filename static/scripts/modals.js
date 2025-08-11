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
        }
    });
  }

  if (messageContainer) {
    const messageModal = document.getElementById("messageModal");
    const messageModalContent = messageModal.querySelector(".modal-content");

    messageContainer.addEventListener("contextmenu", function (e) {
      if (e.target.closest(".message")) {
        e.preventDefault();
        openModalAt(e.clientX, e.clientY, messageModal, messageModalContent);
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

function setChatModalData(chatId) {
  document.getElementById('chat-leave-btn').dataset.chatid = chatId;
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
