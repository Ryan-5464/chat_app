function addChatModalListenerToChatContainer() {
    const modal = ConfigureChatModal();
    const container = document.getElementById('chats-container');
    const configureEditButton = ConfigureEditButton(modal.__controller);
    const configureLeaveButton = ConfigureLeaveButton(modal.__controller);
    
    container.addEventListener("contextmenu", (e) => {
        e.preventDefault();
        e.stopPropagation();
        const chatId = e.target.closest('[data-chatid]')?.getAttribute('data-chatid');
        if (!chatId) return;
        modal.__controller.OpenAt(e.clientX, e.clientY);
        configureEditButton(chatId);
        configureLeaveButton(chatId);
    });
};

function ConfigureEditButton(chatModalController) {
    const editNameButton = document.getElementById('chat-edit-btn');
    let currentChatId = null;
    editNameButton.addEventListener('click', () => {
        if (!currentChatId) return;
        chatModalController.EditChatName(currentChatId);
    });
    return (chatId) => { currentChatId = chatId };
};

function ConfigureLeaveButton(chatModalController) {
    const leaveButton = document.getElementById('chat-leave-btn');
    let currentChatId = null;
    leaveButton.addEventListener('click', () => {
        if (!currentChatId) return;
        chatModalController.LeaveChat(currentChatId);
    });
    return (chatId) => { currentChatId = chatId };
};

function ConfigureChatModal() {
    const modal = document.getElementById('chatModal');
    modal.__controller = {
        Close: () => CloseModal(modal),
        OpenAt: (clientX, clientY) => OpenModalAt(modal, clientX, clientY),
        EditChatName: (chatId) => EditChatName(chatId, () => CloseModal(modal)),
        LeaveChat: (chatId) => LeaveChat(chatId, () => CloseModal(modal)),
    };
    return modal;
};

function EditChatName(chatId, closeModal) {
    const openInput = document.getElementById('chat-name-input')
    if (openInput) {
        const openInputChatName = document.createElement('div')
        openInputChatName.classList.add('chat-name')
        openInputChatName.innerHTML = openInput.__oldtext
        openInput.replaceWith(openInputChatName)
    }
    if (openInput) { openInput.remove()}
    const chat = document.querySelector(`[data-chatid="${chatId}"]`);
    const chatName = chat.querySelector('.chat-name');
    const input = replaceWithInput(chatName, "Enter new name", 'chat-name-input');
    input.focus();
    closeModal();
    input.addEventListener('keydown', (e) => {
        if (e.key === "Enter") {
            EditChatNameHandler(input.value, chatId);
        };
    });
};

function LeaveChat(chatId, closeModal) {
    const activeChat = document.querySelector('.active');
    const activeChatId = activeChat?.getAttribute("data-chatid");
    const isActive = chatId === activeChatId;
    closeModal();
    console.log("is active", isActive)
    LeaveChatHandler(chatId, isActive);
};