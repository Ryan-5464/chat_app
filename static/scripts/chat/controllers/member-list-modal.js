function ConfigureMemberModal() {
    const modal = document.getElementById('chatMemberModal');
    const content = document.getElementById('memberModal-content');

    modal.__controller = {
        Open: () => {
            modal.style.display = "flex"; // enable flex centering here
            content.classList.add("opening");
            content.classList.remove("closing");
        },
        Close: () => {
            content.classList.remove("opening");
            content.classList.add("closing");
            setTimeout(() => {
                modal.style.display = "none"; // hide again
            }, 300);
        },
        AddMemberToChat: (email, chatId) => AddMemberToChat(email, chatId)
    } 

    // 🟢 Close when clicking the backdrop (but not the modal content)
    modal.addEventListener("click", (e) => {
        if (e.target === modal) {
            modal.__controller.Close();
        }
    });

    const addMemberInput = document.getElementById('add-member-input')
    addMemberInput.addEventListener('keydown', (e) => {
        if (e.key !== 'Enter') return;
        e.preventDefault()
        const chatId = addMemberInput.getAttribute('data-chatid')
        modal.__controller.AddMemberToChat(addMemberInput.value, chatId)
        addMemberInput.value = ''
    })

    return modal;
}

function AddMemberToChat(email, chatId) {
    AddMemberToChatHandler(email, chatId);
}