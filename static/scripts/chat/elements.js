function ChatElement(chat) {
    const chatElem = document.createElement('div')
    chatElem.classList.add('chat', 'sidebar-elem')
    chatElem.setAttribute('data-chatid', chat.Id)
    chatElem.setAttribute('data-adminid', chat.AdminId)
    chatElem.setAttribute('data-createdat', chat.CreatedAt)

    const chatHeader = document.createElement('div')
    chatHeader.classList.add('chat-header')
    chatElem.appendChild(chatHeader)

    const chatName = document.createElement('div')
    chatName.id = `chat${chat.Id}`
    chatName.classList.add('chat-name')
    chatName.innerHTML = chat.Name
    chatHeader.appendChild(chatName)

    const adminName = document.createElement('div')
    adminName.classList.add('chat-admin-name')
    adminName.innerHTML = chat.AdminName
    chatHeader.appendChild(adminName)

    const unreadMessageCount = document.createElement('div')
    unreadMessageCount.classList.add('chat-unread-message-count')
    unreadMessageCount.innerHTML = chat.UnreadMessageCount
    chatHeader.appendChild(unreadMessageCount)

    const chatFooter = document.createElement('div')
    chatFooter.classList.add('chat-footer')
    chatElem.appendChild(chatFooter)

    return chatElem
}


function ContactElement(contact) {
    const contactElem = document.createElement('div')
    contactElem.classList.add('contact', 'sidebar-elem')
    contactElem.setAttribute('data-contactid', contact.Id)
    contactElem.setAttribute('data-contactchatid', contact.ContactChatId)

    const contactHeader = document.createElement('div')
    contactHeader.classList.add('contact-header')
    contactElem.appendChild(contactHeader)

    const name = document.createElement('div')
    name.classList.add('contact-name')
    name.innerHTML = contact.Name
    contactHeader.appendChild(name)

    const email = document.createElement('div')
    email.classList.add('contact-email')
    email.innerHTML = contact.Email
    contactHeader.appendChild(email)

    const messageFooter = document.createElement('div')
    messageFooter.classList.add('contact-footer')
    contactElem.appendChild(messageFooter)

    const knownSince = document.createElement('div')
    knownSince.classList.add('contact-since')
    knownSince.innerHTML = contact.KnownSince
    messageFooter.appendChild(knownSince)

    const onlineStatus = document.createElement('div')
    onlineStatus.classList.add('contact-status')
    onlineStatus.innerHTML = contact.OnlineStatus
    contactHeader.appendChild(onlineStatus)

    return contactElem

}

function MessageElement(message) {
    const messageElem = document.createElement('div')
    messageElem.classList.add('message')
    console.log("is user message", message.IsUserMessage)
    if (message.IsUserMessage) {
        messageElem.classList.add('me')
    }
    messageElem.setAttribute('data-messageid', message.Id)
    messageElem.setAttribute('data-chatid', message.ChatId)
    messageElem.setAttribute('data-userid', message.UserId)
    messageElem.setAttribute('data-replyid', message.ReplyId)

    const messageHeader = document.createElement('div')
    messageHeader.classList.add('message-header')
    messageElem.appendChild(messageHeader)

    const author = document.createElement('div')
    author.classList.add('message-author')
    author.innerHTML = message.Author
    messageHeader.appendChild(author)

    const formatCreatedAt = new Date(message.CreatedAt)
    const formatLastEditAt = new Date(message.LastEditAt)

    if (formatCreatedAt < formatLastEditAt) {
        const lastEditAt = document.createElement('div')
        lastEditAt.classList.add('message-lasteditat')
        lastEditAt.innerHTML = `Edited: ${formatLastEditAt.toLocaleString()}`
        messageHeader.appendChild(lastEditAt)
    } else {
        const createdAt = document.createElement('div')
        createdAt.classList.add('message-createdat')
        createdAt.innerHTML = `Sent: ${formatCreatedAt.toLocaleString()}`
        messageHeader.appendChild(createdAt)
    }

    const messgeFooter = document.createElement('div')
    messgeFooter.classList.add('message-footer')
    messageElem.appendChild(messgeFooter)

    const messageText = document.createElement('div')
    messageText.classList.add('message-text')
    messageText.innerHTML = message.Text
    messageHeader.appendChild(messageText)

    return messageElem

}

function ChatNameElement(name) {
    const chatNameElem = document.createElement('div');
    chatNameElem.classList.add('chat-name');
    chatNameElem.innerHTML = name;
    return chatNameElem;
};