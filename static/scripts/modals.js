document.addEventListener("DOMContentLoaded", function () {
    addMessageModelEventListener()
    addChatModelEventListener()
    addContactModelEventListener()
});

function addChatModelEventListener() {
    const chats = document.querySelectorAll('.chat');
    const modal = document.getElementById("chatModal");
    const modalContent = modal.querySelector(".modal-content");

    chats.forEach((chat) => {
        chat.addEventListener("contextmenu", function (e) {
            e.preventDefault();
            openModalAt(e.clientX, e.clientY, modal, modalContent);
        });
    });
}

function addMessageModelEventListener() {
    const messages = document.querySelectorAll('.message');
    const modal = document.getElementById("messageModal");
    const modalContent = modal.querySelector(".modal-content");

    messages.forEach((message) => {
        message.addEventListener("contextmenu", function (e) {
            e.preventDefault();
            openModalAt(e.clientX, e.clientY, modal, modalContent);
        });
    });
}

function addContactModelEventListener() {
    const chats = document.querySelectorAll('.contact');
    const modal = document.getElementById("contactModal");
    const modalContent = modal.querySelector(".modal-content");

    chats.forEach((chat) => {
        chat.addEventListener("contextmenu", function (e) {
            e.preventDefault();
            openModalAt(e.clientX, e.clientY, modal, modalContent);
        });
        addEditChatEventListener(chat)
    });
}

function addEditChatEventListener(chat) {
    const modalOptions = chat.querySelector('#chatModel-options')
    const editButton = modalOptions.querySelector('#edit-btn')
    const leaveButton = modalOptions.querySelector('#leave-btn')
    const membersButton = modalOptions.querySelector('#members-btn')

    editButton.addEventListener("click", function(e) {
        e.preventDefault()
        const chatName = chat.querySelector('#chat-name')
        chatName.innerHTML = ""
        const newChatNameInput = document.createElement('input')
        newChatNameInput.classList.add('input-elem')
        newChatNameInput.type('text')
        newChatNameInput.placeholder('Type new name.')
        chatName.appendChild(newChatNameInput)
    })

}

function addLeaveChatEventListener(chat) {

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