function ChatElement(chat) {
    const chatElem = CreateElement({
        classes:[APP.CLS.CHAT.TAG, APP.CLS.SIDEBAR_ELEM], 
        data:{
            [APP.DATA.CHAT.ID]: contact.Id,
            [APP.DATA.CHAT.ADMINID]: contact.AdminId,
            [APP.DATA.CREATED]: contact.CreatedAt,
        },
    });
    if (chat.UserIsAdmin) {
        chatElem.classList.add(APP.CLS.ME);
    };
    const chatHeader = CreateElement({classes:[APP.CLS.CHAT.HEADER]});
    const chatFooter = CreateElement({classes:[APP.CLS.CHAT.FOOTER]});
    const chatName = CreateElement({id: `chat${chat.Id}`, classes:[APP.CLS.CHAT.NAME], innerHTML:chat.Name});
    const unreadMessageCount = CreateElement({classes:[APP.CLS.CHAT.UNREAD_MSG_CNT], innerHTML:chat.UnreadMessageCount});
    if (chat.UnreadMessageCount == 0) {
        unreadMessageCount.classList.add(APP.CLS.HIDDEN);
    };
    chatElem.appendChild(chatHeader, chatFooter);
    chatHeader.appendChild(chatName, unreadMessageCount);
    return chatElem;
};


function ContactElement(contact) {
    const contactElem = CreateElement({
        classes:[APP.CLS.CONTACT.TAG, APP.CLS.SIDEBAR_ELEM], 
        data:{
            [APP.DATA.CONTACT.ID]: contact.Id,
            [APP.DATA.CONTACT.CHATID]: contact.ContactChatId,
        },
    });
    const contactHeader = CreateElement({classes:[APP.CLS.CONTACT.HEADER]});
    const contactFooter = CreateElement({classes:[APP.CLS.CONTACT.FOOTER]});
    const name = CreateElement({classes:[APP.CLS.CONTACT.NAME], innerHTML:contact.Name});
    const email = CreateElement({classes:[APP.CLS.CONTACT.EMAIL], innerHTML:contact.Email});
    const onlineStatus = CreateElement({classes:[APP.CLS.CONTACT.STATUS], innerHTML:contact.OnlineStatus});
    contactElem.appendChild(contactHeader, contactFooter);
    contactHeader.appendChild(name, email, onlineStatus);
    return contactElem;
};

function MessageElement(message) {
    const messageElem = CreateElement({
        classes:[APP.CLS.MESSAGE.TAG], 
        data:{
            [APP.DATA.MESSAGE.ID]: message.Id,
            [APP.DATA.CHAT.ID]: message.ChatId,
            [APP.DATA.USER.ID]: message.UserId,
        },
    });
    if (message.IsUserMessage) {
        messageElem.classList.add(APP.CLS.ME);
    };
    const messageHeader = CreateElement({classes:[APP.CLS.MESSAGE.HEADER]});
    const messageFooter = CreateElement({classes:[APP.CLS.MESSAGE.FOOTER]});
    const author = CreateElement({classes:[APP.CLS.MESSAGE.AUTHOR], innerHTML:message.Author});
    const messageText = CreateElement({classes:[APP.CLS.MESSAGE.TEXT], innerHTML:message.Text});
    const createdAt = new Date(message.CreatedAt);
    const lastEditAt = new Date(message.LastEditAt);
    let date;
    if (createdAt < lastEditAt) {
        date = CreateElement({classes:[APP.CLS.MESSAGE.LAST_EDIT], innerHTML:`Edited: ${FormatDate(message.LastEditAt)}`});
    } else {
        date = CreateElement({classes:[APP.CLS.MESSAGE.CREATED], innerHTML:`Sent: ${FormatDate(message.CreatedAt)}`});
    };
    messageElem.appendChild(messageHeader, messageFooter);
    messageHeader.appendChild(author, date, messageText);
    return messageElem;
};

function MemberElement(member) {
    const memberElem = CreateElement({
        classes:[APP.CLS.MEMBER.TAG], 
        data:{
            [APP.DATA.CHAT.ID]:member.ChatId,
            [APP.DATA.USER.ID]:member.UserId,
        },
    });
    const name = CreateElement({classes:[APP.CLS.MEMBER.NAME], innerHTML:member.Name})
    const email = CreateElement({classes:[APP.CLS.MEMBER.EMAIL], innerHTML:member.Email})
    memberElem.appendChild(name, email);
    return memberElem;
};

function ChatNameElement(name) {
    return CreateElement({classes:[APP.CLS.CHAT.NAME], innerHTML:name});
};

function MessageTextElement(text) {
    return CreateElement({classes:[APP.CLS.MESSAGE.TEXT], innerHTML:text});
};

function CreateElement(
    {
        elemType='div', 
        id='', 
        classes=[], 
        innerHTML='', 
        data={}
    } = {}) {
    const elem = document.createElement(elemType);
    elem.id = id;
    elem.classList.add(...classes);
    elem.innerHTML = innerHTML;
    for (const [key, value] of Object.entries(data)) {
        elem.dataset[key] = value;
    };
    return elem;
};

function FormatDate(date) {
    const formatDate = new Date(date);
    return `${formatDate.toLocaleString()}`;
};