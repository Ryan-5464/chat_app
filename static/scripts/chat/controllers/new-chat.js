function addNewChatEventListenerToNewChatInput() {
    const newChatInput = ConfigureNewChatInput();
    newChatInput.addEventListener('keydown', (e) => {
        if (e.key !== 'Enter' || !newChatInput.value.trim()) return;
        console.log("new chat name: ", newChatInput.value)
        newChatInput.__controller.NewChat(newChatInput.value);
        newChatInput.value = ''
    });
};

function ConfigureNewChatInput () {
    const newChatInput = document.getElementById('new-chat-input');
    newChatInput.__controller = {
        NewChat: (input) => NewChatHandler(input),
    };
    return newChatInput;
};

