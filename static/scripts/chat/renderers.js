class Renderer {
    activate() {

    }

    _render(container, elemFactory, data, overwrite) {
        if (overwrite == true) {
            container.innerHTML = ''
        }
        
        Object.values(data).forEach(obj => {
            container.appendChild(elemFactory(obj))
        })
    }

}

class GroupChatRenderer extends Renderer {
    render(chats, overwrite) {
        const chatContainer = document.getElementById('chats-container')
        this._render(chatContainer, ChatElement, chats, overwrite)
    }
}

class ContactChatRenderer extends Renderer {
    render(contacts, overwrite) {
        const contactsContainer = document.getElementById('contacts-container')
        this._render(contactsContainer, ContactElement, contacts, overwrite)
    }
}

class ChatMessageRenderer extends Renderer {
    render(messages, overwrite) {
        const messagesContainer = document.getElementById('messages-container')
        this._render(messagesContainer, MessageElement, messages, overwrite)
    }
}
