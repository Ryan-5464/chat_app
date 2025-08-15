async function EditChatNameRequest(name, chatId) {
    return fetch(EDIT_CHAT_NAME_ENDPOINT, POST({ Name: name, ChatId: chatId }))
        .then(response => {
            if (!response.ok) throw new Error("Network response was not ok");
            return response.json();
        })
        .then(data => {
            return data.Name
        })
        .catch(error => {
            console.log('Error:', error);
        });
}